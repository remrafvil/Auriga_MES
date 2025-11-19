package rAcxcL8

import "github.com/remrafvil/Auriga_API/internal/repositories/rAssets/rA_Types"

var AssetsToCreate_CXC_L8_48_FB_Die = []rA_Types.AssetCreationData{

	// Feed Block & Die
	{
		ProductName:    "3_FB_Die",
		ParentTechCode: ptrUint(25083000),
		Location:       "CXC_L08_FB_Die",
		TechCode:       2508480000,
		Code:           "FB_Die",
		Sn:             "Kuhne_141225-70-01A1.0.00.0.3_FB_Die",
		SapCode:        "",
	},

	// Feed Block
	{
		ProductName:    "KU5",
		ParentTechCode: ptrUint(2508480000),
		Location:       "CXC_L08_FB",
		TechCode:       25084810000,
		Code:           "FB",
		Sn:             "Kuhne_141225-70-01A1.0.00.0.3_KU5",
		SapCode:        "",
	},
	// Feed Block Heating
	{
		ProductName:    "3_ExtruderHeating",
		ParentTechCode: ptrUint(25084810000),
		Location:       "CXC_L08_FB_HT",
		TechCode:       25084811000,
		Code:           "HT",
		Sn:             "25084811000_XX",
		SapCode:        "",
	},
	// Z02_01
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084811000),
		Location:       "CXC_L08_FB_HT_Z02_01",
		TechCode:       25084812100,
		Code:           "Z02_01",
		Sn:             "25084812100_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084812100),
		Location:       "CXC_L08_FB_HT_Z02_01_TS1",
		TechCode:       25084812105,
		Code:           "TS1",
		Sn:             "25084812105_XX",
		SapCode:        "",
	},

	// Z02_02
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084811000),
		Location:       "CXC_L08_FB_HT_Z02_02",
		TechCode:       25084812200,
		Code:           "Z02_02",
		Sn:             "25084812200_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084812200),
		Location:       "CXC_L08_FB_HT_Z02_02_TS1",
		TechCode:       25084812205,
		Code:           "TS1",
		Sn:             "25084812205_XX",
		SapCode:        "",
	},

	// Z02_03
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084811000),
		Location:       "CXC_L08_FB_HT_Z02_03",
		TechCode:       25084812300,
		Code:           "Z02_03",
		Sn:             "25084812300_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084812300),
		Location:       "CXC_L08_FB_HT_Z02_03_TS1",
		TechCode:       25084812305,
		Code:           "TS1",
		Sn:             "25084812305_XX",
		SapCode:        "",
	},

	// Die
	{
		ProductName:    "BDF117_601-17.A1.0.00.2.1",
		ParentTechCode: ptrUint(2508480000),
		Location:       "CXC_L08_FB",
		TechCode:       25084830000,
		Code:           "Die",
		Sn:             "Kuhne_141225-70-01A1.0.00.0.3_BDF117",
		SapCode:        "",
	},
	// Die Heating
	{
		ProductName:    "3_ExtruderHeating",
		ParentTechCode: ptrUint(25084830000),
		Location:       "CXC_L08_Die_HT",
		TechCode:       25084831000,
		Code:           "HT",
		Sn:             "25084831000_XX",
		SapCode:        "",
	},

	// Z03_01
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084831000),
		Location:       "CXC_L08_FB_HT_Z03_01",
		TechCode:       25084833100,
		Code:           "Z03_01",
		Sn:             "25084833100_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084833100),
		Location:       "CXC_L08_FB_HT_Z03_01_TS1",
		TechCode:       25084833105,
		Code:           "TS1",
		Sn:             "25084833105_XX",
		SapCode:        "",
	},
	// Z03_02
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084831000),
		Location:       "CXC_L08_FB_HT_Z03_02",
		TechCode:       25084833200,
		Code:           "Z03_02",
		Sn:             "25084833200_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084833200),
		Location:       "CXC_L08_FB_HT_Z03_02_TS1",
		TechCode:       25084833205,
		Code:           "TS1",
		Sn:             "25084833205_XX",
		SapCode:        "",
	},
	// Z03_03
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084831000),
		Location:       "CXC_L08_FB_HT_Z03_03",
		TechCode:       25084833300,
		Code:           "Z03_03",
		Sn:             "25084833300_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084833300),
		Location:       "CXC_L08_FB_HT_Z03_03_TS1",
		TechCode:       25084833305,
		Code:           "TS1",
		Sn:             "25084833305_XX",
		SapCode:        "",
	},
	// Z03_04
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084831000),
		Location:       "CXC_L08_FB_HT_Z03_04",
		TechCode:       25084833400,
		Code:           "Z03_04",
		Sn:             "25084833400_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084833400),
		Location:       "CXC_L08_FB_HT_Z03_04_TS1",
		TechCode:       25084833405,
		Code:           "TS1",
		Sn:             "25084833405_XX",
		SapCode:        "",
	},
	// Z03_05
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084831000),
		Location:       "CXC_L08_FB_HT_Z03_05",
		TechCode:       25084833500,
		Code:           "Z03_05",
		Sn:             "25084833500_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084833500),
		Location:       "CXC_L08_FB_HT_Z03_05_TS1",
		TechCode:       25084833505,
		Code:           "TS1",
		Sn:             "25084833505_XX",
		SapCode:        "",
	},
	// Z03_06
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084831000),
		Location:       "CXC_L08_FB_HT_Z03_06",
		TechCode:       25084833600,
		Code:           "Z03_06",
		Sn:             "25084833600_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084833600),
		Location:       "CXC_L08_FB_HT_Z03_06_TS1",
		TechCode:       25084833605,
		Code:           "TS1",
		Sn:             "25084833605_XX",
		SapCode:        "",
	},
	// Z03_07
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084831000),
		Location:       "CXC_L08_FB_HT_Z03_07",
		TechCode:       25084833700,
		Code:           "Z03_07",
		Sn:             "25084833700_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084833700),
		Location:       "CXC_L08_FB_HT_Z03_07_TS1",
		TechCode:       25084833705,
		Code:           "TS1",
		Sn:             "25084833705_XX",
		SapCode:        "",
	},
	// Z03_08
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084831000),
		Location:       "CXC_L08_FB_HT_Z03_08",
		TechCode:       25084833800,
		Code:           "Z03_08",
		Sn:             "25084833800_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084833800),
		Location:       "CXC_L08_FB_HT_Z03_08_TS1",
		TechCode:       25084833805,
		Code:           "TS1",
		Sn:             "25084833805_XX",
		SapCode:        "",
	},
	// Z03_09
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084831000),
		Location:       "CXC_L08_FB_HT_Z03_09",
		TechCode:       25084833900,
		Code:           "Z03_09",
		Sn:             "25084833900_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084833900),
		Location:       "CXC_L08_FB_HT_Z03_09_TS1",
		TechCode:       25084833905,
		Code:           "TS1",
		Sn:             "25084833905_XX",
		SapCode:        "",
	},
	// Z03_10
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25084831000),
		Location:       "CXC_L08_FB_HT_Z03_10",
		TechCode:       25084834000,
		Code:           "Z03_10",
		Sn:             "25084834000_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25084834000),
		Location:       "CXC_L08_FB_HT_Z03_10_TS1",
		TechCode:       25084834005,
		Code:           "TS1",
		Sn:             "25084834005_XX",
		SapCode:        "",
	},
}
