package sProducts

import (
	"context"
	"errors"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ProductTypeResponse - Respuesta para listado básico
type ProductTypeResponse struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ProductTypeDetailResponse - Respuesta con características específicas
type ProductTypeDetailResponse struct {
	ID       uint                     `json:"id"`
	Name     string                   `json:"name"`
	Features []ProductFeatureResponse `json:"features"`
}

type ProductFeatureResponse struct {
	ID        uint              `json:"id"`
	Name      string            `json:"name"`
	Symbol    string            `json:"symbol"`
	Unit      string            `json:"unit"`
	ValueType rModels.ValueType `json:"value_type"`
}

func (s *service) GetProductTypesList(ctx context.Context) ([]ProductTypeResponse, error) {
	s.logger.Info("Getting product types list")

	productTypes, err := s.repository.ProductTypeFindAll(ctx)
	if err != nil {
		s.logger.Error("Failed to get product types list", zap.Error(err))
		return nil, err
	}

	response := make([]ProductTypeResponse, len(productTypes))
	for i, pt := range productTypes {
		response[i] = ProductTypeResponse{
			ID:   pt.ID,
			Name: pt.Name,
		}
	}

	return response, nil
}

func (s *service) GetProductTypeByID(ctx context.Context, id uint) (*ProductTypeDetailResponse, error) {
	s.logger.Info("Getting product type by ID", zap.Uint("id", id))

	productType, err := s.repository.ProductTypeFindByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Warn("Product type not found", zap.Uint("id", id))
			return nil, errors.New("product type not found")
		}
		s.logger.Error("Failed to get product type", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}

	return s.mapToDetailResponse(productType), nil
}

func (s *service) GetAllProductTypesWithFeatures(ctx context.Context) ([]ProductTypeDetailResponse, error) {
	s.logger.Info("Getting all product types with features")

	productTypes, err := s.repository.ProductTypeFindAllWithFeatures(ctx)
	if err != nil {
		s.logger.Error("Failed to get product types with features", zap.Error(err))
		return nil, err
	}

	response := make([]ProductTypeDetailResponse, len(productTypes))
	for i, pt := range productTypes {
		response[i] = *s.mapToDetailResponse(&pt)
	}

	return response, nil
}

func (s *service) mapToDetailResponse(pt *rModels.MrProductType) *ProductTypeDetailResponse {
	features := make([]ProductFeatureResponse, len(pt.Features))
	for i, f := range pt.Features {
		features[i] = ProductFeatureResponse{
			ID:        f.ID,
			Name:      f.Name,
			Symbol:    f.Symbol,
			Unit:      f.Unit,
			ValueType: f.ValueType,
		}
	}

	return &ProductTypeDetailResponse{
		ID:       pt.ID,
		Name:     pt.Name,
		Features: features,
	}
}
