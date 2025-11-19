package hAssets

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type responseMessage struct {
	Message string `json:"message"`
}

type Handler interface {
	AssetShow(c echo.Context) error
	AssetList(c echo.Context) error
}

type mhAsset struct {
	ID       uint   `json:"id" form:"id" query:"id"`
	ParentID uint   `json:"parent_id" form:"parent_id" query:"parent_id"`
	Code     string `json:"code" form:"code" query:"code"`
}

func (h *handler) AssetShow(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	log.Println(id)
	use, err := h.service.AssetInfo(uint(id))

	log.Println("Error aqui:", err)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Activo no existeeee"})
	}
	return c.JSON(http.StatusOK, use)
}

func (h *handler) AssetShowDetail(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	log.Println(id)
	use, err := h.service.GetAssetWithSimplifiedProduct(uint(id))

	log.Println("Error aqui:", err)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Activo no existeeee"})
	}
	return c.JSON(http.StatusOK, use)
}

func (h *handler) AssetShowHierarchi(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	log.Println(id)
	use, err := h.service.GetAssetHierarchy(uint(id))

	log.Println("Error aqui:", err)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Activo no existeeee"})
	}
	return c.JSON(http.StatusOK, use)
}

func (h *handler) AssetList(c echo.Context) error {
	// ✅ Obtener el contexto estándar ya preparado por el middleware
	ctx := c.Request().Context()

	// ✅ Llamar directamente al servicio - el contexto ya tiene toda la información
	list, err := h.service.AssetList(ctx)
	if err != nil {
		h.logger.Error("Service call failed",
			zap.Error(err),
			zap.String("handler", "AssetList"))
		return err
	}

	h.logger.Info("Asset list completed successfully",
		zap.Int("count", len(list)),
		zap.String("handler", "AssetList"))

	return c.JSON(http.StatusOK, list)
}
