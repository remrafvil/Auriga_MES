package rEvents

// lo llamaremso repositories

import (
	"context"
	"errors"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

func (r *repository) FindCategoriesWithEventTypes(ctx context.Context) ([]rModels.MrEventCategory, error) {
	var categories []rModels.MrEventCategory

	err := r.db.WithContext(ctx).
		Preload("EventTypes").
		Find(&categories).Error

	if err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *repository) FindCategoryWithEventTypesByName(ctx context.Context, name string) (*rModels.MrEventCategory, error) {
	var category rModels.MrEventCategory

	err := r.db.WithContext(ctx).
		Preload("EventTypes").
		Where("name = ?", name).
		First(&category).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Categoría no encontrada
		}
		return nil, err
	}

	return &category, nil
}

func DB_InitEventsCategory(db *gorm.DB) {
	// Create Event Categories with Types
	categories := []rModels.MrEventCategory{
		{
			Name:        "Maintenance",
			Description: "Events related to preventive and corrective maintenance",
			EventTypes: []rModels.MrEventType{
				{Name: "Preventive Maintenance", Description: "Scheduled inspection and maintenance tasks"},
				{Name: "Corrective Maintenance", Description: "Repairs after a failure occurs"},
				{Name: "Conditional Monitoring", Description: "Monitoring for early detection of potential failures"},
				{Name: "Lubrication", Description: "Routine lubrication of components"},
				{Name: "Alignment Verification", Description: "Ensuring correct alignment of machine parts"},
			},
		},
		{
			Name:        "Production",
			Description: "Events related to production operations and performance",
			EventTypes: []rModels.MrEventType{
				{Name: "Startup Operation", Description: "Starting the production line"},
				{Name: "Scheduled Stop", Description: "Planned stop of operations"},
				{Name: "Process Adjustment", Description: "Adjustments to optimize production processes"},
				{Name: "Material Delay", Description: "Delays in material supply to the production line"},
			},
		},
		{
			Name:        "Safety",
			Description: "Events related to workplace safety and incidents",
			EventTypes: []rModels.MrEventType{
				{Name: "Incident Report", Description: "Report of an incident in the workplace"},
				{Name: "Emergency Drill", Description: "Scheduled emergency training drill"},
				{Name: "Spill Containment", Description: "Management of hazardous material spills"},
				{Name: "Protective Equipment Check", Description: "Inspection of safety equipment and PPE"},
			},
		},
		{
			Name:        "Energy and Resources",
			Description: "Events related to energy consumption and resource utilization",
			EventTypes: []rModels.MrEventType{
				{Name: "Energy Spike", Description: "Sudden increase in energy usage"},
				{Name: "Supply Outage", Description: "Interruption in energy or resource supply"},
				{Name: "Air Pressure Drop", Description: "Decrease in compressed air supply pressure"},
				{Name: "Water Leak", Description: "Detected water leak in the system"},
			},
		},
		{
			Name:        "Logistics",
			Description: "Events related to material handling and inventory",
			EventTypes: []rModels.MrEventType{
				{Name: "Raw Material Delay", Description: "Delay in raw material delivery"},
				{Name: "Inventory Shortage", Description: "Insufficient inventory for production"},
				{Name: "Overproduction", Description: "Production exceeds planned output"},
			},
		},
		{
			Name:        "Quality and Compliance",
			Description: "Events related to audits, certifications, and quality control",
			EventTypes: []rModels.MrEventType{
				{Name: "Audit Inspection", Description: "Regulatory or internal audit of processes"},
				{Name: "ISO Compliance Review", Description: "Review for ISO certification compliance"},
				{Name: "Quality Control Failure", Description: "Detection of quality control non-conformities"},
			},
		},
		{
			Name:        "Technology and Automation",
			Description: "Events related to control systems and automation",
			EventTypes: []rModels.MrEventType{
				{Name: "PLC Error", Description: "Error detected in the PLC system"},
				{Name: "HMI Update", Description: "Update to the HMI software"},
				{Name: "Sensor Out of Range", Description: "Sensor readings outside acceptable range"},
				{Name: "Actuator Fault", Description: "Fault detected in an actuator"},
			},
		},
		{
			Name:        "Emergencies",
			Description: "Events related to emergencies and disasters",
			EventTypes: []rModels.MrEventType{
				{Name: "Power Outage", Description: "Interruption of electrical power"},
				{Name: "Earthquake Stop", Description: "Emergency stop due to seismic activity"},
				{Name: "Critical Equipment Failure", Description: "Complete failure of critical equipment"},
			},
		},
		{
			Name:        "Documentation and Administration",
			Description: "Events related to administrative tasks and documentation",
			EventTypes: []rModels.MrEventType{
				{Name: "Maintenance Log", Description: "Log entry for completed maintenance"},
				{Name: "Failure Report", Description: "Detailed report of a failure"},
				{Name: "Work Order Creation", Description: "Creation of a new work order for tasks"},
			},
		},
	}

	// Insert categories into the database
	for _, category := range categories {
		db.Create(&category)
	}
}
