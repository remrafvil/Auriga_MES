package rModels

import (
	"time"

	"gorm.io/gorm"
)

type MrRawEvents struct {
	Time      time.Time      `gorm:"column:_time"`
	ESW       int            `gorm:"column:\"ESW\";type:int4"`
	EL_Lv0    string         `gorm:"column:\"EL_Lv0\""`
	Name      string         `gorm:"column:name"`
	EL_Lv1    string         `gorm:"column:\"EL_Lv1\""`
	EL_Lv2    string         `gorm:"column:\"EL_Lv2\""`
	EL_Lv3    string         `gorm:"column:\"EL_Lv3\""`
	Host      string         `gorm:"column:host"`
	SlaveID   string         `gorm:"column:slave_id"`
	Type      string         `gorm:"column:type"`
	ID        uint           `gorm:"primaryKey"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type MrCommitEvents struct {
	EventTime     time.Time
	EventType     string `gorm:"type:text"`
	EventCategory string `gorm:"type:text"`
	Factory       string `gorm:"type:text"`
	ProdLine      string `gorm:"type:text"`
	System        string `gorm:"type:text"`
	Machine       string `gorm:"type:text"`
	Part          string `gorm:"type:text"`
	gorm.Model
}

// Event Category
type MrEventCategory struct {
	ID           uint            `gorm:"primaryKey"`
	Name         string          `gorm:"size:255;not null"` // Name of the category
	Description  string          `gorm:"type:text"`         // Description of the category
	EventTypes   []MrEventType   `gorm:"foreignKey:event_category_id"`
	ProductTypes []MrProductType `gorm:"many2many:mr_product_type_event_categories"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// Event Type
type MrEventType struct {
	ID              uint            `gorm:"primaryKey"`
	Name            string          `gorm:"size:255;not null"` // Name of the event type
	Description     string          `gorm:"type:text"`         // Description of the event type
	EventCategoryID uint            `gorm:"not null"`          // Foreign key to EventCategory
	Category        MrEventCategory `gorm:"foreignKey:event_category_id"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
