package rProducts

import (
	"fmt"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"github.com/remrafvil/Auriga_API/internal/repositories/rProducts/rP_Types"

	"gorm.io/gorm"
)

func CreateMultipleProductsWithValues(db *gorm.DB, products []rP_Types.ProductCreationData) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// Pre-cache de tipos de producto para optimización
		productTypes := make(map[string]rModels.MrProductType)

		// Pre-cache de características para optimización
		featureTypes := make(map[string]rModels.MrProductFeatureType)

		// 1. Primera pasada: validar todos los datos antes de crear cualquier registro
		for _, productData := range products {
			// Verificar tipo de producto
			if _, exists := productTypes[productData.ProductType]; !exists {
				var pt rModels.MrProductType
				if err := tx.Where("name = ?", productData.ProductType).First(&pt).Error; err != nil {
					return fmt.Errorf("no se encontró el tipo de producto '%s': %v", productData.ProductType, err)
				}
				productTypes[productData.ProductType] = pt
			}

			// Verificar características
			productType := productTypes[productData.ProductType]
			for featureName := range productData.FeatureValues {
				if _, exists := featureTypes[featureName]; !exists {
					var ft rModels.MrProductFeatureType
					if err := tx.Where("name = ?", featureName).First(&ft).Error; err != nil {
						return fmt.Errorf("no se encontró la característica '%s'", featureName)
					}
					featureTypes[featureName] = ft

					// Verificar relación entre característica y tipo de producto
					var relationExists int64
					if err := tx.Model(&rModels.MrProductFeatureTypeRelation{}).
						Where("mr_product_type_id = ? AND mr_product_feature_type_id = ?",
							productType.ID, ft.ID).
						Count(&relationExists).Error; err != nil || relationExists == 0 {
						return fmt.Errorf("la característica '%s' no está asociada al tipo de producto '%s'",
							featureName, productData.ProductType)
					}
				}
			}
		}

		// 2. Segunda pasada: crear todos los productos y sus valores
		for _, productData := range products {
			productType := productTypes[productData.ProductType]

			// Crear el producto
			product := rModels.MrProduct{
				Name:          productData.Name,
				ProductTypeID: productType.ID,
				Manufacturer:  productData.Manufacturer,
				Family:        productData.Family,
				Description:   productData.Description,
			}
			if err := tx.Create(&product).Error; err != nil {
				return fmt.Errorf("error al crear el producto '%s': %v", productData.Name, err)
			}

			// Crear los valores de las características
			for featureName, value := range productData.FeatureValues {
				if value == "" {
					continue
				}

				featureType := featureTypes[featureName]
				featureValue := rModels.MrProductFeatureValue{
					MrProductID:            product.ID,
					MrProductFeatureTypeID: featureType.ID,
				}

				if featureType.ValueType == rModels.ValueTypeNumber {
					var numValue float64
					if _, err := fmt.Sscanf(value, "%f", &numValue); err != nil {
						return fmt.Errorf("el valor '%s' para la característica '%s' no es un número válido",
							value, featureName)
					}
					featureValue.NumberValue = numValue
				} else {
					featureValue.StringValue = value
				}

				if err := tx.Create(&featureValue).Error; err != nil {
					return fmt.Errorf("error al crear el valor para la característica '%s' del producto '%s': %v",
						featureName, productData.Name, err)
				}
			}
		}

		return nil
	})
}

// Función  para crear productos con valores
func createProductWithValues(db *gorm.DB, productName, productTypeName, manufacturer, family,
	description string, featureValues map[string]string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. Obtener el tipo de producto
		var productType rModels.MrProductType
		if err := tx.Where("name = ?", productTypeName).First(&productType).Error; err != nil {
			return fmt.Errorf("no se encontró el tipo de producto '%s': %v", productTypeName, err)
		}

		// 2. Crear el producto
		product := rModels.MrProduct{
			Name:          productName,
			ProductTypeID: productType.ID,
			Manufacturer:  manufacturer,
			Family:        family,
			Description:   description,
		}
		if err := tx.Create(&product).Error; err != nil {
			return fmt.Errorf("error al crear el producto: %v", err)
		}

		// 3. Para cada característica en el mapa de valores, crear el valor asociado
		for featureName, value := range featureValues {
			if value == "" {
				continue // Saltar valores vacíos
			}

			// Obtener la definición de la característica
			var featureType rModels.MrProductFeatureType
			if err := tx.Where("name = ?", featureName).First(&featureType).Error; err != nil {
				return fmt.Errorf("no se encontró la característica '%s' para el producto '%s': %v", featureName, productName, err)
			}

			// Verificar que la característica pertenece al tipo de producto
			var relationExists int64
			if err := tx.Model(&rModels.MrProductFeatureTypeRelation{}).
				Where("mr_product_type_id = ? AND mr_product_feature_type_id = ?", productType.ID, featureType.ID).
				Count(&relationExists).Error; err != nil || relationExists == 0 {
				return fmt.Errorf("la característica '%s' no está asociada al tipo de producto '%s'", featureName, productTypeName)
			}

			// Crear el valor de la característica según su tipo
			featureValue := rModels.MrProductFeatureValue{
				MrProductID:            product.ID,
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
	})
}
