package hw09structvalidator

import (
	"reflect"
)

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	var validationErrors ValidationErrors

	t := val.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := val.Field(i)

		tag := field.Tag.Get("validate")
		if tag == "" {
			continue
		}

		err := validateField(field.Name, fieldValue, tag)
		if err != nil {
			if errs, ok := err.(ValidationErrors); ok { //nolint:errorlint
				validationErrors = append(validationErrors, errs...)
			} else {
				return err
			}
		}
	}

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}
