package hShifts

import (
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers"
	"github.com/remrafvil/Auriga_API/internal/httpapi/middlewares"
	"github.com/remrafvil/Auriga_API/internal/services/sAssets"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type handler struct {
	service        sAssets.Service
	authMiddleware *middlewares.AuthMiddleware
	logger         *zap.Logger
}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

type Params struct {
	fx.In

	Service        sAssets.Service
	AuthMiddleware *middlewares.AuthMiddleware
	Logger         *zap.Logger
}

/* func New(p Params) Result {
	return Result{
		Handler: &handler{
			service:        p.Service,
			authMiddleware: p.AuthMiddleware,
			logger:         p.Logger,
		},
	}
} */
