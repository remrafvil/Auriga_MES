package rLabor

import (
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
)

func (r *repository) GetAllTeams() ([]rModels.MrTeam, error) {
	var teams []rModels.MrTeam
	err := r.db.Preload("Factory").Preload("Department").Preload("Leader").
		Where("active = ?", true).Find(&teams).Error
	return teams, err
}

func (r *repository) GetTeamByID(id uint) (*rModels.MrTeam, error) {
	var team rModels.MrTeam
	err := r.db.Preload("Factory").Preload("Department").Preload("Leader").
		Preload("TeamMembers.Employee").
		First(&team, id).Error
	return &team, err
}

func (r *repository) CreateTeam(team *rModels.MrTeam) error {
	return r.db.Create(team).Error
}

func (r *repository) UpdateTeam(team *rModels.MrTeam) error {
	return r.db.Save(team).Error
}

func (r *repository) DeleteTeam(id uint) error {
	return r.db.Model(&rModels.MrTeam{}).Where("id = ?", id).Update("active", false).Error
}

func (r *repository) GetTeamsByFactory(factoryID uint) ([]rModels.MrTeam, error) {
	var teams []rModels.MrTeam
	err := r.db.Preload("Factory").Preload("Department").Preload("Leader").
		Where("factory_id = ? AND active = ?", factoryID, true).Find(&teams).Error
	return teams, err
}

func (r *repository) GetTeamsByDepartment(departmentID uint) ([]rModels.MrTeam, error) {
	var teams []rModels.MrTeam
	err := r.db.Preload("Factory").Preload("Department").Preload("Leader").
		Where("department_id = ? AND active = ?", departmentID, true).Find(&teams).Error
	return teams, err
}

func (r *repository) AddTeamMember(teamMember *rModels.MrTeamMember) error {
	return r.db.Create(teamMember).Error
}

func (r *repository) RemoveTeamMember(teamID, employeeID uint) error {
	return r.db.Where("team_id = ? AND employee_id = ?", teamID, employeeID).
		Delete(&rModels.MrTeamMember{}).Error
}

func (r *repository) GetTeamMembers(teamID uint) ([]rModels.MrTeamMember, error) {
	var members []rModels.MrTeamMember
	err := r.db.Preload("Employee").Preload("Team").
		Where("team_id = ? AND active = ?", teamID, true).Find(&members).Error
	return members, err
}
