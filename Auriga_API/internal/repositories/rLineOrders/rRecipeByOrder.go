package rLineOrders

// lo llamaremso repositories

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
	"gorm.io/gorm"
)

/******************************      COMPONENT       ******************************/

func (m *repository) ComponentList(ctx context.Context) []rModels.MrComponent { // Lista todos
	var data = []rModels.MrComponent{}

	log.Println("*****************      Llego por aquí REPOSITORY *****************")
	err := m.db.Find(&data).Error
	if err != nil {
		log.Fatalln(err)
	}
	return data
}

func (m *repository) ComponentInfo(id uint) (rModels.MrComponent, error) { // Encuentra uno
	var data rModels.MrComponent
	rs := m.db.Where("id  = ?", id).Find(&data)
	if rs.RowsAffected == 0 {
		log.Println("no existe registro con ese identificador")
		return rModels.MrComponent{}, errors.New("no existe registro con ese id")
	}
	if rs.Error != nil {

		return rModels.MrComponent{}, errors.New("fallo registro con ese id")
	}

	return data, nil
}

func (m *repository) ComponentAdd(SapCode string, Description string) (rModels.MrComponent, error) { // Añade uno
	var data rModels.MrComponent

	data.SapCode = SapCode
	data.Description = Description

	if err := m.db.Create(&data).Error; err != nil {
		return rModels.MrComponent{}, err
	} else {

		return data, nil
	}
}

func (m *repository) ComponentUpdate(id uint, SapCode string, Description string) (rModels.MrComponent, error) { // Actualiza uno
	var data rModels.MrComponent

	if err := m.db.Where("id = ? ", id).First(&data).Error; err != nil {
		return rModels.MrComponent{}, errors.New("no hay registro con ese id")
	}
	data.SapCode = SapCode
	data.Description = Description

	if err := m.db.Save(&data).Error; err != nil {
		return rModels.MrComponent{}, errors.New("no se pudo actualizar")
	}
	return data, nil

}

func (m *repository) ComponentDel(id uint) error { // Borra uno
	var data rModels.MrComponent
	if err := m.db.Where("id = ?", id).Delete(&data).Error; err != nil {
		return err
	}
	return nil
}

func (m *repository) ComponentUpdateOrInsert(components []rModels.MrComponent) ([]rModels.MrComponent, error) {
	var data []rModels.MrComponent

	for _, p := range components {
		if err := m.db.FirstOrCreate(&p, rModels.MrComponent{SapCode: p.SapCode}).Error; err != nil {
			return data, err
		}
	}
	return data, nil
}

/******************************        RECIPE        ******************************/

func (m *repository) RecipeList(ctx context.Context) []rModels.MrRecipe { // Lista todos
	var data = []rModels.MrRecipe{}

	log.Println("*****************      Llego por aquí REPOSITORY *****************")
	err := m.db.Find(&data).Error
	if err != nil {
		log.Fatalln(err)
	}
	return data
}

func (m *repository) RecipeInfo(id uint) (rModels.MrRecipe, error) { // Encuentra uno
	var data rModels.MrRecipe
	rs := m.db.Where("id  = ?", id).Find(&data)
	if rs.RowsAffected == 0 {
		log.Println("no existe registro con ese identificador")
		return rModels.MrRecipe{}, errors.New("no existe registro con ese id")
	}
	if rs.Error != nil {

		return rModels.MrRecipe{}, errors.New("fallo registro con ese id")
	}

	return data, nil
}

func (m *repository) RecipeAdd(SapCode string, Description string) (rModels.MrRecipe, error) { // Añade uno
	var data rModels.MrRecipe

	data.SapCode = SapCode
	data.Description = Description

	if err := m.db.Create(&data).Error; err != nil {
		return rModels.MrRecipe{}, err
	} else {

		return data, nil
	}
}

func (m *repository) RecipeUpdate(id uint, SapCode string, Description string) (rModels.MrRecipe, error) { // Actualiza uno
	var data rModels.MrRecipe

	if err := m.db.Where("id = ? ", id).First(&data).Error; err != nil {
		return rModels.MrRecipe{}, errors.New("no hay registro con ese id")
	}
	data.SapCode = SapCode
	data.Description = Description

	if err := m.db.Save(&data).Error; err != nil {
		return rModels.MrRecipe{}, errors.New("no se pudo actualizar")
	}
	return data, nil

}

func (m *repository) RecipeDel(id uint) error { // Borra uno
	var data rModels.MrRecipe
	if err := m.db.Where("id = ?", id).Delete(&data).Error; err != nil {
		return err
	}
	return nil
}

func (m *repository) RecipeUpdateOrInsert(SapCode string, Description string) (rModels.MrRecipe, error) {
	var data rModels.MrRecipe
	var ui_recipe rModels.MrRecipe

	ui_recipe.SapCode = SapCode
	ui_recipe.Description = Description

	if err := m.db.Where(rModels.MrRecipe{SapCode: SapCode}).FirstOrCreate(&ui_recipe).Error; err != nil {
		return data, err
	}
	log.Println("*****************      Llego por aquí REPOSITORY RECIPE Update o Create *****************")

	log.Println(ui_recipe.SapCode)
	log.Println(ui_recipe.Description)
	m.db.Where("sap_code = ?", ui_recipe.SapCode).Find(&data)
	return data, nil
}

/******************************   COMPONENT / RECIPE   ******************************/

func (m *repository) RecipCompUpdateOrInsert(recipe rModels.MrRecipe, recComp []rModels.MrRecipeComponent) error {
	// Usar una transacción para crear los datos
	err := m.db.Transaction(func(tx *gorm.DB) error {
		// Cargar datos de la receta
		recipe1 := rModels.MrRecipe{SapCode: recipe.SapCode, Description: recipe.Description}
		if err := tx.FirstOrCreate(&recipe1).Error; err != nil {
			return fmt.Errorf("failed to create recipe: %w", err)
		}
		// Cargar datos de los componentes
		for _, component := range recipe.Components {
			if err := tx.FirstOrCreate(&component, rModels.MrComponent{SapCode: component.SapCode}).Error; err != nil {
				return fmt.Errorf("failed to create/find component: %w", err)
			}

			log.Println("Este es el dato updateData2", component.SapCode)
		}

		// Crear relaciones en la tabla intermedia
		for _, relation := range recComp {
			if err := tx.Save(&relation).Error; err != nil {
				return fmt.Errorf("failed to create recipe-component relation: %w", err)
			}
		}
		// Si todo va bien, la transacción se confirma
		return nil
	})

	// Manejar errores de transacción
	if err != nil {
		fmt.Println("Transaction failed RecipCompUpdateOrInsert:", err)
		return err
	}
	fmt.Println("Transaction completed successfully RecipCompUpdateOrInsert")
	return nil
}

func (m *repository) LineRecipCompFind(sapOrderCode string) (rModels.MrRecipe, error) {
	var recipe rModels.MrRecipe

	err := m.db.Preload("Components", func(db *gorm.DB) *gorm.DB {
		return db.Order("mr_components.sap_code")
	}).
		Preload("Relations", func(db *gorm.DB) *gorm.DB {
			return db.Order("mr_component_sap_code")
		}).
		First(&recipe, "sap_code = ?", sapOrderCode).Error

	if err != nil {
		fmt.Println("Error fetching recipe components:", err)
		return recipe, err
	}

	return recipe, nil
}

func DB_InitComponents(c *gorm.DB) { //
	// c.Create(&rModels.MrComponent{ComponentName: "RM00006718", MaterialDescription: "HIPS NATURAL 4400 PAMPA", RequiredQuantity: "33803.200", MeeasurenmentUnitRQ: "KG", CommitedQuantity: "24539.000", MeeasurenmentUnitCQ: "KG", WithDrawnQuantity: "24539.000"})
	// c.Create(&rModels.MrComponent{ComponentName: "RM00006597", MaterialDescription: "GPPS CRISTAL HH101 PAMPA", RequiredQuantity: "4568.000", MeeasurenmentUnitRQ: "KG", CommitedQuantity: "2635.490", MeeasurenmentUnitCQ: "KG", WithDrawnQuantity: "0.000"})
	// c.Create(&rModels.MrComponent{ComponentName: "RM00006678", MaterialDescription: "MASTER BATCH PS AMARILLO 1512/362316", RequiredQuantity: "913.600", MeeasurenmentUnitRQ: "KG", CommitedQuantity: "91.000", MeeasurenmentUnitCQ: "KG", WithDrawnQuantity: "272.000"})
	// c.Create(&rModels.MrComponent{ComponentName: "RE00003548", MaterialDescription: "RE RECUP. CREMA (1311) COEXPAN", RequiredQuantity: "4568.000", MeeasurenmentUnitRQ: "ZPKG01", CommitedQuantity: "1084.000<", MeeasurenmentUnitCQ: "KG", WithDrawnQuantity: "1084.000"})
}
