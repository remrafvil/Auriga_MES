package rInfluxQuery

import (
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

/* func InfluxQueryInitToCreate(c *gorm.DB) { //
	c.Create(&rModels.MrInfluxGrafanaQuery{
		Name:        "FactoryEnergyConps",
		Description: "Factory Energy Compsumption",
		Query:       `from(bucket: "FSP") |> range(start: v.timeRangeStart, stop: v.timeRangeStop) |> filter(fn: (r) => r["name"] == "Line_02") |> filter(fn: (r) => r["EL_Lv1"] == "2_Dosing") |> filter(fn: (r) => r["EL_Lv2"] == "Doser_01") |> filter(fn: (r) => r["EL_Lv3"] == "C1" or r["EL_Lv3"] == "C2" or r["EL_Lv3"] == "C3") |> filter(fn: (r) => r["_field"] == "HopperNr" or r["_field"] == "TotalDispensed_kg") |> aggregateWindow(every: v.windowPeriod, fn: mean, createEmpty: false) |> yield(name: "mean")`,
	})
}*/

func InfluxQueryInitToCreate(c *gorm.DB) { //
	c.Create(&rModels.MrInfluxGrafanaQuery{
		Name:        "CXC_Line_08_Energy",
		Description: "CXC Line 08 Energy Compsumption",
		Query: `from(bucket: "CXC")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["name"] == "Electrical")
  |> filter(fn: (r) => r["EL_Lv1"] == "TR5")
  |> filter(fn: (r) => r["EL_Lv2"] == "L8")
  |> filter(fn: (r) => r["EL_Lv3"] == "Energy")
  |> filter(fn: (r) => r["_field"] == "TotalApparentEnergyPositive_KVAh" or r["_field"] == "ActiveEnergyPositive_kWh")
  |> aggregateWindow(every: v.windowPeriod, fn: mean, createEmpty: false)
  |> yield(name: "mean")`,
	})
}

/* func InfluxQueryInitToCreate(c *gorm.DB) { //
	c.Create(&rModels.MrInfluxGrafanaQuery{
		Name:        "CXC_Line_08_Power",
		Description: "CXC Line 08 Power Consumption",
		Query: `from(bucket: "CXC")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r["name"] == "Electrical")
  |> filter(fn: (r) => r["EL_Lv1"] == "TR5")
  |> filter(fn: (r) => r["EL_Lv2"] == "L8")
  |> filter(fn: (r) => r["EL_Lv3"] == "Total")
  |> filter(fn: (r) => r["_field"] == "ApparentPower_kVA" or r["_field"] == "ActivePower_kW")
  |> aggregateWindow(every: v.windowPeriod, fn: mean, createEmpty: false)
  |> yield(name: "mean")`,
	})
} */
