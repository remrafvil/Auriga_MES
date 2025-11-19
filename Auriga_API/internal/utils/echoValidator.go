// internal/utils/validator.go
package utils

import (
	"reflect"
	"time"

	"github.com/go-playground/validator/v10"
)

// CustomValidator estructura para el validador personalizado
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator crea una nueva instancia del validador con validaciones personalizadas
func NewCustomValidator() *CustomValidator {
	v := validator.New()

	// Registrar validación personalizada para time.Time
	v.RegisterValidation("notzerotime", func(fl validator.FieldLevel) bool {
		field := fl.Field()

		if field.Type() == reflect.TypeOf(time.Time{}) {
			t := field.Interface().(time.Time)
			return !t.IsZero()
		}

		return false
	})

	return &CustomValidator{validator: v}
}

// Validate implementa la interfaz de Echo
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// Helper para obtener mensajes de error más legibles
func (cv *CustomValidator) GetValidationErrors(err error) map[string]string {
	validationErrors := make(map[string]string)

	if ve, ok := err.(validator.ValidationErrors); ok {
		for _, fe := range ve {
			field := fe.Field()
			tag := fe.Tag()

			switch tag {
			case "required":
				validationErrors[field] = "Este campo es requerido"
			case "min":
				validationErrors[field] = "El valor es demasiado corto"
			case "max":
				validationErrors[field] = "El valor es demasiado largo"
			case "email":
				validationErrors[field] = "Formato de email inválido"
			case "notzerotime":
				validationErrors[field] = "La fecha/hora es requerida"
			default:
				validationErrors[field] = "Valor inválido"
			}
		}
	}

	return validationErrors
}
