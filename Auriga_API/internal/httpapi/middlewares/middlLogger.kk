package middlewares

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/remrafvil/Auriga_API/internal/utils"
	"go.uber.org/zap"
)

// RequestLoggerMiddleware crea middleware de logging con contexto
func RequestLoggerMiddleware(logger *utils.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			// Generar request ID
			requestID := fmt.Sprintf("req-%d", time.Now().UnixNano())

			// Crear contexto con request ID
			ctx := context.WithValue(c.Request().Context(), utils.RequestIDKey, requestID)

			// Establecer operación basada en ruta y método
			operation := getOperationFromPath(c.Path(), c.Request().Method)
			ctx = utils.WithOperation(ctx, operation)

			c.SetRequest(c.Request().WithContext(ctx))

			// Logger con contexto
			requestLogger := logger.WithContext(ctx)

			// Registrar inicio
			requestLogger.Debug("Request started",
				zap.String("method", c.Request().Method),
				zap.String("path", c.Path()),
				zap.String("remote_ip", c.RealIP()),
				zap.String("user_agent", c.Request().UserAgent()),
			)

			// Procesar request
			err := next(c)

			// Registrar fin
			latency := time.Since(start)
			status := c.Response().Status

			logFields := []zap.Field{
				zap.Int("status", status),
				zap.String("method", c.Request().Method),
				zap.String("path", c.Path()),
				zap.String("remote_ip", c.RealIP()),
				zap.Duration("latency", latency),
				zap.Int64("latency_ms", latency.Milliseconds()),
				zap.String("user_agent", c.Request().UserAgent()),
			}

			if err != nil {
				logFields = append(logFields, zap.Error(err))
			}

			switch {
			case status >= 500:
				requestLogger.Error("Request completed with server error", logFields...)
			case status >= 400:
				requestLogger.Warn("Request completed with client error", logFields...)
			default:
				requestLogger.Info("Request completed successfully", logFields...)
			}

			return err
		}
	}
}

func getOperationFromPath(path, method string) string {
	// Mapear rutas a nombres de operación legibles
	switch {
	case path == "/api/v1/users" && method == "POST":
		return "create_user"
	case path == "/api/v1/auth/login" && method == "POST":
		return "user_login"
	case strings.HasPrefix(path, "/api/v1/assets"):
		return "asset_operation"
	case strings.HasPrefix(path, "/api/v1/orders"):
		return "order_operation"
	default:
		return fmt.Sprintf("%s_%s", strings.ToLower(method), sanitizePath(path))
	}
}

func sanitizePath(path string) string {
	// Convertir /api/v1/users/123 -> users
	parts := strings.Split(path, "/")
	for i, part := range parts {
		if part == "api" && i+2 < len(parts) {
			return parts[i+2] // Devuelve el segmento después de /api/v1/
		}
	}
	if len(parts) > 0 && parts[len(parts)-1] != "" {
		return parts[len(parts)-1]
	}
	return "root"
}
