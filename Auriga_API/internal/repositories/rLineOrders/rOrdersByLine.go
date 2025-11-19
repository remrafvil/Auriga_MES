package rLineOrders

// lo llamaremso repositories

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

func (m *repository) LineOrdersList() ([]rModels.MrProductionOrder, error) { // Lista todos
	var data = []rModels.MrProductionOrder{}

	log.Println("*****************      Llego por aquí REPOSITORY *****************")
	err := m.db.Find(&data).Error
	if err != nil {
		log.Fatalln(err)
		return data, err
	}
	return data, nil
}

func (m *repository) LineOrdersInfo(id uint) (rModels.MrProductionOrder, error) { // Encuentra uno
	var data rModels.MrProductionOrder
	rs := m.db.Where("id  = ?", id).Find(&data)
	if rs.RowsAffected == 0 {
		log.Println("no existe registro con ese identificador")
		return rModels.MrProductionOrder{}, errors.New("no existe registro con ese id")
	}
	if rs.Error != nil {

		return rModels.MrProductionOrder{}, errors.New("fallo registro con ese id")
	}

	return data, nil
}

func (m *repository) LineOrdersAdd(Factory string, ProdLine string, OrderNumber string, OrderNType string, ProductName string, ProductDescription string, QuantityToProduce string, QuantityProduced string, QuantityRemainedToProduce string, MeasurementUnit string, StarteddAt time.Time, FinishedAt time.Time) (rModels.MrProductionOrder, error) { // Añade uno
	var data rModels.MrProductionOrder

	data.Factory = Factory
	data.ProdLine = ProdLine
	data.OrderNumber = OrderNumber
	data.OrderNType = OrderNType
	data.ProductName = ProductName
	data.ProductDescription = ProductDescription
	data.QuantityToProduce = QuantityToProduce
	data.QuantityProduced = QuantityProduced
	data.QuantityRemainedToProduce = QuantityRemainedToProduce
	data.MeasurementUnit = MeasurementUnit
	data.StarteddAt = StarteddAt
	data.FinishedAt = FinishedAt

	if err := m.db.Create(&data).Error; err != nil {
		return rModels.MrProductionOrder{}, err
	} else {

		return data, nil
	}
}

func (m *repository) LineOrdersUpdate(id uint, Factory string, ProdLine string, OrderNumber string, OrderNType string, ProductName string, ProductDescription string, QuantityToProduce string, QuantityProduced string, QuantityRemainedToProduce string, MeasurementUnit string, StarteddAt time.Time, FinishedAt time.Time) (rModels.MrProductionOrder, error) { // Actualiza uno
	var data rModels.MrProductionOrder

	if err := m.db.Where("id = ? ", id).First(&data).Error; err != nil {
		return rModels.MrProductionOrder{}, errors.New("no hay registro con ese id")
	}
	data.Factory = Factory
	data.ProdLine = ProdLine
	data.OrderNumber = OrderNumber
	data.OrderNType = OrderNType
	data.ProductName = ProductName
	data.ProductDescription = ProductDescription
	data.QuantityToProduce = QuantityToProduce
	data.QuantityProduced = QuantityProduced
	data.QuantityRemainedToProduce = QuantityRemainedToProduce
	data.MeasurementUnit = MeasurementUnit
	data.StarteddAt = StarteddAt
	data.FinishedAt = FinishedAt

	if err := m.db.Save(&data).Error; err != nil {
		return rModels.MrProductionOrder{}, errors.New("no se pudo actualizar")
	}

	log.Println("*****************      Llego por aquí REPOSITORY *****************")
	return data, nil

}

func (m *repository) LineOrdersDel(id uint) error { // Borra uno
	var data rModels.MrProductionOrder
	if err := m.db.Where("id = ?", id).Delete(&data).Error; err != nil {
		return err
	}
	return nil
}

func (m *repository) LineOrdersUpdateOrInsert(lineOrders []rModels.MrProductionOrder) error {
	for _, p := range lineOrders {
		if err := m.db.FirstOrCreate(&p, rModels.MrProductionOrder{OrderNumber: p.OrderNumber}).Error; err != nil {
			return err
		}
	}
	log.Println("*****************      Llego por aquí REPOSITORY LINE Update o Create *****************")
	return nil
}

func (m *repository) LineOrdersFind(factory string, location string) ([]rModels.MrProductionOrder, error) {
	var data []rModels.MrProductionOrder
	if err := m.db.Where("factory = ? AND prod_line = ?", factory, location).Find(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}

func (m *repository) LineOrderStartFinish(factory string, prodLine string, orderNumber string, startFinish string) error { // Borra uno
	var startFinishLabel string

	updateOrder := rModels.MrProductionOrder{
		Factory:     factory,
		ProdLine:    prodLine,
		OrderNumber: orderNumber,
	}
	if startFinish == "Start" {
		startFinishLabel = "startedd_at"
	} else if startFinish == "Finish" {
		startFinishLabel = "finished_at"
	}

	result := m.db.Model(&rModels.MrProductionOrder{}).
		Where("factory = ? AND prod_line = ? AND order_number = ?",
			updateOrder.Factory,
			updateOrder.ProdLine,
			updateOrder.OrderNumber,
		).
		Update(startFinishLabel, time.Now())

	if result.Error != nil {
		fmt.Println("Error updating record:", result.Error)
		return result.Error
	} else if result.RowsAffected == 0 {
		fmt.Println("No record found to update")
	} else {
		fmt.Println("Record updated successfully")
	}
	return nil
}

func (m *repository) LineOrderUpdateTime(factory string, prodLine string, orderNumber string, starteddAt time.Time, finishedAt time.Time) error { // Borra uno

	result := m.db.Model(&rModels.MrProductionOrder{}).
		Where("factory = ? AND prod_line = ? AND order_number = ?", factory, prodLine, orderNumber).
		Updates(map[string]interface{}{
			"startedd_at": starteddAt,
			"finished_at": finishedAt,
		})

	if result.Error != nil {
		fmt.Println("Error updating record:", result.Error)
		return result.Error
	} else if result.RowsAffected == 0 {
		fmt.Println("No record found to update")
	} else {
		fmt.Println("Record updated successfully")
	}
	return nil
}

func (m *repository) LineOrdersFindStartFinish(factory string, location string, orderNumber string) (time.Time, time.Time, error) {
	var data rModels.MrProductionOrder
	var starteddAt, finishedAt time.Time
	if err := m.db.Where("factory = ? AND prod_line = ? AND order_number = ?", factory, location, orderNumber).Find(&data).Error; err != nil {
		return starteddAt, finishedAt, err
	}
	starteddAt = data.StarteddAt
	finishedAt = data.FinishedAt

	return starteddAt, finishedAt, nil
}

func DB_InitLineOrders(c *gorm.DB) { //
	c.Create(&rModels.MrProductionOrder{Factory: "FSP", ProdLine: "Line_02", OrderNumber: "000010068196", OrderNType: "ZP01", ProductName: "EX00026361", ProductDescription: "BR.CREMA 457 X 1.10", QuantityToProduce: "40000.000", QuantityProduced: "39242.000", QuantityRemainedToProduce: "758.000", MeasurementUnit: "KG", StarteddAt: time.Now(), FinishedAt: time.Now()})
	c.Create(&rModels.MrProductionOrder{Factory: "FSP", ProdLine: "Line_02", OrderNumber: "000010068197", OrderNType: "ZP01", ProductName: "EX00026361", ProductDescription: "BR.CREMA 457 X 1.10", QuantityToProduce: "40000.000", QuantityProduced: "39242.000", QuantityRemainedToProduce: "758.000", MeasurementUnit: "KG", StarteddAt: time.Now(), FinishedAt: time.Now()})
	c.Create(&rModels.MrProductionOrder{Factory: "FSP", ProdLine: "Line_02", OrderNumber: "000010068198", OrderNType: "ZP01", ProductName: "EX00026361", ProductDescription: "BR.CREMA 457 X 1.10", QuantityToProduce: "40000.000", QuantityProduced: "39242.000", QuantityRemainedToProduce: "758.000", MeasurementUnit: "KG", StarteddAt: time.Now(), FinishedAt: time.Now()})
	c.Create(&rModels.MrProductionOrder{Factory: "FSP", ProdLine: "Line_02", OrderNumber: "000010068199", OrderNType: "ZP01", ProductName: "EX00026361", ProductDescription: "BR.CREMA 457 X 1.10", QuantityToProduce: "40000.000", QuantityProduced: "39242.000", QuantityRemainedToProduce: "758.000", MeasurementUnit: "KG", StarteddAt: time.Now(), FinishedAt: time.Now()})

}
