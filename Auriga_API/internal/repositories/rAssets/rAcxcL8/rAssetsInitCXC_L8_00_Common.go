package rAcxcL8

import "github.com/remrafvil/Auriga_API/internal/repositories/rAssets/rA_Types"

// Función auxiliar para convertir uint a *uint
func ptrUint(i uint) *uint {
	return &i
}

/* var assetsToCreate = []rA_Types.AssetCreationData{
	// Activos raíz (sin padre)
	{
		ProductName:    "CXC",
		ParentTechCode: nil,
		Location:       "Chile",
		TechCode:       25000000,
		Code:           "CXC",
		Sn:             "02020202",
		SapCode:        "1010",
	},
	{
		ProductName:    "CXC_L08",
		ParentTechCode: ptrUint(25000000),
		Location:       "CXC_L08",
		TechCode:       25080000,
		Code:           "Line_08",
		Sn:             "02_08",
		SapCode:        "W1010188",
	},
} */

var AssetsToCreate_CXC_L8_00_Common = []rA_Types.AssetCreationData{

	// Subsistemas (nivel 2)
	{
		ProductName:    "0_Common",
		ParentTechCode: ptrUint(25080000),
		Location:       "CXC_L08",
		TechCode:       25080100,
		Code:           "0_Common",
		Sn:             "25_08_01",
		SapCode:        "",
	},
	{
		ProductName:    "1_Transport",
		ParentTechCode: ptrUint(25080000),
		Location:       "CXC_L08",
		TechCode:       25081000,
		Code:           "1_Transport",
		Sn:             "25_08_02",
		SapCode:        "",
	},
	{
		ProductName:    "2_Dosing",
		ParentTechCode: ptrUint(25080000),
		Location:       "CXC_L08",
		TechCode:       25082000,
		Code:           "2_Dosing",
		Sn:             "25_08_03",
		SapCode:        "",
	},
	{
		ProductName:    "3_Extrusion",
		ParentTechCode: ptrUint(25080000),
		Location:       "CXC_L08",
		TechCode:       25083000,
		Code:           "3_Extrusion",
		Sn:             "25_08_04",
		SapCode:        "",
	},
	{
		ProductName:    "5_Stack",
		ParentTechCode: ptrUint(25080000),
		Location:       "CXC_L08",
		TechCode:       25085000,
		Code:           "5_Stack",
		Sn:             "25_08_05",
		SapCode:        "",
	},
	{
		ProductName:    "8_Winders",
		ParentTechCode: ptrUint(25080000),
		Location:       "CXC_L08",
		TechCode:       25088000,
		Code:           "8_Winders",
		Sn:             "25_08_08",
		SapCode:        "",
	},
}
