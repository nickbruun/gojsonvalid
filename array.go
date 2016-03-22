package jsonvalid

import (
	"fmt"
)

type ArrayValidator struct {
	required bool
	minLen int
	itemValidator Validator
}

func (v *ArrayValidator) clone() *ArrayValidator {
	return &ArrayValidator{
		required: v.required,
		minLen: v.minLen,
		itemValidator: v.itemValidator,
	}
}

func (v *ArrayValidator) Required() *ArrayValidator {
	nv := v.clone()
	nv.required = true
	return nv
}

func (v *ArrayValidator) Of(validator Validator) *ArrayValidator {
	nv := v.clone()
	nv.itemValidator = validator
	return nv
}

func (v *ArrayValidator) Validate(path Path, value interface{}) (interface{}, error) {
	// Test if the value is nil, in which case we can short-circuit to checking
	// if the value is required.
	if value == nil {
		if v.required {
			return nil, ValidationErrorAtPath(path, ValueErrorRequired)
		}

		return []interface{}{}, nil
	}

	// Test if the value is an array.
	arrValue, ok := value.([]interface{})
	if !ok {
		return nil, ValidationErrorAtPath(path, ValueError{
			Code: "invalid_type",
			Message: "Value must be an array",
		})
	}

	if len(arrValue) == 0 && v.required {
		return nil, ValidationErrorAtPath(path, ValueErrorRequired)
	}

	// Construct a result.
	var err *ValidationError
	result := make([]interface{}, 0, len(arrValue))

	for i, elemValue := range arrValue {
		resultValue, resultErr := v.itemValidator.Validate(path.Elem(i), elemValue)

		if resultErr != nil {
			if resultValidationErr, ok := resultErr.(*ValidationError); ok {
				if err == nil {
					err = resultValidationErr
				} else {
					err = err.Concat(resultValidationErr)
				}
			} else {
				return nil, resultErr
			}
		}

		result = append(result, resultValue)
	}

	if len(result) < v.minLen {
		return nil, ValidationErrorAtPath(path, ValueError{
			Code: "invalid",
			Message: fmt.Sprintf("Value must be an array of at least %d element(s)", v.minLen),
		})
	}

	if err == nil {
		return result, nil
	} else {
		return result, err
	}
}

func (v *ArrayValidator) MinLen(minLen int) *ArrayValidator {
	nv := v.clone()
	nv.minLen = minLen
	return nv
}

func Array() *ArrayValidator {
	return &ArrayValidator{}
}

func ArrayOf(of Validator) *ArrayValidator {
	return Array().Of(of)
}
