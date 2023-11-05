package validation

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Error struct {
	Field string
	Err   error
}

type Errors []Error

var (
	ErrWrongValue      = errors.New("error during parsing input")
	ErrUnsupportedType = errors.New("unsupported type for validation")
)

var (
	ErrLenNotEqual       = errors.New("len doesn't satisfy the rule")
	ErrOutOfSet          = errors.New("value doesn't belong to the set")
	ErrValueOutOfBondary = errors.New("values is less or over the limit")
	ErrNotRegexp         = errors.New("value doesn't satisfy the regexp")
)

func (v Errors) Error() string {
	var err string
	for _, s := range v {
		err += fmt.Sprintf("field: %s, error: %s ", s.Field, s.Err.Error())
	}
	return err
}

func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	valType := val.Type()

	if val.Kind() != reflect.Struct {
		return ErrUnsupportedType
	}

	valErrors := make(Errors, 0, val.NumField())

	for i := 0; i < val.NumField(); i++ {
		fval := val.Field(i)
		if fval.Kind() == reflect.Pointer {
			fval = fval.Elem()
		}
		if fval.Kind() == reflect.Struct {
			continue
		}

		value, ok := valType.Field(i).Tag.Lookup("validate")
		if !ok {
			continue
		}

		if err := VerifyRules(fval, value); err != nil {
			if errors.Is(err, ErrWrongValue) || errors.Is(err, ErrUnsupportedType) {
				return err
			}
			valErrors = append(valErrors, Error{
				Field: valType.Field(i).Name,
				Err:   err,
			})
		}
	}
	if len(valErrors) == 0 {
		return nil
	}

	return valErrors
}

func VerifyRules(v reflect.Value, rule string) error {
	rules := strings.Split(rule, "|")

	for _, r := range rules {
		switch {
		case strings.HasPrefix(r, "len:"):
			if err := verifyLen(v, r); err != nil {
				return err
			}
		case strings.HasPrefix(r, "in:"):
			if err := verifyIn(v, r); err != nil {
				return err
			}
		case strings.HasPrefix(r, "min:") || strings.HasPrefix(r, "max:"):
			if err := VerifyMinOrMax(v, r); err != nil {
				return err
			}
		case strings.HasPrefix(r, "regexp:"):
			if err := verifyRegExp(v, r); err != nil {
				return err
			}
		}
	}

	return nil
}

// VerifyLen tells whether value's length is equal to the given length or not
// Works for strings and slices of strings
// For slices works for each value in it.
func verifyLen(v reflect.Value, r string) error {
	parts := strings.SplitN(r, ":", 2)
	if len(parts) != 2 {
		return ErrWrongValue
	}

	switch kind := v.Kind(); kind {
	case reflect.String:
		l, err := strconv.Atoi(parts[1])
		if err != nil {
			return ErrWrongValue
		}

		if l != v.Len() {
			return ErrLenNotEqual
		}
	case reflect.Slice:
		if stringSlice, ok := v.Interface().([]string); ok {
			for _, v := range stringSlice {
				err := verifyLen(reflect.ValueOf(v), r)
				if err != nil {
					return err
				}
			}
			return nil
		}
		return ErrUnsupportedType
	default:
		return ErrUnsupportedType
	}

	return nil
}

// VerifyRegExp tells whether string matches given pattern or not
// Works for strings and slices of strings
// For slices works for each value in it.
func verifyRegExp(v reflect.Value, r string) error {
	kind := v.Kind()

	parts := strings.SplitN(r, ":", 2)
	if len(parts) != 2 {
		return ErrWrongValue
	}

	rExp, err := regexp.Compile(parts[1])
	if err != nil {
		return ErrWrongValue
	}

	switch {
	case kind == reflect.String:
		strVal := v.String()
		ok := rExp.MatchString(strVal)
		if !ok {
			return ErrNotRegexp
		}
	case kind == reflect.Slice:
		if stringSlice, ok := v.Interface().([]string); ok {
			for _, v := range stringSlice {
				if err := verifyRegExp(reflect.ValueOf(v), r); err != nil {
					return err
				}
			}
			return nil
		}
		return ErrUnsupportedType
	default:
		return ErrUnsupportedType
	}

	return nil
}

// verifyInInt tells whether the value belongs to the set.

func verifyInInt(v reflect.Value, r string) error {
	set := strings.Split(r[3:], ",")
	intVal := v.Int()

	for _, v := range set {
		i, err := strconv.ParseInt(v, 10, 0)
		if err != nil {
			return ErrWrongValue
		}
		if intVal == i {
			return nil
		}
	}

	return ErrOutOfSet
}

// verifyInString tells whether the value belongs to the set
func verifyInString(v reflect.Value, r string) error {
	set := strings.Split(r[3:], ",")
	strVal := v.String()

	for _, s := range set {
		if s == strVal {
			return nil
		}
	}

	return ErrOutOfSet
}

// verifyIn checks if the value belongs to the set based on its type
func verifyIn(v reflect.Value, r string) error {
	kind := v.Kind()

	switch kind {
	case reflect.String:
		return verifyInString(v, r)
	case reflect.Int:
		return verifyInInt(v, r)
	case reflect.Slice:
		if stringSlice, ok := v.Interface().([]string); ok {
			for _, v := range stringSlice {
				if err := verifyInString(reflect.ValueOf(v), r); err != nil {
					return err
				}
			}
			return nil
		}
		if intSlice, ok := v.Interface().([]int); ok {
			for _, v := range intSlice {
				if err := verifyInInt(reflect.ValueOf(v), r); err != nil {
					return err
				}
			}
			return nil
		}
	default:
		return ErrUnsupportedType
	}
	return ErrUnsupportedType
}

func VerifyMinOrMax(v reflect.Value, r string) error {
	kind := v.Kind()

	switch {
	case kind == reflect.Int:
		num, err := strconv.ParseInt(r[4:], 0, 0)
		if err != nil {
			return ErrWrongValue
		}
		switch {
		case r[:3] == "min":
			if v.Int() < num {
				return ErrValueOutOfBondary
			}
			return nil
		case r[:3] == "max":
			if v.Int() > num {
				return ErrValueOutOfBondary
			}
			return nil
		}
	case kind == reflect.Slice:
		if intSlice, ok := v.Interface().([]int); ok {
			for _, v := range intSlice {
				if err := VerifyMinOrMax(reflect.ValueOf(v), r); err != nil {
					return err
				}
			}
			return nil
		}
	}

	return ErrUnsupportedType
}
