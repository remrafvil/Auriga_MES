package rLabor

import (
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
)

func (r *repository) GetAssignmentsByEmployee(employeeID uint) ([]rModels.MrShiftAssignment, error) {
	var assignments []rModels.MrShiftAssignment
	err := r.db.Preload("Shift").Preload("Employee").
		Where("employee_id = ? AND active = ?", employeeID, true).
		Find(&assignments).Error
	return assignments, err
}

func (r *repository) GetAssignmentsByTeam(teamID uint) ([]rModels.MrShiftAssignment, error) {
	var assignments []rModels.MrShiftAssignment
	// Obtener assignments a través de los miembros del equipo
	err := r.db.Preload("Shift").Preload("Employee").
		Joins("JOIN mr_team_members ON mr_team_members.employee_id = mr_shift_assignments.employee_id").
		Where("mr_team_members.team_id = ? AND mr_shift_assignments.active = ? AND mr_team_members.active = ?",
			teamID, true, true).
		Find(&assignments).Error
	return assignments, err
}

func (r *repository) GetAssignmentsByShift(shiftID uint) ([]rModels.MrShiftAssignment, error) {
	var assignments []rModels.MrShiftAssignment
	err := r.db.Preload("Shift").Preload("Employee").
		Where("shift_id = ? AND active = ?", shiftID, true).
		Find(&assignments).Error
	return assignments, err
}

func (r *repository) CreateAssignment(assignment *rModels.MrShiftAssignment) error {
	return r.db.Create(assignment).Error
}

func (r *repository) UpdateAssignment(assignment *rModels.MrShiftAssignment) error {
	return r.db.Save(assignment).Error
}

func (r *repository) DeleteAssignment(id uint) error {
	return r.db.Model(&rModels.MrShiftAssignment{}).Where("id = ?", id).Update("active", false).Error
}

func (r *repository) GetCurrentAssignment(employeeID uint, date time.Time) (*rModels.MrShiftAssignment, error) {
	var assignment rModels.MrShiftAssignment
	err := r.db.Preload("Shift").Preload("Employee").
		Where("employee_id = ? AND active = ? AND start_date <= ? AND (end_date IS NULL OR end_date >= ?)",
			employeeID, true, date, date).
		First(&assignment).Error
	return &assignment, err
}

func (r *repository) GetAssignmentsByDateRange(startDate, endDate time.Time) ([]rModels.MrShiftAssignment, error) {
	var assignments []rModels.MrShiftAssignment
	err := r.db.Preload("Shift").Preload("Employee").
		Where("active = ? AND ((start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?) OR (start_date <= ? AND (end_date IS NULL OR end_date >= ?)))",
			true, startDate, endDate, startDate, endDate, startDate, endDate).
		Find(&assignments).Error
	return assignments, err
}
