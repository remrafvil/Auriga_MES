package rModels

type MrInfluxGrafanaQuery struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:text;uniqueIndex" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Query       string `gorm:"type:text;not null" json:"query"`
}
