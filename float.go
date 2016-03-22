package jsonvalid

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Float64ValueValidator func(value float64) (float64, *ValueError)

type Float64Validator struct {
	defaultValue float64
	required bool
	valueValidators []Float64ValueValidator
}

func (v *Float64Validator) clone() *Float64Validator {
	return &Float64Validator{
		defaultValue: v.defaultValue,
		required: v.required,
		valueValidators: v.valueValidators,
	}
}

func (v *Float64Validator) Required() *Float64Validator {
	nv := v.clone()
	nv.required = true
	return nv
}

func (v *Float64Validator) Default(value float64) *Float64Validator {
	nv := v.clone()
	nv.defaultValue = value
	return nv
}

func (v *Float64Validator) Validate(path Path, value interface{}) (interface{}, error) {
	// Test if the value is nil, in which case we can short-circuit to checking
	// if the value is required.
	if value == nil {
		if v.required {
			return nil, ValidationErrorAtPath(path, ValueErrorRequired)
		}

		return v.defaultValue, nil
	}

	// Test if the value is a floating point number.
	var floatValue float64

	switch tv := value.(type) {
	case float64:
		floatValue = tv
	case json.Number:
		var err error
		floatValue, err = tv.Float64()
		if err != nil {
			return nil, ValidationErrorAtPath(path, ValueError{
				Code: "invalid_type",
				Message: "Value must be a floating point number",
			})
		}
	case string:
		var err error
		if floatValue, err = strconv.ParseFloat(tv, 64); err != nil {
			return nil, ValidationErrorAtPath(path, ValueError{
				Code: "invalid_type",
				Message: "Value must be a floating point number",
			})
		}
	default:
		return nil, ValidationErrorAtPath(path, ValueError{
			Code: "invalid_type",
			Message: "Value must be a floating point number",
		})
	}

	// Validate the value.
	for _, valueValidator := range v.valueValidators {
		var err *ValueError
		if floatValue, err = valueValidator(floatValue); err != nil {
			return floatValue, ValidationErrorAtPath(path, *err)
		}
	}

	return floatValue, nil
}

func (v *Float64Validator) ValidateValue(svv Float64ValueValidator) *Float64Validator {
	nv := v.clone()
	nv.valueValidators = append(nv.valueValidators, svv)
	return nv
}

func (v *Float64Validator) Min(minValue float64) *Float64Validator {
	return v.ValidateValue(func(value float64) (float64, *ValueError) {
		if value < minValue {
			return value, &ValueError{
				Code: "invalid",
				Message: fmt.Sprintf("Value must be at least %d", minValue),
			}
		}

		return value, nil
	})
}

func (v *Float64Validator) Max(maxValue float64) *Float64Validator {
	return v.ValidateValue(func(value float64) (float64, *ValueError) {
		if value > maxValue {
			return value, &ValueError{
				Code: "invalid",
				Message: fmt.Sprintf("Value must be at most %d", maxValue),
			}
		}

		return value, nil
	})
}

func Float64() *Float64Validator {
	return &Float64Validator{}
}
