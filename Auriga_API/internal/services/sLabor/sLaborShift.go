package sLabor

import (
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
)

func (s *service) GetAllShifts() ([]rModels.MrShift, error) {
	return s.repository.GetAllShifts()
}

func (s *service) GetShiftByID(id uint) (*rModels.MrShift, error) {
	return s.repository.GetShiftByID(id)
}

func (s *service) CreateShift(req CreateShiftRequest) (*rModels.MrShift, error) {
	shift := &rModels.MrShift{
		Name:        req.Name,
		StartTime:   req.StartTime,
		EndTime:     req.EndTime,
		Description: req.Description,
		Tolerance:   req.Tolerance,
		Active:      true,
	}

	if err := s.repository.CreateShift(shift); err != nil {
		return nil, err
	}

	return shift, nil
}

func (s *service) UpdateShift(id uint, req UpdateShiftRequest) (*rModels.MrShift, error) {
	shift, err := s.repository.GetShiftByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		shift.Name = req.Name
	}
	if !req.StartTime.IsZero() {
		shift.StartTime = req.StartTime
	}
	if !req.EndTime.IsZero() {
		shift.EndTime = req.EndTime
	}
	if req.Description != "" {
		shift.Description = req.Description
	}
	if req.Tolerance != 0 {
		shift.Tolerance = req.Tolerance
	}
	if req.Active != nil {
		shift.Active = *req.Active
	}

	if err := s.repository.UpdateShift(shift); err != nil {
		return nil, err
	}

	return shift, nil
}

func (s *service) DeleteShift(id uint) error {
	return s.repository.DeleteShift(id)
}
