package cores

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type StructValidator struct {
	validator *validator.Validate
}

func NewStructValidator() *StructValidator {
	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		for _, tagKey := range []string{"json", "msgpack"} {
			if tagVal := fld.Tag.Get(tagKey); tagVal != "" {
				name := strings.SplitN(tagVal, ",", 2)[0]
				if name != "" && name != "-" {
					return name
				}
			}
		}
		return fld.Name
	})
	return &StructValidator{
		validator: v,
	}
}

// Validate satisfies the fiber.StructValidator interface.
func (v *StructValidator) Validate(out any) error {
	return v.validator.Struct(out)
}
