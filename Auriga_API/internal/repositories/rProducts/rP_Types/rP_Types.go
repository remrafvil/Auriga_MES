package rP_Types

import "github.com/remrafvil/Auriga_API/internal/repositories/rModels"

type ProductCreationData struct {
	Name          string
	ProductType   string
	Manufacturer  string
	Family        string
	Description   string
	FeatureValues map[string]string
}
type FeatureData struct {
	Name      string
	Symbol    string
	Unit      string
	Order     uint
	ValueType rModels.ValueType
}

type ProductTypeData struct {
	Name     string
	Features []FeatureData
}
