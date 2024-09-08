package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
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
			res.WriteString(", ")
		}
	}
	return res.String()
}

func (v ValidationErrors) Unwrap() (errs []error) {
	if len(v) == 0 {
		return nil
	}
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
		field := vv.Field(i)
		fieldName := vt.Field(i).Name
		rules, err := parseRules(tag)
		if err != nil {
			res = append(res, ValidationError{
				Field: fieldName,
				Err:   fmt.Errorf("parsing %v: %w", vt.Name(), err),
			})
		}
		if field.Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				item := field.Index(j)
				res = checkRule(rules, item, res, fieldName)
			}
		} else {
			res = checkRule(rules, field, res, fieldName)
		}
	}
	if len(res) > 0 {
		return res
	}
	return nil
}

func checkRule(rules []Rule, field reflect.Value, res ValidationErrors, fieldName string) ValidationErrors {
	var errValidate error
	for _, rule := range rules {
		if field.Kind() == reflect.String {
			errValidate = validateString(field.String(), rule)
		}
		if field.Kind() == reflect.Int {
			errValidate = validateInt(field.Interface().(int), rule)
		}
		if errValidate != nil {
			res = append(res, ValidationError{
				fieldName,
				errValidate,
			})
		}
	}
	return res
}

func validateString(field string, rule Rule) error {
	switch rule.name {
	case "len":
		if len(field) != rule.paramNum {
			return fmt.Errorf("len must be %v", rule.paramNum)
		}
	case "in":
		paramList := strings.Split(rule.param, ",")
		for _, param := range paramList {
			if field == param {
				return nil
			}
		}
		return fmt.Errorf("%v in [%v] is required", field, rule.param)
	case "regexp":
		matched, errReg := regexp.MatchString(rule.param, field)
		if errReg != nil {
			return fmt.Errorf("regexp err: %w", errReg)
		}
		if !matched {
			return fmt.Errorf("'%v' not matched", field)
		}
	default:
		return fmt.Errorf("invalid rule: %v", rule.name)
	}
	return nil
}

func validateInt(field int, rule Rule) error {
	switch rule.name {
	case "in":
		paramList := strings.Split(rule.param, ",")
		for _, param := range paramList {
			paramInt, paramErr := strconv.Atoi(param)
			if paramErr != nil {
				return fmt.Errorf("invalid int tag element: %v", param)
			}
			if field == paramInt {
				return nil
			}
		}
		return fmt.Errorf("%v in [%v] is required", field, rule.param)
	case "min":
		if field < rule.paramNum {
			return fmt.Errorf("must be minimum %v", rule.paramNum)
		}
	case "max":
		if field > rule.paramNum {
			return fmt.Errorf("must be maximum %v", rule.paramNum)
		}
	default:
		return fmt.Errorf("invalid rule: %v", rule.name)
	}
	return nil
}

type Rule struct {
	name     string
	param    string
	paramNum int
}

// var ErrInvalidTag = errors.New("invalid tag")

func parseRules(tag string) (res []Rule, err error) {
	for _, rule := range strings.Split(tag, "|") {
		parts := strings.Split(rule, ":")
		switch len(parts) {
		case 1:
			res = append(res, Rule{strings.TrimSpace(parts[0]), "", 0})
		case 2:
			name := strings.TrimSpace(parts[0])
			param := strings.TrimSpace(parts[1])
			var paramNum int
			if name == "len" || name == "min" || name == "max" {
				paramNum, err = strconv.Atoi(param)
				if err != nil {
					return nil, fmt.Errorf("number param for %s: %w", name, err)
				}
			}
			res = append(res, Rule{name, param, paramNum})
		default:
			return nil, fmt.Errorf("invalid tag: %v", tag)
		}
	}

	return res, nil
}
