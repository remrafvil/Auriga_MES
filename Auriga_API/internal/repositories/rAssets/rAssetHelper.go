package rAssets

import (
	"fmt"

	"github.com/lib/pq"
	"github.com/remrafvil/Auriga_API/internal/repositories/rAssets/rA_Types"
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

func CreateMultipleAssets(db *gorm.DB, assets []rA_Types.AssetCreationData) error {
	// Iniciar una transacción
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Primera pasada: crear todos los activos raíz (sin padre)
	for _, assetData := range assets {
		if assetData.ParentTechCode == nil {
			if _, err := createAssetFromData(tx, assetData, nil); err != nil {
				tx.Rollback()
				return fmt.Errorf("error creando activo raíz: %v", err)
			}
		}
	}

	// Segunda pasada: crear activos con padres
	for _, assetData := range assets {
		if assetData.ParentTechCode != nil {
			if _, err := createAssetFromData(tx, assetData, assetData.ParentTechCode); err != nil {
				tx.Rollback()
				return fmt.Errorf("error creando activo con padre: %v", err)
			}
		}
	}

	// Confirmar la transacción
	return tx.Commit().Error
}

// createAssetFromData es una función auxiliar que maneja la lógica de creación de activos
func createAssetFromData(tx *gorm.DB, assetData rA_Types.AssetCreationData, parentTechCode *uint) (*rModels.MrAsset, error) {
	// Primero creamos el activo base
	asset, err := createAsset(tx, assetData, parentTechCode)
	if err != nil {
		return nil, err
	}

	var hierarchicalLevel pq.StringArray

	if parentTechCode != nil {
		// Obtenemos el padre
		var parent rModels.MrAsset
		if err := tx.Where("tech_code = ?", *parentTechCode).First(&parent).Error; err != nil {
			return nil, fmt.Errorf("error buscando padre con tech_code %d: %v", *parentTechCode, err)
		}

		// Usamos el hierarchical_level del padre como base
		if len(parent.HierarchicalLevel) > 0 {
			hierarchicalLevel = make(pq.StringArray, len(parent.HierarchicalLevel))
			copy(hierarchicalLevel, parent.HierarchicalLevel)
		}

		// Añadimos el código del padre solo si no está ya en el hierarchical_level
		if len(hierarchicalLevel) == 0 || hierarchicalLevel[len(hierarchicalLevel)-1] != parent.Code {
			hierarchicalLevel = append(hierarchicalLevel, parent.Code)
		}
	}

	// Añadimos el código del activo actual
	hierarchicalLevel = append(hierarchicalLevel, asset.Code)

	// Actualizamos el activo
	asset.HierarchicalLevel = hierarchicalLevel
	if err := tx.Model(asset).Update("hierarchical_level", hierarchicalLevel).Error; err != nil {
		return nil, fmt.Errorf("error actualizando hierarchical_level: %v", err)
	}

	return asset, nil
}

// createAsset es la función base que crea un activo en la base de datos
func createAsset(tx *gorm.DB, assetData rA_Types.AssetCreationData, parentTechCode *uint) (*rModels.MrAsset, error) {
	// Buscar o crear el producto
	var product rModels.MrProduct
	if err := tx.Where("name = ?", assetData.ProductName).First(&product).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			product = rModels.MrProduct{Name: assetData.ProductName}
			if err := tx.Create(&product).Error; err != nil {
				return nil, fmt.Errorf("error creando producto: %v", err)
			}
		} else {
			return nil, fmt.Errorf("error buscando producto: %v", err)
		}
	}

	// Preparar el activo
	asset := &rModels.MrAsset{
		ProductID: product.ID,
		Location:  assetData.Location,
		TechCode:  assetData.TechCode,
		Code:      assetData.Code,
		Sn:        assetData.Sn,
		SapCode:   assetData.SapCode,
	}

	// Manejar el padre si se especifica
	if parentTechCode != nil {
		var parentAsset rModels.MrAsset
		if err := tx.Where("tech_code = ?", *parentTechCode).First(&parentAsset).Error; err != nil {
			return nil, fmt.Errorf("error buscando activo padre: %v", err)
		}
		asset.ParentID = &parentAsset.ID
	}

	// Crear el activo
	if err := tx.Create(asset).Error; err != nil {
		return nil, fmt.Errorf("error creando activo: %v", err)
	}

	//fmt.Printf("Activo creado: %s (TechCode: %d)\n", assetData.ProductName, assetData.TechCode)
	return asset, nil
}
