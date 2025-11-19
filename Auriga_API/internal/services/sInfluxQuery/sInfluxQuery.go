package sInfluxQuery

import (
	"context"

	"github.com/remrafvil/Auriga_API/internal/repositories/rInfluxQuery"
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"go.uber.org/zap"
)

type Service interface {
	GetQueryByName(ctx context.Context, name string) (*rModels.MrInfluxGrafanaQuery, error)
}

// service implementación
type service struct {
	repository rInfluxQuery.Repository
	logger     *zap.Logger
}

func New(repository rInfluxQuery.Repository, logger *zap.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

func (s *service) GetQueryByName(ctx context.Context, name string) (*rModels.MrInfluxGrafanaQuery, error) {
	s.logger.Info("Getting influx query by name", zap.String("name", name))
	return s.repository.FindByName(ctx, name)
}
