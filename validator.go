package jsonvalid

// Validator.
type Validator interface {
	Validate(path Path, value interface{}) (interface{}, error)
}
