package rP_General

import (
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"github.com/remrafvil/Auriga_API/internal/repositories/rProducts/rP_Types"
)

var ProductTypesData = []rP_Types.ProductTypeData{
	{
		Name: "Drive_AC",
		Features: []rP_Types.FeatureData{
			{"Nominal Current", "IN", "Aac", 1, rModels.ValueTypeNumber},
			{"Maximum Current", "Imax", "Aac", 2, rModels.ValueTypeNumber},
			{"Nominal Power", "PN", "kW", 3, rModels.ValueTypeNumber},
			{"Heavy Duty Current", "IHD", "Aac", 4, rModels.ValueTypeNumber},
			{"Heavy Duty Power", "PHD", "kW", 5, rModels.ValueTypeNumber},
			{"Dissipated Power", "Pd", "kW", 6, rModels.ValueTypeNumber},
			{"Frame", "", "", 7, rModels.ValueTypeString},
		},
	},

	{
		Name: "Rectifier_ISU",
		Features: []rP_Types.FeatureData{
			{"Nominal Current", "IN", "Adc", 1, rModels.ValueTypeNumber},
			{"Maximum Current", "Imax", "Adc", 2, rModels.ValueTypeNumber},
			{"Nominal Power", "PN", "kW", 3, rModels.ValueTypeNumber},
			{"Heavy Duty Current", "IHD", "Adc", 4, rModels.ValueTypeNumber},
			{"Heavy Duty Power", "PHD", "kW", 5, rModels.ValueTypeNumber},
			{"Dissipated Power", "Pd", "kW", 6, rModels.ValueTypeNumber},
			{"Frame", "", "", 7, rModels.ValueTypeString},
		},
	},

	{
		Name: "LCL_Filter",
		Features: []rP_Types.FeatureData{
			{"Nominal Current", "IN", "Aac", 1, rModels.ValueTypeNumber},
			{"Dissipated Power", "Pd", "kW", 6, rModels.ValueTypeNumber},
		},
	},
	{
		Name: "Motor_AC",
		Features: []rP_Types.FeatureData{
			{"Nominal Power", "PN", "kW", 1, rModels.ValueTypeNumber},
			{"Stator Voltage", "Vs", "Vac", 2, rModels.ValueTypeNumber},
			{"Stator Current", "Is", "Aac", 3, rModels.ValueTypeNumber},
			{"Nominal Speed", "Nn", "rpm", 4, rModels.ValueTypeNumber},
			{"Frequency", "f", "Hz", 5, rModels.ValueTypeNumber},
			{"Power Factor", "PF", "", 6, rModels.ValueTypeNumber},
			{"Efficiency", "η", "%", 7, rModels.ValueTypeNumber},
			{"Torque", "T", "Nm", 8, rModels.ValueTypeNumber},
			{"Maximum Speed", "Nmax", "rpm", 9, rModels.ValueTypeNumber},
			{"Water Flow", "Q", "m3/h", 10, rModels.ValueTypeNumber},
			{"Frame", "", "", 11, rModels.ValueTypeString},
			{"Weight", "m", "kg", 12, rModels.ValueTypeNumber},
			{"Cooling", "", "", 13, rModels.ValueTypeString},
			{"TempSensor", "", "", 14, rModels.ValueTypeString},
			{"MountingPosition", "", "", 15, rModels.ValueTypeString},
			{"InsulationClass", "", "", 16, rModels.ValueTypeString},
			{"EfficiencyClass", "", "", 17, rModels.ValueTypeString},
			{"Bearing_DE", "", "", 18, rModels.ValueTypeString},
			{"Bearing_NDE", "", "", 19, rModels.ValueTypeString},
		},
	},

	{
		Name: "TransportSystem",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "HopperLoader",
		Features: []rP_Types.FeatureData{
			{"Volume", "vol.", "l", 1, rModels.ValueTypeNumber},
		},
	},

	{
		Name: "Blender",
		Features: []rP_Types.FeatureData{
			{"Throughput", "TP", "kg/h", 1, rModels.ValueTypeNumber},
			{"Number Components", "N. Comp.", "NC", 2, rModels.ValueTypeNumber},
		},
	},
	{
		Name: "DosingVolu",
		Features: []rP_Types.FeatureData{
			{"Throughput", "TP", "kg/h", 1, rModels.ValueTypeNumber},
			{"Number Components", "N. Comp.", "NC", 2, rModels.ValueTypeNumber},
		},
	},
	{
		Name: "ComponentHopper",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},

	{
		Name: "Factory",
		Features: []rP_Types.FeatureData{
			{"Site Code", "SC", "", 1, rModels.ValueTypeNumber},
			{"Latitude", " lat.", "°", 2, rModels.ValueTypeNumber},
			{"Longitude", "long", "°", 3, rModels.ValueTypeNumber},
			{"Country", "", "", 4, rModels.ValueTypeString},
			{"State", "", "", 5, rModels.ValueTypeString},
			{"City", "", "", 6, rModels.ValueTypeString},
			{"Postal Code", "", "", 7, rModels.ValueTypeString},
			{"Address", "", "", 8, rModels.ValueTypeString},
		},
	},
	{
		Name: "ExtrusionLine",
		Features: []rP_Types.FeatureData{
			{"Throughput", "TP", "kg/h", 1, rModels.ValueTypeNumber},
		},
	},
	{
		Name: "LineCommon",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "LineTransport",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "LineDosing",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "LineExtrusion",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "FB_Die",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "ExtruderDrive",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "ExtruderHeating",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "ExtruderPressure",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "LineStack",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "LineStack_Common",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "LineStack_Cooling",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "LineStack_Cyl",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "LineStack_CylGAP",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},

	{
		Name: "LineWinders",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "LineWinders_Common",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "LineWinders_WS",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "LineWinders_W",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "Extruder",
		Features: []rP_Types.FeatureData{
			{"Type", "Typ.", "", 1, rModels.ValueTypeString},
			{"Diameter", "D", "mm", 2, rModels.ValueTypeNumber},
			{"Throughput", "TP", "kg/h", 3, rModels.ValueTypeNumber},
			{"L/D Ratio", "L/D", "", 4, rModels.ValueTypeNumber}},
	},
	{
		Name: "Feed-Block",
		Features: []rP_Types.FeatureData{
			{"Type", "Typ.", "", 1, rModels.ValueTypeString},
			{"Weight", "m", "kg", 2, rModels.ValueTypeNumber},
			{"Dimensions", "Dim.", "mm", 3, rModels.ValueTypeString},
			{"Laminar Flow Geometry", "LFG", "", 4, rModels.ValueTypeString},
			{"Laminar Flow Dimensions", "LFD", "mm", 5, rModels.ValueTypeString},
			{"Joins Pieces", "JP", "N", 6, rModels.ValueTypeNumber},
			{"Extruders Connected", "EC", "N", 7, rModels.ValueTypeNumber},
			{"Number of Layers", "N", "N", 8, rModels.ValueTypeNumber},
			{"Layer Configuration", "LC", "", 9, rModels.ValueTypeString},
			{"Layer Sequencer", "LS", "", 10, rModels.ValueTypeString},
			{"Layer Profiling Tool", "LPT", "", 11, rModels.ValueTypeString},
			{"Heating Zones", "HZ", "N", 12, rModels.ValueTypeNumber},
			{"Pressure Sensor", "PS", "", 13, rModels.ValueTypeNumber},
		},
	},
	{
		Name: "Die",
		Features: []rP_Types.FeatureData{
			{"Type", "Typ.", "", 1, rModels.ValueTypeString},
			{"Weight", "m", "kg", 2, rModels.ValueTypeNumber},
			{"Dimensions", "Dim.", "mm", 3, rModels.ValueTypeString},
			{"Laminar Flow Geometry", "LFG", "", 4, rModels.ValueTypeString},
			{"Laminar Flow Dimensions", "LFD", "mm", 5, rModels.ValueTypeString},
			{"Lip Width", "LW", "mm", 6, rModels.ValueTypeNumber},
			{"Lip Thickness Range", "LTR", "", 7, rModels.ValueTypeString},
			{"Lip Static Interchangeable", "ISL", "", 8, rModels.ValueTypeString},
			{"Deckling", "D", "", 9, rModels.ValueTypeString},
			{"Gap Adjustment", "GA", "", 10, rModels.ValueTypeString},
			{"Bolt Number", "BN", "N", 11, rModels.ValueTypeNumber},
			{"Bolt Diameter", "BCD", "mm", 12, rModels.ValueTypeNumber},
			{"Flow restrictor bar", "FRB", "", 13, rModels.ValueTypeString},
			{"Flow restrictor bar Adjustment", "FRBA", "", 14, rModels.ValueTypeString},
			{"Flow restrictor bar Bolts", "FRBB", "", 15, rModels.ValueTypeNumber},
			{"Flow restrictor bar Bolts Diameter", "FRBBD", "mm", 16, rModels.ValueTypeNumber},
			{"Heating Zones", "HZ", "N", 17, rModels.ValueTypeNumber},
		},
	},
	{
		Name: "Gearbox",
		Features: []rP_Types.FeatureData{
			{"Type", "Typ.", "", 1, rModels.ValueTypeString},
			{"Ratio", "i", "", 2, rModels.ValueTypeNumber},
			{"Nominal Torque", "Tn", "Nm", 3, rModels.ValueTypeNumber},
			{"Maximum Torque", "Tmax", "Nm", 4, rModels.ValueTypeNumber},
			{"Motor Power", "PM", "kW", 5, rModels.ValueTypeNumber},
			{"Motor Speed", "NM", "rpm", 6, rModels.ValueTypeNumber},
			{"Output Speed", "Nout", "rpm", 7, rModels.ValueTypeNumber},
			{"Weight", "m", "kg", 8, rModels.ValueTypeNumber},
			{"Oil Type", "", "", 9, rModels.ValueTypeString},
			{"Oil Volume", "V", "l", 10, rModels.ValueTypeNumber},
			{"Service Factor", "SF", "", 11, rModels.ValueTypeNumber},
			{"Mounting Position", "", "", 12, rModels.ValueTypeString},
		},
	},
	{
		Name: "HeatingZone",
		Features: []rP_Types.FeatureData{
			{"Feature", "", "", 1, rModels.ValueTypeString},
		},
	},
	{
		Name: "HeatingResistance",
		Features: []rP_Types.FeatureData{
			{"Nominal Power", "PN", "kW", 1, rModels.ValueTypeNumber},
			{"Nominal Voltage", "VN", "Vac", 2, rModels.ValueTypeNumber},
			{"Nominal Current", "IN", "Aac", 3, rModels.ValueTypeNumber},
			{"Resistance", "R", "Ohm", 4, rModels.ValueTypeNumber},
			{"Width", "W", "mm", 5, rModels.ValueTypeNumber},
			{"Length", "L", "mm", 6, rModels.ValueTypeNumber},
			{"Type", "Typ.", "", 7, rModels.ValueTypeString},
		},
	},
	{
		Name: "Blower",
		Features: []rP_Types.FeatureData{
			{"Nominal Power", "PN", "kW", 1, rModels.ValueTypeNumber},
			{"Nominal Voltage", "VN", "Vac", 2, rModels.ValueTypeNumber},
			{"Nominal Current", "IN", "Aac", 3, rModels.ValueTypeNumber},
			{"Nominal Speed", "NN", "rpm", 4, rModels.ValueTypeNumber},
			{"Frequency", "f", "Hz", 5, rModels.ValueTypeNumber},
			{"Air Flow", "Q", "m3/h", 6, rModels.ValueTypeNumber},
		},
	},
	{
		Name: "TemperatureSensor",
		Features: []rP_Types.FeatureData{
			{"Type", "Typ.", "", 1, rModels.ValueTypeString},
			{"Measuring Range", "MR", "°C", 2, rModels.ValueTypeString},
		},
	},
	{
		Name: "PressureSensor",
		Features: []rP_Types.FeatureData{
			{"Type", "Typ.", "", 1, rModels.ValueTypeString},
			{"Measuring Range", "MR", "bar", 2, rModels.ValueTypeString},
			{"Nominal Voltage", "VN", "Vdc", 3, rModels.ValueTypeNumber},
			{"Signal Type", "ST", "", 4, rModels.ValueTypeString},
			{"Calibration", "Cal.", "", 5, rModels.ValueTypeString},
		},
	},
	{
		Name: "PositionSensorsMagnetostrictive",
		Features: []rP_Types.FeatureData{
			{"Type", "Typ.", "", 1, rModels.ValueTypeString},
			{"Bit", "Bit", "", 2, rModels.ValueTypeString},
			{"Measuring Range", "Range", "mm", 3, rModels.ValueTypeNumber},
			{"Accuracy", "Acc", "mm", 4, rModels.ValueTypeNumber},
			{"Gradient", "G", "mm/s", 5, rModels.ValueTypeNumber},
			{"Protocol", "Prot.", "", 6, rModels.ValueTypeString},
		},
	},

	{
		Name: "ScreenChanger",
		Features: []rP_Types.FeatureData{
			{"Type", "Typ.", "", 1, rModels.ValueTypeString},
			{"Dimension", "D", "mm", 2, rModels.ValueTypeNumber},
			{"Number of Plates", "NPl", "N", 3, rModels.ValueTypeNumber},
			{"Max. Pressure", "Pmax", "bar", 4, rModels.ValueTypeNumber},
			{"Max. Temperature", "Tmax", "°C", 5, rModels.ValueTypeNumber},
			{"Max. Flow", "Qmax", "kg/h", 6, rModels.ValueTypeNumber},
		},
	},

	{
		Name: "MeltPump",
		Features: []rP_Types.FeatureData{
			{"Dimension", "D", "mm", 1, rModels.ValueTypeNumber},
			{"Volume", "V", "cm3/rev", 2, rModels.ValueTypeNumber},
			{"Max. Pressure Drop", "Pmax", "bar", 3, rModels.ValueTypeNumber},
			{"Max. Flow", "Qmax", "kg/h", 5, rModels.ValueTypeNumber},
		},
	},
	{
		Name: "DegassingSystem",
		Features: []rP_Types.FeatureData{
			{"Type", "Typ.", "", 1, rModels.ValueTypeString},
			{"Max. Flow", "Qmax", "m3/h", 2, rModels.ValueTypeNumber},
			{"Max. Vacuum", "Vmax", "mbar", 3, rModels.ValueTypeNumber},
		},
	},
	{
		Name: "VacuumPump",
		Features: []rP_Types.FeatureData{
			{"Type", "Typ.", "", 1, rModels.ValueTypeString},
			{"Max. Flow", "Qmax", "m3/h", 3, rModels.ValueTypeNumber},
			{"Max. Vacuum", "Vmax", "mbar", 4, rModels.ValueTypeNumber},
		},
	},
	{
		Name: "WaterCoolingUnit",
		Features: []rP_Types.FeatureData{
			{"Type", "Typ.", "", 1, rModels.ValueTypeString},
			{"Max. Pressure", "Pmax", "bar", 2, rModels.ValueTypeNumber},
			{"Max. Temperature", "Tmax", "°C", 3, rModels.ValueTypeNumber},
			{"Min. Temperature", "Tmin", "°C", 4, rModels.ValueTypeNumber},
			{"Nominal Power", "PN", "kW", 5, rModels.ValueTypeNumber},
			{"Nominal Current", "IN", "Aac", 6, rModels.ValueTypeNumber},
			{"Nominal Voltage", "VN", "Vac", 7, rModels.ValueTypeNumber},
			{"Control Voltage", "VC", "Vdc", 8, rModels.ValueTypeNumber},
		},
	},
}
