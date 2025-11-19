package hInfluxQuery

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type mhInfluxDBSQuerySelector struct {
	Factory   string `json:"factory" 			form:"factory" 				query:"factory"`
	ProdLine  string `json:"prodline" 			form:"prodline" 			query:"prodline"`
	System    string `json:"system" 				form:"system" 				query:"system"`
	Machine   string `json:"Machine" 			form:"Machine" 				query:"Machine"`
	Part      string `json:"part" 				form:"part" 				query:"part"`
	Subpart1  string `json:"subpart1" 			form:"subpart1" 			query:"subpart1"`
	Subpart2  string `json:"subpart2" 			form:"subpart2" 			query:"subpart2"`
	Subpart3  string `json:"subpart3" 			form:"subpart3" 			query:"subpart3"`
	QueryType string `json:"querytype" 			form:"querytype" 			query:"querytype"`
}

func (h *handler) GetQueryByName(c echo.Context) error {
	ctx := c.Request().Context()
	u := new(mhInfluxDBSQuerySelector)

	// Obtener valores de las cabeceras
	u.Factory = strings.TrimSpace(c.Request().Header.Get("Factory"))
	u.ProdLine = strings.TrimSpace(c.Request().Header.Get("ProdLine"))
	u.System = strings.TrimSpace(c.Request().Header.Get("System"))
	u.Machine = strings.TrimSpace(c.Request().Header.Get("Machine"))
	u.Part = strings.TrimSpace(c.Request().Header.Get("Part"))
	u.Subpart1 = strings.TrimSpace(c.Request().Header.Get("Subpart1"))
	u.Subpart2 = strings.TrimSpace(c.Request().Header.Get("Subpart2"))
	u.Subpart3 = strings.TrimSpace(c.Request().Header.Get("Subpart3"))
	u.QueryType = strings.TrimSpace(c.Request().Header.Get("QueryType"))

	// Construir el nombre compuesto
	name := buildCompositeName(u)

	if name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "At least one header parameter is required to build the query name",
		})
	}

	query, err := h.service.GetQueryByName(ctx, name)
	if err != nil {
		h.logger.Error("Failed to get query by name",
			zap.Error(err),
			zap.String("name", name),
		)
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Query not found for name: " + name,
		})
	}

	// Si solo quieren el query directamente (para copiar/pegar)
	if strings.TrimSpace(c.Request().Header.Get("raw")) == "true" {
		return c.String(http.StatusOK, query.Query)
	}

	// Para uso programático - devolver el query como string ejecutable
	response := map[string]interface{}{
		"id":          query.ID,
		"name":        query.Name,
		"description": query.Description,
		"query":       query.Query, // Esto debería ser el string con saltos de línea reales
	}

	return c.JSON(http.StatusOK, response)
}

// buildCompositeName construye el nombre compuesto basado en los campos no vacíos
func buildCompositeName(u *mhInfluxDBSQuerySelector) string {
	var parts []string

	// Agregar campos no vacíos en el orden específico
	if u.Factory != "" {
		parts = append(parts, u.Factory)
	}
	if u.ProdLine != "" {
		parts = append(parts, u.ProdLine)
	}
	if u.System != "" {
		parts = append(parts, u.System)
	}
	if u.Machine != "" {
		parts = append(parts, u.Machine)
	}
	if u.Part != "" {
		parts = append(parts, u.Part)
	}
	if u.Subpart1 != "" {
		parts = append(parts, u.Subpart1)
	}
	if u.Subpart2 != "" {
		parts = append(parts, u.Subpart2)
	}
	if u.Subpart3 != "" {
		parts = append(parts, u.Subpart3)
	}
	if u.QueryType != "" {
		parts = append(parts, u.QueryType)
	}

	// Unir todas las partes con underscore
	return strings.Join(parts, "_")
}
