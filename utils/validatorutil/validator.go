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

// NewCommonValidator ...
func NewCommonValidator() CommonValidator {
	return CommonValidator{
		validator: validator.New(),
	}
}

// ValidateStruct метод валидации структуры с общей обработкой
func (v *CommonValidator) ValidateStruct(strct struct{}, handler ValidationHandler) error {
	err := v.validator.Struct(strct)
	if err != nil {
		return handler(err.(validator.ValidationErrors))
	}
	return nil
}
