package sEvents

import (
	"log"
	"time"
)

type msRawEvents struct {
	EventTime time.Time `json:"EventTime"`
	Event     int       `json:"Event"`
	Factory   string    `json:"Factory"`
	ProdLine  string    `json:"ProdLine"`
	System    string    `json:"System"`
	Machine   string    `json:"Machine"`
	Part      string    `json:"Part"`
	ID        uint      `json:"ID"`
}

func (s *service) EventsRawByLineList(factory string, lineNumber string) ([]msRawEvents, error) {
	var data = []msRawEvents{}

	// Leemos el listado de eventos registrados automáticamente por la línea desde POSTGRES
	rawEventData, err := s.repositoryEven.EventsRawByLineList(factory, lineNumber)
	if err != nil {
		log.Println("Error lectura listado de Eventos Raw Service EventsRawByLineList:", err)
		return data, err
	}
	for _, p := range rawEventData {
		data = append(data, msRawEvents{
			EventTime: p.Time,
			Event:     p.ESW,
			Factory:   p.EL_Lv0,
			ProdLine:  p.Name,
			System:    p.EL_Lv1,
			Machine:   p.EL_Lv2,
			Part:      p.EL_Lv3,
			ID:        p.ID,
		})
	}

	return data, nil
}

func (s *service) EventsRawToCommitLine(id uint, eventTime time.Time, factory string, prodline string, system string, machine string, part string, eventTypt string) ([]msCommitEvents, error) {
	var data = []msCommitEvents{}

	eventsCommit, err := s.repositoryEven.EventsRawToCommitLine(id, eventTime, factory, prodline, system, machine, part, eventTypt)
	if err != nil {
		log.Println("Error actualizar Service EventsCommitByLineUpdate:", err)
		return data, err
	}

	log.Println("Events Commit:", eventsCommit)

	return data, nil
}

func (s *service) EventsRawByLineDel(id uint) ([]msRawEvents, error) {
	var data = []msRawEvents{}

	eventsRaw, err := s.repositoryEven.EventsRawByLineDel(id)
	if err != nil {
		log.Println("Error borrar Service EventsCommitByLineDel:", err)
		return data, err
	}

	log.Println("Events Commit:", eventsRaw)

	return data, nil
}
