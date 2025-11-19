package sEvents

import (
	"context"
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rAssets"
	"github.com/remrafvil/Auriga_API/internal/repositories/rEvents"
	"github.com/remrafvil/Auriga_API/internal/repositories/rLineOrders"
	"github.com/remrafvil/Auriga_API/internal/repositories/riInfluxdb"
	"github.com/remrafvil/Auriga_API/internal/repositories/rsSap"
)

type Service interface {
	EventsRawByLineList(factory string, lineNumber string) ([]msRawEvents, error)
	EventsRawToCommitLine(id uint, eventTime time.Time, factory string, prodline string, system string, machine string, part string, eventTypt string) ([]msCommitEvents, error)
	EventsRawByLineDel(id uint) ([]msRawEvents, error)

	EventsCommitByLineList(factory string, lineNumber string) ([]msCommitEvents, error)
	EventsCommitByLineAdd(eventTime time.Time, factory string, prodline string, system string, machine string, part string, eventTypt string, eventCategory string) ([]msCommitEvents, error)
	EventsCommitByLineUpdate(id uint, eventTime time.Time, factory string, prodline string, system string, machine string, part string, eventTypt string, eventCategory string) ([]msCommitEvents, error)
	EventsCommitByLineDel(id uint) ([]msCommitEvents, error)

	EventsSapByLineDel(id uint) ([]msCommitEvents, error)

	GetAllCategoriesWithEventTypes(ctx context.Context) ([]msEventCategoryDTO, error)
	GetCategoryWithEventTypesByName(ctx context.Context, name string) (*msEventCategoryDTO, error)
}

type service struct {
	repositoryEven   rEvents.Repository
	repositoryAss    rAssets.Repository
	repositoryOrd    rLineOrders.Repository
	repositorySap    rsSap.Repository
	repositoryInflux riInfluxdb.Repository
}

func New(repositoryEven rEvents.Repository, repositoryAss rAssets.Repository, repositoryOrd rLineOrders.Repository, repositorySap rsSap.Repository, repositoryInflux riInfluxdb.Repository) Service {
	return &service{
		repositoryEven:   repositoryEven,
		repositoryAss:    repositoryAss,
		repositoryOrd:    repositoryOrd,
		repositorySap:    repositorySap,
		repositoryInflux: repositoryInflux,
	}
}
