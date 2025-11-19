package rModels

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type ValueType string

const (
	ValueTypeString ValueType = "string"
	ValueTypeNumber ValueType = "number"
)

type MrProductType struct {
	ID        uint                           `gorm:"primaryKey"`
	Name      string                         `gorm:"unique;not null"`
	Features  []MrProductFeatureType         `gorm:"many2many:mr_product_feature_type_relations"`
	Relations []MrProductFeatureTypeRelation `gorm:"foreignKey:MrProductTypeID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type MrProductFeatureType struct {
	ID        uint                           `gorm:"primaryKey"`
	Name      string                         `gorm:"not null;uniqueIndex"`
	Symbol    string                         `gorm:"not null"`
	Unit      string                         `gorm:"not null"`
	ValueType ValueType                      `gorm:"type:varchar(10);not null;default:'string'"`
	Relations []MrProductFeatureTypeRelation `gorm:"foreignKey:MrProductFeatureTypeID"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type MrProductFeatureTypeRelation struct {
	MrProductTypeID        uint `gorm:"primaryKey"`
	MrProductFeatureTypeID uint `gorm:"primaryKey"`
	VisuOrder              uint `gorm:"type:integer"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt `gorm:"index"`
}

type MrProduct struct {
	ID            uint                    `gorm:"primaryKey"`
	Name          string                  `gorm:"not null"`
	ProductTypeID uint                    `gorm:"not null"`
	ProductType   MrProductType           `gorm:"foreignKey:ProductTypeID"`
	Manufacturer  string                  `gorm:"not null"`
	Family        string                  `gorm:"not null"`
	Description   string                  `gorm:"type:text"`
	FeatureValues []MrProductFeatureValue `gorm:"foreignKey:MrProductID"`
	Documents     []MrDocuments           `gorm:"many2many:mr_product_documents;joinForeignKey:ProductID;joinReferences:DocumentID"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type MrProductFeatureValue struct {
	ID                     uint                 `gorm:"primaryKey"`
	MrProductID            uint                 `gorm:"not null;uniqueIndex:idx_product_feature"`
	MrProductFeatureTypeID uint                 `gorm:"not null;uniqueIndex:idx_product_feature"`
	ProductFeature         MrProductFeatureType `gorm:"foreignKey:MrProductFeatureTypeID"`
	StringValue            string               `gorm:"type:text"`
	NumberValue            float64              `gorm:"type:decimal(20,6)"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
	DeletedAt              gorm.DeletedAt `gorm:"index"`
}

type MrAsset struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	ProductID         uint           `gorm:"not null;index" json:"product_id"`
	Product           MrProduct      `gorm:"foreignKey:ProductID" json:"mr_products"`
	ParentID          *uint          `gorm:"index" json:"parent_id"`
	Parent            *MrAsset       `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children          []MrAsset      `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Location          string         `gorm:"type:text;not null" json:"location"`
	TechCode          uint           `gorm:"not null;uniqueIndex" json:"tech_code"`
	Code              string         `gorm:"type:text;not null" json:"code"`
	Sn                string         `gorm:"type:text;not null;uniqueIndex" json:"sn"`
	SapCode           string         `gorm:"type:text" json:"sap_code"`
	Documents         []MrDocuments  `gorm:"many2many:mr_asset_documents;joinForeignKey:AssetID;joinReferences:DocumentID"`
	HierarchicalLevel pq.StringArray `gorm:"type:text[]" json:"hierarchical_level"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// NUEVA RELACIÓN: Un MrAsset puede tener una MrFactory
	MrFactory *MrFactory `gorm:"foreignKey:AssetID" json:"mr_factory,omitempty"`
}

type MrDocuments struct {
	ID          uint `gorm:"primaryKey" json:"id"`
	Nombre      string
	Descripcion string
	URL         string
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
