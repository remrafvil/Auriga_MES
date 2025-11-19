package hSap

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

type mhOrderConsumption struct {
	Factory          string `json:"factory" 			form:"factory" 				query:"factory"`
	ProdLine         string `json:"prodline" 			form:"prodline" 			query:"prodline"`
	System           string `json:"system" 				form:"system" 				query:"system"`
	Machine          string `json:"Machine" 			form:"Machine" 				query:"Machine"`
	Part             string `json:"part" 				form:"part" 				query:"part"`
	SapOrderCode     string `json:"sapordercode" 		form:"sapordercode"			query:"sapordercode"`
	SapRequest       string `json:"saprequest" 			form:"saprequest" 			query:"saprequest"`
	SapComponentCode string `json:"sapcomponentcode" 	form:"sapcomponentcode" 	query:"sapcomponentcode"`
	SapCode          string `json:"sapcode" 			form:"sapcode" 			query:"sapcode"`
}

func (h *handler) OrderConsumption(c echo.Context) error {
	u := new(mhOrderConsumption)
	u.Factory = c.Request().Header.Get("Factory")
	u.ProdLine = c.Request().Header.Get("ProdLine")
	u.System = c.Request().Header.Get("System")
	u.SapOrderCode = c.Request().Header.Get("SapOrderCode")
	u.SapRequest = c.Request().Header.Get("SapRequest")

	log.Println("Factory", u.Factory)
	log.Println("ProdLine", u.ProdLine)
	log.Println("System", u.System)
	log.Println("SapOrderCode", u.SapOrderCode)
	log.Println("SapRequest", u.SapRequest)

	use, err := h.service.DosingConsumptionList(u.Factory, u.ProdLine, u.System, u.SapOrderCode, u.SapRequest)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusForbidden, responseMessage{Message: "Registro no existeeee"})
	}

	//log.Println("Estos son los datos:", use)
	return c.JSON(http.StatusOK, use)
}

func (h *handler) OrderConsumptionAdd(c echo.Context) error {
	u := new(mhOrderConsumption)
	u.Factory = c.Request().Header.Get("Factory")
	u.ProdLine = c.Request().Header.Get("ProdLine")
	u.System = c.Request().Header.Get("System")
	u.Machine = c.Request().Header.Get("Machine")
	u.Part = c.Request().Header.Get("Part")
	u.SapOrderCode = c.Request().Header.Get("SapOrderCode")
	u.SapComponentCode = c.Request().Header.Get("SapComponentCode")

	log.Println("Factory", u.Factory)
	log.Println("ProdLine", u.ProdLine)
	log.Println("System", u.System)
	log.Println("Machine", u.Machine)
	log.Println("Part", u.Part)
	log.Println("SapOrderCode", u.SapOrderCode)
	log.Println("SapComponentCode", u.SapComponentCode)

	use, err := h.service.DosingConsumptionAdd(u.Factory, u.ProdLine, u.System, u.Machine, u.Part, u.SapOrderCode, u.SapComponentCode)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Registro no insertado"})
	}

	//log.Println("Estos son los datos:", use)
	return c.JSON(http.StatusOK, use)
}

func (h *handler) OrderConsumptionDel(c echo.Context) error {
	u := new(mhOrderConsumption)
	u.Factory = c.Request().Header.Get("Factory")
	u.ProdLine = c.Request().Header.Get("ProdLine")
	u.System = c.Request().Header.Get("System")
	u.Machine = c.Request().Header.Get("Machine")
	u.Part = c.Request().Header.Get("Part")
	u.SapOrderCode = c.Request().Header.Get("SapOrderCode")
	u.SapComponentCode = c.Request().Header.Get("SapComponentCode")

	log.Println("Factory", u.Factory)
	log.Println("ProdLine", u.ProdLine)
	log.Println("System", u.System)
	log.Println("Machine", u.Machine)
	log.Println("Part", u.Part)
	log.Println("SapOrderCode", u.SapOrderCode)
	log.Println("SapComponentCode", u.SapComponentCode)

	use, err := h.service.DosingConsumptionDel(u.Factory, u.ProdLine, u.System, u.Machine, u.Part, u.SapOrderCode, u.SapComponentCode)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Registro no borrado"})
	}

	//log.Println("Estos son los datos:", use)
	return c.JSON(http.StatusOK, use)
}

func (h *handler) OrderConsumptionUpdate(c echo.Context) error {
	u := new(mhOrderConsumption)
	u.Factory = c.Request().Header.Get("Factory")
	u.ProdLine = c.Request().Header.Get("ProdLine")
	u.System = c.Request().Header.Get("System")
	u.Machine = c.Request().Header.Get("Machine")
	u.Part = c.Request().Header.Get("Part")
	u.SapOrderCode = c.Request().Header.Get("SapOrderCode")
	u.SapComponentCode = c.Request().Header.Get("SapComponentCode")

	log.Println("Factory", u.Factory)
	log.Println("ProdLine", u.ProdLine)
	log.Println("System", u.System)
	log.Println("Machine", u.Machine)
	log.Println("Part", u.Part)
	log.Println("SapOrderCode", u.SapOrderCode)
	log.Println("SapComponentCode", u.SapComponentCode)

	use, err := h.service.DosingConsumptionUpdate(u.Factory, u.ProdLine, u.System, u.Machine, u.Part, u.SapOrderCode, u.SapComponentCode)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Registro no actualizado"})
	}

	//log.Println("Estos son los datos:", use)
	return c.JSON(http.StatusOK, use)
}

func (h *handler) OrderConsumptionCalculate(c echo.Context) error {
	u := new(mhOrderConsumption)
	u.Factory = c.Request().Header.Get("Factory")
	u.ProdLine = c.Request().Header.Get("ProdLine")
	u.SapCode = c.Request().Header.Get("SapCode")
	u.SapOrderCode = c.Request().Header.Get("SapOrderCode")

	log.Println("Factory", u.Factory)
	log.Println("ProdLine", u.ProdLine)
	log.Println("SapCode", u.SapCode)
	log.Println("SapOrderCode", u.SapOrderCode)

	use, err := h.service.DosingConsumptionCalculate(u.Factory, u.ProdLine, u.SapOrderCode)
	if err != nil {
		fmt.Println("registo no: %w", err)
		return c.JSON(http.StatusInternalServerError, responseMessage{Message: "Registro no actualizado"})
	}

	//log.Println("Estos son los datos:", use)
	return c.JSON(http.StatusOK, use)
}

func (h *handler) OrderConsumptionSummaryToSAP(c echo.Context) error {
	u := new(mhOrderConsumption)
	u.Factory = c.Request().Header.Get("Factory")
	u.ProdLine = c.Request().Header.Get("ProdLine")
	u.SapCode = c.Request().Header.Get("SapCode")
	u.SapOrderCode = c.Request().Header.Get("SapOrderCode")

	log.Println("Factory", u.Factory)
	log.Println("ProdLine", u.ProdLine)
	log.Println("SapCode", u.SapCode)
	log.Println("SapOrderCode", u.SapOrderCode)

	//log.Println("Estos son los datos:", use)
	return c.JSON(http.StatusOK, u)
}
