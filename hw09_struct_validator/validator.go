package hw09structvalidator

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("%v: %v", e.Field, e.Err)
}

func (e *ValidationError) Unwrap() error { return e.Err }

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	if len(v) == 0 {
		return ""
	}
	var res strings.Builder
	for i, err := range v {
		res.WriteString(err.Error())

		if i < len(v)-1 {
			res.WriteString(",")
		}
	}
	return res.String()
}

func (v ValidationErrors) Unwrap() []error {
	if len(v) == 0 {
		return nil
	}
	var errs []error
	for _, err := range v {
		errs = append(errs, err.Err)
	}
	return errs
}

func Validate(v interface{}) error {
	vt := reflect.TypeOf(v)
	vv := reflect.ValueOf(v)
	if vt.Kind() != reflect.Struct {
		return fmt.Errorf("%v is not a struct", vt.Name())
	}
	res := ValidationErrors{}
	for i := 0; i < vt.NumField(); i++ {

		tag := vt.Field(i).Tag.Get("validate")
		if tag == "" {
			continue
		}
		rules := parseRules(tag)
		for _, rule := range rules {
			switch rule.name {
			case "len":
				if err := lenValidate(vv.Field(i), rule.param); err != nil {
					res = append(res, ValidationError{vt.Field(i).Name, err})
				}
			}
		}
	}
	if len(res) > 0 {
		return res
	}
	return nil
}

type Rule struct {
	name  string
	param string
}

// var ErrInvalidTag = errors.New("invalid tag")

func parseRules(tag string) []Rule {
	var res []Rule
	for _, rule := range strings.Split(tag, "|") {
		parts := strings.Split(rule, ":")
		if len(parts) != 2 {
			continue
		}
		res = append(res, Rule{parts[0], parts[1]})
	}

	return res
}

func lenValidate(v reflect.Value, param string) error {
	lenParam, err := strconv.Atoi(param)
	if err != nil {
		return err
	}
	switch v.Kind() {
	case reflect.String:
		if v.Len() != lenParam {
			return fmt.Errorf("len must be %v", lenParam)
		}
		return nil
	default:
		return fmt.Errorf("%v is not a valid type", v.Kind())

	}
}
