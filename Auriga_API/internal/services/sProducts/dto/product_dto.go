// dto/product_dto.go   (Data Transfer Objects)
package dto

type CreateProductRequest struct {
	Name          string            `json:"name" validate:"required"`
	ProductType   string            `json:"product_type" validate:"required"`
	Manufacturer  string            `json:"manufacturer" validate:"required"`
	Family        string            `json:"family" validate:"required"`
	Description   string            `json:"description"`
	FeatureValues map[string]string `json:"feature_values"`
}

type UpdateProductRequest struct {
	Name          string            `json:"name"`
	ProductType   string            `json:"product_type"`
	Manufacturer  string            `json:"manufacturer"`
	Family        string            `json:"family"`
	Description   string            `json:"description"`
	FeatureValues map[string]string `json:"feature_values"`
}
