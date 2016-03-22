package jsonvalid

import (
	"testing"
)

var (
	optionalIntOptionalArrayValidator = ArrayOf(Int())
	optionalIntRequiredArrayValidator = ArrayOf(Int()).Required()
)

func AssertArrayValidationFails(t *testing.T, valueDesc, validatorDesc string, validator Validator, value interface{}) {
	_, err := validator.Validate("", value)

	if err == nil {
		t.Errorf("expected error form validating %s with %s", valueDesc, validatorDesc)
	}
}

func AssertArrayValidationResult(t *testing.T, valueDesc, validatorDesc string, validator Validator, value interface{}, expected []interface{}) {
	result, err := validator.Validate("", value)

	if err != nil {
		t.Errorf("unexpected error validating %s with %s: %v", valueDesc, validatorDesc, err)
		return
	}

	array, ok := result.([]interface{})
	if !ok {
		t.Errorf("return value from validating %s with %s is not an array: %v", valueDesc, validatorDesc, result)
		return
	}

	if len(array) != len(expected) {
		t.Errorf("expected return value from validating %s with %s to be %v but it is: %v", valueDesc, validatorDesc, expected, array)
		return
	}
}

func TestArray(t *testing.T) {
	// Test validating null value.
	AssertArrayValidationResult(
		t,
		"null value",
		"optional array validator",
		optionalIntOptionalArrayValidator,
		nil,
		[]interface{}{},
	)
	AssertArrayValidationFails(
		t,
		"null value",
		"required array validator",
		optionalIntRequiredArrayValidator,
		nil,
	)

	// Test validating empty array.
	AssertArrayValidationResult(
		t,
		"empty array",
		"optional array validator",
		optionalIntOptionalArrayValidator,
		[]interface{}{},
		[]interface{}{},
	)
	AssertArrayValidationFails(
		t,
		"empty array",
		"required array validator",
		optionalIntRequiredArrayValidator,
		[]interface{}{},
	)

	// Test validating array with values.
	AssertArrayValidationResult(
		t,
		"array with values",
		"optional array validator",
		optionalIntOptionalArrayValidator,
		[]interface{}{nil, 1, 5},
		[]interface{}{nil, 1, 5},
	)
	AssertArrayValidationResult(
		t,
		"array with values",
		"required array validator",
		optionalIntOptionalArrayValidator,
		[]interface{}{nil, 1, 5},
		[]interface{}{nil, 1, 5},
	)
}
