package rEvents

// lo llamaremso repositories

import (
	"fmt"
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

func (m *repository) EventsCommitByLineList(factory string, location string) ([]rModels.MrCommitEvents, error) {
	var data []rModels.MrCommitEvents
	if err := m.db.Where("factory = ? AND prod_line = ?", factory, location).Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (m *repository) EventsCommitByLineAdd(eventTime time.Time, factory string, prodLine string, system string, machine string, part string, eventTypt string, eventCategory string) ([]rModels.MrCommitEvents, error) {
	var eventsCommit []rModels.MrCommitEvents

	newEventCommit := rModels.MrCommitEvents{
		EventTime:     eventTime,
		Factory:       factory,
		ProdLine:      prodLine,
		System:        system,
		Machine:       machine,
		Part:          part,
		EventType:     eventTypt,
		EventCategory: eventCategory,
	}
	// Insertar el registro en la base de datos
	err := m.db.Create(&newEventCommit).Error
	if err != nil {
		if m.db.Error != nil && gorm.ErrDuplicatedKey == err {
			fmt.Println("Error: Duplicate record")
		} else {
			fmt.Println("Error inserting data:", err)
		}
		return eventsCommit, err
	} else {
		fmt.Println("Record inserted successfully:", newEventCommit)
	}
	return eventsCommit, nil
}

func (m *repository) EventsCommitByLineUpdate(id uint, eventTime time.Time, factory string, prodLine string, system string, machine string, part string, eventTypt string, eventCategory string) ([]rModels.MrCommitEvents, error) {
	var eventsCommit []rModels.MrCommitEvents

	updateEventCommit := rModels.MrCommitEvents{
		// MrComponentSapCode: sapComponentCode,
		EventTime:     eventTime,
		Factory:       factory,
		ProdLine:      prodLine,
		System:        system,
		Machine:       machine,
		Part:          part,
		EventType:     eventTypt,
		EventCategory: eventCategory,
	}
	result := m.db.Model(&rModels.MrCommitEvents{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"event_time":     updateEventCommit.EventTime,
			"factory":        updateEventCommit.Factory,
			"prod_line":      updateEventCommit.ProdLine,
			"system":         updateEventCommit.System,
			"machine":        updateEventCommit.Machine,
			"part":           updateEventCommit.Part,
			"event_type":     updateEventCommit.EventType,
			"event_category": updateEventCommit.EventCategory,
		})

	if result.Error != nil {
		fmt.Println("Error updating Event Commit:", result.Error)
		return eventsCommit, result.Error
	} else if result.RowsAffected == 0 {
		fmt.Println("No Event Commit found to Update")
	} else {
		fmt.Println("Event Commit Updated Successfully")
	}
	return eventsCommit, nil
}

func (m *repository) EventsCommitByLineDel(id uint) ([]rModels.MrCommitEvents, error) {
	var eventsCommit []rModels.MrCommitEvents

	err := m.db.Delete(&rModels.MrCommitEvents{}, id).Error
	if err != nil {
		panic("failed to delete record")
	}

	// Confirmar que se eliminó el registro
	var count int64
	m.db.Model(&rModels.MrCommitEvents{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		println("Registro eliminado con éxito")
	} else {
		println("El registro aún existe")
	}
	return eventsCommit, nil
}
func (m *repository) EventsCommitByLineFind(id uint) (rModels.MrCommitEvents, string, error) {
	var eventsCommit rModels.MrCommitEvents
	var order rModels.MrProductionOrder
	err := m.db.First(&eventsCommit, id).Error

	if err != nil {
		return rModels.MrCommitEvents{}, "", fmt.Errorf("evento no encontrado: %w", err)
	}

	err = m.db.
		Where("? BETWEEN startedd_at AND finished_at", eventsCommit.CreatedAt).
		Where("factory = ? AND prod_line = ?", eventsCommit.Factory, eventsCommit.ProdLine).
		First(&order).Error

	if err != nil {
		fmt.Println("No se encontró la orden de fabricación para el evento:", err)
		return rModels.MrCommitEvents{}, "", err
	}

	return eventsCommit, order.OrderNumber, nil
}

func DB_InitEventsCommit(c *gorm.DB) { //
	timeString := "2024-05-195 02:30:45"
	theTime, _ := time.Parse("2006-01-02 03:04:05", timeString)
	c.Create(&rModels.MrCommitEvents{EventTime: theTime, EventType: "Fault", EventCategory: "Trip", Factory: "FSP", ProdLine: "Line_02", System: "Extrusion", Machine: "ExtruderA", Part: "Z01_01"})
}
