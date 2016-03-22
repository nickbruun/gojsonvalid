package jsonvalid

import (
	"errors"
)

// Value error.
type ValueError struct {
	// Error code.
	Code string `json:"code"`

	// Humanly readable error message.
	Message string `json:"message"`
}

var (
	// Required value error.
	ValueErrorRequired = ValueError{
		Code: "required",
		Message: "This field is required",
	}

	// Invalid value error.
	ValueErrorInvalidValue = ValueError{
		Code: "invalid",
		Message: "Invalid value",
	}

	// Parse error.
	//
	// Represents an error parsing a requests' body.
	ErrParseError = errors.New("parse error")
)

// Field error.
//
// Pathed value error.
type FieldError struct {
	Path Path `json:"path"`
	ValueError
}

// Validation error.
type ValidationError struct {
	Fields []FieldError
}

// Concatenate with other validation error.
func (e *ValidationError) Concat(o *ValidationError) *ValidationError {
	return &ValidationError{
		Fields: append(e.Fields, o.Fields...),
	}
}

func (e *ValidationError) Error() string {
	return "validation error"
}

// Validation error at path.
//
// Short hand for creating a validation error with a specific value error at
// a given path.
func ValidationErrorAtPath(path Path, valueError ValueError) *ValidationError {
	return &ValidationError{
		Fields: []FieldError{{path, valueError}},
	}
}
