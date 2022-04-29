package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrTag        = errors.New("validate tag ill-formed")
	ErrNotAStruct = errors.New("interface is not a structure")
)

/* Validation errors. */
var (
	ErrMin    = errors.New("value is less than minimum")
	ErrMax    = errors.New("value is greater than maximum")
	ErrIn     = errors.New("no values match")
	ErrLen    = errors.New("string len doesn't fit")
	ErrRegexp = errors.New("string doesn't match regexp")
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errsLine := "ValidationErrors:"
	for i := range v {
		errsLine += fmt.Sprintf(" %s", v[i])
	}

	return errsLine
}

func ValidateInt(intValue int, tag string) error {
	args := strings.FieldsFunc(tag, func(r rune) bool { return r == ':' })
	if len(args) != 2 {
		return ErrTag
	}

	switch {
	case args[0] == "min":
		i, err := strconv.Atoi(args[1])
		if err != nil {
			return ErrTag
		}
		if i >= intValue {
			return ErrMin
		}
	case args[0] == "max":
		i, err := strconv.Atoi(args[1])
		if err != nil {
			return ErrTag
		}
		if i <= intValue {
			return ErrMax
		}
	case args[0] == "in":
		values := strings.FieldsFunc(args[1], func(r rune) bool { return r == ',' })
		inValues := make([]int, 0, len(values))
		for i := range values {
			v, err := strconv.Atoi(values[i])
			if err != nil {
				return ErrTag
			}
			inValues = append(inValues, v)
		}
		for i := range inValues {
			if inValues[i] == intValue {
				return nil
			}
		}
		return ErrIn
	default:
		return ErrTag
	}

	return nil
}

func ValidateString(str string, tag string) error {
	args := strings.FieldsFunc(tag, func(r rune) bool { return r == ':' })
	if len(args) != 2 {
		return ErrTag
	}

	switch {
	case args[0] == "len":
		i, err := strconv.Atoi(args[1])
		if err != nil {
			return ErrTag
		}
		if len(str) != i {
			return ErrLen
		}
	case args[0] == "regexp":
		r, err := regexp.Compile(args[1])
		if err != nil {
			return ErrTag
		}

		if !r.MatchString(str) {
			return ErrRegexp
		}
	case args[0] == "in":
		values := strings.FieldsFunc(args[1], func(r rune) bool { return r == ',' })
		for i := range values {
			if values[i] == str {
				return nil
			}
		}
		return ErrIn
	default:
		return ErrTag
	}

	return nil
}

func ValidateCondition(v reflect.Value, condition string) (err error) {
	switch v.Kind() {
	case reflect.String:
		err = ValidateString(v.String(), condition)
	case reflect.Int:
		err = ValidateInt(int(v.Int()), condition)
	case reflect.Slice:
		for i := 0; i < v.Len(); i++ {
			err = ValidateCondition(v.Index(i), condition)
			if err != nil {
				break
			}
		}
	case reflect.Struct:
		if condition == "nested" && v.CanInterface() {
			/* Recursion */
			err = Validate(v.Interface())
		}
	default:
		fmt.Printf("not implemented tag: %v\n", v.Kind())
	}

	return
}

func Validate(vinterface interface{}) error {
	var errorsArray ValidationErrors

	v := reflect.ValueOf(vinterface)
	if v.Kind() != reflect.Struct {
		return ErrNotAStruct
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		validateTag, ok := ft.Tag.Lookup("validate")
		if !ok {
			continue
		}

		conditions := strings.FieldsFunc(validateTag, func(r rune) bool { return r == '|' })
		if len(conditions) == 0 {
			return ErrTag
		}

		for j := range conditions {
			err := ValidateCondition(fv, conditions[j])

			var errs ValidationErrors
			switch {
			case errors.Is(err, ErrTag):
				return err
			case errors.As(err, &errs):
				errorsArray = append(errorsArray, errs...)
			case err != nil:
				errorsArray = append(errorsArray, ValidationError{Err: err, Field: ft.Name})
			}
		}
	}

	if len(errorsArray) == 0 {
		return nil
	}
	return errorsArray
}
