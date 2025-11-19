package sAuth

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// ValidateTokenWithClaims valida un token y retorna los claims
func (s *service) ValidateTokenWithClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := s.validator.ValidateToken(tokenString)
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok {
		return *claims, nil
	} else if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, fmt.Errorf("could not extract claims from token")
}

// GetAuthURL genera la URL de autorización
func (s *service) GetAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// ExchangeCode intercambia el código por tokens
func (s *service) ExchangeCode(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.oauthConfig.Exchange(ctx, code)
}

// GetUserInfo obtiene la información del usuario desde Authentik
func (s *service) GetUserInfo(ctx context.Context, token *oauth2.Token) (*rModels.AuthentikUserInfo, error) {
	client := s.oauthConfig.Client(ctx, token)

	userInfoURL := s.config.AuthConfig.Issuer + "/application/o/userinfo/"
	resp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("userinfo endpoint returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var userInfo rModels.AuthentikUserInfo
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal user info: %w", err)
	}

	return &userInfo, nil
}

// SyncUser sincroniza el usuario con la tabla MrEmployee
func (s *service) SyncUser(userInfo *rModels.AuthentikUserInfo) (*rModels.MrEmployee, error) {
	employee, err := s.repository.SyncUser(userInfo)
	if err != nil {
		return nil, err
	}

	// LOG REDUCIDO: Solo información esencial
	s.logger.Info("Employee synchronized",
		zap.String("authentik_id", userInfo.Sub),
		zap.Uint("employee_id", employee.ID),
		zap.String("email", employee.Email))

	return employee, nil
}

// FindCurrentUserInfo - Buscar empleado por ID de Authentik
func (s *service) FindCurrentUserInfo(userID string, ctx context.Context) (*rModels.MrEmployee, error) {
	employee, _, err := s.GetEmployeeWithOrganization(userID, ctx)
	return employee, err
}

// Nuevo método que combina empleado BD + organización del JWT (sin caché)
func (s *service) GetEmployeeWithOrganization(userID string, ctx context.Context) (*rModels.MrEmployee, map[string]interface{}, error) {
	// Obtener empleado de la base de datos directamente (sin caché)
	employee, err := s.repository.FindCurrentUserInfo(userID, ctx)
	if err != nil {
		return nil, nil, err
	}

	// La organización viene del JWT, no de la BD
	organization := make(map[string]interface{})

	// LOG REDUCIDO: Solo en caso de error o para información importante
	if employee == nil {
		s.logger.Warn("Employee not found in database", zap.String("user_id", userID))
	}

	return employee, organization, nil
}

// ValidateToken valida un token usando el validador local
func (s *service) ValidateToken(tokenString string) (*jwt.Token, error) {
	return s.validator.ValidateToken(tokenString)
}

// GetValidator retorna el validador JWKS para operaciones administrativas
func (s *service) GetValidator() *JWKSValidator {
	return s.validator
}

// ForceRefreshJWKS fuerza la actualización de las claves JWKS
func (s *service) ForceRefreshJWKS() error {
	return s.validator.ForceRefresh()
}

// Métodos de gestión de tokens
func (s *service) AddToBlacklist(token string, expiresAt time.Time) {
	s.blacklist.Add(token, expiresAt)
	// LOG REDUCIDO: Solo información importante
	s.logger.Info("Token added to blacklist", zap.Time("expires_at", expiresAt))
}

func (s *service) IsTokenRevoked(token string) bool {
	return s.blacklist.IsRevoked(token)
}

func (s *service) ClearUserCache(userID string) {
	// Método mantenido por compatibilidad, pero no hace nada
	// LOG ELIMINADO: No es necesario loggear esta operación ahora
}

func (s *service) CleanupBlacklist() {
	s.blacklist.Cleanup()
}
