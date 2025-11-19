package rProducts

import (
	"context"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

func (r *repository) ProductTypeFindAll(ctx context.Context) ([]rModels.MrProductType, error) {
	var productTypes []rModels.MrProductType
	err := r.db.WithContext(ctx).
		Select("id", "name").
		Find(&productTypes).Error
	if err != nil {
		return nil, err
	}
	return productTypes, nil
}

func (r *repository) ProductTypeFindByID(ctx context.Context, id uint) (*rModels.MrProductType, error) {
	var productType rModels.MrProductType
	err := r.db.WithContext(ctx).
		Preload("Features", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "symbol", "unit", "value_type")
		}).
		First(&productType, id).Error
	if err != nil {
		return nil, err
	}
	return &productType, nil
}

func (r *repository) ProductTypeFindAllWithFeatures(ctx context.Context) ([]rModels.MrProductType, error) {
	var productTypes []rModels.MrProductType
	err := r.db.WithContext(ctx).
		Preload("Features", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name", "symbol", "unit", "value_type")
		}).
		Find(&productTypes).Error
	if err != nil {
		return nil, err
	}
	return productTypes, nil
}
