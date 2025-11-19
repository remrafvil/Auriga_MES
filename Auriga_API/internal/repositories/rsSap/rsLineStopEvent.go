package rsSap

import (
	"encoding/xml"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

//API GET para enviar los datos de una parada de una máquina SAP

type ZPP_IOT_PARADA_LOG struct {
	XMLName   xml.Name `xml:"ns1:ZPP_IOT_PARADA_LOG"`
	Xmlns     string   `xml:"xmlns:ns1,attr"`
	Maquina   string   `xml:"IV_MAQUINA"`
	Timestamp string   `xml:"IV_TIMESTAMP"`
	Estado    string   `xml:"IV_ESTADO"`
	Motivo    string   `xml:"IV_MOTIVO"`
	OF        string   `xml:"IV_OF"`
	Operario  string   `xml:"IV_OPERARIO"`
}

func (m *repository) RsLineStopEvent(maquina string, timeEvent time.Time, estado string, motivo string, of string, operario string) error {
	url := m.config.Sap.Url
	auth := m.config.Sap.Auth
	// url := "https://l20163-iflmap.hcisbp.eu1.hana.ondemand.com/http/SAP_RFC_IOT_DHM"
	// auth := "czAwMjA2OTUzMjI6MUxhbnRlcltd"

	p := ZPP_IOT_PARADA_LOG{
		Xmlns:     "urn:sap-com:document:sap:rfc:functions",
		Maquina:   maquina,
		Timestamp: timeEvent.Format("2006-01-02T15:04:05"),
		Estado:    estado,
		Motivo:    motivo,
		OF:        of,
		Operario:  operario,
	}
	xmlBytes, err := xml.Marshal(p)
	if err != nil {
		fmt.Println(err)
		return err
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
		return err
	}

	fmt.Println("final body-", resp.String())
	fmt.Println("status code-", resp.StatusCode())

	return err
}
