package rEvents

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

func (m *repository) EventsRawByLineList(factory string, location string) ([]rModels.MrRawEvents, error) {
	var data []rModels.MrRawEvents
	if err := m.db.Where("\"EL_Lv0\" = ? AND name = ?", factory, location).Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (m *repository) EventsRawByLineDel(id uint) ([]rModels.MrRawEvents, error) {
	var eventsRaw []rModels.MrRawEvents

	err := m.db.Delete(&rModels.MrRawEvents{}, id).Error
	if err != nil {
		panic("failed to delete record")
	}

	// Confirmar que se eliminó el registro
	var count int64
	m.db.Model(&rModels.MrRawEvents{}).Where("id = ?", id).Count(&count)
	if count == 0 {
		println("Registro eliminado con éxito")
	} else {
		println("El registro aún existe")
	}
	return eventsRaw, nil
}

func (m *repository) EventsRawToCommitLine(id uint, eventTime time.Time, factory string, prodline string, system string, machine string, part string, eventTypt string) ([]rModels.MrCommitEvents, error) {
	var eventsCommit []rModels.MrCommitEvents

	err := m.db.Transaction(func(tx *gorm.DB) error {
		var rawEvent rModels.MrRawEvents
		// Leer el registro de rModels.MrRawEvents
		log.Println("Leer el registro de rModels.MrRawEvents:    KKKKKKKKKKKKKKKKKKKKKKKKKKKKKKK ")
		if err := tx.First(&rawEvent, id).Error; err != nil {
			return fmt.Errorf("error fetching rModels.MrRawEvents: %w", err)
		}
		esw_String := strconv.Itoa(rawEvent.ESW)
		log.Println("rawEvent.Time: ", rawEvent.Time)
		// eventTimeDB, _ := time.Parse("2006-01-02 15:04:05 -0700 MST", rawEvent.Time)
		// log.Println("esw_String: ", eventTimeDB)

		// Crear un registro en rModels.MrCommitEvents con los datos obtenidos
		commitEvent := rModels.MrCommitEvents{
			EventTime: rawEvent.Time, // Puedes mapear los datos según lo necesario
			EventType: esw_String,
			Factory:   rawEvent.EL_Lv0,
			ProdLine:  rawEvent.Name,
			System:    rawEvent.EL_Lv1,
			Machine:   rawEvent.EL_Lv2,
			Part:      rawEvent.EL_Lv3,
		}
		if err := tx.Create(&commitEvent).Error; err != nil {
			return fmt.Errorf("error creating rModels.MrCommitEvents: %w", err)
		}

		// Borrar el registro de rModels.MrRawEvents
		if err := tx.Delete(&rawEvent).Error; err != nil {
			return fmt.Errorf("error deleting rModels.MrRawEvents: %w", err)
		}

		// Si todo se ejecuta sin errores, se confirma la transacción
		return nil
	})

	// Manejo de errores fuera de la transacción
	if err != nil {
		fmt.Printf("Transaction failed: %v\n", err)
	} else {
		fmt.Println("Transaction succeeded")
	}
	return eventsCommit, nil
}
