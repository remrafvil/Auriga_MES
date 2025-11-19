package httpapi

import (
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers/hAssets"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers/hAuth"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers/hEvents"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers/hInfluxQuery"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers/hLabor"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers/hLabor_KKKK/hEmployee"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers/hProducts"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers/hSap"
	"github.com/remrafvil/Auriga_API/internal/httpapi/middlewares"
	"go.uber.org/fx"
)

var Module = fx.Module("httpapi", fx.Provide(
	middlewares.NewAuthMiddleware,
	hAssets.New,
	hAuth.New,
	hProducts.New,
	hSap.New,
	hEvents.New,
	hInfluxQuery.New,
	hEmployee.New,
	hLabor.New,
))
