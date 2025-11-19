package rLabor

import (
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
)

func (r *repository) GetAllShifts() ([]rModels.MrShift, error) {
	var shifts []rModels.MrShift
	err := r.db.Where("active = ?", true).Find(&shifts).Error
	return shifts, err
}

func (r *repository) GetShiftByID(id uint) (*rModels.MrShift, error) {
	var shift rModels.MrShift
	err := r.db.First(&shift, id).Error
	return &shift, err
}

func (r *repository) CreateShift(shift *rModels.MrShift) error {
	return r.db.Create(shift).Error
}

func (r *repository) UpdateShift(shift *rModels.MrShift) error {
	return r.db.Save(shift).Error
}

func (r *repository) DeleteShift(id uint) error {
	return r.db.Model(&rModels.MrShift{}).Where("id = ?", id).Update("active", false).Error
}

func (r *repository) GetShiftsByFactory(factoryID uint) ([]rModels.MrShift, error) {
	var shifts []rModels.MrShift
	// Esta consulta necesitaría ajustarse según tu lógica de negocio
	err := r.db.Where("active = ?", true).Find(&shifts).Error
	return shifts, err
}
