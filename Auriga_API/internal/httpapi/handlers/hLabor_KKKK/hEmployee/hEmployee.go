package hEmployee

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/remrafvil/Auriga_API/config"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers"
	"github.com/remrafvil/Auriga_API/internal/httpapi/middlewares"
	"github.com/remrafvil/Auriga_API/internal/services/sLabor1"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type handler struct {
	service        sLabor1.Service
	authMiddleware *middlewares.AuthMiddleware
	logger         *zap.Logger
}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

type Params struct {
	fx.In

	Service        sLabor1.Service
	AuthMiddleware *middlewares.AuthMiddleware
	Logger         *zap.Logger
}

func New(p Params) Result {
	return Result{
		Handler: &handler{
			service:        p.Service,
			authMiddleware: p.AuthMiddleware,
			logger:         p.Logger,
		},
	}
}

// ejemplo para usar en main.go	routes.Asset(e)
func (h *handler) RegisterRoutes(e *echo.Echo, s *config.Settings) {
	r := e.Group("/employees")

	// Si necesitas middleware, descomenta estas líneas:
	// r.Use(h.authMiddleware.CombinedMiddleware())
	// r.Use(h.authMiddleware.Handler())

	// Rutas dentro del grupo /employees
	r.POST("/sync", h.SyncEmployees)
	r.GET("/count", h.GetEmployeeCount)
	r.GET("/:code", h.GetEmployeeByCode)
}

func (h *handler) SyncEmployees(c echo.Context) error {
	h.logger.Info("Received request to sync employees")

	err := h.service.SyncEmployees()
	if err != nil {
		h.logger.Error("Failed to sync employees", zap.Error(err))
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to sync employees: " + err.Error(),
		})
	}

	count, err := h.service.GetEmployeeCount()
	if err != nil {
		h.logger.Error("Failed to get employee count", zap.Error(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":         "Employees synchronized successfully",
		"total_employees": count,
	})
}

func (h *handler) GetEmployeeCount(c echo.Context) error {
	count, err := h.service.GetEmployeeCount()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to get employee count",
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"total_employees": count,
	})
}

func (h *handler) GetEmployeeByCode(c echo.Context) error {
	code := c.Param("code")
	if code == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Employee code is required",
		})
	}

	employee, err := h.service.GetEmployeeByCode(code)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Employee not found",
		})
	}

	return c.JSON(http.StatusOK, employee)
}
