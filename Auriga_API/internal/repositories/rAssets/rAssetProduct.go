package rAssets

import (
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
)

// GetAssetWithProductDetails obtiene un activo con los detalles del producto
// incluyendo las características y documentos asociados.

func (m *repository) GetAssetWithProductDetails(assetID uint) (rModels.MrAsset, error) {
	var asset rModels.MrAsset

	err := m.db.Preload("Product").
		Preload("Product.ProductType").
		Preload("Product.ProductType.Features").
		Preload("Product.ProductType.Relations").
		Preload("Product.FeatureValues").
		Preload("Product.FeatureValues.ProductFeature").
		Preload("Product.Documents").
		Preload("Documents").
		//	Preload("Children").
		First(&asset, assetID).Error

	if err != nil {
		return rModels.MrAsset{}, err
	}

	return asset, nil
}
