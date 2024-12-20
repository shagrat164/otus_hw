package hw09structvalidator

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func validateLength(fieldValue reflect.Value, param string) error {
	expectedLen, _ := strconv.Atoi(param)

	switch fieldValue.Kind() { //nolint:exhaustive
	case reflect.String:
		if len(fieldValue.String()) != expectedLen {
			return ErrInvalidLength
		}
	case reflect.Slice:
		if fieldValue.Len() != expectedLen {
			return ErrInvalidLength
		}
	}

	return nil
}

func validateMin(fieldValue reflect.Value, param string) error {
	min, _ := strconv.Atoi(param)

	if fieldValue.Kind() == reflect.Int && int(fieldValue.Int()) < min {
		return ErrInvalidMin
	}

	return nil
}

func validateMax(fieldValue reflect.Value, param string) error {
	max, _ := strconv.Atoi(param)

	if fieldValue.Kind() == reflect.Int && int(fieldValue.Int()) > max {
		return ErrInvalidMax
	}

	return nil
}

func validateRegexp(fieldValue reflect.Value, param string) error {
	re, err := regexp.Compile(param)
	if err != nil {
		return ErrInvalidTag
	}

	if fieldValue.Kind() == reflect.String && !re.MatchString(fieldValue.String()) {
		return ErrInvalidRegexp
	}

	return nil
}

func validateIn(fieldValue reflect.Value, param string) error {
	allowedValues := strings.Split(param, ",")

	switch fieldValue.Kind() { //nolint:exhaustive
	case reflect.String:
		for _, v := range allowedValues {
			if fieldValue.String() == v {
				return nil
			}
		}
	case reflect.Int:
		for _, v := range allowedValues {
			allowedInt, _ := strconv.Atoi(v)
			if int(fieldValue.Int()) == allowedInt {
				return nil
			}
		}
	}

	return ErrInvalidIn
}
