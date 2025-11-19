package sSap

import (
	"log"

	"github.com/remrafvil/Auriga_API/internal/repositories/rModels"
)

type msComponent struct {
	SapCode           string `json:"SapCode"`
	Description       string `json:"Description"`
	RequiredQuantity  string `json:"RequiredQuantity"`
	MeasurementUnitRQ string `json:"MeasurementUnitRQ"`
	CommittedQuantity string `json:"CommittedQuantity"`
	MeasurementUnitCQ string `json:"MeasurementUnitCQ"`
	WithDrawnQuantity string `json:"WithDrawnQuantity"`
}

type msRecipe struct {
	SapCode     string `json:"SapCode"`
	Description string `json:"Description"`
	Components  []msComponent
}

func (s *service) RecipeComponentList(sapOrderCode string, sapRequest string) (msRecipe, error) {
	var data = msRecipe{}
	var recipeR = rModels.MrRecipe{}
	var recCompR = []rModels.MrRecipeComponent{}

	// Actualizamos o no desde SAP, dependiendo de la solicitud
	if sapRequest == "true" {
		rsData, err := s.repositorySap.RsLineRecipe(sapOrderCode)
		if err != nil {
			log.Println("Error service RsLineRecipe:", err)
			return data, err
		}

		recipeR.SapCode = sapOrderCode
		recipeR.Description = "Esta es la descripcion"
		for _, p := range rsData {
			recipeR.Components = append(recipeR.Components, rModels.MrComponent{
				SapCode:     p.SapCode,
				Description: p.Description,
			})
		}
		for _, p := range rsData {
			recCompR = append(recCompR, rModels.MrRecipeComponent{
				MrRecipeSapCode:    sapOrderCode,
				MrComponentSapCode: p.SapCode,
				RequiredQuantity:   p.RequiredQuantity,
				MeasurementUnitRQ:  p.MeasurementUnitRQ,
				CommittedQuantity:  p.CommittedQuantity,
				MeasurementUnitCQ:  p.MeasurementUnitCQ,
				WithDrawnQuantity:  p.WithDrawnQuantity,
			})
		}

		err = s.repositoryOrd.RecipCompUpdateOrInsert(recipeR, recCompR)
		if err != nil {
			log.Println("Error service RecipCompUpdateOrInsert:", err)
			return data, err
		}
	}

	// Leemos la receta de una orden de fabricación desde POSTGRES
	recipe, err := s.repositoryOrd.LineRecipCompFind(sapOrderCode)
	if err != nil {
		log.Println("Error lectura listado de ordenes Postgrest LineOrderList:", err)
		return data, err
	}
	data.SapCode = recipe.SapCode
	data.Description = recipe.Description
	for i, p := range recipe.Relations {
		data.Components = append(data.Components, msComponent{
			RequiredQuantity:  p.RequiredQuantity,
			MeasurementUnitRQ: p.MeasurementUnitRQ,
			CommittedQuantity: p.CommittedQuantity,
			MeasurementUnitCQ: p.MeasurementUnitCQ,
			WithDrawnQuantity: p.WithDrawnQuantity,
		})
		data.Components[i].SapCode = recipe.Components[i].SapCode
		data.Components[i].Description = recipe.Components[i].Description
	}

	return data, nil
}
