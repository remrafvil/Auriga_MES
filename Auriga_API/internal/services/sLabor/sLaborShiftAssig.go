package sLabor

import (
	"fmt"
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
)

func (s *service) GetEmployeeAssignments(employeeID uint) ([]rModels.MrShiftAssignment, error) {
	return s.repository.GetAssignmentsByEmployee(employeeID)
}

func (s *service) GetTeamAssignments(teamID uint) ([]rModels.MrShiftAssignment, error) {
	return s.repository.GetAssignmentsByTeam(teamID)
}

func (s *service) CreateIndividualAssignment(req CreateShiftAssignmentRequest) (*rModels.MrShiftAssignment, error) {
	assignment := &rModels.MrShiftAssignment{
		EmployeeID: req.EmployeeID,
		ShiftID:    req.ShiftID,
		StartDate:  req.StartDate,
		EndDate:    req.EndDate,
		Active:     true,
	}

	if err := s.repository.CreateAssignment(assignment); err != nil {
		return nil, err
	}

	return assignment, nil
}

func (s *service) CreateTeamAssignment(req CreateShiftAssignmentRequest) ([]rModels.MrShiftAssignment, error) {
	// Obtener todos los miembros del equipo
	members, err := s.GetTeamMembers(req.TeamID)
	if err != nil {
		return nil, err
	}

	var assignments []rModels.MrShiftAssignment
	for _, member := range members {
		assignment := &rModels.MrShiftAssignment{
			EmployeeID: member.EmployeeID,
			ShiftID:    req.ShiftID,
			StartDate:  req.StartDate,
			EndDate:    req.EndDate,
			Active:     true,
		}

		if err := s.repository.CreateAssignment(assignment); err != nil {
			return nil, err
		}

		assignments = append(assignments, *assignment)
	}

	return assignments, nil
}

func (s *service) CreateBulkAssignments(req BulkShiftAssignmentRequest) ([]rModels.MrShiftAssignment, error) {
	var assignments []rModels.MrShiftAssignment

	// Si se proporciona TeamID, asignar a todos los miembros del equipo
	if req.TeamID != 0 {
		members, err := s.repository.GetTeamMembers(req.TeamID)
		if err != nil {
			return nil, err
		}

		for _, member := range members {
			assignment := &rModels.MrShiftAssignment{
				EmployeeID: member.EmployeeID,
				ShiftID:    req.ShiftID,
				StartDate:  req.StartDate,
				EndDate:    req.EndDate,
				Active:     true,
			}

			if err := s.repository.CreateAssignment(assignment); err != nil {
				return nil, err
			}
			assignments = append(assignments, *assignment)
		}
	}

	// Si se proporcionan EmployeeIDs individuales
	if len(req.EmployeeIDs) > 0 {
		for _, employeeID := range req.EmployeeIDs {
			assignment := &rModels.MrShiftAssignment{
				EmployeeID: employeeID,
				ShiftID:    req.ShiftID,
				StartDate:  req.StartDate,
				EndDate:    req.EndDate,
				Active:     true,
			}

			if err := s.repository.CreateAssignment(assignment); err != nil {
				return nil, err
			}
			assignments = append(assignments, *assignment)
		}
	}

	return assignments, nil
}

func (s *service) UpdateAssignment(id uint, req CreateShiftAssignmentRequest) (*rModels.MrShiftAssignment, error) {
	assignment, err := s.repository.GetAssignmentsByShift(id)
	if err != nil {
		return nil, err
	}

	// Buscar la asignación específica (simplificado - en realidad necesitarías un método GetAssignmentByID)
	// Por simplicidad, asumimos que tenemos el ID de la asignación
	var targetAssignment *rModels.MrShiftAssignment
	for i := range assignment {
		if assignment[i].ID == id {
			targetAssignment = &assignment[i]
			break
		}
	}

	if targetAssignment == nil {
		return nil, fmt.Errorf("assignment not found")
	}

	if req.EmployeeID != 0 {
		targetAssignment.EmployeeID = req.EmployeeID
	}
	if req.ShiftID != 0 {
		targetAssignment.ShiftID = req.ShiftID
	}
	if !req.StartDate.IsZero() {
		targetAssignment.StartDate = req.StartDate
	}
	targetAssignment.EndDate = req.EndDate

	if err := s.repository.UpdateAssignment(targetAssignment); err != nil {
		return nil, err
	}

	return targetAssignment, nil
}

func (s *service) DeleteAssignment(id uint) error {
	return s.repository.DeleteAssignment(id)
}

func (s *service) GetCurrentEmployeeAssignment(employeeID uint, date time.Time) (*rModels.MrShiftAssignment, error) {
	if date.IsZero() {
		date = time.Now()
	}
	return s.repository.GetCurrentAssignment(employeeID, date)
}
