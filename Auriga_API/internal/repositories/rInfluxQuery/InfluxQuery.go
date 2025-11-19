package rInfluxQuery

import (
	"context"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	FindByName(ctx context.Context, name string) (*rModels.MrInfluxGrafanaQuery, error)
}

type repository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func New(db *gorm.DB, logger *zap.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}

func (r *repository) FindByName(ctx context.Context, name string) (*rModels.MrInfluxGrafanaQuery, error) {
	var query rModels.MrInfluxGrafanaQuery
	err := r.db.WithContext(ctx).Where("name = ?", name).First(&query).Error
	if err != nil {
		return nil, err
	}
	return &query, nil
}
