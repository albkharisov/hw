package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var TagError = errors.New("validate tag ill-formed")

// Validation errors
var (
	MinError    = errors.New("value is less than minimum")
	MaxError    = errors.New("value is greater than maximum")
	InError     = errors.New("no values match")
	LenError    = errors.New("string len doesn't fit")
	RegexpError = errors.New("string doesn't match regexp")
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
	fmt.Printf("ValidateInt: '%v', tag: '%v'\n", intValue, tag)

	args := strings.FieldsFunc(tag, func(r rune) bool { return r == ':' })
	if len(args) != 2 {
		return TagError
	}

	switch {
	case args[0] == "min":
		i, err := strconv.Atoi(args[1])
		if err != nil {
			return TagError
		}
		if i >= intValue {
			return MinError
		}
	case args[0] == "max":
		i, err := strconv.Atoi(args[1])
		if err != nil {
			return TagError
		}
		if i <= intValue {
			return MaxError
		}
	case args[0] == "in":
		values := strings.FieldsFunc(args[1], func(r rune) bool { return r == ',' })
		for i := range values {
			v, err := strconv.Atoi(values[i])
			if err != nil {
				return TagError
			}

			if v == intValue {
				return nil
			}
		}
		return InError
	default:
		return TagError
	}

	return nil
}

func ValidateString(str string, tag string) error {
	fmt.Printf("ValidateString: '%v', tag: '%v'\n", str, tag)

	args := strings.FieldsFunc(tag, func(r rune) bool { return r == ':' })
	if len(args) != 2 {
		return TagError
	}

	switch {
	case args[0] == "len":
		i, err := strconv.Atoi(args[1])
		if err != nil {
			return TagError
		}
		if len(str) != i {
			return LenError
		}
	case args[0] == "regexp":
		r, err := regexp.Compile(args[1])
		if err != nil {
			return TagError
		}

		if !r.MatchString(str) {
			return RegexpError
		}
	case args[0] == "in":
		values := strings.FieldsFunc(args[1], func(r rune) bool { return r == ',' })
		for i := range values {
			if values[i] == str {
				return nil
			}
		}
		return InError
	default:
		return TagError
	}

	return nil
}

func ValidateCondition(v reflect.Value, condition string) (err error) {
	fmt.Println("ValidateCondition")
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
	default:
		panic(fmt.Errorf("not implemented tag, %T, %v", v.Kind(), v.Kind()))
	}

	return
}

func Validate(vinterface interface{}) error {

	var errorsArray ValidationErrors

	v := reflect.ValueOf(vinterface)
	if v.Kind() != reflect.Struct {
		return errors.New("interface is not a structure")
	}

	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		ft := t.Field(i)
		fv := v.Field(i)

		validateTag, ok := ft.Tag.Lookup("validate")
		if !ok {
			fmt.Printf("%v: Lookup don't found 'validate' in tag: %v\n", ft.Name, ft.Tag)
			continue
		}

		conditions := strings.FieldsFunc(validateTag, func(r rune) bool { return r == '|' })
		for j := range conditions {
			err := ValidateCondition(fv, conditions[j])
			if err == TagError {
				return err
			}
			if err != nil {
				errorsArray = append(errorsArray, ValidationError{Err: err, Field: ft.Name})
			}
		}
	}

	if len(errorsArray) == 0 {
		return nil
	}
	return errorsArray
}
