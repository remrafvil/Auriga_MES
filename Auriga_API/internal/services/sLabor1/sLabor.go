package sLabor1

import (
	"fmt"

	"github.com/remrafvil/Auriga_API/internal/repositories/rLabor_KKK"
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"github.com/remrafvil/Auriga_API/internal/repositories/rwWorkera"
	"go.uber.org/zap"
)

type Service interface {
	SyncEmployees() error
	GetEmployeeCount() (int64, error)
	GetEmployeeByCode(code string) (*rModels.Employee, error)
}

// service implementación
type service struct {
	workeraRepo rwWorkera.Repository
	laborRepo   rLabor_KKK.Repository
	logger      *zap.Logger
}

func New(workeraRepo rwWorkera.Repository, laborRepo rLabor_KKK.Repository, logger *zap.Logger) Service {
	return &service{
		workeraRepo: workeraRepo,
		laborRepo:   laborRepo,
		logger:      logger,
	}
}

func (s *service) SyncEmployees() error {
	s.logger.Info("Starting employee synchronization process")

	// 1. Obtener empleados de Workera (repositorio Workera)
	s.logger.Info("Fetching employees from Workera API")
	employees, err := s.workeraRepo.GetAllEmployees()
	if err != nil {
		s.logger.Error("Failed to get employees from Workera API", zap.Error(err))
		return fmt.Errorf("failed to get employees from Workera: %w", err)
	}

	s.logger.Info("Retrieved employees from Workera API",
		zap.Int("count", len(employees)))

	// 2. Guardar empleados en la base de datos local (repositorio Labor)
	s.logger.Info("Saving employees to local database")
	err = s.laborRepo.SaveOrUpdateEmployees(employees)
	if err != nil {
		s.logger.Error("Failed to save employees to local database", zap.Error(err))
		return fmt.Errorf("failed to save employees to database: %w", err)
	}

	s.logger.Info("Successfully synchronized employees",
		zap.Int("totalProcessed", len(employees)))
	return nil
}

func (s *service) GetEmployeeCount() (int64, error) {
	count, err := s.laborRepo.GetEmployeeCount()
	if err != nil {
		s.logger.Error("Failed to get employee count from database", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (s *service) GetEmployeeByCode(code string) (*rModels.Employee, error) {
	employee, err := s.laborRepo.GetEmployeeByCode(code)
	if err != nil {
		s.logger.Error("Failed to get employee by code",
			zap.String("code", code),
			zap.Error(err))
		return nil, err
	}
	return employee, nil
}
