package httputils

import (
	"github.com/go-playground/validator/v10"
)

type Validator struct {
	delegate *validator.Validate
}

func NewValidator() *Validator {
	return NewValidatorWithDelegator(validator.New())
}

func NewValidatorWithDelegator(validator *validator.Validate) *Validator {
	return &Validator{delegate: validator}
}

func (v *Validator) Validate(i interface{}) error {
	return v.delegate.Struct(i)
}
