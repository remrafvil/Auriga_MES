package hEvents

import (
	"github.com/labstack/echo/v4"
	"github.com/remrafvil/Auriga_API/config"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers"
	"github.com/remrafvil/Auriga_API/internal/httpapi/middlewares"
	"github.com/remrafvil/Auriga_API/internal/services/sEvents"
	"go.uber.org/fx"
)

type handler struct {
	service        sEvents.Service
	authMiddleware *middlewares.AuthMiddleware
}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

type Params struct {
	fx.In

	Service        sEvents.Service
	AuthMiddleware *middlewares.AuthMiddleware
}

func New(p Params) Result {
	return Result{
		Handler: &handler{
			service:        p.Service,
			authMiddleware: p.AuthMiddleware,
		},
	}
}

// ejemplo para usar en main.go	routes.Land(e)
func (h *handler) RegisterRoutes(e *echo.Echo, s *config.Settings) {
	//jwtconfig := middlewares.JwtInitConfig()
	//middlewares.MainMiddlewares(e, s)

	r := e.Group("/events")
	/*middlewares*/
	//r.Use(middleware.JWTWithConfig(jwtconfig))

	r.GET("/raw", h.EventsRawByLineList)
	r.GET("/raw/tocommit", h.EventsRawToCommitLine)
	r.GET("/raw/del", h.EventsRawByLineDel)

	r.GET("/commit", h.EventsCommitByLineList)
	r.GET("/commit/add", h.EventsCommitByLineAdd)
	r.GET("/commit/update", h.EventsCommitByLineUpdate)
	r.GET("/commit/del", h.EventsCommitByLineDel)

	r.GET("/sapCommit", h.EventsSapCommit)

	r.GET("/categories", h.GetAllCategoriesWithEventTypes)

	r.GET("/typeByCat", h.GetCategoryWithEventTypesByName)
}
