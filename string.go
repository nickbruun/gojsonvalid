package jsonvalid

import (
	"github.com/nickbruun/goinput"
	"strings"
	"regexp"
)

type StringValueValidator func(value string) (string, *ValueError)

type StringValidator struct {
	required bool
	valueValidators []StringValueValidator
}

func (v *StringValidator) clone() *StringValidator {
	return &StringValidator{
		required: v.required,
		valueValidators: v.valueValidators,
	}
}

func (v *StringValidator) ValidateValue(svv StringValueValidator) *StringValidator {
	nv := v.clone()
	nv.valueValidators = append(nv.valueValidators, svv)
	return nv
}

func (v *StringValidator) Required() *StringValidator {
	nv := v.clone()
	nv.required = true
	return nv
}

func (v *StringValidator) Strip() *StringValidator {
	return v.ValidateValue(func(value string) (string, *ValueError) {
		return strings.TrimSpace(value), nil
	})
}

func (v *StringValidator) StripSingleLine() *StringValidator {
	return v.ValidateValue(func(value string) (string, *ValueError) {
		return input.TrimWhitespaceNormalizeLine(value), nil
	})
}

func (v *StringValidator) OneOf(values ...string) *StringValidator {
	valueSet := make(map[string]struct{}, len(values))

	for _, value := range values {
		valueSet[value] = struct{}{}
	}

	return v.ValidateValue(func(value string) (string, *ValueError) {
		if _, ok := valueSet[value]; !ok {
			return value, &ValueError{
				Code: "invalid",
				Message: "Invalid value",
			}
		}

		return value, nil
	})
}

func (v *StringValidator) Matches(expr string) *StringValidator {
	return v.MatchesRegexp(regexp.MustCompile(expr))
}

func (v *StringValidator) MatchesRegexp(re *regexp.Regexp) *StringValidator {
	return v.ValidateValue(func(value string) (string, *ValueError) {
		if !re.MatchString(value) {
			return value, &ValueError{
				Code: "invalid",
				Message: "Invalid value",
			}
		}

		return value, nil
	})
}

func (v *StringValidator) Validate(path Path, value interface{}) (interface{}, error) {
	// Test if the value is nil, in which case we can short-circuit to checking
	// if the value is required.
	if value == nil {
		if v.required {
			return "", ValidationErrorAtPath(path, ValueErrorRequired)
		}

		return "", nil
	}

	// Test if the value is a string.
	strValue, ok := value.(string)
	if !ok {
		return "", ValidationErrorAtPath(path, ValueError{
			Code: "invalid_type",
			Message: "Value must be a string",
		})
	}

	// Apply the validators.
	for _, svv := range v.valueValidators {
		var valErr *ValueError
		if strValue, valErr = svv(strValue); valErr != nil {
			return "", ValidationErrorAtPath(path, *valErr)
		}
	}

	// Reevaluate if the cleaned value is required.
	if strValue == "" && v.required {
		return "", ValidationErrorAtPath(path, ValueErrorRequired)
	}

	return strValue, nil
}

func String() *StringValidator {
	return &StringValidator{}
}
