package sEvents

import (
	"fmt"
	"log"
)

func (s *service) EventsSapByLineDel(id uint) ([]msCommitEvents, error) {
	var data = []msCommitEvents{}

	// Leo el evento por ID de la base de datos
	eventCommit, orderNumber, err := s.repositoryEven.EventsCommitByLineFind(id)
	if err != nil {
		log.Println("Error borrar Service EventsCommitByLineDel:", err)
		return data, err
	}
	fmt.Println("Events Commit Time:", eventCommit.EventTime)
	fmt.Println("Events Commit Factory:", eventCommit.Factory)
	fmt.Println("Events Commit ProdLine:", eventCommit.ProdLine)
	fmt.Println("Events Commit System:", eventCommit.System)
	fmt.Println("Events Commit Machine:", eventCommit.Machine)
	fmt.Println("Events Commit Part:", eventCommit.Part)
	fmt.Println("Events Commit EventType:", eventCommit.EventType)
	fmt.Println("Events Commit EventCategory:", eventCommit.EventCategory)
	fmt.Println("Events Commit Production Order:", orderNumber)

	// Obtenemos los datos del activo por la linea de la fabrica
	assetOrder, err := s.repositoryAss.AssetByFactLine(eventCommit.Factory, eventCommit.ProdLine)
	if err != nil {
		log.Println("Activo no detectado Service DosingConsumptionList:", err)
		return data, err
	}

	fmt.Println("El código SAP de la línea es:", assetOrder.SapCode)

	// Aquí implemento la función de envío a SAP
	err = s.repositorySap.RsLineStopEvent(assetOrder.SapCode, eventCommit.EventTime, eventCommit.EventType, eventCommit.EventCategory, orderNumber, "operario")
	if err != nil {
		log.Println("Error lectura listado de ordenes SAP LineOrderList:", err)
		return data, err
	}

	// Borro el evento de la base de datos

	eventsCommitReturn, err := s.repositoryEven.EventsCommitByLineDel(id)
	if err != nil {
		log.Println("Error borrar Service EventsCommitByLineDel:", err)
		return data, err
	}

	log.Println("Events Commit:", eventsCommitReturn)

	return data, nil
}
