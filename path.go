package jsonvalid

import (
	"fmt"
)

// Path.
//
// Represents a field path in a JSON structure.
type Path string

// Property of path.
//
// Returns a path to a property of the path.
func (p Path) Prop(name string) Path {
	if p == "" {
		return Path(name)
	}

	return Path(string(p) + "." + name)
}

// Element of path.
//
// Returns a path to an element of the path, ie. an array index.
func (p Path) Elem(index int) Path {
	return Path(fmt.Sprintf("%s[%d]", p, index))
}
