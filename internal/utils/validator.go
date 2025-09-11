package utils

import (
	"regexp"
	"slices"
)

type ValidationErrors map[string]string

type Validator struct {
	Errors ValidationErrors
}

func NewValidator() *Validator {
	return &Validator{
		Errors: make(map[string]string),
	}
}

func (v *Validator) add(field, errorMsg string) {
	if _, ok := v.Errors[field]; !ok {
		v.Errors[field] = errorMsg
	}
}

func (v *Validator) Must(ok bool, field, errorMsg string) {
	if !ok {
		v.add(field, errorMsg)
	}
}

func (v *Validator) In(text string, values []string, field, errorMsg string) {
	if !slices.Contains(values, text) {
		v.add(field, errorMsg)
	}
}

func (v *Validator) Regex(text string, pattern *regexp.Regexp, field, errorMsg string) {
	if !pattern.MatchString(text) {
		v.add(field, errorMsg)
	}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}
