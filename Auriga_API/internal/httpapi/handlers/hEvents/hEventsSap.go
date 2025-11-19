package hEvents

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h *handler) EventsSapCommit(c echo.Context) error {
	u := new(mhLineEvents)
	u.ID = c.Request().Header.Get("ID")

	log.Println("ID", u.ID)
	log.Println("Vamos a enviar este evento a SAP")
	fmt.Println("Vamos a enviar este evento a SAP")

	uint64_ID, _ := strconv.ParseUint(u.ID, 10, 32)
	uint_ID := uint(uint64_ID)

	use, err := h.service.EventsSapByLineDel(uint_ID)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Registro no Borrado Handler EventsCommitByLineDel"})
	}
	return c.JSON(http.StatusOK, use)
}
