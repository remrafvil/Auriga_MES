package hProducts

import (
	"github.com/labstack/echo/v4"
	"github.com/remrafvil/Auriga_API/config"
	"github.com/remrafvil/Auriga_API/internal/httpapi/handlers"
	"github.com/remrafvil/Auriga_API/internal/httpapi/middlewares"
	"github.com/remrafvil/Auriga_API/internal/services/sProducts"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type handler struct {
	service        sProducts.Service
	authMiddleware *middlewares.AuthMiddleware
	logger         *zap.Logger
}

type Result struct {
	fx.Out

	Handler handlers.Handler `group:"handlers"`
}

type Params struct {
	fx.In

	Service        sProducts.Service
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

// ejemplo para usar en main.go	routes.Product(e)
func (h *handler) RegisterRoutes(e *echo.Echo, s *config.Settings) {
	//jwtconfig := middlewares.JwtInitConfig()

	r := e.Group("/product")
	/*middlewares*/
	//r.Use(middleware.JWTWithConfig(jwtconfig))

	// r.POST("/add", h.ProductAdd)
	// r.PUT("/update/:id", h.ProductUpdate)
	// r.DELETE("/del/:id", h.ProductDel)
	// r.GET("/show/:id", h.ProductShow)
	// r.GET("/list", h.ProductList)

	// Listado de productos
	r.GET("/list", h.GetProductsList)

	// Producto específico por ID
	r.GET("/show/:id", h.GetProductByID)

	// Listado básico
	r.GET("/typeslist/sort", h.GetProductTypesList)
	// Product type completo por ID
	r.GET("/typedetail/:id", h.GetProductTypeByID)
	// Todos los product types con características
	r.GET("/typeslist/complete", h.GetAllProductTypesWithFeatures)
	// Crear producto
	r.POST("/add", h.CreateProduct)
	// Actualizar producto
	r.PUT("/update/:id", h.UpdateProduct)
	// Eliminar producto
	r.DELETE("/delete/:id", h.DeleteProduct)
}
