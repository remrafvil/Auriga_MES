package rLabor

import (
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	GetAllShifts() ([]rModels.MrShift, error)
	GetShiftByID(id uint) (*rModels.MrShift, error)
	CreateShift(shift *rModels.MrShift) error
	UpdateShift(shift *rModels.MrShift) error
	DeleteShift(id uint) error
	GetShiftsByFactory(factoryID uint) ([]rModels.MrShift, error)

	GetAllTeams() ([]rModels.MrTeam, error)
	GetTeamByID(id uint) (*rModels.MrTeam, error)
	CreateTeam(team *rModels.MrTeam) error
	UpdateTeam(team *rModels.MrTeam) error
	DeleteTeam(id uint) error
	GetTeamsByFactory(factoryID uint) ([]rModels.MrTeam, error)
	GetTeamsByDepartment(departmentID uint) ([]rModels.MrTeam, error)
	AddTeamMember(teamMember *rModels.MrTeamMember) error
	RemoveTeamMember(teamID, employeeID uint) error
	GetTeamMembers(teamID uint) ([]rModels.MrTeamMember, error)

	GetAssignmentsByEmployee(employeeID uint) ([]rModels.MrShiftAssignment, error)
	GetAssignmentsByTeam(teamID uint) ([]rModels.MrShiftAssignment, error)
	GetAssignmentsByShift(shiftID uint) ([]rModels.MrShiftAssignment, error)
	CreateAssignment(assignment *rModels.MrShiftAssignment) error
	UpdateAssignment(assignment *rModels.MrShiftAssignment) error
	DeleteAssignment(id uint) error
	GetCurrentAssignment(employeeID uint, date time.Time) (*rModels.MrShiftAssignment, error)
	GetAssignmentsByDateRange(startDate, endDate time.Time) ([]rModels.MrShiftAssignment, error)
}

type repository struct {
	db     *gorm.DB
	logger *zap.Logger
}

func New(db *gorm.DB, logger *zap.Logger) Repository {
	return &repository{
		db:     db,
		logger: logger,
	}
}
