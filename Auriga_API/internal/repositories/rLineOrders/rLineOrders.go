package rLineOrders

// lo llamaremso repositories

import (
	"context"
	"time"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

type Repository interface {
	LineOrdersList() ([]rModels.MrProductionOrder, error)
	LineOrdersInfo(id uint) (rModels.MrProductionOrder, error)
	LineOrdersAdd(Factory string, ProdLine string, OrderNumber string, OrderNType string, ProductName string, ProductDescription string, QuantityToProduce string, QuantityProduced string, QuantityRemainedToProduce string, MeasurementUnit string, StarteddAt time.Time, FinishedAt time.Time) (rModels.MrProductionOrder, error)
	LineOrdersUpdate(id uint, Factory string, ProdLine string, OrderNumber string, OrderNType string, ProductName string, ProductDescription string, QuantityToProduce string, QuantityProduced string, QuantityRemainedToProduce string, MeasurementUnit string, StarteddAt time.Time, FinishedAt time.Time) (rModels.MrProductionOrder, error)
	LineOrdersDel(id uint) error
	LineOrdersUpdateOrInsert(lineOrders []rModels.MrProductionOrder) error
	LineOrdersFind(factory string, location string) ([]rModels.MrProductionOrder, error)
	LineOrderStartFinish(factory string, prodLine string, orderNumber string, startFinish string) error
	LineOrderUpdateTime(factory string, prodLine string, orderNumber string, starteddAt time.Time, finishedAt time.Time) error
	LineOrdersFindStartFinish(factory string, location string, orderNumber string) (time.Time, time.Time, error)

	ComponentList(ctx context.Context) []rModels.MrComponent
	ComponentInfo(id uint) (rModels.MrComponent, error)
	ComponentAdd(SapCode string, Description string) (rModels.MrComponent, error)
	ComponentUpdate(id uint, SapCode string, Description string) (rModels.MrComponent, error)
	ComponentDel(id uint) error
	ComponentUpdateOrInsert(components []rModels.MrComponent) ([]rModels.MrComponent, error)

	RecipeList(ctx context.Context) []rModels.MrRecipe
	RecipeInfo(id uint) (rModels.MrRecipe, error)
	RecipeAdd(SapCode string, Description string) (rModels.MrRecipe, error)
	RecipeUpdate(id uint, SapCode string, Description string) (rModels.MrRecipe, error)
	RecipeDel(id uint) error
	RecipeUpdateOrInsert(SapCode string, Description string) (rModels.MrRecipe, error)

	RecipCompUpdateOrInsert(recipe rModels.MrRecipe, recComp []rModels.MrRecipeComponent) error
	LineRecipCompFind(sapOrderCode string) (rModels.MrRecipe, error)

	ConsumptionByOrder(recipeSapCode string, factory string, prodLine string) ([]rModels.MrConsumption, error)
	ConsumptionComponentAdd(recipeSapCode string, sapComponentCode string, factory string, prodLine string, dosingUnit string, dosingComponent string) ([]rModels.MrConsumption, error)
	ConsumptionComponentDel(recipeSapCode string, sapComponentCode string, factory string, prodLine string, dosingUnit string, dosingComponent string) ([]rModels.MrConsumption, error)
	ConsumptionComponentUpdate(recipeSapCode string, sapComponentCode string, factory string, prodLine string, dosingUnit string, dosingComponent string) ([]rModels.MrConsumption, error)
	ConsumptionAddReal(realConsumptions []rModels.MrConsumption) error
}

type repository struct {
	db *gorm.DB
}

func New(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}
