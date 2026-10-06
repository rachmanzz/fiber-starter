package cores

import (
	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	validator *validator.Validate
}

func NewStructValidator() *StructValidator {
	return &StructValidator{
		validator: validator.New(),
	}
}

// Validate satisfies the fiber.StructValidator interface.
func (v *StructValidator) Validate(out any) error {
	return v.validator.Struct(out)
}
