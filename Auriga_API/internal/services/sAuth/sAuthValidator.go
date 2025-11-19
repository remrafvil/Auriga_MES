package sAuth

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/remrafvil/Auriga_API/config"
	"go.uber.org/zap"
)

type JWKSValidator struct {
	jwksURL   string
	logger    *zap.Logger
	keys      map[string]*rsa.PublicKey
	keysMutex sync.RWMutex
	lastFetch time.Time
	cacheTTL  time.Duration
}

func NewJWKSValidator(config *config.Settings, logger *zap.Logger) *JWKSValidator {
	validator := &JWKSValidator{
		jwksURL:  config.AuthConfig.JWKSEndpoint,
		logger:   logger,
		keys:     make(map[string]*rsa.PublicKey),
		cacheTTL: config.AuthConfig.JWKSCacheTTL,
	}

	// Cargar claves al iniciar
	if err := validator.fetchKeys(); err != nil {
		logger.Warn("Failed to fetch initial JWKS keys", zap.Error(err))
	}

	return validator
}

// ValidateToken valida el token LOCALMENTE usando claves cacheadas
func (v *JWKSValidator) ValidateToken(tokenString string) (*jwt.Token, error) {
	// Verificar si necesitamos actualizar las claves
	if v.shouldRefreshKeys() {
		if err := v.fetchKeys(); err != nil {
			v.logger.Warn("Failed to refresh JWKS keys", zap.Error(err))
			// Continuar con claves cacheadas si hay error
		}
	}

	// Usar jwt.ParseWithClaims para obtener claims específicos
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, v.getKeyFunc())
	if err != nil {
		return nil, fmt.Errorf("token parsing failed: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	// Validar claims adicionales
	if claims, ok := token.Claims.(*jwt.MapClaims); ok {
		if err := v.validateClaims(*claims); err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("invalid token claims type")
	}

	return token, nil
}

// validateClaims valida los claims del token
func (v *JWKSValidator) validateClaims(claims jwt.MapClaims) error {
	// Verificar expiración usando jwt.RegisteredClaims para mejor manejo
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return fmt.Errorf("failed to get expiration time: %w", err)
	}

	if exp == nil || exp.Time.Before(time.Now()) {
		return fmt.Errorf("token expired")
	}

	// Verificar que tenga subject
	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return fmt.Errorf("sub claim is required")
	}

	// Verificar issuer (opcional, pero recomendado)
	if iss, ok := claims["iss"].(string); ok {
		v.logger.Debug("Token issuer", zap.String("issuer", iss))
		// Puedes agregar validación específica del issuer aquí si lo deseas
	}

	return nil
}

// getKeyFunc retorna la función para obtener la clave de firma
func (v *JWKSValidator) getKeyFunc() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		// Verificar el algoritmo
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		kid, ok := token.Header["kid"].(string)
		if !ok {
			return nil, fmt.Errorf("kid not found in token header")
		}

		v.keysMutex.RLock()
		defer v.keysMutex.RUnlock()

		key, exists := v.keys[kid]
		if !exists {
			return nil, fmt.Errorf("key not found for kid: %s", kid)
		}

		return key, nil
	}
}

// shouldRefreshKeys verifica si es necesario actualizar las claves JWKS
func (v *JWKSValidator) shouldRefreshKeys() bool {
	v.keysMutex.RLock()
	defer v.keysMutex.RUnlock()

	return time.Since(v.lastFetch) > v.cacheTTL || len(v.keys) == 0
}

// fetchKeys obtiene y actualiza las claves JWKS desde Authentik
func (v *JWKSValidator) fetchKeys() error {
	v.logger.Debug("Fetching JWKS from Authentik", zap.String("url", v.jwksURL))

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", v.jwksURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch JWKS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("JWKS endpoint returned status: %d", resp.StatusCode)
	}

	var jwks struct {
		Keys []struct {
			Kid string `json:"kid"`
			Kty string `json:"kty"`
			Use string `json:"use"`
			Alg string `json:"alg"`
			N   string `json:"n"`
			E   string `json:"e"`
		} `json:"keys"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&jwks); err != nil {
		return fmt.Errorf("failed to decode JWKS: %w", err)
	}

	newKeys := make(map[string]*rsa.PublicKey)

	for _, key := range jwks.Keys {
		if key.Kty == "RSA" && key.Use == "sig" {
			rsaKey, err := convertJWKToRSAPublicKey(key.N, key.E)
			if err != nil {
				v.logger.Warn("Failed to parse RSA key",
					zap.String("kid", key.Kid),
					zap.Error(err))
				continue
			}
			newKeys[key.Kid] = rsaKey
			v.logger.Debug("Added RSA key to cache", zap.String("kid", key.Kid))
		}
	}

	if len(newKeys) == 0 {
		return fmt.Errorf("no valid RSA keys found in JWKS")
	}

	v.keysMutex.Lock()
	defer v.keysMutex.Unlock()

	v.keys = newKeys
	v.lastFetch = time.Now()

	v.logger.Info("JWKS keys updated",
		zap.Int("key_count", len(v.keys)),
		zap.Time("last_fetch", v.lastFetch))
	return nil
}

// convertJWKToRSAPublicKey convierte una clave JWK a RSA Public Key
func convertJWKToRSAPublicKey(nStr, eStr string) (*rsa.PublicKey, error) {
	// Decodificar base64 URL encoding
	nBytes, err := base64.RawURLEncoding.DecodeString(nStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modulus: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(eStr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode exponent: %w", err)
	}

	modulus := new(big.Int).SetBytes(nBytes)
	exponent := new(big.Int).SetBytes(eBytes)

	if !exponent.IsInt64() {
		return nil, fmt.Errorf("exponent is too large")
	}

	return &rsa.PublicKey{
		N: modulus,
		E: int(exponent.Int64()),
	}, nil
}

// ForceRefresh fuerza la actualización de las claves JWKS
func (v *JWKSValidator) ForceRefresh() error {
	v.keysMutex.Lock()
	v.lastFetch = time.Time{} // Reset para forzar refresh
	v.keysMutex.Unlock()

	return v.fetchKeys()
}

// GetKeysCount retorna el número de claves cacheadas (para debugging)
func (v *JWKSValidator) GetKeysCount() int {
	v.keysMutex.RLock()
	defer v.keysMutex.RUnlock()
	return len(v.keys)
}
