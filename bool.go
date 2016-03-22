package jsonvalid

type BoolValidator struct {
	defaultValue bool
	required bool
}

func (v *BoolValidator) clone() *BoolValidator {
	return &BoolValidator{
		defaultValue: v.defaultValue,
		required: v.required,
	}
}

func (v *BoolValidator) Required() *BoolValidator {
	nv := v.clone()
	nv.required = true
	return nv
}

func (v *BoolValidator) Default(value bool) *BoolValidator {
	nv := v.clone()
	nv.defaultValue = value
	return nv
}

func (v *BoolValidator) Validate(path Path, value interface{}) (interface{}, error) {
	// Test if the value is nil, in which case we can short-circuit to checking
	// if the value is required.
	if value == nil {
		if v.required {
			return nil, ValidationErrorAtPath(path, ValueErrorRequired)
		}

		return v.defaultValue, nil
	}

	// Test if the value is a boolean.
	boolValue, ok := value.(bool)
	if !ok {
		return nil, ValidationErrorAtPath(path, ValueError{
			Code: "invalid_type",
			Message: "Value must be a boolean",
		})
	}

	return boolValue, nil
}

func Bool() *BoolValidator {
	return &BoolValidator{}
}
