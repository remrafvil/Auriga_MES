package hEvents

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type responseMessage struct {
	Message string `json:"message"`
}

type Handler interface {
	EventsRawByLineList(c echo.Context) error
	EventsRawToCommitLine(c echo.Context) error
	EventsRawByLineDel(c echo.Context) error

	EventsCommitByLineList(c echo.Context) error
	EventsCommitByLineAdd(c echo.Context) error
	EventsCommitByLineUpdate(c echo.Context) error
	EventsCommitByLineDel(c echo.Context) error

	GetAllCategoriesWithEventTypes(c echo.Context) error
	GetCategoryWithEventTypesByName(c echo.Context) error
}

func (h *handler) EventsRawByLineList(c echo.Context) error {
	u := new(mhLineEvents)
	u.Factory = c.Request().Header.Get("Factory")
	u.ProdLine = c.Request().Header.Get("ProdLine")

	log.Println(u.Factory)
	log.Println(u.ProdLine)

	use, err := h.service.EventsRawByLineList(u.Factory, u.ProdLine)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Registro no existeeee"})
	}
	return c.JSON(http.StatusOK, use)
}

func (h *handler) EventsRawToCommitLine(c echo.Context) error {
	u := new(mhLineEvents)
	u.ID = c.Request().Header.Get("ID")
	u.EventTime = c.Request().Header.Get("EventTime")
	u.Factory = c.Request().Header.Get("Factory")
	u.ProdLine = c.Request().Header.Get("ProdLine")
	u.System = c.Request().Header.Get("System")
	u.Machine = c.Request().Header.Get("Machine")
	u.Part = c.Request().Header.Get("Part")
	u.Type = c.Request().Header.Get("Type")

	log.Println("ID", u.ID)
	log.Println("EventTime", u.EventTime)
	log.Println("Factory", u.Factory)
	log.Println("ProdLine", u.ProdLine)
	log.Println("System", u.System)
	log.Println("Machine", u.Machine)
	log.Println("Part", u.Part)
	log.Println("Type", u.Type)

	//u_ID, _ := strconv.Atoi(u.ID)
	uint64_ID, _ := strconv.ParseUint(u.ID, 10, 32)
	uint_ID := uint(uint64_ID)
	eventTime, _ := time.Parse(time.RFC3339, u.EventTime)
	// log.Println(StarteddAtTime)
	// log.Println(FinishedAtTime)

	use, err := h.service.EventsRawToCommitLine(uint_ID, eventTime, u.Factory, u.ProdLine, u.System, u.Machine, u.Part, u.Type)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Registro no Actualizado Handler EventsCommitByLineUpdate"})
	}
	return c.JSON(http.StatusOK, use)
}

func (h *handler) EventsRawByLineDel(c echo.Context) error {
	u := new(mhLineEvents)
	u.ID = c.Request().Header.Get("ID")

	log.Println("ID", u.ID)

	uint64_ID, _ := strconv.ParseUint(u.ID, 10, 32)
	uint_ID := uint(uint64_ID)

	use, err := h.service.EventsRawByLineDel(uint_ID)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Registro no Borrado Handler EventsCommitByLineDel"})
	}
	return c.JSON(http.StatusOK, use)
}
