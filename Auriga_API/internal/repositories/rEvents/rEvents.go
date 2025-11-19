package rEvents

// lo llamaremso repositories

import (
	"context"
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

type Repository interface {
	EventsRawByLineList(factory string, location string) ([]rModels.MrRawEvents, error)
	EventsRawByLineDel(id uint) ([]rModels.MrRawEvents, error)

	EventsRawToCommitLine(id uint, eventTime time.Time, factory string, prodline string, system string, machine string, part string, eventTypt string) ([]rModels.MrCommitEvents, error)

	EventsCommitByLineList(factory string, location string) ([]rModels.MrCommitEvents, error)
	EventsCommitByLineAdd(eventTime time.Time, factory string, prodLine string, system string, machine string, part string, eventTypt string, eventCategory string) ([]rModels.MrCommitEvents, error)
	EventsCommitByLineUpdate(id uint, eventTime time.Time, factory string, prodLine string, system string, machine string, part string, eventTypt string, eventCategory string) ([]rModels.MrCommitEvents, error)
	EventsCommitByLineDel(id uint) ([]rModels.MrCommitEvents, error)
	EventsCommitByLineFind(id uint) (rModels.MrCommitEvents, string, error)

	FindCategoriesWithEventTypes(ctx context.Context) ([]rModels.MrEventCategory, error)
	FindCategoryWithEventTypesByName(ctx context.Context, name string) (*rModels.MrEventCategory, error)
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}
