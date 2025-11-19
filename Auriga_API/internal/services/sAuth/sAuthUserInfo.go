package sAuth

import (
	"context"

	"go.uber.org/zap"
)

// GetUserInfoFromContext obtiene la información del usuario desde el contexto
func (s *service) GetUserInfoFromContext(ctx context.Context) (userID, email, name string, groups []string, organization map[string]interface{}) {
	// Acceder directamente a los valores del contexto estándar
	userID, _ = ctx.Value("user_id").(string)
	email, _ = ctx.Value("user_email").(string)
	name, _ = ctx.Value("user_name").(string)
	organization, _ = ctx.Value("user_organization").(map[string]interface{})

	// Extraer grupos si se necesitan
	if groupsInterface, ok := ctx.Value("user_groups").([]interface{}); ok {
		for _, group := range groupsInterface {
			if groupStr, ok := group.(string); ok {
				groups = append(groups, groupStr)
			}
		}
	}

	s.logger.Debug("User info extracted from context",
		zap.String("user_id", userID),
		zap.String("email", email),
		zap.String("name", name),
		zap.Int("groups_count", len(groups)),
		zap.Bool("has_organization", organization != nil))

	return
}

// GetUserFactoryNames obtiene los nombres de las fábricas del usuario
func (s *service) GetUserFactoryNames(organization map[string]interface{}) []string {
	var factoryNames []string

	if organization == nil {
		return factoryNames
	}

	factories, ok := organization["factories"].(map[string]interface{})
	if !ok {
		return factoryNames
	}

	for factoryName := range factories {
		factoryNames = append(factoryNames, factoryName)
	}

	s.logger.Debug("User factories extracted",
		zap.Strings("factories", factoryNames),
		zap.Int("count", len(factoryNames)))

	return factoryNames
}
