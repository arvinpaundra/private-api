package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Error map[string]string

type Validator struct {
	validate *validator.Validate
}

func NewValidator() *Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			return ""
		}

		return name
	})

	return &Validator{
		validate: v,
	}
}

func (v *Validator) Validate(val any) Error {
	err := v.validate.Struct(val)

	if err != nil {
		return formatError(err)
	}

	return nil
}

func formatError(errs error) Error {
	verrs := make(Error)

	for _, err := range errs.(validator.ValidationErrors) {
		field := err.Field()

		switch err.Tag() {
		case "required", "required_if":
			verrs[field] = "this field is required"
		case "email":
			verrs[field] = "invalid email format"
		case "min":
			verrs[field] = fmt.Sprintf("min length %s characters", err.Param())
		case "max":
			verrs[field] = fmt.Sprintf("max length %s characters", err.Param())
		case "numeric":
			verrs[field] = "only numeric format"
		case "oneof":
			params := strings.ReplaceAll(err.Param(), " ", ", ")
			verrs[field] = fmt.Sprintf("must be one of %s", params)
		default:
			verrs[field] = err.Error()
		}
	}

	return verrs
}
