package hLabor

import (
	"net/http"
	"strconv"
	"time"

	"github.com/remrafvil/Auriga_API/internal/services/sLabor"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *handler) GetEmployeeAssignments(c echo.Context) error {
	employeeID, err := strconv.ParseUint(c.Param("employeeId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid employee ID"})
	}

	assignments, err := h.service.GetEmployeeAssignments(uint(employeeID))
	if err != nil {
		h.logger.Error("Error getting employee assignments", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get assignments"})
	}

	return c.JSON(http.StatusOK, assignments)
}

func (h *handler) GetTeamAssignments(c echo.Context) error {
	teamID, err := strconv.ParseUint(c.Param("teamId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
	}

	assignments, err := h.service.GetTeamAssignments(uint(teamID))
	if err != nil {
		h.logger.Error("Error getting team assignments", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get assignments"})
	}

	return c.JSON(http.StatusOK, assignments)
}

func (h *handler) CreateIndividualAssignment(c echo.Context) error {
	var req sLabor.CreateShiftAssignmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validación usando Echo con CustomValidator
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Solo validaciones de lógica de negocio
	if !req.EndDate.IsZero() && req.StartDate.After(req.EndDate) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Start date must be before end date"})
	}

	assignment, err := h.service.CreateIndividualAssignment(req)
	if err != nil {
		h.logger.Error("Error creating individual assignment", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create assignment"})
	}

	return c.JSON(http.StatusCreated, assignment)
}

func (h *handler) CreateTeamAssignment(c echo.Context) error {
	var req sLabor.CreateShiftAssignmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validación usando Echo con CustomValidator
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Solo validaciones de lógica de negocio
	if !req.EndDate.IsZero() && req.StartDate.After(req.EndDate) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Start date must be before end date"})
	}

	assignments, err := h.service.CreateTeamAssignment(req)
	if err != nil {
		h.logger.Error("Error creating team assignment", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create assignments"})
	}

	return c.JSON(http.StatusCreated, assignments)
}

func (h *handler) CreateBulkAssignments(c echo.Context) error {
	var req sLabor.BulkShiftAssignmentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validación usando Echo con CustomValidator
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Solo validaciones de lógica de negocio
	if !req.EndDate.IsZero() && req.StartDate.After(req.EndDate) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Start date must be before end date"})
	}

	assignments, err := h.service.CreateBulkAssignments(req)
	if err != nil {
		h.logger.Error("Error creating bulk assignments", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create assignments"})
	}

	return c.JSON(http.StatusCreated, assignments)
}

func (h *handler) GetCurrentEmployeeAssignment(c echo.Context) error {
	employeeID, err := strconv.ParseUint(c.Param("employeeId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid employee ID"})
	}

	dateParam := c.QueryParam("date")
	var date time.Time
	if dateParam != "" {
		date, err = time.Parse("2006-01-02", dateParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date format. Use YYYY-MM-DD"})
		}
	}

	assignment, err := h.service.GetCurrentEmployeeAssignment(uint(employeeID), date)
	if err != nil {
		h.logger.Error("Error getting current assignment", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get current assignment"})
	}

	return c.JSON(http.StatusOK, assignment)
}
