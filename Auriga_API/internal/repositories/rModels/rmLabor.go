package rModels

import (
	"time"
)

// MrFactory - Modelo de fábrica con relación explícita
type MrFactory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	AssetID   uint      `gorm:"not null;uniqueIndex" json:"asset_id"` // Foreign key a MrAsset
	AssetCode string    `gorm:"size:50;uniqueIndex" json:"asset_code"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	City      string    `gorm:"size:100" json:"city"`
	Address   string    `gorm:"type:text" json:"address"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Active    bool      `gorm:"default:true" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// RELACIÓN: Una MrFactory pertenece a un MrAsset
	Asset MrAsset `gorm:"foreignKey:AssetID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"asset"`

	// RELACIÓN: Una MrFactory tiene muchos MrFactoryDepartment
	FactoryDepartments []MrFactoryDepartment `gorm:"foreignKey:FactoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"factory_departments,omitempty"`

	// RELACIÓN: Una MrFactory tiene muchos MrTeam
	Teams []MrTeam `gorm:"foreignKey:FactoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teams,omitempty"`
}

// MrDepartment - Departamento
type MrDepartment struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Code        string    `gorm:"size:50;unique" json:"code"`
	Description string    `gorm:"type:text" json:"description"`
	Active      bool      `gorm:"default:true" json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// RELACIÓN: Un MrDepartment tiene muchos MrFactoryDepartment
	FactoryDepartments []MrFactoryDepartment `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"factory_departments,omitempty"`

	// RELACIÓN: Un MrDepartment tiene muchos MrTeam
	Teams []MrTeam `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"teams,omitempty"`
}

// MrFactoryDepartment - Tabla intermedia para relación muchos a muchos
type MrFactoryDepartment struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FactoryID    uint      `gorm:"not null;index:idx_mr_factory_department_ids" json:"factory_id"`
	DepartmentID uint      `gorm:"not null;index:idx_mr_factory_department_ids" json:"department_id"`
	Active       bool      `gorm:"default:true" json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// RELACIONES existentes
	Factory    MrFactory    `gorm:"foreignKey:FactoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"factory"`
	Department MrDepartment `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"department"`

	// NUEVA RELACIÓN: Una combinación Factory-Department puede tener muchos empleados
	EmployeeAssignments []MrEmployeeFactoryDepartment `gorm:"foreignKey:FactoryDepartmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"employee_assignments,omitempty"`
}

// MrEmployee - Empleado con relación a Factory-Department
type MrEmployee struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	IDN             string    `gorm:"size:20" json:"idn"` // id_card (puede ser NULL para nuevos usuarios)
	WorkdayID       string    `gorm:"size:50;uniqueIndex" json:"workday_id"`
	AccessControlID string    `gorm:"size:50" json:"access_control_id"`
	AuthentikID     string    `gorm:"size:100;uniqueIndex" json:"authentik_id"` // ID de Authentik
	FirstName       string    `gorm:"size:255;not null" json:"first_name"`
	LastName        string    `gorm:"size:255;not null" json:"last_name"`
	Email           string    `gorm:"size:255;uniqueIndex" json:"email"`
	Phone           string    `gorm:"size:20" json:"phone"`
	HireDate        time.Time `json:"hire_date"`
	ContractType    string    `gorm:"size:50" json:"contract_type"`
	Salary          float64   `gorm:"type:decimal(10,2)" json:"salary"`
	Active          bool      `gorm:"default:true" json:"active"`
	External        bool      `gorm:"default:false" json:"external"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// RELACIÓN: Un empleado puede tener múltiples asignaciones Factory-Department
	FactoryDepartmentAssignments []MrEmployeeFactoryDepartment `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"factory_department_assignments,omitempty"`

	// RELACIÓN: Asignaciones directas de roles desde Authentik
	AuthentikRoles []MrEmployeeAuthentikRole `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"authentik_roles,omitempty"`

	// Otras relaciones (mantenidas)
	Attendances      []MrAttendance      `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"attendances,omitempty"`
	TeamMemberships  []MrTeamMember      `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"team_memberships,omitempty"`
	ShiftAssignments []MrShiftAssignment `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shift_assignments,omitempty"`
}

// MODELO: Roles de Authentik por fábrica y departamento
type MrEmployeeAuthentikRole struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	EmployeeID   uint      `gorm:"not null;index:idx_employee_factory_dept_role" json:"employee_id"`
	Factory      string    `gorm:"size:50;not null;index:idx_employee_factory_dept_role" json:"factory"`
	Department   string    `gorm:"size:100;not null;index:idx_employee_factory_dept_role" json:"department"`
	Role         string    `gorm:"size:100;not null;index:idx_employee_factory_dept_role" json:"role"`
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	AssignedAt   time.Time `json:"assigned_at"`
	LastSyncedAt time.Time `json:"last_synced_at"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// RELACIÓN
	Employee MrEmployee `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"employee"`
}

// MrEmployeeFactoryDepartment - Tabla intermedia para empleado ↔ (fábrica-departamento)
type MrEmployeeFactoryDepartment struct {
	ID                  uint `gorm:"primaryKey" json:"id"`
	EmployeeID          uint `gorm:"not null;index:idx_employee_factory_dept" json:"employee_id"`
	FactoryDepartmentID uint `gorm:"not null;index:idx_employee_factory_dept" json:"factory_department_id"`

	// Campos adicionales específicos de la asignación empleado
	StartDate           time.Time `json:"start_date"`
	EndDate             time.Time `json:"end_date"` // Para historial
	IsPrimaryAssignment bool      `gorm:"default:false" json:"is_primary_assignment"`
	Active              bool      `gorm:"default:true" json:"active"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// RELACIONES
	Employee          MrEmployee          `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"employee"`
	FactoryDepartment MrFactoryDepartment `gorm:"foreignKey:FactoryDepartmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"factory_department"`
}

// MrShift - Turno
type MrShift struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Description string    `gorm:"type:text" json:"description"`
	Tolerance   int       `gorm:"default:15" json:"tolerance"`
	Active      bool      `gorm:"default:true" json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// RELACIÓN: Un MrShift tiene muchas MrShiftAssignment
	ShiftAssignments []MrShiftAssignment `gorm:"foreignKey:ShiftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shift_assignments,omitempty"`
}

// MrShiftAssignment - Asignación de turnos
type MrShiftAssignment struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	EmployeeID uint      `gorm:"not null;index" json:"employee_id"`
	ShiftID    uint      `gorm:"not null;index" json:"shift_id"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Active     bool      `gorm:"default:true" json:"active"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// RELACIONES
	Employee MrEmployee `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"employee"`
	Shift    MrShift    `gorm:"foreignKey:ShiftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shift"`
}

// MrTeam - Equipo
type MrTeam struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FactoryID    uint      `gorm:"not null;index" json:"factory_id"`
	DepartmentID uint      `gorm:"not null;index" json:"department_id"`
	Name         string    `gorm:"size:255;not null" json:"name"`
	Description  string    `gorm:"type:text" json:"description"`
	LeaderID     uint      `json:"leader_id"`
	Active       bool      `gorm:"default:true" json:"active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// RELACIONES
	Factory    MrFactory    `gorm:"foreignKey:FactoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"factory"`
	Department MrDepartment `gorm:"foreignKey:DepartmentID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"department"`
	Leader     MrEmployee   `gorm:"foreignKey:LeaderID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"leader"`

	// RELACIÓN: Un MrTeam tiene muchos MrTeamMember
	TeamMembers []MrTeamMember `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"team_members,omitempty"`
}

// MrTeamMember - Miembro de equipo
type MrTeamMember struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TeamID     uint      `gorm:"not null;index" json:"team_id"`
	EmployeeID uint      `gorm:"not null;index" json:"employee_id"`
	Role       string    `gorm:"size:100" json:"role"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	Active     bool      `gorm:"default:true" json:"active"`
	CreatedAt  time.Time `json:"created_at"`

	// RELACIONES
	Team     MrTeam     `gorm:"foreignKey:TeamID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"team"`
	Employee MrEmployee `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"employee"`
}

// MrAttendance - Asistencia
type MrAttendance struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	EmployeeID   uint       `gorm:"not null;index:idx_mr_attendance_employee_date" json:"employee_id"`
	Date         time.Time  `gorm:"index:idx_mr_attendance_employee_date" json:"date"`
	ClockIn      time.Time  `json:"clock_in"`
	ClockOut     *time.Time `json:"clock_out"`
	ShiftID      uint       `json:"shift_id"`
	Status       string     `gorm:"size:50" json:"status"`
	HoursWorked  float64    `gorm:"type:decimal(4,2)" json:"hours_worked"`
	Observations string     `gorm:"type:text" json:"observations"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// RELACIONES
	Employee MrEmployee `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"employee"`
	Shift    MrShift    `gorm:"foreignKey:ShiftID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"shift"`
}
