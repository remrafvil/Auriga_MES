package hAssets

import (
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) DosingSystemByLine(c echo.Context) error {

	lineIDStr := c.Request().Header.Get("ProdLine_ID")

	log.Println("ProdLine_ID", lineIDStr)
	//lineIDStr := c.QueryParam("line_id")
	if lineIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "line_id parameter is required",
		})
	}

	lineID, err := strconv.ParseUint(lineIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid line_id parameter",
		})
	}

	use, err := h.service.GetDosingSystemByLine(c.Request().Context(), uint(lineID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get dosing system: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, use)
}

func (h *handler) DoserComponents(c echo.Context) error {

	doserIDStr := c.Request().Header.Get("Doser_ID")
	//doserIDStr := c.QueryParam("doser_id")
	if doserIDStr == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "doser_id parameter is required",
		})
	}

	doserID, err := strconv.ParseUint(doserIDStr, 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid doser_id parameter",
		})
	}

	components, err := h.service.GetDoserComponents(c.Request().Context(), uint(doserID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "failed to get doser components: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, components)
}
