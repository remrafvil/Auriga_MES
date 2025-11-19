package rLabor_KKK

import (
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	SaveOrUpdateEmployees(employees []rModels.Employee) error
	GetEmployeeCount() (int64, error)
	GetEmployeeByCode(code string) (*rModels.Employee, error)
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

func (r *repository) SaveOrUpdateEmployees(employees []rModels.Employee) error {
	r.logger.Info("Starting to save/update employees in database",
		zap.Int("employeesCount", len(employees)))

	// Use transaction for bulk insert/update
	err := r.db.Transaction(func(tx *gorm.DB) error {
		for _, employee := range employees {
			var existingEmployee rModels.Employee
			result := tx.Where("code = ?", employee.Code).First(&existingEmployee)

			if result.Error == gorm.ErrRecordNotFound {
				// Create new employee
				if err := tx.Create(&employee).Error; err != nil {
					r.logger.Error("Failed to create employee in database",
						zap.String("code", employee.Code),
						zap.Error(err))
					return err
				}
				r.logger.Debug("Created employee in database", zap.String("code", employee.Code))
			} else if result.Error == nil {
				// Update existing employee
				employee.ID = existingEmployee.ID
				if err := tx.Save(&employee).Error; err != nil {
					r.logger.Error("Failed to update employee in database",
						zap.String("code", employee.Code),
						zap.Error(err))
					return err
				}
				r.logger.Debug("Updated employee in database", zap.String("code", employee.Code))
			} else {
				r.logger.Error("Failed to check existing employee in database",
					zap.String("code", employee.Code),
					zap.Error(result.Error))
				return result.Error
			}
		}
		return nil
	})

	if err != nil {
		r.logger.Error("Transaction failed while saving employees", zap.Error(err))
		return err
	}

	r.logger.Info("Successfully saved/updated employees in database",
		zap.Int("totalProcessed", len(employees)))
	return nil
}

func (r *repository) GetEmployeeCount() (int64, error) {
	var count int64
	err := r.db.Model(&rModels.Employee{}).Count(&count).Error
	if err != nil {
		r.logger.Error("Failed to get employee count from database", zap.Error(err))
		return 0, err
	}
	return count, nil
}

func (r *repository) GetEmployeeByCode(code string) (*rModels.Employee, error) {
	var employee rModels.Employee
	err := r.db.Where("code = ?", code).First(&employee).Error
	if err != nil {
		return nil, err
	}
	return &employee, nil
}
