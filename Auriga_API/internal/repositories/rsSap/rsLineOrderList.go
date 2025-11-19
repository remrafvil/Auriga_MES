package rsSap

import (
	"log"

	"encoding/xml"

	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
)

type ZPP_IOT_ORDER_SEQUENCE struct {
	XMLName  xml.Name `xml:"ns1:ZPP_IOT_ORDER_SEQUENCE"`
	Xmlns    string   `xml:"xmlns:ns1,attr"`
	IV_LINEA string   `xml:"IV_LINEA"`
}

type OrderSequence struct {
	ET_ORDERS struct {
		Item []struct {
			OrderNumber               string `xml:"AUFNR" json:"AUFNR"`
			OrderNType                string `xml:"AUART" json:"AUART"`
			ProductName               string `xml:"MATNR" json:"MATNR"`
			ProductDescription        string `xml:"MAKTX" json:"MAKTX"`
			StartDateOF               string `xml:"GSTRI" json:"GSTRI"`
			StartTimeOF               string `xml:"GSUZI" json:"GSUZI"`
			EndDateOF                 string `xml:"GETRI" json:"GETRI"`
			EndTimeOF                 string `xml:"GEUZI" json:"GEUZI"`
			QuantityToProduce         string `xml:"GAMNG" json:"GAMNG"`
			QuantityProduced          string `xml:"WEMNG" json:"WEMNG"`
			QuantityRemainedToProduce string `xml:"CANT_RES" json:"CANT_RES"`
			MeasurementUnit           string `xml:"GMEIN" json:"GMEIN"`
		} `xml:"item"`
	} `xml:"ET_ORDERS"`
}

func (m *repository) RsLineOrderList(factory string, lineNumber string, lineSapCode string) ([]rModels.MrProductionOrder, error) {

	var data = make([]rModels.MrProductionOrder, 0)

	url := m.config.Sap.Url
	auth := m.config.Sap.Auth
	// url := "https://l20163-iflmap.hcisbp.eu1.hana.ondemand.com/http/SAP_RFC_IOT_DHM"
	// auth := "czAwMjA2OTUzMjI6MUxhbnRlcltd"

	z := ZPP_IOT_ORDER_SEQUENCE{
		Xmlns:    "urn:sap-com:document:sap:rfc:functions",
		IV_LINEA: lineSapCode,
	}
	xmlBytes, err := xml.Marshal(z)
	if err != nil {
		fmt.Println(err)
		return data, err
	}
	// Add the XML declaration
	xmlDecl := []byte(`<?xml version="1.0" encoding="UTF-8"?>`)
	xmlBytes = append(xmlDecl, xmlBytes...)
	xmlStr := string(xmlBytes)

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
	dataS := &OrderSequence{}
	err = xml.Unmarshal(body, dataS)
	if err != nil {
		fmt.Println("Unmarshal Error - ", err)
		return data, err
	}

	for _, p := range dataS.ET_ORDERS.Item {
		//		StarteddAt, _ := time.Parse("2006-01-02 15:04:05", p.StartDateOF+" "+p.StartTimeOF)
		//		FinishedAt, _ := time.Parse("2006-01-02 15:04:05", p.EndDateOF+" "+p.EndTimeOF)
		data = append(data, rModels.MrProductionOrder{
			Factory:                   factory,
			ProdLine:                  lineNumber,
			OrderNumber:               p.OrderNumber,
			OrderNType:                p.OrderNType,
			ProductName:               p.ProductName,
			ProductDescription:        p.ProductDescription,
			QuantityToProduce:         p.QuantityToProduce,
			QuantityProduced:          p.QuantityProduced,
			QuantityRemainedToProduce: p.QuantityRemainedToProduce,
			MeasurementUnit:           p.MeasurementUnit,
			//			StarteddAt:                StarteddAt,
			//			FinishedAt:                FinishedAt,
		})
	}
	log.Println("*****************      Llego por aquí REPOSITORY SAP  *****************")
	return data, nil
}
