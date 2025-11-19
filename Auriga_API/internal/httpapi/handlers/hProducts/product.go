package hProducts

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/remrafvil/Auriga_API/internal/services/sProducts/dto"
	"go.uber.org/zap"
)

type responseMessage struct {
	Message string `json:"message"`
}

type Handler interface {
	GetProductsList(c echo.Context) error
	GetProductByID(c echo.Context) error
	CreateProduct(c echo.Context) error
	UpdateProduct(c echo.Context) error
	DeleteProduct(c echo.Context) error

	GetProductTypesList(c echo.Context) error
	GetProductTypeByID(c echo.Context) error
	GetAllProductTypesWithFeatures(c echo.Context) error
}

// Request structures para el handler
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

// GetProductsList godoc
// @Summary Get products list
// @Description Get a list of products with basic information
// @Tags products
// @Accept json
// @Produce json
// @Success 200 {array} service.ProductResponse
// @Failure 500 {object} map[string]string
// @Router /products [get]
func (h *handler) GetProductsList(c echo.Context) error {
	ctx := c.Request().Context()

	products, err := h.service.GetProductsList(ctx)
	if err != nil {
		h.logger.Error("Failed to get products list", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get products list",
		})
	}

	return c.JSON(http.StatusOK, products)
}

// GetProductByID godoc
// @Summary Get product by ID
// @Description Get complete product information with feature values
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} service.ProductDetailResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [get]
func (h *handler) GetProductByID(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid product ID", zap.String("id", idStr))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product ID",
		})
	}

	product, err := h.service.GetProductByID(ctx, uint(id))
	if err != nil {
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Product not found",
			})
		}
		h.logger.Error("Failed to get product", zap.Uint("id", uint(id)), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get product",
		})
	}

	return c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with feature values
// @Tags products
// @Accept json
// @Produce json
// @Param product body dto.CreateProductRequest true "Product data"
// @Success 201 {object} service.ProductDetailResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products [post]
func (h *handler) CreateProduct(c echo.Context) error {
	ctx := c.Request().Context()

	var req dto.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Validaciones básicas
	if req.Name == "" || req.ProductType == "" || req.Manufacturer == "" || req.Family == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Name, product_type, manufacturer and family are required",
		})
	}

	product, err := h.service.CreateProduct(ctx, &req)
	if err != nil {
		h.logger.Error("Failed to create product", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, product)
}

// UpdateProduct godoc
// @Summary Update a product
// @Description Update an existing product and its feature values
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body dto.UpdateProductRequest true "Product data"
// @Success 200 {object} service.ProductDetailResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [put]
func (h *handler) UpdateProduct(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid product ID", zap.String("id", idStr))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product ID",
		})
	}

	var req dto.UpdateProductRequest
	if err := c.Bind(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	product, err := h.service.UpdateProduct(ctx, uint(id), &req)
	if err != nil {
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Product not found",
			})
		}
		h.logger.Error("Failed to update product", zap.Uint("id", uint(id)), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, product)
}

// DeleteProduct godoc
// @Summary Delete a product
// @Description Delete a product and its feature values
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /products/{id} [delete]
func (h *handler) DeleteProduct(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid product ID", zap.String("id", idStr))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product ID",
		})
	}

	if err := h.service.DeleteProduct(ctx, uint(id)); err != nil {
		if err.Error() == "product not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Product not found",
			})
		}
		h.logger.Error("Failed to delete product", zap.Uint("id", uint(id)), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
