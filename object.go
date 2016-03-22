package jsonvalid

import (
	"encoding/json"
	"net/http"
	"io/ioutil"
	"reflect"
)

type ObjectValidator struct {
	required bool
	props map[string]*ObjectProp
	targetType reflect.Type
}

func (v *ObjectValidator) clone() *ObjectValidator {
	return &ObjectValidator{
		required: v.required,
		props: v.props,
		targetType: v.targetType,
	}
}

func (v *ObjectValidator) Required() *ObjectValidator {
	nv := v.clone()
	nv.required = true
	return nv
}

func (v *ObjectValidator) UnmarshalTo(typ interface{}) *ObjectValidator {
	nv := v.clone()

	if refTyp, ok := typ.(reflect.Type); ok {
		nv.targetType = refTyp
	} else {
		nv.targetType = reflect.TypeOf(typ)
	}

	return nv
}

func (v *ObjectValidator) ParseAndValidateHttpRequest(req *http.Request) (interface{}, error) {
	// First, read the request body.
	data, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, ErrParseError
	}

	// Parse the form as a JSON object.
	var jsonObj map[string]interface{}
	if err = json.Unmarshal(data, &jsonObj); err != nil {
		return nil, ErrParseError
	}

	return v.Validate("", jsonObj)
}

func (v *ObjectValidator) Validate(path Path, value interface{}) (interface{}, error) {
	if value == nil {
		if v.required {
			return nil, ValidationErrorAtPath(path, ValueErrorRequired)
		}

		return nil, nil
	}

	// Validate that the value is an object.
	jsonObj, ok := value.(map[string]interface{})
	if !ok {
		return nil, ValidationErrorAtPath(path, ValueError{
			Code: "invalid_type",
			Message: "Field must be an object",
		})
	}

	// Validate each provided property.
	var err *ValidationError
	result := make(map[string]interface{})

	for propName, propValue := range jsonObj {
		var resultValue interface{}
		var resultErr error

		prop, ok := v.props[propName]
		if !ok {
			resultErr = ValidationErrorAtPath(path.Prop(propName), ValueError{
				Code: "invalid_property",
				Message: "Invalid property",
			})
		} else {
			resultValue, resultErr = prop.Validator().Validate(path.Prop(propName), propValue)
		}

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

		result[propName] = resultValue
	}

	// Validate all properties not provided in the object.
	for propName, prop := range v.props {
		if _, handled := result[propName]; !handled {
			resultValue, resultErr := prop.Validator().Validate(path.Prop(propName), nil)

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

			result[propName] = resultValue
		}
	}

	// Bail if there are any errors.
	if err != nil {
		return nil, err
	}

	// Unmarshal if necessary.
	if v.targetType == nil {
		return result, nil
	}

	return unmarshalObject(result, v.targetType)
}

func Object(props ...*ObjectProp) *ObjectValidator {
	propMap := make(map[string]*ObjectProp, len(props))

	for _, prop := range props {
		propMap[prop.Name()] = prop
	}
	
	return &ObjectValidator{
		props: propMap,
	}
}

type ObjectProp struct {
	name string
	validator Validator
}

func (p *ObjectProp) Name() string {
	return p.name
}

func (p *ObjectProp) Validator() Validator {
	return p.validator
}

func Prop(name string, validator Validator) *ObjectProp {
	return &ObjectProp{
		name: name,
		validator: validator,
	}
}
