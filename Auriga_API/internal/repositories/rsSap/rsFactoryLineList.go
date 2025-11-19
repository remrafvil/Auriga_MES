package rsSap

import (
	"encoding/json"
	"encoding/xml"
	"fmt"

	"github.com/go-resty/resty/v2"
)

//API POST para obtener las líneas de producción de un centro de fabricación

type ZPP_IOT_LINES_FACTORY struct {
	XMLName  xml.Name `xml:"ns1:ZPP_IOT_LINES_FACTORY"`
	Xmlns    string   `xml:"xmlns:ns1,attr"`
	IV_WERKS string   `xml:"IV_WERKS"` //nº del centro de fabricación
	IV_SPRAS string   `xml:"IV_SPRAS"` //idioma nombre de la línea(E inglés, S español, F francés)
}
type LineFactory struct {
	ET_LINES struct {
		Item []struct {
			IV_LINEA string `xml:"ARBPL" json:"ARBPL"` //puesto de trabajo SAP(identificador único de linea)
			NameLine string `xml:"KTEXT" json:"KTEXT"` //descripción o nombre de la linea
		} `xml:"item"`
	} `xml:"ET_LINES"`
}

func (m *repository) rsFactoryLineList() {
	url := m.config.Sap.Url
	auth := m.config.Sap.Auth
	// url := "https://l20163-iflmap.hcisbp.eu1.hana.ondemand.com/http/SAP_RFC_IOT_DHM"
	// auth := "czAwMjA2OTUzMjI6MUxhbnRlcltd"

	l := ZPP_IOT_LINES_FACTORY{
		Xmlns:    "urn:sap-com:document:sap:rfc:functions",
		IV_WERKS: "1004",
		IV_SPRAS: "E",
	}
	xmlBytes, err := xml.Marshal(l)
	if err != nil {
		fmt.Println(err)
		return
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
		return
	}
	body := resp.Body()
	data := &LineFactory{}
	err = xml.Unmarshal(body, data)
	if err != nil {
		fmt.Println("Unmarshal Error - ", err)
		return
	}
	jsonBytes, _ := json.Marshal(data.ET_LINES.Item)
	fmt.Println("JSON Data:", string(jsonBytes))
}
