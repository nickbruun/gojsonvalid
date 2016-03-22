package jsonvalid

import (
	"fmt"
	"reflect"
	"encoding/json"
)

func unmarshalObject(marshaled map[string]interface{}, typ reflect.Type) (interface{}, error) {
	// Resolve the actual struct type.
	unmarshalToPtr := false

	if typ.Kind() == reflect.Ptr {
		unmarshalToPtr = true
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return nil, fmt.Errorf("cannot unmarshal object to %v", typ)
	}

	// Initialize an object representation.
	unmarshaled := reflect.New(typ)

	// Unmarshal.
	//
	// Really lazy for now.
	marshaledBytes, err := json.Marshal(marshaled)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(marshaledBytes, unmarshaled.Interface()); err != nil {
		return nil, err
	}

	// Return the unmarshaled object.
	if !unmarshalToPtr {
		return unmarshaled.Elem().Interface(), nil
	}

	return unmarshaled.Interface(), nil
}
