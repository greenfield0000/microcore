package validatorutil

import (
	"github.com/go-playground/validator/v10"
)

// ValidationHandler Обработчик валидации
type ValidationHandler func(errors validator.ValidationErrors) error

// Validator метод валидации структуры
type Validator interface {
	ValidateStruct(strct struct{}, handler ValidationHandler) error
}

// CommonValidator реализация валидатора по-умолчанию
type CommonValidator struct {
	validator *validator.Validate
}

func customNotEmpty(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	return value != ""
}

// NewCommonValidator ...
func NewCommonValidator() (CommonValidator, error) {
	v := validator.New()
	err := v.RegisterValidation("notEmpty", customNotEmpty)
	if err != nil {
		return CommonValidator{}, err
	}

	return CommonValidator{validator: v}, nil
}

// ValidateStruct метод валидации структуры с общей обработкой
func (v *CommonValidator) ValidateStruct(strct interface{}, handler ValidationHandler) error {
	err := v.validator.Struct(strct)
	if err != nil {
		return handler(err.(validator.ValidationErrors))
	}
	return nil
}
