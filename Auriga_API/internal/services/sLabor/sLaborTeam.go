package sLabor

import (
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
)

func (s *service) GetAllTeams() ([]rModels.MrTeam, error) {
	return s.repository.GetAllTeams()
}

func (s *service) GetTeamByID(id uint) (*rModels.MrTeam, error) {
	return s.repository.GetTeamByID(id)
}

func (s *service) CreateTeam(req CreateTeamRequest) (*rModels.MrTeam, error) {
	team := &rModels.MrTeam{
		FactoryID:    req.FactoryID,
		DepartmentID: req.DepartmentID,
		Name:         req.Name,
		Description:  req.Description,
		LeaderID:     req.LeaderID,
		Active:       true,
	}

	if err := s.repository.CreateTeam(team); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *service) UpdateTeam(id uint, req UpdateTeamRequest) (*rModels.MrTeam, error) {
	team, err := s.repository.GetTeamByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		team.Name = req.Name
	}
	if req.Description != "" {
		team.Description = req.Description
	}
	if req.LeaderID != 0 {
		team.LeaderID = req.LeaderID
	}
	if req.Active != nil {
		team.Active = *req.Active
	}

	if err := s.repository.UpdateTeam(team); err != nil {
		return nil, err
	}

	return team, nil
}

func (s *service) DeleteTeam(id uint) error {
	return s.repository.DeleteTeam(id)
}

func (s *service) AddTeamMember(req AddTeamMemberRequest) (*rModels.MrTeamMember, error) {
	teamMember := &rModels.MrTeamMember{
		TeamID:     req.TeamID,
		EmployeeID: req.EmployeeID,
		Role:       req.Role,
		StartDate:  req.StartDate,
		Active:     true,
	}

	if err := s.repository.AddTeamMember(teamMember); err != nil {
		return nil, err
	}

	return teamMember, nil
}

func (s *service) RemoveTeamMember(teamID, employeeID uint) error {
	return s.repository.RemoveTeamMember(teamID, employeeID)
}

func (s *service) GetTeamMembers(teamID uint) ([]rModels.MrTeamMember, error) {
	return s.repository.GetTeamMembers(teamID)
}
