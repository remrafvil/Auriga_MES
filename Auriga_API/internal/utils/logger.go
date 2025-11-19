package utils

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type contextKey string

const (
	RequestIDKey contextKey = "request_id"
	OperationKey contextKey = "operation"
)

// Logger wrapper para inyectar contexto
type Logger struct {
	*zap.Logger
}

// NewLogger crea un nuevo logger configurado
func NewLogger(environment string) (*Logger, error) {
	var config zap.Config

	if environment == "production" {
		config = zap.NewProductionConfig()
		config.EncoderConfig.TimeKey = "timestamp"
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.LevelKey = "severity"
	config.DisableStacktrace = false

	// Configuración común
	config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	if environment != "production" {
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return &Logger{logger}, nil
}

// WithContext crea un logger con campos de contexto
func (l *Logger) WithContext(ctx context.Context) *Logger {
	fields := []zap.Field{
		zap.Time("timestamp", time.Now()),
	}

	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		fields = append(fields, zap.String("request_id", requestID))
	}

	if operation, ok := ctx.Value(OperationKey).(string); ok {
		fields = append(fields, zap.String("operation", operation))
	}

	return &Logger{l.Logger.With(fields...)}
}

// SafeFields filtra campos sensibles
func SafeFields(data map[string]interface{}) []zap.Field {
	fields := make([]zap.Field, 0, len(data))
	sensitiveKeys := map[string]bool{
		"password": true, "token": true, "authorization": true,
		"email": true, "phone": true, "salary": true,
		"secret": true, "apikey": true, "credential": true,
	}

	for key, value := range data {
		if sensitiveKeys[key] {
			fields = append(fields, zap.String(key, "***REDACTED***"))
		} else {
			fields = append(fields, zap.Any(key, value))
		}
	}
	return fields
}

// Helper para crear contexto con operación
func WithOperation(ctx context.Context, operation string) context.Context {
	return context.WithValue(ctx, OperationKey, operation)
}

// Helper para obtener request ID del contexto
func GetRequestID(ctx context.Context) string {
	if requestID, ok := ctx.Value(RequestIDKey).(string); ok {
		return requestID
	}
	return ""
}
