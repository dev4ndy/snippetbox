package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) Invalid() bool {
	return !v.Valid()
}

func (v *Validator) AddFieldError(key string, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key string, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func Required(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxLength(value string, max int) bool {
	return utf8.RuneCountInString(value) <= max
}

func AllowedValues(value int, allowedValues ...int) bool {
	for _, val := range allowedValues {
		if value == val {
			return true
		}
	}
	return false
}
