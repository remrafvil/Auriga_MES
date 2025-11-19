package rsSap

import (
	"encoding/xml"
	"fmt"

	"github.com/go-resty/resty/v2"
)

//API POST para obtener los componentes de una orden de fabricación

type ZPP_IOT_ORDER_COMPONENTES struct {
	XMLName  xml.Name `xml:"ns1:ZPP_IOT_ORDER_COMPONENTES"`
	Xmlns    string   `xml:"xmlns:ns1,attr"`
	IV_AUFNR string   `xml:"IV_AUFNR"`
}
type MrsRecipe struct { //ZPP_IOT_ORDER_COMPONENTES
	ET_COMPONENTES struct {
		Item []struct {
			SapCode           string `xml:"COMPONENTE"`           //componente
			Description       string `xml:"MATERIAL_DESCRIPTION"` //descripción del componente
			RequiredQuantity  string `xml:"REQ_QUAN"`             //cantidad necesaria para fabricar
			MeasurementUnitRQ string `xml:"BASE_UOM"`             //unidad de medida de cantidad necesaria
			CommittedQuantity string `xml:"COMMITED_QUANTITY"`    //cantidad confirmada
			MeasurementUnitCQ string `xml:"ENTRY_UOM"`            //unidad de medida de cantidad confirmada
			WithDrawnQuantity string `xml:"WITHDRAWN_QUANTITY"`   //cantidad tomada
		} `xml:"item"`
	} `xml:"ET_COMPONENTES"`
}

type MrsComponent struct {
	SapCode           string
	Description       string
	RequiredQuantity  string
	MeasurementUnitRQ string
	CommittedQuantity string
	MeasurementUnitCQ string
	WithDrawnQuantity string
}

func (m *repository) RsLineRecipe(lineNumber string) ([]MrsComponent, error) {

	var data = make([]MrsComponent, 0)

	url := m.config.Sap.Url
	auth := m.config.Sap.Auth
	// url := "https://l20163-iflmap.hcisbp.eu1.hana.ondemand.com/http/SAP_RFC_IOT_DHM"
	// auth := "czAwMjA2OTUzMjI6MUxhbnRlcltd"

	c := ZPP_IOT_ORDER_COMPONENTES{
		Xmlns:    "urn:sap-com:document:sap:rfc:functions",
		IV_AUFNR: lineNumber,
	}
	xmlBytes, err := xml.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return data, err
	}
	// Add the XML declaration
	xmlDecl := []byte(`<?xml version="1.0" encoding="UTF-8"?>`)
	xmlBytes = append(xmlDecl, xmlBytes...)
	xmlStr := string(xmlBytes)
	fmt.Println("xmlBytes - ", xmlStr)

	client := resty.New()
	resp, err := client.R().
		SetHeader("Content-Type", "application/xml").
		SetHeader("Authorization", "Basic "+auth).
		SetBody(xmlStr).
		Post(url)
	if err != nil {
		fmt.Println("api failure", err)
		return data, err
	}
	body := resp.Body()
	dataS := &MrsRecipe{}
	err = xml.Unmarshal(body, dataS)
	if err != nil {
		fmt.Println("Unmarshal Error - ", err)
		return data, err
	}

	for _, p := range dataS.ET_COMPONENTES.Item {
		data = append(data, MrsComponent{
			SapCode:           p.SapCode,
			Description:       p.Description,
			RequiredQuantity:  p.RequiredQuantity,
			MeasurementUnitRQ: p.MeasurementUnitRQ,
			CommittedQuantity: p.CommittedQuantity,
			MeasurementUnitCQ: p.MeasurementUnitCQ,
			WithDrawnQuantity: p.WithDrawnQuantity,
		})
	}
	return data, nil
}
