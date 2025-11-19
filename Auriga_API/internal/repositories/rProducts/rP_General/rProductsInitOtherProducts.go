package rP_General

import (
	"github.com/remrafvil/Auriga_API/internal/repositories/rProducts/rP_Types"
)

var EquipmentToCreate = []rP_Types.ProductCreationData{
	// Blenders
	{
		Name:         "FGB-5-6/0",
		ProductType:  "Blender",
		Manufacturer: "Ferlin",
		Family:       "FGB",
		Description:  "Batch Blender",
		FeatureValues: map[string]string{
			"Throughput":        "500.0",
			"Number Components": "6",
		},
	},
	{
		Name:         "FGB-10-6",
		ProductType:  "Blender",
		Manufacturer: "Ferlin",
		Family:       "FGB",
		Description:  "Batch Blender",
		FeatureValues: map[string]string{
			"Throughput":        "1000.0",
			"Number Components": "6",
		},
	},
	{
		Name:         "FGB-10-1-2-2",
		ProductType:  "Blender",
		Manufacturer: "Ferlin",
		Family:       "FGB",
		Description:  "Batch Blender",
		FeatureValues: map[string]string{
			"Throughput":        "1000.0",
			"Number Components": "5",
		},
	},
	{
		Name:         "FGB-5-1-2-2",
		ProductType:  "Blender",
		Manufacturer: "Ferlin",
		Family:       "FGB",
		Description:  "Batch Blender",
		FeatureValues: map[string]string{
			"Throughput":        "500.0",
			"Number Components": "5",
		},
	},

	// DosingVolu
	{
		Name:         "Coex_D-1",
		ProductType:  "DosingVolu",
		Manufacturer: "Coexpan",
		Family:       "Coex_Dosing",
		Description:  "Volimetric Blender",
		FeatureValues: map[string]string{
			"Throughput":        "100.0",
			"Number Components": "1",
		},
	},
	{
		Name:         "Coex_D-2",
		ProductType:  "DosingVolu",
		Manufacturer: "Coexpan",
		Family:       "Coex_Dosing",
		Description:  "Volimetric Blender",
		FeatureValues: map[string]string{
			"Throughput":        "200.0",
			"Number Components": "2",
		},
	},

	// ComponentHopper
	{
		Name:          "Coex_CH",
		ProductType:   "ComponentHopper",
		Manufacturer:  "Coexpan",
		Family:        "Coex_Component",
		Description:   "Hopper Loader Gravity Discharge",
		FeatureValues: map[string]string{},
	},

	// Transport Group

	{
		Name:         "Coex_TG",
		ProductType:  "TransportSystem",
		Manufacturer: "Coexpan",
		Family:       "Coex_TranSystem",
		Description:  "Transport Conveyor System",
	},

	// HopperLoader
	{
		Name:         "VC570R66",
		ProductType:  "HopperLoader",
		Manufacturer: "Ferlin",
		Family:       "Movac",
		Description:  "Hopper Loader Gravity Discharge",
		FeatureValues: map[string]string{
			"Volume": "70.0",
		},
	},
	{
		Name:         "VC340R66",
		ProductType:  "HopperLoader",
		Manufacturer: "Ferlin",
		Family:       "Movac",
		Description:  "Hopper Loader Gravity Discharge",
		FeatureValues: map[string]string{
			"Volume": "40.0",
		},
	},
	{
		Name:         "VC325R66",
		ProductType:  "HopperLoader",
		Manufacturer: "Ferlin",
		Family:       "Movac",
		Description:  "Hopper Loader Gravity Discharge",
		FeatureValues: map[string]string{
			"Volume": "25.0",
		},
	},
	{
		Name:         "VC207R66",
		ProductType:  "HopperLoader",
		Manufacturer: "Ferlin",
		Family:       "Movac",
		Description:  "Hopper Loader Gravity Discharge",
		FeatureValues: map[string]string{
			"Volume": "7.0",
		},
	},
}
