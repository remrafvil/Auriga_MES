package rModels

import (
	"time"

	"gorm.io/gorm"
)

// MrMaintenancePlan - Plan de mantenimiento
type MrMaintenancePlan struct {
	ID              uint            `gorm:"primaryKey" json:"id"`
	Name            string          `gorm:"size:255;not null" json:"name"`
	Description     string          `gorm:"type:text" json:"description"`
	AssetID         uint            `gorm:"not null;index:idx_maintenanceplan_asset_active" json:"asset_id"` // Activo al que aplica el plan
	MaintenanceType MaintenanceType `gorm:"type:varchar(20);not null" json:"maintenance_type"`
	Priority        PriorityLevel   `gorm:"type:varchar(20);not null" json:"priority"`
	FrequencyType   FrequencyType   `gorm:"type:varchar(20);not null" json:"frequency_type"`
	FrequencyValue  int             `gorm:"not null" json:"frequency_value"`
	DurationMinutes int             `gorm:"not null" json:"duration_minutes"`
	Active          bool            `gorm:"default:true;index:idx_maintenanceplan_asset_active" json:"active"`
	LastExecuted    *time.Time      `json:"last_executed"`
	NextScheduled   time.Time       `json:"next_scheduled"`

	// Campos para mantenimiento predictivo
	ConditionParameter *string  `gorm:"size:100" json:"condition_parameter"` // ej: "temperature", "vibration"
	ThresholdValue     *float64 `gorm:"type:decimal(10,4)" json:"threshold_value"`
	SensorAssetID      *uint    `gorm:"index" json:"sensor_asset_id"` // Activo sensor relacionado

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// RELACIONES
	Asset              MrAsset                  `gorm:"foreignKey:AssetID" json:"asset"`
	SensorAsset        *MrAsset                 `gorm:"foreignKey:SensorAssetID" json:"sensor_asset,omitempty"`
	Procedures         []MrMaintenanceProcedure `gorm:"foreignKey:MaintenancePlanID" json:"procedures"`
	RequiredSpareParts []MrMaintenanceSparePart `gorm:"foreignKey:MaintenancePlanID" json:"required_spare_parts"`
	WorkOrders         []MrMaintWorkOrder       `gorm:"foreignKey:MaintenancePlanID" json:"work_orders"`
}

// MrMaintenanceProcedure - Procedimiento de mantenimiento
type MrMaintenanceProcedure struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	MaintenancePlanID uint      `gorm:"not null;index" json:"maintenance_plan_id"`
	StepNumber        int       `gorm:"not null" json:"step_number"`
	Title             string    `gorm:"size:255;not null" json:"title"`
	Description       string    `gorm:"type:text" json:"description"`
	ExpectedDuration  int       `gorm:"not null" json:"expected_duration"` // en minutos
	SafetyNotes       string    `gorm:"type:text" json:"safety_notes"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// RELACIONES
	MaintenancePlan MrMaintenancePlan `gorm:"foreignKey:MaintenancePlanID" json:"maintenance_plan"`
}

// MrMaintenanceSparePart - Repuestos requeridos para mantenimiento (referencian PRODUCTOS)
type MrMaintenanceSparePart struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	MaintenancePlanID uint      `gorm:"not null;index" json:"maintenance_plan_id"`
	ProductID         uint      `gorm:"not null;index" json:"product_id"` // Referencia a MrProduct
	QuantityRequired  float64   `gorm:"type:decimal(10,3);not null" json:"quantity_required"`
	Unit              string    `gorm:"size:50;not null" json:"unit"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// RELACIONES
	MaintenancePlan MrMaintenancePlan `gorm:"foreignKey:MaintenancePlanID" json:"maintenance_plan"`
	Product         MrProduct         `gorm:"foreignKey:ProductID" json:"product"`
}

// MrMaintWorkOrder - Orden de trabajo de mantenimiento
type MrMaintWorkOrder struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	MaintenancePlanID *uint           `gorm:"index" json:"maintenance_plan_id"` // Puede ser nulo para OT correctivas no planificadas
	AssetID           uint            `gorm:"not null;index:idx_maintworkorder_asset_status" json:"asset_id"`
	WorkOrderType     WorkOrderType   `gorm:"type:varchar(20);not null" json:"work_order_type"`
	Priority          PriorityLevel   `gorm:"type:varchar(20);not null" json:"priority"`
	Status            WorkOrderStatus `gorm:"type:varchar(20);not null;index:idx_maintworkorder_asset_status" json:"status"`
	Title             string          `gorm:"size:255;not null" json:"title"`
	Description       string          `gorm:"type:text" json:"description"`
	ScheduledDate     time.Time       `json:"scheduled_date"`
	StartDate         *time.Time      `json:"start_date"`
	EndDate           *time.Time      `json:"end_date"`
	EstimatedHours    float64         `gorm:"type:decimal(4,2)" json:"estimated_hours"`
	ActualHours       float64         `gorm:"type:decimal(4,2)" json:"actual_hours"`
	Cost              float64         `gorm:"type:decimal(10,2)" json:"cost"`
	CreatedBy         uint            `gorm:"not null" json:"created_by"` // EmployeeID del creador
	AssignedTeamID    *uint           `gorm:"index" json:"assigned_team_id"`

	// Campos adicionales para trazabilidad
	CompletionNotes *string `gorm:"type:text" json:"completion_notes"` // Observaciones al cerrar
	QualityCheck    *bool   `gorm:"default:null" json:"quality_check"` // Verificación de calidad

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// RELACIONES
	MaintenancePlan *MrMaintenancePlan           `gorm:"foreignKey:MaintenancePlanID" json:"maintenance_plan,omitempty"`
	Asset           MrAsset                      `gorm:"foreignKey:AssetID" json:"asset"`
	Creator         MrEmployee                   `gorm:"foreignKey:CreatedBy" json:"creator"`
	AssignedTeam    *MrTeam                      `gorm:"foreignKey:AssignedTeamID" json:"assigned_team,omitempty"`
	AssignedMembers []MrMaintWorkOrderAssignment `gorm:"foreignKey:WorkOrderID" json:"assigned_members"`
	Tasks           []MrMaintWorkOrderTask       `gorm:"foreignKey:WorkOrderID" json:"tasks"`
	UsedSpareParts  []MrMaintWorkOrderSparePart  `gorm:"foreignKey:WorkOrderID" json:"used_spare_parts"`
	Documents       []MrDocuments                `gorm:"many2many:mr_maint_work_order_documents;joinForeignKey:WorkOrderID;joinReferences:DocumentID" json:"documents"`
}

// MrMaintWorkOrderAssignment - Asignación de empleados a orden de trabajo de mantenimiento
type MrMaintWorkOrderAssignment struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	WorkOrderID uint       `gorm:"not null;index" json:"work_order_id"`
	EmployeeID  uint       `gorm:"not null;index" json:"employee_id"`
	Role        string     `gorm:"size:100;not null" json:"role"` // Líder, Técnico, Ayudante, etc.
	StartTime   *time.Time `json:"start_time"`
	EndTime     *time.Time `json:"end_time"`
	HoursWorked float64    `gorm:"type:decimal(4,2)" json:"hours_worked"`
	CreatedAt   time.Time  `json:"created_at"`

	// RELACIONES
	WorkOrder MrMaintWorkOrder `gorm:"foreignKey:WorkOrderID" json:"work_order"`
	Employee  MrEmployee       `gorm:"foreignKey:EmployeeID" json:"employee"`
}

// MrMaintWorkOrderTask - Tareas específicas de la orden de trabajo de mantenimiento
type MrMaintWorkOrderTask struct {
	ID          uint                `gorm:"primaryKey" json:"id"`
	WorkOrderID uint                `gorm:"not null;index:idx_maintworkordertask_workorder_status" json:"work_order_id"`
	TaskNumber  int                 `gorm:"not null" json:"task_number"`
	Description string              `gorm:"type:text;not null" json:"description"`
	Status      WorkOrderTaskStatus `gorm:"type:varchar(20);not null;index:idx_maintworkordertask_workorder_status" json:"status"`
	CompletedAt *time.Time          `json:"completed_at"`

	// Campo adicional para asignación específica
	AssignedTo *uint `gorm:"index" json:"assigned_to"` // Empleado específico asignado

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// RELACIONES
	WorkOrder MrMaintWorkOrder `gorm:"foreignKey:WorkOrderID" json:"work_order"`
	Assignee  *MrEmployee      `gorm:"foreignKey:AssignedTo" json:"assignee,omitempty"`
}

// MrMaintWorkOrderSparePart - Repuestos utilizados en orden de trabajo de mantenimiento (referencian ACTIVOS del almacén)
type MrMaintWorkOrderSparePart struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	WorkOrderID      uint      `gorm:"not null;index" json:"work_order_id"`
	ProductID        uint      `gorm:"not null;index" json:"product_id"`
	StockID          *uint     `gorm:"index" json:"stock_id"` // Referencia al stock específico
	SparePartAssetID *uint     `gorm:"index" json:"spare_part_asset_id"`
	Quantity         float64   `gorm:"type:decimal(15,6);not null" json:"quantity"`
	UnitCost         float64   `gorm:"type:decimal(15,6)" json:"unit_cost"`
	TotalCost        float64   `gorm:"type:decimal(15,2)" json:"total_cost"`
	AssetMovementID  *uint     `gorm:"index" json:"asset_movement_id"` // Movimiento asociado (antes StockMovementID)
	CreatedAt        time.Time `json:"created_at"`

	// RELACIONES
	WorkOrder      MrMaintWorkOrder         `gorm:"foreignKey:WorkOrderID" json:"work_order"`
	Product        MrProduct                `gorm:"foreignKey:ProductID" json:"product"`
	Stock          *MrSparePartStock        `gorm:"foreignKey:StockID" json:"stock,omitempty"`
	SparePartAsset *MrAsset                 `gorm:"foreignKey:SparePartAssetID" json:"spare_part_asset,omitempty"`
	AssetMovement  *MrAssetRegisterMovement `gorm:"foreignKey:AssetMovementID" json:"asset_movement,omitempty"`
}

// MrSparePartStock - Control de stock para productos no discretos o control avanzado
type MrSparePartStock struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	WarehouseAssetID uint       `gorm:"not null;index:idx_sparepartstock_warehouse_product" json:"warehouse_asset_id"`
	ProductID        uint       `gorm:"not null;index:idx_sparepartstock_warehouse_product" json:"product_id"`
	Quantity         float64    `gorm:"type:decimal(15,6);not null" json:"quantity"`
	MinStock         float64    `gorm:"type:decimal(15,6)" json:"min_stock"`
	MaxStock         float64    `gorm:"type:decimal(15,6)" json:"max_stock"`
	BatchNumber      string     `gorm:"size:100" json:"batch_number"`          // Para trazabilidad
	ExpiryDate       *time.Time `json:"expiry_date"`                           // Para productos perecederos
	Location         string     `gorm:"size:100" json:"location"`              // Ubicación específica
	UnitCost         float64    `gorm:"type:decimal(15,6)" json:"unit_cost"`   // Costo unitario promedio
	TotalValue       float64    `gorm:"type:decimal(15,2)" json:"total_value"` // Valor total
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`

	// RELACIONES
	Product        MrProduct                 `gorm:"foreignKey:ProductID" json:"product"`
	WarehouseAsset MrAsset                   `gorm:"foreignKey:WarehouseAssetID" json:"warehouse_asset"`
	AssetMovements []MrAssetRegisterMovement `gorm:"foreignKey:StockID" json:"asset_movements"`
}

// MrAssetRegisterMovement - Registro de movimientos de inventario de activos
type MrAssetRegisterMovement struct {
	ID            uint              `gorm:"primaryKey" json:"id"`
	StockID       uint              `gorm:"not null;index:idx_assetregistermovement_stock_date" json:"stock_id"`
	MovementType  AssetMovementType `gorm:"type:varchar(20);not null" json:"movement_type"`
	MovementDate  time.Time         `gorm:"index:idx_assetregistermovement_stock_date" json:"movement_date"`
	Quantity      float64           `gorm:"type:decimal(15,6);not null" json:"quantity"`
	UnitCost      float64           `gorm:"type:decimal(15,6)" json:"unit_cost"`
	TotalCost     float64           `gorm:"type:decimal(15,2)" json:"total_cost"`
	PreviousStock float64           `gorm:"type:decimal(15,6)" json:"previous_stock"`
	NewStock      float64           `gorm:"type:decimal(15,6)" json:"new_stock"`

	// Referencias para trazabilidad
	WorkOrderID     *uint `gorm:"index" json:"work_order_id"`     // Si es consumo por OT
	PurchaseOrderID *uint `gorm:"index" json:"purchase_order_id"` // Si es entrada por compra
	RequisitionID   *uint `gorm:"index" json:"requisition_id"`    // Si es por requisición
	AssetID         *uint `gorm:"index" json:"asset_id"`          // Activo específico relacionado

	// Información del movimiento
	ReferenceNumber string    `gorm:"size:100" json:"reference_number"` // Nº documento
	Reason          string    `gorm:"type:text" json:"reason"`          // Motivo del movimiento
	CreatedBy       uint      `gorm:"not null" json:"created_by"`       // Empleado que registra
	CreatedAt       time.Time `json:"created_at"`

	// RELACIONES
	Stock         MrSparePartStock  `gorm:"foreignKey:StockID" json:"stock"`
	WorkOrder     *MrMaintWorkOrder `gorm:"foreignKey:WorkOrderID" json:"work_order,omitempty"`
	PurchaseOrder *MrPurchaseOrder  `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order,omitempty"`
	Asset         *MrAsset          `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	Creator       MrEmployee        `gorm:"foreignKey:CreatedBy" json:"creator"`
}

// MrPurchaseOrder - Orden de compra para repuestos
type MrPurchaseOrder struct {
	ID               uint                `gorm:"primaryKey" json:"id"`
	WarehouseAssetID uint                `gorm:"not null;index" json:"warehouse_asset_id"`
	OrderNumber      string              `gorm:"size:100;not null;unique" json:"order_number"`
	Supplier         string              `gorm:"size:255;not null" json:"supplier"`
	Status           PurchaseOrderStatus `gorm:"type:varchar(20);not null" json:"status"`
	OrderDate        time.Time           `json:"order_date"`
	ExpectedDate     time.Time           `json:"expected_date"`
	ReceivedDate     *time.Time          `json:"received_date"`
	TotalAmount      float64             `gorm:"type:decimal(15,2)" json:"total_amount"`
	CreatedBy        uint                `gorm:"not null" json:"created_by"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`

	// RELACIONES
	WarehouseAsset MrAsset                   `gorm:"foreignKey:WarehouseAssetID" json:"warehouse_asset"`
	Creator        MrEmployee                `gorm:"foreignKey:CreatedBy" json:"creator"`
	OrderItems     []MrPurchaseOrderItem     `gorm:"foreignKey:PurchaseOrderID" json:"order_items"`
	AssetMovements []MrAssetRegisterMovement `gorm:"foreignKey:PurchaseOrderID" json:"asset_movements"`
}

// MrPurchaseOrderItem - Items de orden de compra
type MrPurchaseOrderItem struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	PurchaseOrderID  uint      `gorm:"not null;index" json:"purchase_order_id"`
	ProductID        uint      `gorm:"not null;index" json:"product_id"`
	Quantity         float64   `gorm:"type:decimal(15,6);not null" json:"quantity"`
	UnitPrice        float64   `gorm:"type:decimal(15,6);not null" json:"unit_price"`
	TotalPrice       float64   `gorm:"type:decimal(15,2);not null" json:"total_price"`
	ReceivedQuantity float64   `gorm:"type:decimal(15,6)" json:"received_quantity"`
	CreatedAt        time.Time `json:"created_at"`

	// RELACIONES
	PurchaseOrder MrPurchaseOrder `gorm:"foreignKey:PurchaseOrderID" json:"purchase_order"`
	Product       MrProduct       `gorm:"foreignKey:ProductID" json:"product"`
}

// Enums y tipos personalizados
type MaintenanceType string

const (
	MaintenanceTypePreventive MaintenanceType = "preventive"
	MaintenanceTypeCorrective MaintenanceType = "corrective"
	MaintenanceTypePredictive MaintenanceType = "predictive"
)

type WorkOrderType string

const (
	WorkOrderTypePreventive WorkOrderType = "preventive"
	WorkOrderTypeCorrective WorkOrderType = "corrective"
	WorkOrderTypePredictive WorkOrderType = "predictive"
	WorkOrderTypeInspection WorkOrderType = "inspection"
)

type PriorityLevel string

const (
	PriorityLow      PriorityLevel = "low"
	PriorityMedium   PriorityLevel = "medium"
	PriorityHigh     PriorityLevel = "high"
	PriorityCritical PriorityLevel = "critical"
)

type FrequencyType string

const (
	FrequencyDays   FrequencyType = "days"
	FrequencyWeeks  FrequencyType = "weeks"
	FrequencyMonths FrequencyType = "months"
	FrequencyHours  FrequencyType = "hours"
	FrequencyUses   FrequencyType = "uses"
)

type WorkOrderStatus string

const (
	WorkOrderPending    WorkOrderStatus = "pending"
	WorkOrderScheduled  WorkOrderStatus = "scheduled"
	WorkOrderInProgress WorkOrderStatus = "in_progress"
	WorkOrderCompleted  WorkOrderStatus = "completed"
	WorkOrderCancelled  WorkOrderStatus = "cancelled"
)

type WorkOrderTaskStatus string

const (
	TaskPending    WorkOrderTaskStatus = "pending"
	TaskInProgress WorkOrderTaskStatus = "in_progress"
	TaskCompleted  WorkOrderTaskStatus = "completed"
)

type PurchaseRequisitionStatus string

const (
	RequisitionDraft    PurchaseRequisitionStatus = "draft"
	RequisitionPending  PurchaseRequisitionStatus = "pending"
	RequisitionApproved PurchaseRequisitionStatus = "approved"
	RequisitionRejected PurchaseRequisitionStatus = "rejected"
	RequisitionOrdered  PurchaseRequisitionStatus = "ordered"
)

// Enums para movimientos de activos
type AssetMovementType string

const (
	AssetMovementTypePurchase    AssetMovementType = "purchase"    // Entrada por compra
	AssetMovementTypeConsumption AssetMovementType = "consumption" // Salida por consumo
	AssetMovementTypeAdjustment  AssetMovementType = "adjustment"  // Ajuste de inventario
	AssetMovementTypeTransfer    AssetMovementType = "transfer"    // Transferencia entre almacenes
	AssetMovementTypeInitial     AssetMovementType = "initial"     // Stock inicial
	AssetMovementTypeReturn      AssetMovementType = "return"      // Devolución
)

type PurchaseOrderStatus string

const (
	PurchaseOrderDraft     PurchaseOrderStatus = "draft"
	PurchaseOrderOrdered   PurchaseOrderStatus = "ordered"
	PurchaseOrderReceived  PurchaseOrderStatus = "received"
	PurchaseOrderCancelled PurchaseOrderStatus = "cancelled"
	PurchaseOrderPartial   PurchaseOrderStatus = "partial"
)

/*
FLUJO DE TRABAJO COMPLETO:

1- PLAN:
	- Creación del Plan de Mantenimiento
	- Definición de Procedimientos Estándar
	- Especificación de Repuestos Requeridos

2- PROGRAMACIÓN:
	- Generación Automática de Órdenes de Trabajo
	- Asignación de Recursos
	- Verificación de Disponibilidad de Repuestos

3- EJECUCIÓN:
	- Inicio de la Orden de Trabajo
 	- Ejecución de Tareas Específicas
 	- Consumo de Repuestos
 	- Registro de Tiempos

4- REGISTRO:
	- Completación de la Orden de Trabajo
	- Actualización del Plan de Mantenimiento
	- Actualización de Inventario
	- Generación de Métricas e Informes
	- Aprendizaje y Mejora

JERARQUÍA DE ACTIVOS:
MrFactory (Asset)
└── MrAsset (Almacén)
    ├── MrAsset (Repuesto 1) → MrProduct (Tipo Filtro)
    ├── MrAsset (Repuesto 2) → MrProduct (Tipo Rodamiento)
    └── MrSparePartStock (Producto a granel - Aceite, Grasa, etc.)

CONTROL DE INVENTARIO:
- Productos discretos: Se gestionan como MrAsset individuales
- Productos no discretos: Se gestionan con MrSparePartStock (cantidades)
- Trazabilidad completa: MrAssetRegisterMovement registra todos los movimientos
*/
