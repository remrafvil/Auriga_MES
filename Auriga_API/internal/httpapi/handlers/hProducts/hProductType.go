// handler/mr_product_type_handler.go
package hProducts

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// GetProductTypesList godoc
// @Summary Get product types list
// @Description Get a list of product types with ID and Name only
// @Tags product-types
// @Accept json
// @Produce json
// @Success 200 {array} service.ProductTypeResponse
// @Failure 500 {object} map[string]string
// @Router /product-types [get]
func (h *handler) GetProductTypesList(c echo.Context) error {
	ctx := c.Request().Context()

	log.Println("*************************************************************************llega aqui")
	productTypes, err := h.service.GetProductTypesList(ctx)
	if err != nil {
		h.logger.Error("Failed to get product types list", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get product types list",
		})
	}

	return c.JSON(http.StatusOK, productTypes)
}

// GetProductTypeByID godoc
// @Summary Get product type by ID
// @Description Get complete product type with all features
// @Tags product-types
// @Accept json
// @Produce json
// @Param id path int true "Product Type ID"
// @Success 200 {object} service.ProductTypeDetailResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /product-types/{id} [get]
func (h *handler) GetProductTypeByID(c echo.Context) error {
	ctx := c.Request().Context()

	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		h.logger.Warn("Invalid product type ID", zap.String("id", idStr))
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid product type ID",
		})
	}

	productType, err := h.service.GetProductTypeByID(ctx, uint(id))
	if err != nil {
		if err.Error() == "product type not found" {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Product type not found",
			})
		}
		h.logger.Error("Failed to get product type", zap.Uint("id", uint(id)), zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get product type",
		})
	}

	return c.JSON(http.StatusOK, productType)
}

// GetAllProductTypesWithFeatures godoc
// @Summary Get all product types with features
// @Description Get all product types with complete features information
// @Tags product-types
// @Accept json
// @Produce json
// @Success 200 {array} service.ProductTypeDetailResponse
// @Failure 500 {object} map[string]string
// @Router /product-types/complete [get]
func (h *handler) GetAllProductTypesWithFeatures(c echo.Context) error {
	ctx := c.Request().Context()

	productTypes, err := h.service.GetAllProductTypesWithFeatures(ctx)
	if err != nil {
		h.logger.Error("Failed to get product types with features", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get product types with features",
		})
	}

	return c.JSON(http.StatusOK, productTypes)
}
