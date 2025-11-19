package hLabor

import (
	"net/http"
	"strconv"

	"github.com/remrafvil/Auriga_API/internal/services/sLabor"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func (h *handler) GetAllTeams(c echo.Context) error {
	teams, err := h.service.GetAllTeams()
	if err != nil {
		h.logger.Error("Error getting teams", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}
	return c.JSON(http.StatusOK, teams)
}

func (h *handler) GetTeam(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
	}

	team, err := h.service.GetTeamByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Team not found"})
	}

	return c.JSON(http.StatusOK, team)
}

func (h *handler) CreateTeam(c echo.Context) error {
	var req sLabor.CreateTeamRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validación usando Echo con CustomValidator
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	team, err := h.service.CreateTeam(req)
	if err != nil {
		h.logger.Error("Error creating team", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create team"})
	}

	return c.JSON(http.StatusCreated, team)
}

func (h *handler) UpdateTeam(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
	}

	var req sLabor.UpdateTeamRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validación usando Echo con CustomValidator
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	team, err := h.service.UpdateTeam(uint(id), req)
	if err != nil {
		h.logger.Error("Error updating team", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update team"})
	}

	return c.JSON(http.StatusOK, team)
}

func (h *handler) DeleteTeam(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
	}

	if err := h.service.DeleteTeam(uint(id)); err != nil {
		h.logger.Error("Error deleting team", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete team"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Team deleted successfully"})
}

func (h *handler) AddTeamMember(c echo.Context) error {
	var req sLabor.AddTeamMemberRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request body"})
	}

	// Validación usando Echo con CustomValidator
	if err := c.Validate(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	teamMember, err := h.service.AddTeamMember(req)
	if err != nil {
		h.logger.Error("Error adding team member", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add team member"})
	}

	return c.JSON(http.StatusCreated, teamMember)
}

func (h *handler) RemoveTeamMember(c echo.Context) error {
	teamID, err := strconv.ParseUint(c.Param("teamId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
	}

	employeeID, err := strconv.ParseUint(c.Param("employeeId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid employee ID"})
	}

	if err := h.service.RemoveTeamMember(uint(teamID), uint(employeeID)); err != nil {
		h.logger.Error("Error removing team member", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove team member"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Team member removed successfully"})
}

func (h *handler) GetTeamMembers(c echo.Context) error {
	teamID, err := strconv.ParseUint(c.Param("teamId"), 10, 32)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid team ID"})
	}

	members, err := h.service.GetTeamMembers(uint(teamID))
	if err != nil {
		h.logger.Error("Error getting team members", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get team members"})
	}

	return c.JSON(http.StatusOK, members)
}
