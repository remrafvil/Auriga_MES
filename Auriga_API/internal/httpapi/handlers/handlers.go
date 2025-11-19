package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/remrafvil/Auriga_API/config"
)

type Handler interface {
	RegisterRoutes(e *echo.Echo, s *config.Settings)
}
