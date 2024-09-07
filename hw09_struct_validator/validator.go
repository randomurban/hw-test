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
		field := vv.Field(i)
		rules, err := parseRules(tag)
		if err != nil {
			return fmt.Errorf("parsing %v: %w", vt.Name(), err)
		}
		switch field.Kind() {
		case reflect.String:
			for _, rule := range rules {
				switch rule.name {
				case "len":
					if field.Len() != rule.paramNum {
						res = append(res,
							ValidationError{
								vt.Field(i).Name,
								fmt.Errorf("len must be %v", rule.paramNum),
							})
					}
				default:
					return fmt.Errorf("invalid rule: %v", rule.name)
				}
			}
		default:
			return fmt.Errorf("%v is not supported type", field.Kind())
		}
	}
	if len(res) > 0 {
		return res
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
					return nil, fmt.Errorf("must be a number param for %s: %w", name, err)
				}
			}
			res = append(res, Rule{name, param, paramNum})
		default:
			return nil, fmt.Errorf("invalid tag: %v", tag)
		}
	}

	return res, nil
}
