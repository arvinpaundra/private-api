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
		namespace := err.Namespace()
		parts := strings.SplitN(namespace, ".", 2)

		field := err.Field()

		if len(parts) >= 2 {
			field = parts[1]
		}

		switch err.Tag() {
		case "required", "required_if":
			verrs[field] = "this field is required"
		case "email":
			verrs[field] = "invalid email format"
		case "min":
			verrs[field] = formatMinError(err)
		case "max":
			verrs[field] = formatMaxError(err)
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

func formatMinError(err validator.FieldError) string {
	switch err.Kind() {
	case reflect.String:
		return fmt.Sprintf("min length %s characters", err.Param())
	case reflect.Slice, reflect.Array, reflect.Map:
		return fmt.Sprintf("must contain at least %s item(s)", err.Param())
	default:
		return fmt.Sprintf("min value %s", err.Param())
	}
}

func formatMaxError(err validator.FieldError) string {
	switch err.Kind() {
	case reflect.String:
		return fmt.Sprintf("max length %s characters", err.Param())
	case reflect.Slice, reflect.Array, reflect.Map:
		return fmt.Sprintf("must contain at most %s item(s)", err.Param())
	default:
		return fmt.Sprintf("max value %s", err.Param())
	}
}
