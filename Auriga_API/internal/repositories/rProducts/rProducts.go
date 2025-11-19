package rProducts

// lo llamaremso repositories

import (
	"context"
	"fmt"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	ProductFindAll(ctx context.Context) ([]rModels.MrProduct, error)
	ProductFindByID(ctx context.Context, id uint) (*rModels.MrProduct, error)
	ProductCreate(ctx context.Context, product *rModels.MrProduct, featureValues map[string]string) error
	ProductUpdate(ctx context.Context, id uint, product *rModels.MrProduct, featureValues map[string]string) error
	ProductDelete(ctx context.Context, id uint) error

	ProductTypeFindAll(ctx context.Context) ([]rModels.MrProductType, error)
	ProductTypeFindByID(ctx context.Context, id uint) (*rModels.MrProductType, error)
	ProductTypeFindAllWithFeatures(ctx context.Context) ([]rModels.MrProductType, error)
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
func (r *repository) ProductFindAll(ctx context.Context) ([]rModels.MrProduct, error) {
	var products []rModels.MrProduct
	err := r.db.WithContext(ctx).
		Preload("ProductType").
		Select("id", "name", "product_type_id", "manufacturer", "family", "description").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (r *repository) ProductFindByID(ctx context.Context, id uint) (*rModels.MrProduct, error) {
	var product rModels.MrProduct
	err := r.db.WithContext(ctx).
		Preload("ProductType").
		Preload("FeatureValues.ProductFeature").
		First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}
func (r *repository) ProductCreate(ctx context.Context, product *rModels.MrProduct, featureValues map[string]string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Obtener el tipo de producto
		var productType rModels.MrProductType
		if err := tx.Where("name = ?", product.ProductType.Name).First(&productType).Error; err != nil {
			return fmt.Errorf("no se encontró el tipo de producto '%s': %v", product.ProductType.Name, err)
		}

		// 2. Crear un producto "limpio" sin el objeto ProductType embebido
		newProduct := &rModels.MrProduct{
			Name:          product.Name,
			ProductTypeID: productType.ID, // ← Solo usar el ID
			Manufacturer:  product.Manufacturer,
			Family:        product.Family,
			Description:   product.Description,
			// NO incluir product.ProductType aquí
		}

		// 3. Crear el producto
		if err := tx.Create(newProduct).Error; err != nil {
			return fmt.Errorf("error al crear el producto: %v", err)
		}

		// 4. Actualizar el producto original con el ID generado
		product.ID = newProduct.ID
		product.ProductTypeID = newProduct.ProductTypeID

		// 5. Crear los valores de características
		if err := r.ProductCreateFeatureValues(tx, product.ID, productType.ID, featureValues); err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) ProductUpdate(ctx context.Context, id uint, product *rModels.MrProduct, featureValues map[string]string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Verificar que el producto existe
		var existingProduct rModels.MrProduct
		if err := tx.First(&existingProduct, id).Error; err != nil {
			return fmt.Errorf("producto no encontrado: %v", err)
		}

		// 2. Obtener el tipo de producto si se está cambiando
		var productType rModels.MrProductType
		if err := tx.Where("name = ?", product.ProductType.Name).First(&productType).Error; err != nil {
			return fmt.Errorf("no se encontró el tipo de producto '%s': %v", product.ProductType.Name, err)
		}

		// 3. Actualizar el producto
		product.ID = id
		product.ProductTypeID = productType.ID

		// Crear un producto "limpio" para la actualización
		updateProduct := &rModels.MrProduct{
			Name:          product.Name,
			ProductTypeID: productType.ID,
			Manufacturer:  product.Manufacturer,
			Family:        product.Family,
			Description:   product.Description,
		}

		if err := tx.Model(&existingProduct).Updates(updateProduct).Error; err != nil {
			return fmt.Errorf("error al actualizar el producto: %v", err)
		}

		// 4. ELIMINACIÓN FÍSICA de valores de características existentes
		// Usar Unscoped().Delete() para eliminar físicamente en lugar de soft delete
		if err := tx.Unscoped().Where("mr_product_id = ?", id).Delete(&rModels.MrProductFeatureValue{}).Error; err != nil {
			return fmt.Errorf("error al eliminar valores de características existentes: %v", err)
		}

		// 5. Crear nuevos valores de características
		if err := r.ProductCreateFeatureValues(tx, id, productType.ID, featureValues); err != nil {
			return err
		}

		return nil
	})
}

func (r *repository) ProductDelete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. Verificar que el producto existe
		var product rModels.MrProduct
		if err := tx.First(&product, id).Error; err != nil {
			return fmt.Errorf("producto no encontrado: %v", err)
		}

		// 2. Eliminar valores de características
		if err := tx.Where("mr_product_id = ?", id).Delete(&rModels.MrProductFeatureValue{}).Error; err != nil {
			return fmt.Errorf("error al eliminar valores de características: %v", err)
		}

		// 3. Eliminar el producto
		if err := tx.Delete(&product).Error; err != nil {
			return fmt.Errorf("error al eliminar el producto: %v", err)
		}

		return nil
	})
}

func (r *repository) ProductCreateFeatureValues(tx *gorm.DB, productID uint, productTypeID uint, featureValues map[string]string) error {
	for featureName, value := range featureValues {
		if value == "" {
			continue // Saltar valores vacíos
		}

		// Obtener la definición de la característica
		var featureType rModels.MrProductFeatureType
		if err := tx.Where("name = ?", featureName).First(&featureType).Error; err != nil {
			return fmt.Errorf("no se encontró la característica '%s': %v", featureName, err)
		}

		// Verificar que la característica pertenece al tipo de producto
		var relationExists int64
		if err := tx.Model(&rModels.MrProductFeatureTypeRelation{}).
			Where("mr_product_type_id = ? AND mr_product_feature_type_id = ?", productTypeID, featureType.ID).
			Count(&relationExists).Error; err != nil || relationExists == 0 {
			return fmt.Errorf("la característica '%s' no está asociada al tipo de producto", featureName)
		}

		// Crear el valor de la característica según su tipo
		featureValue := rModels.MrProductFeatureValue{
			MrProductID:            productID,
			MrProductFeatureTypeID: featureType.ID,
		}

		if featureType.ValueType == rModels.ValueTypeNumber {
			var numValue float64
			if _, err := fmt.Sscanf(value, "%f", &numValue); err != nil {
				return fmt.Errorf("el valor '%s' para la característica '%s' no es un número válido", value, featureName)
			}
			featureValue.NumberValue = numValue
		} else {
			featureValue.StringValue = value
		}

		if err := tx.Create(&featureValue).Error; err != nil {
			return fmt.Errorf("error al crear el valor para la característica '%s': %v", featureName, err)
		}
	}

	return nil
}
