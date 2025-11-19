package rsSap

// lo llamaremso repositories

import (
	"time"

	"github.com/remrafvil/Auriga_API/config"
	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

type Repository interface {
	rsFactoryLineList()
	RsLineOrderList(factory string, lineNumber string, lineSapCode string) ([]rModels.MrProductionOrder, error)
	RsLineRecipe(lineNumber string) ([]MrsComponent, error)
	RsLineStopEvent(maquina string, timeEvent time.Time, estado string, motivo string, of string, operario string) error
}

type repository struct {
	db     *gorm.DB
	config *config.Settings
}

func New(db *gorm.DB, cfg *config.Settings) Repository {
	return &repository{
		db:     db,
		config: cfg,
	}
}
