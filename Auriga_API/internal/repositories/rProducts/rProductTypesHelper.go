package rProducts

import (
	"fmt"
	"strings"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"github.com/remrafvil/Auriga_API/internal/repositories/rProducts/rP_Types"
	"gorm.io/gorm"
)

// Helper mejorado para crear/obtener características (con validación implícita)
func createOrGetFeature(db *gorm.DB, name, symbol, unit string, valueType rModels.ValueType) (rModels.MrProductFeatureType, error) {
	feature := rModels.MrProductFeatureType{
		Name:      strings.TrimSpace(name),
		Symbol:    strings.TrimSpace(symbol),
		Unit:      strings.TrimSpace(unit),
		ValueType: valueType,
	}

	if feature.Name == "" {
		return feature, fmt.Errorf("el nombre de la característica no puede estar vacío")
	}

	result := db.FirstOrCreate(&feature, rModels.MrProductFeatureType{Name: feature.Name})
	return feature, result.Error
}

// Función para asociar características con orden visual (con validación implícita)
func associateFeatureWithOrder(db *gorm.DB, productType *rModels.MrProductType, feature rModels.MrProductFeatureType, order uint) error {
	if productType.ID == 0 {
		return fmt.Errorf("ID de tipo de producto no puede ser cero")
	}
	if feature.ID == 0 {
		return fmt.Errorf("ID de característica no puede ser cero")
	}

	relation := rModels.MrProductFeatureTypeRelation{
		MrProductTypeID:        productType.ID,
		MrProductFeatureTypeID: feature.ID,
		VisuOrder:              order,
	}
	return db.Create(&relation).Error
}

func CreateMultipleProductTypes(db *gorm.DB, productTypes []rP_Types.ProductTypeData) error {
	// Validar los datos antes de comenzar
	if err := validateProductTypesData(productTypes); err != nil {
		return fmt.Errorf("validación fallida: %w", err)
	}

	// Ejecutar en una transacción
	return db.Transaction(func(tx *gorm.DB) error {
		for _, ptData := range productTypes {
			// Crear o obtener el tipo de producto
			productType := rModels.MrProductType{Name: ptData.Name}
			if err := tx.FirstOrCreate(&productType, rModels.MrProductType{Name: ptData.Name}).Error; err != nil {
				return fmt.Errorf("error al crear producto %s: %w", ptData.Name, err)
			}

			// Procesar cada característica
			for _, f := range ptData.Features {
				feature, err := createOrGetFeature(tx, f.Name, f.Symbol, f.Unit, f.ValueType)
				if err != nil {
					return fmt.Errorf("error al crear feature %s: %w", f.Name, err)
				}

				if err := associateFeatureWithOrder(tx, &productType, feature, f.Order); err != nil {
					return fmt.Errorf("error al asociar feature %s: %w", f.Name, err)
				}
			}
		}
		return nil
	})
}

// Función de validación
func validateProductTypesData(productTypes []rP_Types.ProductTypeData) error {
	if len(productTypes) == 0 {
		return fmt.Errorf("lista de tipos de productos vacía")
	}

	seenProducts := make(map[string]bool)
	for _, pt := range productTypes {
		// Validar nombre del producto
		if strings.TrimSpace(pt.Name) == "" {
			return fmt.Errorf("nombre de producto no puede estar vacío")
		}

		// Validar duplicados en el lote actual
		if seenProducts[pt.Name] {
			return fmt.Errorf("nombre de producto duplicado: %s", pt.Name)
		}
		seenProducts[pt.Name] = true

		// Validar características
		if len(pt.Features) == 0 {
			return fmt.Errorf("el producto %s no tiene características", pt.Name)
		}

		seenFeatures := make(map[string]bool)
		seenOrders := make(map[uint]bool)
		for _, f := range pt.Features {
			// Validar nombre de característica
			if strings.TrimSpace(f.Name) == "" {
				return fmt.Errorf("nombre de característica no puede estar vacío en producto %s", pt.Name)
			}

			// Validar duplicados
			if seenFeatures[f.Name] {
				return fmt.Errorf("característica duplicada '%s' en producto %s", f.Name, pt.Name)
			}
			seenFeatures[f.Name] = true

			// Validar orden visual único
			if seenOrders[f.Order] {
				return fmt.Errorf("orden visual duplicado %d en producto %s", f.Order, pt.Name)
			}
			seenOrders[f.Order] = true

			// Validar tipo de valor
			if f.ValueType != rModels.ValueTypeString && f.ValueType != rModels.ValueTypeNumber {
				return fmt.Errorf("tipo de valor inválido para característica %s en producto %s", f.Name, pt.Name)
			}
		}
	}
	return nil
}
