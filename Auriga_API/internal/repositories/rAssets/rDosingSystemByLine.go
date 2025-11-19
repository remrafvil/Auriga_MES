package rAssets

import (
	"context"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

func (r *repository) GetDosingSystemByLine(ctx context.Context, lineID uint) ([]rModels.MrAsset, error) {
	var dosingSystems []rModels.MrAsset

	// Buscar el sistema de dosificación (2_Dosing) para la línea específica
	err := r.db.WithContext(ctx).
		Preload("Children", func(db *gorm.DB) *gorm.DB {
			return db.Preload("Children") // Dosers y sus componentes
		}).
		Where("parent_id = ? AND code = ?", lineID, "2_Dosing").
		Find(&dosingSystems).Error

	if err != nil {
		return nil, err
	}

	return dosingSystems, nil
}

func (r *repository) GetDoserComponents(ctx context.Context, doserID uint) ([]rModels.MrAsset, error) {
	var components []rModels.MrAsset

	err := r.db.WithContext(ctx).
		Where("parent_id = ?", doserID).
		Order("code ASC"). // Ordenar por código para consistencia
		Find(&components).Error

	if err != nil {
		return nil, err
	}

	return components, nil
}
