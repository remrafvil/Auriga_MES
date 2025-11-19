package hLabor

import (
	"github.com/labstack/echo/v4"
	"github.com/remrafvil/Auriga_API/config"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers"
	"github.com/remrafvil/Auriga_API/internal/httpapi/middlewares"
	"github.com/remrafvil/Auriga_API/internal/services/sLabor"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type handler struct {
	service        sLabor.Service
	authMiddleware *middlewares.AuthMiddleware
	logger         *zap.Logger
}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

type Params struct {
	fx.In

	Service        sLabor.Service
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
	r := e.Group("/labor")
	/*middlewares*/

	//r.Use(middlewares.RequestLoggerMiddleware(h.logger))
	//r.Use(h.authMiddleware.CombinedMiddleware())

	//r.Use(h.authMiddleware.Handler())

	// Shift routes
	r.GET("/shifts", h.GetAllShifts)
	r.GET("/shifts/:id", h.GetShift)
	r.POST("/shifts", h.CreateShift)
	r.PUT("/shifts/:id", h.UpdateShift)
	r.DELETE("/shifts/:id", h.DeleteShift)

	// Team routes
	r.GET("/teams", h.GetAllTeams)
	r.GET("/teams/:id", h.GetTeam)
	r.POST("/teams", h.CreateTeam)
	r.PUT("/teams/:id", h.UpdateTeam)
	r.DELETE("/teams/:id", h.DeleteTeam)

	// Team member routes
	r.POST("/teams/:teamId/members", h.AddTeamMember)
	r.DELETE("/teams/:teamId/members/:employeeId", h.RemoveTeamMember)
	r.GET("/teams/:teamId/members", h.GetTeamMembers)

	// Shift assignment routes
	r.GET("/shift-assig/employee/:employeeId", h.GetEmployeeAssignments)
	r.GET("/shift-assig/team/:teamId", h.GetTeamAssignments)
	r.GET("/shift-assig/employee/:employeeId/current", h.GetCurrentEmployeeAssignment)
	r.POST("/shift-assig/individual", h.CreateIndividualAssignment)
	r.POST("/shift-assig/team", h.CreateTeamAssignment)
	r.POST("/shift-assig/bulk", h.CreateBulkAssignments)
}
