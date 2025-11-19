package hSap

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type mhOrderRecipe struct {
	SapOrderCode string `json:"sapordercode" 	form:"sapordercode" query:"sapordercode"`
	SapRequest   string `json:"saprequest" 		form:"saprequest" 	query:"saprequest"`
}

func (h *handler) OrderRecipe(c echo.Context) error {
	u := new(mhOrderRecipe)
	u.SapOrderCode = c.Request().Header.Get("SapOrderCode")
	u.SapRequest = c.Request().Header.Get("SapRequest")

	log.Println("Receta para el siguiente SapOrderCode:", u.SapOrderCode)
	log.Println(u.SapRequest)

	use, err := h.service.RecipeComponentList(u.SapOrderCode, u.SapRequest)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Registro no existeeee"})
	}
	return c.JSON(http.StatusOK, use)
}
