package rAcxcL8

import "github.com/remrafvil/Auriga_API/internal/repositories/rAssets/rA_Types"

var AssetsToCreate_CXC_L8_ExtrA = []rA_Types.AssetCreationData{

	// Subsistemas (nivel 2)
	{
		ProductName:    "KS45-24D_B",
		ParentTechCode: ptrUint(25083000),
		Location:       "CXC_L08_ExtrA",
		TechCode:       25083200,
		Code:           "Extr_A",
		Sn:             "Kuhne_141225-45-01E1.0.00.0.3",
		SapCode:        "",
	},

	// Extruder Drive System
	{
		ProductName:    "3_ExtruderDrive",
		ParentTechCode: ptrUint(25083200),
		Location:       "CXC_L08_ExtrA_Extr",
		TechCode:       25083200020,
		Code:           "Extr",
		Sn:             "25083200010_XX",
		SapCode:        "",
	},
	// Extruder VFD
	{
		ProductName:    "ACS800-01-0040-3",
		ParentTechCode: ptrUint(25083200020),
		Location:       "CXC_L08_ExtrA_Extr_Dr",
		TechCode:       25083200021,
		Code:           "Dr",
		Sn:             "1120701872",
		SapCode:        "",
	},
	// Extruder Motor
	{
		ProductName:    "MDFMAXX 180-12",
		ParentTechCode: ptrUint(25083200020),
		Location:       "CXC_L08_ExtrA_Extr_Mt",
		TechCode:       25083200022,
		Code:           "Mt",
		Sn:             "25083208002_1204",
		SapCode:        "",
	},

	// Extruder Gearbox
	{
		ProductName:    "RSG_116443-002",
		ParentTechCode: ptrUint(25083200020),
		Location:       "CXC_L08_ExtrA_Extr_Gr",
		TechCode:       25083200024,
		Code:           "Gr",
		Sn:             "25083208001_2310",
		SapCode:        "",
	},

	// Extruder Heating
	{
		ProductName:    "3_ExtruderHeating",
		ParentTechCode: ptrUint(25083200),
		Location:       "CXC_L08_ExtrA_HT",
		TechCode:       25083201000,
		Code:           "HT",
		Sn:             "25083201000_XX",
		SapCode:        "",
	},
	// Z01_01
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25083201000),
		Location:       "CXC_L08_ExtrA_HT_Z01_01",
		TechCode:       25083201100,
		Code:           "Z01_01",
		Sn:             "25083201100_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25083201100),
		Location:       "CXC_L08_ExtrA_HT_Z01_01_TS1",
		TechCode:       25083201105,
		Code:           "TS1",
		Sn:             "25083201105_XX",
		SapCode:        "",
	},
	// Z01_02
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25083201000),
		Location:       "CXC_L08_ExtrA_HT_Z01_02",
		TechCode:       25083201200,
		Code:           "Z01_02",
		Sn:             "25083201200_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25083201200),
		Location:       "CXC_L08_ExtrA_HT_Z01_02_TS1",
		TechCode:       25083201205,
		Code:           "TS1",
		Sn:             "25083201205_XX",
		SapCode:        "",
	},
	// Z01_03

	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25083201000),
		Location:       "CXC_L08_ExtrA_HT_Z01_03",
		TechCode:       25083201300,
		Code:           "Z01_03",
		Sn:             "25083201300_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25083201300),
		Location:       "CXC_L08_ExtrA_HT_Z01_03_TS1",
		TechCode:       25083201305,
		Code:           "TS1",
		Sn:             "25083201305_XX",
		SapCode:        "",
	},
	// Z01_04

	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25083201000),
		Location:       "CXC_L08_ExtrA_HT_Z01_04",
		TechCode:       25083201400,
		Code:           "Z01_04",
		Sn:             "25083201400_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25083201400),
		Location:       "CXC_L08_ExtrA_HT_Z01_04_TS1",
		TechCode:       25083201405,
		Code:           "TS1",
		Sn:             "25083201405_XX",
		SapCode:        "",
	},
	// Z02_01
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25083201000),
		Location:       "CXC_L08_ExtrA_HT_Z02_01",
		TechCode:       25083202100,
		Code:           "Z02_01",
		Sn:             "25083202100_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25083202100),
		Location:       "CXC_L08_ExtrA_HT_Z02_01_TS1",
		TechCode:       25083202105,
		Code:           "TS1",
		Sn:             "25083202105_XX",
		SapCode:        "",
	},

	// Z02_02
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25083201000),
		Location:       "CXC_L08_ExtrA_HT_Z02_02",
		TechCode:       25083202200,
		Code:           "Z02_02",
		Sn:             "25083202200_XX",
		SapCode:        "",
	},

	// Z03_01
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25083201000),
		Location:       "CXC_L08_ExtrA_HT_Z03_01",
		TechCode:       25083203100,
		Code:           "Z03_01",
		Sn:             "25083203100_XX",
		SapCode:        "",
	},

	// Z03_02
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25083201000),
		Location:       "CXC_L08_ExtrA_HT_Z03_02",
		TechCode:       25083203200,
		Code:           "Z03_02",
		Sn:             "25083203200_XX",
		SapCode:        "",
	},

	// Z03_03
	{
		ProductName:    "HeatingZone",
		ParentTechCode: ptrUint(25083201000),
		Location:       "CXC_L08_ExtrA_HT_Z03_03",
		TechCode:       25083203300,
		Code:           "Z03_03",
		Sn:             "25083203300_XX",
		SapCode:        "",
	},

	// Extruder Pressure & Melt
	{
		ProductName:    "3_ExtruderPressure",
		ParentTechCode: ptrUint(25083200),
		Location:       "CXC_L08_ExtrA_PR",
		TechCode:       25083204000,
		Code:           "PR",
		Sn:             "25083204000_XX",
		SapCode:        "",
	},
	// PT1
	{
		ProductName:    "IE0-A-7-M-B05C-1-4-0-P-E-D141430B000X0",
		ParentTechCode: ptrUint(25083204000),
		Location:       "CXC_L08_ExtrA_PR_PT1",
		TechCode:       25083204101,
		Code:           "PT1",
		Sn:             "Gef-18410174",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25083204000),
		Location:       "CXC_L08_ExtrA_PR_Mt1",
		TechCode:       25083204401,
		Code:           "Mt1",
		Sn:             "25083204401_XX",
		SapCode:        "",
	},
	{
		ProductName:    "Thermocouple_J",
		ParentTechCode: ptrUint(25083204000),
		Location:       "CXC_L08_ExtrA_PR_Mt2",
		TechCode:       25083204501,
		Code:           "Mt2",
		Sn:             "25083204501_XX",
		SapCode:        "",
	},

	// Screen Changer
	{
		ProductName:    "T005616",
		ParentTechCode: ptrUint(25083200),
		Location:       "CXC_L08_ExtrA",
		TechCode:       25083205000,
		Code:           "SC",
		Sn:             "25083205000_001",
		SapCode:        "",
	},
}
