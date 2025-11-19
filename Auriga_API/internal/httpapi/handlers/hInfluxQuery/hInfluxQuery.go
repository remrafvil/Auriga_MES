package hInfluxQuery

import (
	"github.com/labstack/echo/v4"
	"github.com/remrafvil/Auriga_API/config"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers"
	"github.com/remrafvil/Auriga_API/internal/httpapi/middlewares"
	"github.com/remrafvil/Auriga_API/internal/services/sInfluxQuery"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type handler struct {
	service        sInfluxQuery.Service
	authMiddleware *middlewares.AuthMiddleware
	logger         *zap.Logger
}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

type Params struct {
	fx.In

	Service        sInfluxQuery.Service
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
	//middlewares.MainMiddlewares(e, s)
	//jwtconfig := middlewares.JwtInitConfig()

	r := e.Group("/influxquery")
	/*middlewares*/
	//r.Use(middleware.JWTWithConfig(jwtconfig))

	//r.Use(h.authMiddleware.Handler())

	r.GET("/energy", h.GetQueryByName)
}
