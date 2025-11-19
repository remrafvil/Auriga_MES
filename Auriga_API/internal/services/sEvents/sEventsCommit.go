package sEvents

import (
	"log"
	"time"
)

type msCommitEvents struct {
	EventTime time.Time `json:"EventTime"`
	Type      string    `json:"Type"`
	Category  string    `json:"Category"`
	Factory   string    `json:"Factory"`
	ProdLine  string    `json:"ProdLine"`
	System    string    `json:"System"`
	Machine   string    `json:"Machine"`
	Part      string    `json:"Part"`
	ID        uint      `json:"ID"`
}

func (s *service) EventsCommitByLineList(factory string, lineNumber string) ([]msCommitEvents, error) {
	var data = []msCommitEvents{}

	// Leemos el listado de eventos registrados automáticamente por la línea desde POSTGRES
	commitEventData, err := s.repositoryEven.EventsCommitByLineList(factory, lineNumber)
	if err != nil {
		log.Println("Error lectura listado de Eventos Raw Service EventsCommitByLineList:", err)
		return data, err
	}
	for _, p := range commitEventData {
		//log.Println("EventTime:", p.EventTime)
		data = append(data, msCommitEvents{
			EventTime: p.EventTime,
			Type:      p.EventType,
			Category:  p.EventCategory,
			Factory:   p.Factory,
			ProdLine:  p.ProdLine,
			System:    p.System,
			Machine:   p.Machine,
			Part:      p.Part,
			ID:        p.ID,
		})
	}

	return data, nil
}

func (s *service) EventsCommitByLineAdd(eventTime time.Time, factory string, prodline string, system string, machine string, part string, eventTypt string, eventCategory string) ([]msCommitEvents, error) {
	var data = []msCommitEvents{}

	eventsCommit, err := s.repositoryEven.EventsCommitByLineAdd(eventTime, factory, prodline, system, machine, part, eventTypt, eventCategory)
	if err != nil {
		log.Println("Error añadir Service EventsCommitByLineAdd:", err)
		return data, err
	}

	log.Println("Events Commit:", eventsCommit)

	return data, nil
}

func (s *service) EventsCommitByLineUpdate(id uint, eventTime time.Time, factory string, prodline string, system string, machine string, part string, eventTypt string, eventCategory string) ([]msCommitEvents, error) {
	var data = []msCommitEvents{}

	eventsCommit, err := s.repositoryEven.EventsCommitByLineUpdate(id, eventTime, factory, prodline, system, machine, part, eventTypt, eventCategory)
	if err != nil {
		log.Println("Error actualizar Service EventsCommitByLineUpdate:", err)
		return data, err
	}

	log.Println("Events Commit:", eventsCommit)

	return data, nil
}

func (s *service) EventsCommitByLineDel(id uint) ([]msCommitEvents, error) {
	var data = []msCommitEvents{}

	eventsCommit, err := s.repositoryEven.EventsCommitByLineDel(id)
	if err != nil {
		log.Println("Error borrar Service EventsCommitByLineDel:", err)
		return data, err
	}

	log.Println("Events Commit:", eventsCommit)

	return data, nil
}
