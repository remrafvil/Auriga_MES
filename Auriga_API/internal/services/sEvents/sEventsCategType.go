package sEvents

import (
	"context"
)

// DTOs para la respuesta
type msEventTypeDTO struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type msEventCategoryDTO struct {
	ID          uint             `json:"id"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	EventTypes  []msEventTypeDTO `json:"eventTypes"`
}

func (s *service) GetAllCategoriesWithEventTypes(ctx context.Context) ([]msEventCategoryDTO, error) {
	categories, err := s.repositoryEven.FindCategoriesWithEventTypes(ctx)

	if err != nil {
		return []msEventCategoryDTO{}, err
	}

	// Convertimos los modelos a DTOs
	var categoryDTOs []msEventCategoryDTO
	for _, category := range categories {
		var eventTypeDTOs []msEventTypeDTO
		for _, eventType := range category.EventTypes {
			eventTypeDTOs = append(eventTypeDTOs, msEventTypeDTO{
				ID:          eventType.ID,
				Name:        eventType.Name,
				Description: eventType.Description,
			})
		}

		categoryDTOs = append(categoryDTOs, msEventCategoryDTO{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			EventTypes:  eventTypeDTOs,
		})
	}

	return categoryDTOs, nil
}

func (s *service) GetCategoryWithEventTypesByName(ctx context.Context, name string) (*msEventCategoryDTO, error) {
	category, err := s.repositoryEven.FindCategoryWithEventTypesByName(ctx, name)
	if err != nil {
		return nil, err
	}

	if category == nil {
		return nil, nil // Categoría no encontrada
	}

	// Convertir los tipos de evento a DTOs
	var eventTypeDTOs []msEventTypeDTO
	for _, et := range category.EventTypes {
		eventTypeDTOs = append(eventTypeDTOs, msEventTypeDTO{
			ID:          et.ID,
			Name:        et.Name,
			Description: et.Description,
		})
	}

	return &msEventCategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		EventTypes:  eventTypeDTOs,
	}, nil
}
