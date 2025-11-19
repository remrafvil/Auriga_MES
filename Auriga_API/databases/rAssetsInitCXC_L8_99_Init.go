package databases

import (
	"fmt"

	"github.com/remrafvil/Auriga_API/internal/repositories/rAssets"
	"github.com/remrafvil/Auriga_API/internal/repositories/rAssets/rA_Types"
	"github.com/remrafvil/Auriga_API/internal/repositories/rAssets/rAcxcL8"
	"gorm.io/gorm"
)

func CreateMultipleAssets_CXC_L8(db *gorm.DB) error {
	assetGroups := map[string][]rA_Types.AssetCreationData{
		"Common":     rAcxcL8.AssetsToCreate_CXC_L8_00_Common,
		"Transport":  rAcxcL8.AssetsToCreate_CXC_L8_10_Transport,
		"Dosing":     rAcxcL8.AssetsToCreate_CXC_L8_20_Dosing,
		"ExtrA":      rAcxcL8.AssetsToCreate_CXC_L8_ExtrA,
		"ExtrB":      rAcxcL8.AssetsToCreate_CXC_L8_ExtrB,
		"ExtrC":      rAcxcL8.AssetsToCreate_CXC_L8_ExtrC,
		"48_FB_Die":  rAcxcL8.AssetsToCreate_CXC_L8_48_FB_Die,
		"50_Stack":   rAcxcL8.AssetsToCreate_CXC_L8_50_Stack,
		"80_Winders": rAcxcL8.AssetsToCreate_CXC_L8_80_Winders,
	}

	for name, assetsGroup := range assetGroups {
		if err := rAssets.CreateMultipleAssets(db, assetsGroup); err != nil {
			fmt.Printf("❌ Error al crear activos %s: %v\n", name, err)
		} else {
			fmt.Printf("✅ Activos %s creados exitosamente\n", name)
		}
	}

	return nil
}
