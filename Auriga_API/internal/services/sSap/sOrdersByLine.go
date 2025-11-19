package sSap

import (
	"log"
	"time"
)

type msProductionOrder struct {
	Factory                   string    `json:"Factory"`
	ProdLine                  string    `json:"ProdLine"`
	OrderNumber               string    `json:"OrderNumber"`
	OrderNType                string    `json:"OrderNType"`
	ProductName               string    `json:"ProductName"`
	ProductDescription        string    `json:"ProductDescription"`
	QuantityToProduce         string    `json:"QuantityToProduce"`
	QuantityProduced          string    `json:"QuantityProduced"`
	QuantityRemainedToProduce string    `json:"QuantityRemainedToProduce"`
	MeasurementUnit           string    `json:"MeasurementUnit"`
	StarteddAt                time.Time `json:"StarteddAt"`
	FinishedAt                time.Time `json:"FinishedAt"`
}

func (s *service) LineOrderList(factory string, lineNumber string, lineSapCode string, sapRequest string) ([]msProductionOrder, error) {
	var data = []msProductionOrder{}

	// Obtenemos el activo por la linea de la fabrica
	assetOrder, err := s.repositoryAss.AssetByFactLine(factory, lineNumber)
	if err != nil {
		log.Println("Activo no detectado Service LineOrderList:", err)
		return data, err
	}
	log.Println("Este es el dato assetOrder", assetOrder)

	// Actualizamos o no desde SAP, dependiendo de la solicitud
	if sapRequest == "true" {
		// Actualizamos las ordenes de fabricación desde SAP
		rsData, err := s.repositorySap.RsLineOrderList(factory, lineNumber, assetOrder.SapCode)
		if err != nil {
			log.Println("Error lectura listado de ordenes SAP LineOrderList:", err)
			return data, err
		}

		log.Println("Este es el dato rsData", rsData)

		err = s.repositoryOrd.LineOrdersUpdateOrInsert(rsData)
		if err != nil {
			log.Println("Error actualización listado de ordenes en Postgrest LineOrderList:", err)
			return data, err
		}
	}

	// Leemos las ordenes de fabricación desde POSTGRES
	ordersData, err := s.repositoryOrd.LineOrdersFind(factory, lineNumber)
	if err != nil {
		log.Println("Error lectura listado de ordenes Postgrest LineOrderList:", err)
		return data, err
	}
	for _, p := range ordersData {
		data = append(data, msProductionOrder{
			Factory:                   p.Factory,
			ProdLine:                  p.ProdLine,
			OrderNumber:               p.OrderNumber,
			OrderNType:                p.OrderNType,
			ProductName:               p.ProductName,
			ProductDescription:        p.ProductDescription,
			QuantityToProduce:         p.QuantityToProduce,
			QuantityProduced:          p.QuantityProduced,
			QuantityRemainedToProduce: p.QuantityRemainedToProduce,
			MeasurementUnit:           p.MeasurementUnit,
			StarteddAt:                p.StarteddAt.Truncate(time.Second),
			FinishedAt:                p.FinishedAt.Truncate(time.Second),
		})
	}

	return data, nil
}

func (s *service) LineOrderStartFinish(factory string, lineNumber string, lineSapCode string, sapRequest string, orderNumber string, startFinish string) ([]msProductionOrder, error) {
	var data = []msProductionOrder{}

	err := s.repositoryOrd.LineOrderStartFinish(factory, lineNumber, orderNumber, startFinish)
	if err != nil {
		log.Println("Error Service LineOrderStartFinish:", err)
		return data, err
	}
	// Leemos las ordenes de fabricación desde POSTGRES
	ordersData, err := s.repositoryOrd.LineOrdersFind(factory, lineNumber)
	if err != nil {
		log.Println("Error lectura listado de ordenes Postgrest LineOrderList:", err)
		return data, err
	}
	for _, p := range ordersData {
		data = append(data, msProductionOrder{
			Factory:                   p.Factory,
			ProdLine:                  p.ProdLine,
			OrderNumber:               p.OrderNumber,
			OrderNType:                p.OrderNType,
			ProductName:               p.ProductName,
			ProductDescription:        p.ProductDescription,
			QuantityToProduce:         p.QuantityToProduce,
			QuantityProduced:          p.QuantityProduced,
			QuantityRemainedToProduce: p.QuantityRemainedToProduce,
			MeasurementUnit:           p.MeasurementUnit,
			StarteddAt:                p.StarteddAt,
			FinishedAt:                p.FinishedAt,
		})
	}

	return data, nil
}

func (s *service) LineOrderUpdateTime(factory string, lineNumber string, lineSapCode string, sapRequest string, orderNumber string, starteddAt time.Time, finishedAt time.Time) ([]msProductionOrder, error) {
	var data = []msProductionOrder{}

	// Obtenemos el activo por la linea de la fabrica
	// assetOrder, err := s.repositoryAss.AssetByFactLine(factory, lineNumber)
	// if err != nil {
	// 	log.Println("Activo no detectado Service LineOrderList:", err)
	// 	return data, err
	// }

	err := s.repositoryOrd.LineOrderUpdateTime(factory, lineNumber, orderNumber, starteddAt, finishedAt)
	if err != nil {
		log.Println("Error Service LineOrderStartFinish:", err)
		return data, err
	}
	// Leemos las ordenes de fabricación desde POSTGRES
	ordersData, err := s.repositoryOrd.LineOrdersFind(factory, lineNumber)
	if err != nil {
		log.Println("Error lectura listado de ordenes Postgrest LineOrderList:", err)
		return data, err
	}
	for _, p := range ordersData {
		data = append(data, msProductionOrder{
			Factory:                   p.Factory,
			ProdLine:                  p.ProdLine,
			OrderNumber:               p.OrderNumber,
			OrderNType:                p.OrderNType,
			ProductName:               p.ProductName,
			ProductDescription:        p.ProductDescription,
			QuantityToProduce:         p.QuantityToProduce,
			QuantityProduced:          p.QuantityProduced,
			QuantityRemainedToProduce: p.QuantityRemainedToProduce,
			MeasurementUnit:           p.MeasurementUnit,
			StarteddAt:                p.StarteddAt,
			FinishedAt:                p.FinishedAt,
		})
	}

	return data, nil
}
