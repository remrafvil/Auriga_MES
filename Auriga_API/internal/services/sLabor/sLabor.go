package sLabor

import (
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rLabor"
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"go.uber.org/zap"
)

type Service interface {
	GetAllShifts() ([]rModels.MrShift, error)
	GetShiftByID(id uint) (*rModels.MrShift, error)
	CreateShift(req CreateShiftRequest) (*rModels.MrShift, error)
	UpdateShift(id uint, req UpdateShiftRequest) (*rModels.MrShift, error)
	DeleteShift(id uint) error

	GetAllTeams() ([]rModels.MrTeam, error)
	GetTeamByID(id uint) (*rModels.MrTeam, error)
	CreateTeam(req CreateTeamRequest) (*rModels.MrTeam, error)
	UpdateTeam(id uint, req UpdateTeamRequest) (*rModels.MrTeam, error)
	DeleteTeam(id uint) error
	AddTeamMember(req AddTeamMemberRequest) (*rModels.MrTeamMember, error)
	RemoveTeamMember(teamID, employeeID uint) error
	GetTeamMembers(teamID uint) ([]rModels.MrTeamMember, error)

	GetEmployeeAssignments(employeeID uint) ([]rModels.MrShiftAssignment, error)
	GetTeamAssignments(teamID uint) ([]rModels.MrShiftAssignment, error)
	CreateIndividualAssignment(req CreateShiftAssignmentRequest) (*rModels.MrShiftAssignment, error)
	CreateTeamAssignment(req CreateShiftAssignmentRequest) ([]rModels.MrShiftAssignment, error)
	CreateBulkAssignments(req BulkShiftAssignmentRequest) ([]rModels.MrShiftAssignment, error)
	UpdateAssignment(id uint, req CreateShiftAssignmentRequest) (*rModels.MrShiftAssignment, error)
	DeleteAssignment(id uint) error
	GetCurrentEmployeeAssignment(employeeID uint, date time.Time) (*rModels.MrShiftAssignment, error)
}

// service implementación
type service struct {
	repository rLabor.Repository
	logger     *zap.Logger
}

func New(repository rLabor.Repository, logger *zap.Logger) Service {
	return &service{
		repository: repository,
		logger:     logger,
	}
}

// SHIFT DTOs
type CreateShiftRequest struct {
	Name        string    `json:"name" validate:"required,min=1"`
	StartTime   time.Time `json:"start_time" validate:"required,notzerotime"`
	EndTime     time.Time `json:"end_time" validate:"required,notzerotime"`
	Description string    `json:"description"`
	Tolerance   int       `json:"tolerance" validate:"required,min=0,max=60"`
}

type UpdateShiftRequest struct {
	Name        string    `json:"name" validate:"omitempty,min=1"`
	StartTime   time.Time `json:"start_time" validate:"omitempty,notzerotime"`
	EndTime     time.Time `json:"end_time" validate:"omitempty,notzerotime"`
	Description string    `json:"description"`
	Tolerance   int       `json:"tolerance" validate:"omitempty,min=0,max=60"`
	Active      *bool     `json:"active"`
}

// TEAM DTOs
type CreateTeamRequest struct {
	FactoryID    uint   `json:"factory_id" validate:"required,min=1"`
	DepartmentID uint   `json:"department_id" validate:"required,min=1"`
	Name         string `json:"name" validate:"required,min=1"`
	Description  string `json:"description"`
	LeaderID     uint   `json:"leader_id" validate:"required,min=1"`
}

type UpdateTeamRequest struct {
	Name        string `json:"name" validate:"omitempty,min=1"`
	Description string `json:"description"`
	LeaderID    uint   `json:"leader_id" validate:"omitempty,min=1"`
	Active      *bool  `json:"active"`
}

type AddTeamMemberRequest struct {
	TeamID     uint      `json:"team_id" validate:"required,min=1"`
	EmployeeID uint      `json:"employee_id" validate:"required,min=1"`
	Role       string    `json:"role" validate:"required,min=1"`
	StartDate  time.Time `json:"start_date" validate:"required,notzerotime"`
}

// SHIFT ASSIGNMENT DTOs
type CreateShiftAssignmentRequest struct {
	EmployeeID uint      `json:"employee_id" validate:"required_without=TeamID,min=1"`
	TeamID     uint      `json:"team_id" validate:"required_without=EmployeeID,min=1"`
	ShiftID    uint      `json:"shift_id" validate:"required,min=1"`
	StartDate  time.Time `json:"start_date" validate:"required,notzerotime"`
	EndDate    time.Time `json:"end_date" validate:"omitempty,notzerotime"`
}

type BulkShiftAssignmentRequest struct {
	TeamID      uint      `json:"team_id" validate:"required_without=EmployeeIDs,min=1"`
	ShiftID     uint      `json:"shift_id" validate:"required,min=1"`
	StartDate   time.Time `json:"start_date" validate:"required,notzerotime"`
	EndDate     time.Time `json:"end_date" validate:"omitempty,notzerotime"`
	EmployeeIDs []uint    `json:"employee_ids" validate:"required_without=TeamID,min=1"`
}
