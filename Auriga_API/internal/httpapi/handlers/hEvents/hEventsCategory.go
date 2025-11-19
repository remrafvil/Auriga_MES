package hEvents

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h *handler) GetAllCategoriesWithEventTypes(c echo.Context) error {
	ctx := c.Request().Context()

	response, err := h.service.GetAllCategoriesWithEventTypes(ctx)

	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get event categories")
	}

	return c.JSON(http.StatusOK, response)
}
func (h *handler) GetCategoryWithEventTypesByName(c echo.Context) error {
	ctx := c.Request().Context()
	u := new(mhLineEvents)
	u.Category = c.Request().Header.Get("Category")

	if strings.TrimSpace(u.Category) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "category name cannot be empty")
	}

	categoryResponse, err := h.service.GetCategoryWithEventTypesByName(ctx, u.Category)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get category with event types")
	}

	if categoryResponse == nil {
		return echo.NewHTTPError(http.StatusNotFound, "category not found")
	}

	return c.JSON(http.StatusOK, categoryResponse)
}
