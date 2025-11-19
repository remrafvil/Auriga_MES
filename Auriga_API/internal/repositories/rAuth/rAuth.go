package rAuth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Repository interface {
	SyncUser(userInfo *rModels.AuthentikUserInfo) (*rModels.MrEmployee, error)
	FindCurrentUserInfo(userID string, ctx context.Context) (*rModels.MrEmployee, error)
	FindEmployeeByEmail(email string, ctx context.Context) (*rModels.MrEmployee, error)
	FindEmployeeByWorkdayID(workdayID string, ctx context.Context) (*rModels.MrEmployee, error)
	FindEmployeeByAuthentikID(authentikID string, ctx context.Context) (*rModels.MrEmployee, error)
	FindEmployeeByIDN(idn string, ctx context.Context) (*rModels.MrEmployee, error)
	SyncAuthentikRoles(employeeID uint, organization map[string]interface{}) error
	GetEmployeeAuthentikRoles(employeeID uint, ctx context.Context) ([]rModels.MrEmployeeAuthentikRole, error)
	EnsureFactoryAndDepartment(factoryName, departmentName string) (uint, error)
	GetOrCreateFactory(factoryName string) (*rModels.MrFactory, error)
	GetOrCreateDepartment(departmentName string) (*rModels.MrDepartment, error)
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

// SyncUser sincroniza el usuario con la tabla MrEmployee y sus roles de Authentik
func (m *repository) SyncUser(userInfo *rModels.AuthentikUserInfo) (*rModels.MrEmployee, error) {
	var employee rModels.MrEmployee

	// 1. Buscar por AuthentikID primero (campo principal)
	result := m.db.Where("authentik_id = ?", userInfo.Sub).First(&employee)

	now := time.Now()

	if result.Error == gorm.ErrRecordNotFound {
		// 2. Si no existe por AuthentikID, buscar por email
		result = m.db.Where("email = ?", userInfo.Email).First(&employee)

		if result.Error == gorm.ErrRecordNotFound {
			// Crear NUEVO empleado
			employee = rModels.MrEmployee{
				AuthentikID: userInfo.Sub,
				IDN:         extractIDNFromOrganization(userInfo.Organization),
				WorkdayID:   extractWorkdayIDFromOrganization(userInfo.Organization),
				FirstName:   extractFirstName(userInfo.Name),
				LastName:    extractLastName(userInfo.Name),
				Email:       userInfo.Email,
				Active:      true,
				External:    false,
				HireDate:    now,
				CreatedAt:   now,
				UpdatedAt:   now,
			}

			if err := m.db.Create(&employee).Error; err != nil {
				return nil, fmt.Errorf("failed to create employee: %w", err)
			}

			m.logger.Info("New employee created from Authentik",
				zap.String("authentik_id", userInfo.Sub),
				zap.String("email", employee.Email),
				zap.Uint("db_id", employee.ID),
				zap.String("idn", employee.IDN),
				zap.String("workday_id", employee.WorkdayID))

		} else if result.Error == nil {
			// Empleado existe por email pero no tiene AuthentikID - ACTUALIZAR
			updateData := map[string]interface{}{
				"authentik_id": userInfo.Sub,
				"first_name":   extractFirstName(userInfo.Name),
				"last_name":    extractLastName(userInfo.Name),
				"active":       true,
				"updated_at":   now,
			}

			if newIDN := extractIDNFromOrganization(userInfo.Organization); newIDN != "" {
				updateData["idn"] = newIDN
			}

			if newWorkdayID := extractWorkdayIDFromOrganization(userInfo.Organization); newWorkdayID != "" {
				updateData["workday_id"] = newWorkdayID
			}

			if err := m.db.Model(&employee).Updates(updateData).Error; err != nil {
				return nil, fmt.Errorf("failed to update employee with authentik_id: %w", err)
			}

			m.logger.Info("Existing employee updated with AuthentikID",
				zap.Uint("employee_id", employee.ID),
				zap.String("email", employee.Email),
				zap.String("authentik_id", userInfo.Sub),
				zap.String("idn", employee.IDN),
				zap.String("workday_id", employee.WorkdayID))
		} else {
			return nil, result.Error
		}

	} else if result.Error == nil {
		// Empleado existe por AuthentikID - ACTUALIZAR
		updateData := map[string]interface{}{
			"first_name": extractFirstName(userInfo.Name),
			"last_name":  extractLastName(userInfo.Name),
			"email":      userInfo.Email,
			"active":     true,
			"updated_at": now,
		}

		if employee.IDN == "" {
			if newIDN := extractIDNFromOrganization(userInfo.Organization); newIDN != "" {
				updateData["idn"] = newIDN
			}
		}

		if employee.WorkdayID == "" {
			if newWorkdayID := extractWorkdayIDFromOrganization(userInfo.Organization); newWorkdayID != "" {
				updateData["workday_id"] = newWorkdayID
			}
		}

		if err := m.db.Model(&employee).Updates(updateData).Error; err != nil {
			return nil, err
		}

		m.logger.Debug("Employee updated from Authentik",
			zap.Uint("employee_id", employee.ID),
			zap.String("email", employee.Email),
			zap.String("idn", employee.IDN),
			zap.String("workday_id", employee.WorkdayID))
	} else {
		return nil, result.Error
	}

	// 3. SINCRONIZAR ROLES DE AUTHENTIK
	if err := m.SyncAuthentikRoles(employee.ID, userInfo.Organization); err != nil {
		m.logger.Error("Failed to sync Authentik roles",
			zap.Uint("employee_id", employee.ID),
			zap.Error(err))
		// No retornamos error aquí para no fallar el login completo
	}

	return &employee, nil
}

// SyncAuthentikRoles sincroniza los roles de Authentik del empleado
func (m *repository) SyncAuthentikRoles(employeeID uint, organization map[string]interface{}) error {
	if organization == nil {
		return fmt.Errorf("organization data is nil")
	}

	factories, ok := organization["factories"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("no factories found in organization data")
	}

	now := time.Now()

	return m.db.Transaction(func(tx *gorm.DB) error {
		// 1. Desactivar todos los roles existentes
		if err := tx.Model(&rModels.MrEmployeeAuthentikRole{}).
			Where("employee_id = ? AND is_active = ?", employeeID, true).
			Update("is_active", false).Error; err != nil {
			return fmt.Errorf("failed to deactivate existing roles: %w", err)
		}

		var rolesToCreate []rModels.MrEmployeeAuthentikRole
		var rolesCount int

		for factoryName, factoryData := range factories {
			factoryMap, ok := factoryData.(map[string]interface{})
			if !ok {
				m.logger.Warn("Invalid factory data format",
					zap.String("factory", factoryName),
					zap.Any("factory_data", factoryData))
				continue
			}

			departments, ok := factoryMap["departments"].(map[string]interface{})
			if !ok {
				m.logger.Warn("No departments found for factory",
					zap.String("factory", factoryName))
				continue
			}

			for deptName, deptData := range departments {
				deptMap, ok := deptData.(map[string]interface{})
				if !ok {
					m.logger.Warn("Invalid department data format",
						zap.String("factory", factoryName),
						zap.String("department", deptName))
					continue
				}

				rolesInterface, ok := deptMap["roles"].([]interface{})
				if !ok {
					m.logger.Warn("No roles found for department",
						zap.String("factory", factoryName),
						zap.String("department", deptName))
					continue
				}

				// Asegurar que existe la fábrica y departamento
				factoryDeptID, err := m.EnsureFactoryAndDepartment(factoryName, deptName)
				if err != nil {
					m.logger.Error("Failed to ensure factory and department",
						zap.String("factory", factoryName),
						zap.String("department", deptName),
						zap.Error(err))
					continue // Continuar con el siguiente departamento
				}

				// También crear la asignación empleado ↔ factory-department
				if err := m.ensureEmployeeFactoryDepartment(employeeID, factoryDeptID); err != nil {
					m.logger.Warn("Failed to ensure employee factory-department assignment",
						zap.Uint("employee_id", employeeID),
						zap.Uint("factory_dept_id", factoryDeptID),
						zap.Error(err))
					// No fallar la transacción por esto
				}

				for _, roleInterface := range rolesInterface {
					role, ok := roleInterface.(string)
					if !ok || role == "" {
						continue
					}

					// Verificar si el rol ya existe
					var existingRole rModels.MrEmployeeAuthentikRole
					result := tx.Where("employee_id = ? AND factory = ? AND department = ? AND role = ?",
						employeeID, factoryName, deptName, role).First(&existingRole)

					if result.Error == nil {
						// Rol existe - reactivarlo y actualizar timestamp
						if err := tx.Model(&existingRole).Updates(map[string]interface{}{
							"is_active":      true,
							"last_synced_at": now,
							"updated_at":     now,
						}).Error; err != nil {
							return fmt.Errorf("failed to reactivate role: %w", err)
						}
					} else if result.Error == gorm.ErrRecordNotFound {
						// Rol no existe - crear nuevo
						newRole := rModels.MrEmployeeAuthentikRole{
							EmployeeID:   employeeID,
							Factory:      factoryName,
							Department:   deptName,
							Role:         role,
							IsActive:     true,
							AssignedAt:   now,
							LastSyncedAt: now,
							CreatedAt:    now,
							UpdatedAt:    now,
						}
						rolesToCreate = append(rolesToCreate, newRole)
					} else {
						return fmt.Errorf("failed to check existing role: %w", result.Error)
					}

					rolesCount++
				}
			}
		}

		// Crear nuevos roles en lote
		if len(rolesToCreate) > 0 {
			if err := tx.CreateInBatches(rolesToCreate, 100).Error; err != nil {
				return fmt.Errorf("failed to create new roles: %w", err)
			}
		}

		// Limpiar roles que ya no existen (más de 30 días inactivos)
		if err := tx.Where("employee_id = ? AND is_active = ? AND last_synced_at < ?",
			employeeID, false, now.AddDate(0, 0, -30)).Delete(&rModels.MrEmployeeAuthentikRole{}).Error; err != nil {
			m.logger.Warn("Failed to clean up old roles", zap.Error(err))
			// No fallar la transacción por esto
		}

		m.logger.Info("Authentik roles synchronized successfully",
			zap.Uint("employee_id", employeeID),
			zap.Int("total_roles", rolesCount),
			zap.Int("new_roles", len(rolesToCreate)))

		return nil
	})
}

// EnsureFactoryAndDepartment crea o obtiene la fábrica, departamento y su relación
func (m *repository) EnsureFactoryAndDepartment(factoryName, departmentName string) (uint, error) {
	// 1. Crear o obtener la fábrica
	factory, err := m.GetOrCreateFactory(factoryName)
	if err != nil {
		return 0, fmt.Errorf("failed to ensure factory: %w", err)
	}

	// 2. Crear o obtener el departamento
	department, err := m.GetOrCreateDepartment(departmentName)
	if err != nil {
		return 0, fmt.Errorf("failed to ensure department: %w", err)
	}

	// 3. Crear o obtener la relación factory-department
	factoryDept, err := m.getOrCreateFactoryDepartment(factory.ID, department.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to ensure factory-department relation: %w", err)
	}

	return factoryDept.ID, nil
}

// GetOrCreateFactory busca o crea una fábrica, mapeando con assets existentes por código
func (m *repository) GetOrCreateFactory(factoryName string) (*rModels.MrFactory, error) {
	var factory rModels.MrFactory

	// Buscar por nombre (case insensitive)
	err := m.db.Where("LOWER(name) = LOWER(?)", factoryName).First(&factory).Error

	if err == gorm.ErrRecordNotFound {
		// Buscar asset por código que coincida con la fábrica
		asset, err := m.findAssetByFactoryCode(factoryName)
		if err != nil {
			return nil, fmt.Errorf("failed to find asset for factory %s: %w", factoryName, err)
		}

		// Crear nueva fábrica con el asset_id correspondiente
		factory = rModels.MrFactory{
			AssetID:   asset.ID,
			AssetCode: asset.Code,
			Name:      factoryName,
			Active:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := m.db.Create(&factory).Error; err != nil {
			return nil, fmt.Errorf("failed to create factory %s: %w", factoryName, err)
		}

		m.logger.Info("New factory created automatically with asset mapping",
			zap.String("factory_name", factoryName),
			zap.Uint("factory_id", factory.ID),
			zap.Uint("asset_id", asset.ID),
			zap.String("asset_code", asset.Code))
	} else if err != nil {
		return nil, fmt.Errorf("failed to get factory %s: %w", factoryName, err)
	}

	return &factory, nil
}

// findAssetByFactoryCode busca el asset que corresponde a una fábrica por su código
func (m *repository) findAssetByFactoryCode(factoryName string) (*rModels.MrAsset, error) {
	var asset rModels.MrAsset

	// Mapeo de fábricas a códigos de asset
	factoryToAssetMap := map[string]string{
		"CXC": "CXC", // Chile -> CXC
		"CXM": "CXM", // Mexico -> CXM
		"CXB": "CXB", // Brasil
		"CXD": "CXD", // Alemania
		"CXE": "CXE", // Alcala
		"CXF": "CXF", // Beaucouzé
		"EXT": "EXT", // Rusia
		"FPC": "FPC", // Chantrans
		"FPL": "FPL", // Beaucouzé
		"FSP": "FSP", // Roye
		"MNT": "MNT", // Sumirago
		"RTP": "RTP", // Rio Tinto
		"ITC": "ITC", // Alcala
	}

	// Obtener el código del asset para esta fábrica
	assetCode, exists := factoryToAssetMap[factoryName]
	if !exists {
		// Si no hay mapeo específico, usar el nombre de la fábrica como código
		assetCode = factoryName
	}

	// Buscar el asset por código
	err := m.db.Where("code = ?", assetCode).First(&asset).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("no asset found with code %s for factory %s", assetCode, factoryName)
		}
		return nil, fmt.Errorf("failed to find asset with code %s: %w", assetCode, err)
	}

	m.logger.Debug("Asset found for factory",
		zap.String("factory_name", factoryName),
		zap.String("asset_code", asset.Code),
		zap.Uint("asset_id", asset.ID))

	return &asset, nil
}

// GetOrCreateDepartment busca o crea un departamento
func (m *repository) GetOrCreateDepartment(departmentName string) (*rModels.MrDepartment, error) {
	var department rModels.MrDepartment

	// Buscar por nombre (case insensitive)
	err := m.db.Where("LOWER(name) = LOWER(?)", departmentName).First(&department).Error

	if err == gorm.ErrRecordNotFound {
		// Crear nuevo departamento
		department = rModels.MrDepartment{
			Name:        departmentName,
			Code:        generateDepartmentCode(departmentName),
			Description: fmt.Sprintf("Departamento %s - Creado automáticamente", departmentName),
			Active:      true,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		if err := m.db.Create(&department).Error; err != nil {
			return nil, fmt.Errorf("failed to create department %s: %w", departmentName, err)
		}

		m.logger.Info("New department created automatically",
			zap.String("department_name", departmentName),
			zap.Uint("department_id", department.ID))
	} else if err != nil {
		return nil, fmt.Errorf("failed to get department %s: %w", departmentName, err)
	}

	return &department, nil
}

// getOrCreateFactoryDepartment busca o crea la relación factory-department
func (m *repository) getOrCreateFactoryDepartment(factoryID, departmentID uint) (*rModels.MrFactoryDepartment, error) {
	var factoryDept rModels.MrFactoryDepartment

	err := m.db.Where("factory_id = ? AND department_id = ?", factoryID, departmentID).First(&factoryDept).Error

	if err == gorm.ErrRecordNotFound {
		// Crear nueva relación
		factoryDept = rModels.MrFactoryDepartment{
			FactoryID:    factoryID,
			DepartmentID: departmentID,
			Active:       true,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := m.db.Create(&factoryDept).Error; err != nil {
			return nil, fmt.Errorf("failed to create factory-department relation: %w", err)
		}

		m.logger.Info("New factory-department relation created",
			zap.Uint("factory_id", factoryID),
			zap.Uint("department_id", departmentID),
			zap.Uint("relation_id", factoryDept.ID))
	} else if err != nil {
		return nil, fmt.Errorf("failed to get factory-department relation: %w", err)
	}

	return &factoryDept, nil
}

// ensureEmployeeFactoryDepartment crea la asignación empleado ↔ factory-department
func (m *repository) ensureEmployeeFactoryDepartment(employeeID, factoryDeptID uint) error {
	var assignment rModels.MrEmployeeFactoryDepartment

	err := m.db.Where("employee_id = ? AND factory_department_id = ?", employeeID, factoryDeptID).
		First(&assignment).Error

	if err == gorm.ErrRecordNotFound {
		// Crear nueva asignación
		assignment = rModels.MrEmployeeFactoryDepartment{
			EmployeeID:          employeeID,
			FactoryDepartmentID: factoryDeptID,
			StartDate:           time.Now(),
			IsPrimaryAssignment: false,
			Active:              true,
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		}

		if err := m.db.Create(&assignment).Error; err != nil {
			return fmt.Errorf("failed to create employee factory-department assignment: %w", err)
		}

		m.logger.Debug("Employee factory-department assignment created",
			zap.Uint("employee_id", employeeID),
			zap.Uint("factory_dept_id", factoryDeptID))
	} else if err != nil {
		return fmt.Errorf("failed to get employee factory-department assignment: %w", err)
	}

	return nil
}

// FindCurrentUserInfo busca empleado por AuthentikID (primary) o fallback a email
func (m *repository) FindCurrentUserInfo(userID string, ctx context.Context) (*rModels.MrEmployee, error) {
	var employee rModels.MrEmployee

	// Buscar por AuthentikID primero
	err := m.db.Preload("FactoryDepartmentAssignments.FactoryDepartment.Factory").
		Preload("FactoryDepartmentAssignments.FactoryDepartment.Department").
		Preload("AuthentikRoles", "is_active = ?", true).
		Where("authentik_id = ?", userID).First(&employee).Error

	if err == gorm.ErrRecordNotFound {
		// Fallback: buscar por email (para compatibilidad)
		m.logger.Debug("Employee not found by authentik_id, trying email",
			zap.String("authentik_id", userID))
		err = m.db.Preload("FactoryDepartmentAssignments.FactoryDepartment.Factory").
			Preload("FactoryDepartmentAssignments.FactoryDepartment.Department").
			Preload("AuthentikRoles", "is_active = ?", true).
			Where("email = ?", userID).First(&employee).Error
	}

	if err != nil {
		m.logger.Warn("Employee not found in database",
			zap.String("user_id", userID),
			zap.Error(err))
		return nil, err
	}

	return &employee, nil
}

// FindEmployeeByAuthentikID busca empleado específicamente por Authentik ID
func (m *repository) FindEmployeeByAuthentikID(authentikID string, ctx context.Context) (*rModels.MrEmployee, error) {
	var employee rModels.MrEmployee
	if err := m.db.Preload("FactoryDepartmentAssignments.FactoryDepartment.Factory").
		Preload("FactoryDepartmentAssignments.FactoryDepartment.Department").
		Preload("AuthentikRoles", "is_active = ?", true).
		Where("authentik_id = ?", authentikID).First(&employee).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

// FindEmployeeByIDN busca empleado por IDN (id_card)
func (m *repository) FindEmployeeByIDN(idn string, ctx context.Context) (*rModels.MrEmployee, error) {
	var employee rModels.MrEmployee
	if err := m.db.Preload("FactoryDepartmentAssignments.FactoryDepartment.Factory").
		Preload("FactoryDepartmentAssignments.FactoryDepartment.Department").
		Preload("AuthentikRoles", "is_active = ?", true).
		Where("idn = ?", idn).First(&employee).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

// FindEmployeeByEmail busca empleado por email
func (m *repository) FindEmployeeByEmail(email string, ctx context.Context) (*rModels.MrEmployee, error) {
	var employee rModels.MrEmployee
	if err := m.db.Preload("FactoryDepartmentAssignments.FactoryDepartment.Factory").
		Preload("FactoryDepartmentAssignments.FactoryDepartment.Department").
		Preload("AuthentikRoles", "is_active = ?", true).
		First(&employee, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

// FindEmployeeByWorkdayID busca empleado por Workday ID
func (m *repository) FindEmployeeByWorkdayID(workdayID string, ctx context.Context) (*rModels.MrEmployee, error) {
	var employee rModels.MrEmployee
	if err := m.db.Preload("FactoryDepartmentAssignments.FactoryDepartment.Factory").
		Preload("FactoryDepartmentAssignments.FactoryDepartment.Department").
		Preload("AuthentikRoles", "is_active = ?", true).
		First(&employee, "workday_id = ?", workdayID).Error; err != nil {
		return nil, err
	}
	return &employee, nil
}

// GetEmployeeAuthentikRoles obtiene todos los roles de Authentik de un empleado
func (m *repository) GetEmployeeAuthentikRoles(employeeID uint, ctx context.Context) ([]rModels.MrEmployeeAuthentikRole, error) {
	var roles []rModels.MrEmployeeAuthentikRole
	if err := m.db.Where("employee_id = ? AND is_active = ?", employeeID, true).
		Order("factory, department, role").
		Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// Helper functions
func extractFirstName(fullName string) string {
	names := splitName(fullName)
	if len(names) > 0 {
		return names[0]
	}
	return fullName
}

func extractLastName(fullName string) string {
	names := splitName(fullName)
	if len(names) > 1 {
		return names[len(names)-1]
	}
	return ""
}

func splitName(fullName string) []string {
	var names []string
	last := 0
	for i, c := range fullName {
		if c == ' ' {
			if i > last {
				names = append(names, fullName[last:i])
			}
			last = i + 1
		}
	}
	if last < len(fullName) {
		names = append(names, fullName[last:])
	}
	return names
}

func extractIDNFromOrganization(organization map[string]interface{}) string {
	if organization != nil {
		if idn, ok := organization["idn"].(string); ok && idn != "" {
			return idn
		}
	}
	return ""
}

func extractWorkdayIDFromOrganization(organization map[string]interface{}) string {
	if organization != nil {
		if workdayID, ok := organization["workday_id"].(string); ok && workdayID != "" {
			return workdayID
		}
	}
	return ""
}

func generateDepartmentCode(departmentName string) string {
	// Ejemplo: "Quality" -> "QUAL", "Maintenance" -> "MAINT"
	if len(departmentName) <= 4 {
		return strings.ToUpper(departmentName)
	}
	return strings.ToUpper(departmentName[:4])
}
