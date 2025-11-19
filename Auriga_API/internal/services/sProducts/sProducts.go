package sProducts

import (
	"context"
	"errors"
	"fmt"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"github.com/remrafvil/Auriga_API/internal/repositories/rProducts"
	"github.com/remrafvil/Auriga_API/internal/services/sProducts/dto"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Service interface {
	GetProductsList(ctx context.Context) ([]ProductResponse, error)
	GetProductByID(ctx context.Context, id uint) (*ProductDetailResponse, error)
	CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*ProductDetailResponse, error)
	UpdateProduct(ctx context.Context, id uint, req *dto.UpdateProductRequest) (*ProductDetailResponse, error)
	DeleteProduct(ctx context.Context, id uint) error

	GetProductTypesList(ctx context.Context) ([]ProductTypeResponse, error)
	GetProductTypeByID(ctx context.Context, id uint) (*ProductTypeDetailResponse, error)
	GetAllProductTypesWithFeatures(ctx context.Context) ([]ProductTypeDetailResponse, error)
}

type service struct {
	repository rProducts.Repository
	logger     *zap.Logger
}

func New(repository rProducts.Repository, logger *zap.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

/* // Request structures
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
} */

// Response structures (las mismas de antes)
type ProductResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	ProductType  string `json:"product_type"`
	Manufacturer string `json:"manufacturer"`
	Family       string `json:"family"`
	Description  string `json:"description"`
}

type ProductDetailResponse struct {
	ID            uint              `json:"id"`
	Name          string            `json:"name"`
	ProductType   string            `json:"product_type"`
	Manufacturer  string            `json:"manufacturer"`
	Family        string            `json:"family"`
	Description   string            `json:"description"`
	FeatureValues map[string]string `json:"feature_values"`
}

func (s *service) GetProductsList(ctx context.Context) ([]ProductResponse, error) {
	s.logger.Info("Getting products list")

	products, err := s.repository.ProductFindAll(ctx)
	if err != nil {
		s.logger.Error("Failed to get products list", zap.Error(err))
		return nil, err
	}

	response := make([]ProductResponse, len(products))
	for i, product := range products {
		response[i] = ProductResponse{
			ID:           product.ID,
			Name:         product.Name,
			ProductType:  product.ProductType.Name,
			Manufacturer: product.Manufacturer,
			Family:       product.Family,
			Description:  product.Description,
		}
	}

	return response, nil
}

func (s *service) GetProductByID(ctx context.Context, id uint) (*ProductDetailResponse, error) {
	s.logger.Info("Getting product by ID", zap.Uint("id", id))

	product, err := s.repository.ProductFindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Product not found", zap.Uint("id", id))
			return nil, errors.New("product not found")
		}
		s.logger.Error("Failed to get product", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	featureValues := s.extractFeatureValues(product)

	return &ProductDetailResponse{
		ID:            product.ID,
		Name:          product.Name,
		ProductType:   product.ProductType.Name,
		Manufacturer:  product.Manufacturer,
		Family:        product.Family,
		Description:   product.Description,
		FeatureValues: featureValues,
	}, nil
}

func (s *service) CreateProduct(ctx context.Context, req *dto.CreateProductRequest) (*ProductDetailResponse, error) {
	s.logger.Info("Creating new product", zap.String("name", req.Name))

	product := &rModels.MrProduct{
		Name:         req.Name,
		Manufacturer: req.Manufacturer,
		Family:       req.Family,
		Description:  req.Description,
		ProductType: rModels.MrProductType{
			Name: req.ProductType,
		},
	}

	if err := s.repository.ProductCreate(ctx, product, req.FeatureValues); err != nil {
		s.logger.Error("Failed to create product", zap.String("name", req.Name), zap.Error(err))
		return nil, fmt.Errorf("failed to create product: %v", err)
	}

	// Obtener el producto creado para devolverlo
	return s.GetProductByID(ctx, product.ID)
}

func (s *service) UpdateProduct(ctx context.Context, id uint, req *dto.UpdateProductRequest) (*ProductDetailResponse, error) {
	s.logger.Info("Updating product", zap.Uint("id", id))

	product := &rModels.MrProduct{
		Name:         req.Name,
		Manufacturer: req.Manufacturer,
		Family:       req.Family,
		Description:  req.Description,
		ProductType: rModels.MrProductType{
			Name: req.ProductType,
		},
	}

	if err := s.repository.ProductUpdate(ctx, id, product, req.FeatureValues); err != nil {
		s.logger.Error("Failed to update product", zap.Uint("id", id), zap.Error(err))
		return nil, fmt.Errorf("failed to update product: %v", err)
	}

	// Obtener el producto actualizado para devolverlo
	return s.GetProductByID(ctx, id)
}

func (s *service) DeleteProduct(ctx context.Context, id uint) error {
	s.logger.Info("Deleting product", zap.Uint("id", id))

	if err := s.repository.ProductDelete(ctx, id); err != nil {
		s.logger.Error("Failed to delete product", zap.Uint("id", id), zap.Error(err))
		return fmt.Errorf("failed to delete product: %v", err)
	}

	return nil
}

func (s *service) extractFeatureValues(product *rModels.MrProduct) map[string]string {
	featureValues := make(map[string]string)

	for _, fv := range product.FeatureValues {
		featureName := fv.ProductFeature.Name
		var value string

		switch fv.ProductFeature.ValueType {
		case rModels.ValueTypeNumber:
			// Para valores numéricos, formatear sin decimales si es entero
			if fv.NumberValue == float64(int(fv.NumberValue)) {
				value = fmt.Sprintf("%.0f", fv.NumberValue)
			} else {
				value = fmt.Sprintf("%.1f", fv.NumberValue)
			}
		case rModels.ValueTypeString:
			value = fv.StringValue
		default:
			// Por defecto usar string value
			value = fv.StringValue
		}

		featureValues[featureName] = value
	}

	return featureValues
}
