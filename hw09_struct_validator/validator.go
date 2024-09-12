package hw09structvalidator

import (
	"errors"
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
	var checkRes ValidationError
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
			return fmt.Errorf("parsing %v: %w", vt.Name(), err)
		}
		if field.Kind() == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				item := field.Index(j)
				checkRes, err = checkRule(rules, item, fieldName)
				if err != nil {
					return fmt.Errorf("parsing %v: %w", vt.Name(), err)
				}
				res = append(res, checkRes)
			}
		} else {
			checkRes, err = checkRule(rules, field, fieldName)
			if err != nil {
				return fmt.Errorf("parsing %v: %w", vt.Name(), err)
			}
			res = append(res, checkRes)
		}
	}
	if len(res) > 0 {
		return res
	}
	return nil
}

func checkRule(rules []Rule, field reflect.Value, fieldName string) (ValidationError, error) {
	var errValidate error
	var errParsing ParsingError
	var res ValidationError
	for _, rule := range rules {
		if field.Kind() == reflect.String {
			errValidate = validateString(field.String(), rule)
		}
		if field.Kind() == reflect.Int {
			errValidate = validateInt(field.Interface().(int), rule)
		}
		if errValidate != nil {
			if errors.As(errValidate, &errParsing) {
				return ValidationError{}, errParsing
			}
			res = ValidationError{
				fieldName,
				errValidate,
			}
		}
	}
	return res, nil
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
		if !rule.re.MatchString(field) {
			return fmt.Errorf("'%v' not matched", field)
		}
	default:
		return ParsingError{"illegal tag in string: " + rule.name}
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
				return ParsingError{"invalid int tag element: " + param}
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
		return ParsingError{"illegal tag in int: " + rule.name}
	}
	return nil
}

type Rule struct {
	name     string
	param    string
	paramNum int
	re       *regexp.Regexp
}

type ParsingError struct {
	s string
}

func (e ParsingError) Error() string {
	return e.s
}

func parseRules(tag string) (res []Rule, err error) {
	for _, rule := range strings.Split(tag, "|") {
		parts := strings.Split(rule, ":")
		name := strings.TrimSpace(parts[0])
		var param string
		var re *regexp.Regexp
		var paramNum int

		if len(parts) != 2 {
			return nil, ParsingError{"invalid tag: " + tag}
		}
		param = strings.TrimSpace(parts[1])

		switch name {
		case "in":

		case "len", "min", "max":
			paramNum, err = strconv.Atoi(param)
			if err != nil {
				return nil, ParsingError{"invalid number param for " + name + ": " + err.Error()}
			}

		case "regexp":
			var reErr error
			re, reErr = regexp.Compile(param)
			if reErr != nil {
				return nil, reErr
			}

		default:
			return nil, ParsingError{"invalid tag: " + tag}
		}
		res = append(res, Rule{name, param, paramNum, re})
	}

	return res, nil
}
