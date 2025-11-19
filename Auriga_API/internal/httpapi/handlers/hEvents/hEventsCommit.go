package hEvents

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

type mhLineEvents struct {
	EventTime string `json:"eventtime" 	form:"eventtime" 	query:"eventtime"`
	Factory   string `json:"factory" 	form:"factory" 		query:"factory"`
	ProdLine  string `json:"prodline" 	form:"prodline" 	query:"prodline"`
	System    string `json:"system" 	form:"system" 		query:"system"`
	Machine   string `json:"machine" 	form:"machine" 	query:"machine"`
	Part      string `json:"part" 		form:"part" 		query:"part"`
	Type      string `json:"type" 		form:"type" 		query:"type"`
	Category  string `json:"category" 	form:"category" 	query:"category"`
	ID        string `json:"id" 		form:"id" 		query:"id"`
}

func (h *handler) EventsCommitByLineList(c echo.Context) error {
	u := new(mhLineEvents)
	u.Factory = c.Request().Header.Get("Factory")
	u.ProdLine = c.Request().Header.Get("ProdLine")

	log.Println(u.Factory)
	log.Println(u.ProdLine)

	use, err := h.service.EventsCommitByLineList(u.Factory, u.ProdLine)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Listado no Mostrado Handler  EventsCommitByLineList"})
	}
	return c.JSON(http.StatusOK, use)
}

func (h *handler) EventsCommitByLineAdd(c echo.Context) error {
	u := new(mhLineEvents)
	u.EventTime = c.Request().Header.Get("EventTime")
	u.Factory = c.Request().Header.Get("Factory")
	u.ProdLine = c.Request().Header.Get("ProdLine")
	u.System = c.Request().Header.Get("System")
	u.Machine = c.Request().Header.Get("Machine")
	u.Part = c.Request().Header.Get("Part")
	u.Type = c.Request().Header.Get("Type")
	u.Category = c.Request().Header.Get("Category")

	log.Println("EventTime", u.EventTime)
	log.Println("Factory", u.Factory)
	log.Println("ProdLine", u.ProdLine)
	log.Println("System", u.System)
	log.Println("Machine", u.Machine)
	log.Println("Part", u.Part)
	log.Println("Type", u.Type)
	log.Println("Category", u.Category)

	eventTime, _ := time.Parse(time.RFC3339, u.EventTime)

	use, err := h.service.EventsCommitByLineAdd(eventTime, u.Factory, u.ProdLine, u.System, u.Machine, u.Part, u.Type, u.Category)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Registro no Añadido Handler  EventsCommitByLineAdd"})
	}

	//log.Println("Estos son los datos:", use)
	return c.JSON(http.StatusOK, use)
}

func (h *handler) EventsCommitByLineUpdate(c echo.Context) error {
	u := new(mhLineEvents)
	u.ID = c.Request().Header.Get("ID")
	u.EventTime = c.Request().Header.Get("EventTime")
	u.Factory = c.Request().Header.Get("Factory")
	u.ProdLine = c.Request().Header.Get("ProdLine")
	u.System = c.Request().Header.Get("System")
	u.Machine = c.Request().Header.Get("Machine")
	u.Part = c.Request().Header.Get("Part")
	u.Type = c.Request().Header.Get("Type")
	u.Category = c.Request().Header.Get("Category")

	log.Println("ID", u.ID)
	log.Println("EventTime", u.EventTime)
	log.Println("Factory", u.Factory)
	log.Println("ProdLine", u.ProdLine)
	log.Println("System", u.System)
	log.Println("Machine", u.Machine)
	log.Println("Part", u.Part)
	log.Println("Type", u.Type)
	log.Println("Category", u.Category)

	//u_ID, _ := strconv.Atoi(u.ID)
	uint64_ID, _ := strconv.ParseUint(u.ID, 10, 32)
	uint_ID := uint(uint64_ID)
	eventTime, _ := time.Parse(time.RFC3339, u.EventTime)
	// log.Println(StarteddAtTime)
	// log.Println(FinishedAtTime)

	use, err := h.service.EventsCommitByLineUpdate(uint_ID, eventTime, u.Factory, u.ProdLine, u.System, u.Machine, u.Part, u.Type, u.Category)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Registro no Actualizado Handler EventsCommitByLineUpdate"})
	}
	return c.JSON(http.StatusOK, use)
}

func (h *handler) EventsCommitByLineDel(c echo.Context) error {
	u := new(mhLineEvents)
	u.ID = c.Request().Header.Get("ID")

	log.Println("ID", u.ID)

	uint64_ID, _ := strconv.ParseUint(u.ID, 10, 32)
	uint_ID := uint(uint64_ID)

	use, err := h.service.EventsCommitByLineDel(uint_ID)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Registro no Borrado Handler EventsCommitByLineDel"})
	}
	return c.JSON(http.StatusOK, use)
}
