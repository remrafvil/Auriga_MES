package hSap

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type responseMessage struct {
	Message string `json:"message"`
}

type Handler interface {
	LineOrder(c echo.Context) error
	OrderRecipe(c echo.Context) error
}

type mhLineOrders struct {
	Factory     string `json:"factory" 	form:"factory" 		query:"factory"`
	Line        string `json:"line" 		form:"line" 		query:"line"`
	SapCode     string `json:"sapcode" 	form:"sapcode" 		query:"sapcode"`
	SapRequest  string `json:"saprequest"form:"saprequest" 	query:"saprequest"`
	OrderNumber string `json:"ordernumber"form:"ordernumber" query:"ordernumber"`
	StartFinish string `json:"startfinish"form:"startfinish" 	query:"startfinish"`
	StarteddAt  string `json:"startedat" form:"startedat" 	query:"startedat"`
	FinishedAt  string `json:"finishedat"form:"finishedat" 	query:"finishedat"`
}

func (h *handler) LineOrders(c echo.Context) error {
	u := new(mhLineOrders)
	u.Factory = c.Request().Header.Get("Factory")
	u.Line = c.Request().Header.Get("ProdLine")
	u.SapCode = c.Request().Header.Get("SapCode")
	u.SapRequest = c.Request().Header.Get("SapRequest")

	log.Println(u.Factory)
	log.Println(u.Line)
	log.Println(u.SapCode)
	log.Println(u.SapRequest)

	use, err := h.service.LineOrderList(u.Factory, u.Line, u.SapCode, u.SapRequest)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Registro no existeeee"})
	}
	return c.JSON(http.StatusOK, use)
}

func (h *handler) LineOrderStartFinish(c echo.Context) error {
	u := new(mhLineOrders)
	u.Factory = c.Request().Header.Get("Factory")
	u.Line = c.Request().Header.Get("ProdLine")
	u.SapCode = c.Request().Header.Get("SapCode")
	u.SapRequest = c.Request().Header.Get("SapRequest")
	u.OrderNumber = c.Request().Header.Get("OrderNumber")
	u.StartFinish = c.Request().Header.Get("StartFinish")

	// log.Println(u.Factory)
	// log.Println(u.Line)
	// log.Println(u.SapCode)
	// log.Println(u.SapRequest)
	// log.Println(u.OrderNumber)
	// log.Println(u.StartFinish)

	use, err := h.service.LineOrderStartFinish(u.Factory, u.Line, u.SapCode, u.SapRequest, u.OrderNumber, u.StartFinish)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Registro no existeeee"})
	}
	return c.JSON(http.StatusOK, use)
}

func (h *handler) LineOrderUpdate(c echo.Context) error {
	u := new(mhLineOrders)
	u.Factory = c.Request().Header.Get("Factory")
	u.Line = c.Request().Header.Get("ProdLine")
	u.SapCode = c.Request().Header.Get("SapCode")
	u.SapRequest = c.Request().Header.Get("SapRequest")
	u.OrderNumber = c.Request().Header.Get("OrderNumber")
	//u.StartFinish = c.Request().Header.Get("StartFinish")
	u.StarteddAt = c.Request().Header.Get("StarteddAt")
	u.FinishedAt = c.Request().Header.Get("FinishedAt")

	// log.Println(u.Factory)
	// log.Println(u.Line)
	// log.Println(u.SapCode)
	// log.Println(u.SapRequest)
	// log.Println(u.OrderNumber)
	// log.Println(u.StartFinish)
	// log.Println(u.StarteddAt)
	// log.Println(u.FinishedAt)

	StarteddAtTime, _ := time.Parse(time.RFC3339, u.StarteddAt)
	FinishedAtTime, _ := time.Parse(time.RFC3339, u.FinishedAt)
	// log.Println(StarteddAtTime)
	// log.Println(FinishedAtTime)

	use, err := h.service.LineOrderUpdateTime(u.Factory, u.Line, u.SapCode, u.SapRequest, u.OrderNumber, StarteddAtTime, FinishedAtTime)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Registro no existeeee"})
	}
	return c.JSON(http.StatusOK, use)
}
