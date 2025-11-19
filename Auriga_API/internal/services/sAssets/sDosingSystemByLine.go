package sAssets

import (
	"context"
)

type msComponetUnit struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

// msDosingUnit representa una unidad de dosificación
type msDosingUnit struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Location    string           `json:"location"`
	Components  []msComponetUnit `json:"components"`
}

func (s *service) GetDosingSystemByLine(ctx context.Context, lineID uint) ([]msDosingUnit, error) {
	rData, err := s.repository.GetDosingSystemByLine(ctx, lineID)
	if err != nil {
		return nil, err
	}

	var result []msDosingUnit

	for _, dosingSystem := range rData {
		// Para cada sistema de dosificación, procesar sus dosers
		for _, doser := range dosingSystem.Children {
			dosingUnit := msDosingUnit{
				ID:          doser.ID,
				Name:        doser.Code,
				Description: doser.Code, // Puedes ajustar esto según tus necesidades
				Location:    doser.Location,
			}

			// Procesar componentes de cada doser
			for _, component := range doser.Children {
				componentUnit := msComponetUnit{
					ID:       component.ID,
					Name:     component.Code,
					Location: component.Location,
				}
				dosingUnit.Components = append(dosingUnit.Components, componentUnit)
			}

			result = append(result, dosingUnit)
		}
	}

	return result, nil
}

func (s *service) GetDoserComponents(ctx context.Context, doserID uint) ([]msComponetUnit, error) {
	components, err := s.repository.GetDoserComponents(ctx, doserID)
	if err != nil {
		return nil, err
	}

	var result []msComponetUnit

	for _, component := range components {
		componentUnit := msComponetUnit{
			ID:       component.ID,
			Name:     component.Code,
			Location: component.Location,
		}
		result = append(result, componentUnit)
	}

	return result, nil
}
