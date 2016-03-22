package jsonvalid

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type IntValueValidator func(value int) (int, *ValueError)

type IntValidator struct {
	defaultValue int
	required bool
	valueValidators []IntValueValidator
}

func (v *IntValidator) clone() *IntValidator {
	return &IntValidator{
		defaultValue: v.defaultValue,
		required: v.required,
		valueValidators: v.valueValidators,
	}
}

func (v *IntValidator) Required() *IntValidator {
	nv := v.clone()
	nv.required = true
	return nv
}

func (v *IntValidator) Default(value int) *IntValidator {
	nv := v.clone()
	nv.defaultValue = value
	return nv
}

func (v *IntValidator) Validate(path Path, value interface{}) (interface{}, error) {
	// Test if the value is nil, in which case we can short-circuit to checking
	// if the value is required.
	if value == nil {
		if v.required {
			return nil, ValidationErrorAtPath(path, ValueErrorRequired)
		}

		return v.defaultValue, nil
	}

	// Test if the value is a integer.
	var intValue int

	switch tv := value.(type) {
	case float64:
		intValue = int(tv)
	case json.Number:
		int64Value, err := tv.Int64()
		if err != nil {
			return nil, ValidationErrorAtPath(path, ValueError{
				Code: "invalid_type",
				Message: "Value must be an integer",
			})
		}

		intValue = int(int64Value) // TODO: overflow checking.
	case string:
		var err error
		if intValue, err = strconv.Atoi(tv); err != nil {
			return nil, ValidationErrorAtPath(path, ValueError{
				Code: "invalid_type",
				Message: "Value must be an integer",
			})
		}
	default:
		return nil, ValidationErrorAtPath(path, ValueError{
			Code: "invalid_type",
			Message: "Value must be an integer",
		})
	}

	// Validate the value.
	for _, valueValidator := range v.valueValidators {
		var err *ValueError
		if intValue, err = valueValidator(intValue); err != nil {
			return intValue, ValidationErrorAtPath(path, *err)
		}
	}

	return intValue, nil
}

func (v *IntValidator) ValidateValue(svv IntValueValidator) *IntValidator {
	nv := v.clone()
	nv.valueValidators = append(nv.valueValidators, svv)
	return nv
}

func (v *IntValidator) OneOf(values ...int) *IntValidator {
	valueSet := make(map[int]struct{}, len(values))

	for _, value := range values {
		valueSet[value] = struct{}{}
	}

	return v.ValidateValue(func(value int) (int, *ValueError) {
		if _, ok := valueSet[value]; !ok {
			return value, &ValueError{
				Code: "invalid",
				Message: "Invalid value",
			}
		}

		return value, nil
	})
}

func (v *IntValidator) Min(minValue int) *IntValidator {
	return v.ValidateValue(func(value int) (int, *ValueError) {
		if value < minValue {
			return value, &ValueError{
				Code: "invalid",
				Message: fmt.Sprintf("Value must be at least %d", minValue),
			}
		}

		return value, nil
	})
}

func (v *IntValidator) Max(maxValue int) *IntValidator {
	return v.ValidateValue(func(value int) (int, *ValueError) {
		if value > maxValue {
			return value, &ValueError{
				Code: "invalid",
				Message: fmt.Sprintf("Value must be at most %d", maxValue),
			}
		}

		return value, nil
	})
}

func Int() *IntValidator {
	return &IntValidator{}
}
