package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/remrafvil/Auriga_API/internal/services/sAuth"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

const CacheTTL = 5 * time.Minute

type CachedValidationResult struct {
	Valid   bool
	Claims  jwt.MapClaims
	Expires time.Time
}

type AuthMiddleware struct {
	validator   *sAuth.JWKSValidator
	logger      *zap.Logger
	authService sAuth.Service
	tokenCache  map[string]CachedValidationResult
	cacheMutex  sync.RWMutex
	cacheTTL    time.Duration
}

type AuthMiddlewareParams struct {
	fx.In
	Validator   *sAuth.JWKSValidator
	Logger      *zap.Logger
	AuthService sAuth.Service
}

func NewAuthMiddleware(p AuthMiddlewareParams) *AuthMiddleware {
	return &AuthMiddleware{
		validator:   p.Validator,
		logger:      p.Logger,
		authService: p.AuthService,
		tokenCache:  make(map[string]CachedValidationResult),
		cacheTTL:    CacheTTL,
	}
}

func (m *AuthMiddleware) CombinedMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			m.logger.Debug("🔐 CombinedMiddleware executing",
				zap.String("path", c.Path()),
				zap.String("method", c.Request().Method))

			tokenString, source := m.extractToken(c)
			if tokenString == "" {
				m.logger.Debug("Authentication required",
					zap.String("path", c.Path()))
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			if m.authService.IsTokenRevoked(tokenString) {
				m.logger.Warn("Attempt to use revoked token",
					zap.String("path", c.Path()))
				return echo.NewHTTPError(http.StatusUnauthorized, "Token has been revoked")
			}

			if claims, found := m.getFromCache(tokenString); found {
				m.setUserContext(c, claims)
				return next(c)
			}

			token, err := m.validator.ValidateToken(tokenString)
			if err != nil {
				m.logger.Warn("Token validation failed",
					zap.Error(err),
					zap.String("path", c.Path()))
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token: "+err.Error())
			}

			// ✅ CORREGIDO: Manejar ambos tipos de claims
			var claims jwt.MapClaims
			if tokenClaims, ok := token.Claims.(*jwt.MapClaims); ok {
				claims = *tokenClaims
			} else if tokenClaims, ok := token.Claims.(jwt.MapClaims); ok {
				claims = tokenClaims
			} else {
				return echo.NewHTTPError(http.StatusUnauthorized, "Invalid token claims format")
			}

			m.setToCache(tokenString, claims)
			m.setUserContext(c, claims)

			if userID, ok := claims["sub"].(string); ok {
				m.logger.Debug("Token validated",
					zap.String("user_id", userID),
					zap.String("source", source),
					zap.String("path", c.Path()))
			}

			return next(c)
		}
	}
}

func (m *AuthMiddleware) extractToken(c echo.Context) (string, string) {
	// 1. Intentar obtener token del header Authorization
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.Split(authHeader, " ")
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1], "header"
		}
	}

	// 2. Si no hay header, buscar en cookies
	cookie, err := c.Cookie("auth_token")
	if err == nil && cookie.Value != "" {
		return cookie.Value, "cookie"
	}

	return "", ""
}

func (m *AuthMiddleware) getFromCache(tokenString string) (jwt.MapClaims, bool) {
	m.cacheMutex.RLock()
	defer m.cacheMutex.RUnlock()

	if result, exists := m.tokenCache[tokenString]; exists {
		if time.Now().Before(result.Expires) {
			return result.Claims, true
		}
		delete(m.tokenCache, tokenString)
	}
	return nil, false
}

func (m *AuthMiddleware) setToCache(tokenString string, claims jwt.MapClaims) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	m.tokenCache[tokenString] = CachedValidationResult{
		Valid:   true,
		Claims:  claims,
		Expires: time.Now().Add(m.cacheTTL),
	}

	if len(m.tokenCache) > 1000 {
		m.cleanupCache()
	}
}

func (m *AuthMiddleware) cleanupCache() {
	now := time.Now()
	for token, result := range m.tokenCache {
		if now.After(result.Expires) {
			delete(m.tokenCache, token)
		}
	}
}

func (m *AuthMiddleware) setUserContext(c echo.Context, claims jwt.MapClaims) {
	// 1. Poner en Echo context (para compatibilidad)
	c.Set("user_claims", claims)

	// 2. Crear nuevo contexto estándar con toda la información
	ctx := context.Background()
	ctx = context.WithValue(ctx, "user_claims", claims)

	// Información básica del usuario
	if sub, ok := claims["sub"].(string); ok {
		ctx = context.WithValue(ctx, "user_id", sub)
		c.Set("user_id", sub)
	}
	if email, ok := claims["email"].(string); ok {
		ctx = context.WithValue(ctx, "user_email", email)
		c.Set("user_email", email)
	}
	if name, ok := claims["name"].(string); ok {
		ctx = context.WithValue(ctx, "user_name", name)
		c.Set("user_name", name)
	}

	// Grupos
	if groups, ok := claims["groups"].([]interface{}); ok {
		var groupStrings []string
		for _, group := range groups {
			if groupStr, ok := group.(string); ok {
				groupStrings = append(groupStrings, groupStr)
			}
		}
		ctx = context.WithValue(ctx, "user_groups", groupStrings)
		c.Set("user_groups", groups)
	}

	// Organización y fábricas
	if organization, ok := claims["organization"]; ok {
		// ✅ CORREGIDO: Manejar diferentes tipos de organización
		var orgMap map[string]interface{}

		if originalOrg, ok := organization.(map[string]interface{}); ok {
			orgMap = originalOrg
		} else if interfaceOrg, ok := organization.(map[interface{}]interface{}); ok {
			orgMap = convertMiddlewareMap(interfaceOrg)
		}

		if orgMap != nil {
			ctx = context.WithValue(ctx, "user_organization", orgMap)
			c.Set("user_organization", orgMap)

			// Extraer nombres de fábricas
			if factories, ok := orgMap["factories"].(map[string]interface{}); ok {
				var factoryNames []string
				for factoryName := range factories {
					factoryNames = append(factoryNames, factoryName)
				}
				ctx = context.WithValue(ctx, "factory_names", factoryNames)
				c.Set("factory_names", factoryNames)
			}
		} else {
			m.logger.Debug("Unknown organization type in JWT",
				zap.String("type", fmt.Sprintf("%T", organization)))
			// Poner igualmente en contexto aunque no sea el tipo esperado
			ctx = context.WithValue(ctx, "user_organization", organization)
			c.Set("user_organization", organization)
		}
	}

	// Actualizar el request con el nuevo contexto
	c.SetRequest(c.Request().WithContext(ctx))

	m.logger.Debug("✅ User context set",
		zap.String("user_id", claims["sub"].(string)),
		zap.Bool("has_organization", ctx.Value("user_organization") != nil),
		zap.Bool("has_factories", ctx.Value("factory_names") != nil))
}

// ✅ Función helper para convertir map[interface{}]interface{}
func convertMiddlewareMap(original map[interface{}]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range original {
		if strKey, ok := key.(string); ok {
			if nestedMap, ok := value.(map[interface{}]interface{}); ok {
				result[strKey] = convertMiddlewareMap(nestedMap)
			} else if nestedSlice, ok := value.([]interface{}); ok {
				result[strKey] = convertMiddlewareSlice(nestedSlice)
			} else {
				result[strKey] = value
			}
		}
	}
	return result
}

func convertMiddlewareSlice(original []interface{}) []interface{} {
	result := make([]interface{}, len(original))
	for i, item := range original {
		if nestedMap, ok := item.(map[interface{}]interface{}); ok {
			result[i] = convertMiddlewareMap(nestedMap)
		} else {
			result[i] = item
		}
	}
	return result
}

// Funciones helper esenciales que SÍ se usan

func GetUserClaims(c echo.Context) (jwt.MapClaims, error) {
	if claims := c.Get("user_claims"); claims != nil {
		if userClaims, ok := claims.(jwt.MapClaims); ok {
			return userClaims, nil
		}
		return nil, fmt.Errorf("invalid user_claims type in context")
	}
	return nil, fmt.Errorf("user_claims not found in context")
}

func GetUserOrganization(c echo.Context) (map[string]interface{}, error) {
	if org := c.Get("user_organization"); org != nil {
		if organization, ok := org.(map[string]interface{}); ok {
			return organization, nil
		}
		return nil, fmt.Errorf("invalid user_organization type in context")
	}
	return nil, fmt.Errorf("user_organization not found in context")
}

// Middlewares de autorización esenciales

func RequireGroup(group string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			claims, err := GetUserClaims(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "Authentication required")
			}

			groups := extractGroupsFromClaims(claims)
			if !hasGroup(groups, group) {
				return echo.NewHTTPError(http.StatusForbidden,
					"Insufficient permissions. Required group: "+group)
			}

			return next(c)
		}
	}
}

func RequireFactory(factory string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			organization, err := GetUserOrganization(c)
			if err != nil {
				return echo.NewHTTPError(http.StatusForbidden, "Organization information not found")
			}

			if !hasFactoryAccess(organization, factory) {
				return echo.NewHTTPError(http.StatusForbidden,
					"Insufficient permissions. Required factory access: "+factory)
			}

			return next(c)
		}
	}
}

// Helpers internos

func extractGroupsFromClaims(claims jwt.MapClaims) []string {
	var groups []string
	if groupsInterface, ok := claims["groups"].([]interface{}); ok {
		for _, group := range groupsInterface {
			if groupStr, ok := group.(string); ok {
				groups = append(groups, groupStr)
			}
		}
	}
	return groups
}

func hasGroup(groups []string, targetGroup string) bool {
	for _, group := range groups {
		if group == targetGroup {
			return true
		}
	}
	return false
}

func hasFactoryAccess(organization map[string]interface{}, factory string) bool {
	factories, ok := organization["factories"].(map[string]interface{})
	if !ok {
		return false
	}
	_, exists := factories[factory]
	return exists
}
