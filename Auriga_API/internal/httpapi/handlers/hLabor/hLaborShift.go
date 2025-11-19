package hLabor

import (
	"net/http"
	"strconv"

	"github.com/remrafvil/Auriga_API/internal/services/sLabor"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *handler) GetAllShifts(c echo.Context) error {
	shifts, err := h.service.GetAllShifts()
	if err != nil {
		h.logger.Error("Error getting shifts", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, shifts)
}

func (h *handler) GetShift(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shift ID"})
	}

	shift, err := h.service.GetShiftByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Shift not found"})
	}

	return c.JSON(http.StatusOK, shift)
}

func (h *handler) CreateShift(c echo.Context) error {
	var req sLabor.CreateShiftRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validación usando Echo con CustomValidator
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Solo validaciones de lógica de negocio
	if req.StartTime.After(req.EndTime) || req.StartTime.Equal(req.EndTime) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Start time must be before end time"})
	}

	shift, err := h.service.CreateShift(req)
	if err != nil {
		h.logger.Error("Error creating shift", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create shift"})
	}

	return c.JSON(http.StatusCreated, shift)
}

func (h *handler) UpdateShift(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shift ID"})
	}

	var req sLabor.UpdateShiftRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validación usando Echo con CustomValidator
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Solo validaciones de lógica de negocio
	if !req.StartTime.IsZero() && !req.EndTime.IsZero() &&
		(req.StartTime.After(req.EndTime) || req.StartTime.Equal(req.EndTime)) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Start time must be before end time"})
	}

	shift, err := h.service.UpdateShift(uint(id), req)
	if err != nil {
		h.logger.Error("Error updating shift", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update shift"})
	}

	return c.JSON(http.StatusOK, shift)
}

func (h *handler) DeleteShift(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid shift ID"})
	}

	if err := h.service.DeleteShift(uint(id)); err != nil {
		h.logger.Error("Error deleting shift", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete shift"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Shift deleted successfully"})
}
