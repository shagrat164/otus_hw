package hw09structvalidator

import (
	"reflect"
	"strings"
)

func validateField(fieldName string, fieldValue reflect.Value, tag string) error {
	var validationErrors ValidationErrors

	rules := strings.Split(tag, "|")
	for _, rule := range rules {
		parts := strings.SplitN(rule, ":", 2)
		if len(parts) != 2 {
			return ErrInvalidTag
		}

		switch fieldValue.Kind() { //nolint:exhaustive
		case reflect.Slice:
			for i := 0; i < fieldValue.Len(); i++ {
				elem := fieldValue.Index(i)
				err := applyRule(elem, parts[0], parts[1])
				if err != nil {
					validationErrors = append(validationErrors, ValidationError{
						Field: fieldName,
						Err:   err,
					})
				}
			}
		default:
			err := applyRule(fieldValue, parts[0], parts[1])
			if err != nil {
				validationErrors = append(validationErrors, ValidationError{
					Field: fieldName,
					Err:   err,
				})
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func applyRule(value reflect.Value, rule, param string) error {
	switch rule {
	case "len":
		return validateLength(value, param)
	case "min":
		return validateMin(value, param)
	case "max":
		return validateMax(value, param)
	case "regexp":
		return validateRegexp(value, param)
	case "in":
		return validateIn(value, param)
	default:
		return ErrInvalidTag
	}
}
