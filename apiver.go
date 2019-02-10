package apiver

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

// InternalError represents an code error, such as malformed struct tags, or passing a value instead of a pointer.
// For public services, such as HTTP services, these errors should typically be logged and considered Internal Server Errors, rather than returned to the user, for security reasons.
type InternalError struct {
	Msg string
}

func (e InternalError) Error() string { return e.Msg }

// UserError is a decode error caused by invalid input. These error messages are designed to be user-friendly and safe to return to users. JSON tags rather than internal struct names are returned, and passed user data is not reflected back.
type UserError struct {
	Msg string
}

func (e UserError) Error() string { return e.Msg }

// RejectUnknownFields is whether to fail to parse JSON with unknown fields. This includes fields which exist in the struct at a later version than is being parsed.
const RejectUnknownFields = true // TODO implement

// TagName is the name of the tag to use for parsing versions and string-primitives.
const TagName = `api`

// TagPropertyStr is the name of the tag property to use for parsing numbers and booleans as strings.
const TagPropertyStr = `str`

// UnmarshalJSON parses JSON for the given object.
// bts is the JSON bytes.
// realObj is the object to unmarshal into.
// version is the object version being used. Fields in the object with a newer version than version must be pointers, and will not be deserialized into, even if the field exists in the JSON in bts. This is to preserve Semantic Versioning.
func UnmarshalJSON(bts []byte, realObj interface{}, version float64) error {
	// TODO add option to reject any realObj with a field missing a tc:version tag
	// TODO add option to reject any bts with fields not in realObj - https://golang.org/pkg/encoding/json/#Decoder.DisallowUnknownFields

	obj := reflect.ValueOf(realObj)
	if obj.Kind() != reflect.Ptr {
		return InternalError{"object must be a pointer"}
	}
	if obj.IsNil() {
		return InternalError{"object must not be nil"}
	}

	obj = reflect.Indirect(obj)

	if obj.Type().Kind() != reflect.Struct {
		return InternalError{"object must be a pointer to a struct"}
	}

	newVal := BuildUnmarshalObj(obj, version, true)

	newValI := newVal.Addr().Interface()

	if err := json.Unmarshal(bts, newValI); err != nil {
		return err
	}

	if err := FromUnmarshalObj(newVal, realObj); err != nil {
		return err
	}

	return nil
}

type TagProperties struct {
	// Version is the Traffic Ops API Version. If no version was present, this will be 0
	Version float64
	// Str is whether "str" existed, which indicates that a string should be parsed as a boolean or number.
	Str bool
}

// GetTagProperties returns the properties from the given tag. An empty string may be passed, which will indicate no version (therefore, all versions), and that the field should not accept a string for a number or boolean.
func GetTagProperties(tag string) TagProperties {
	props := TagProperties{}
	for _, prop := range strings.Split(tag, ",") {
		switch prop {
		case TagPropertyStr:
			props.Str = true
		default:
			if f, err := strconv.ParseFloat(prop, 64); err == nil {
				props.Version = f
			} else {
				// TODO log? Return error?
			}
		}
	}
	return props
}

// BuildUnmarshalObj creates an object to be serialized or deserialized into, from the given val, omitting versions newer than version, and dynamically creating types which will deserialize from strings for fields with TagName TagPropertyStr.
// The strTypes should always be false to build an object for unmarshalling into, and should always be true for building an object to marshal into bytes. This parameter exists, because the 'str' types use the largest possible type, and can lead to precision loss for smaller types like float32.
func BuildUnmarshalObj(val reflect.Value, version float64, strTypes bool) reflect.Value {
	newTyp := BuildUnmarshalType(val.Type(), version, strTypes)
	return reflect.New(newTyp).Elem()
}

// Create the object to be passed to encoding/json.Unmarshal.
// This creates a new struct which:
// 1. converts all values to pointers, so we can distinguish missing from default values
// 2. removes any fields newer than version
// 3. converts "str" fields to types which will deserialize as strings or their real type (int,float.bool)
//
func BuildUnmarshalType(typ reflect.Type, version float64, strTypes bool) reflect.Type {
	// TODO error if val has non-pointer fields newer than version (which can never be filled, but must be filled - ergo all non-base versions must be pointers to make any sense)

	newTypeFields := []reflect.StructField{}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		props := GetTagProperties(field.Tag.Get(TagName))
		if props.Version > version {
			continue
		}

		newField := reflect.StructField{}
		newField.Name = field.Name
		newField.Type = field.Type
		if newField.Type.Kind() == reflect.Struct {
			newField.Type = BuildUnmarshalType(newField.Type, version, strTypes)
		}

		if newField.Type.Kind() != reflect.Ptr {
			// convert all fields to pointers
			// this lets us later verify value=required fields exist, and return an error if any value field is nil.
			// Without this, we can't distinguish empty from missing values.
			newField.Type = reflect.PtrTo(newField.Type)
		}

		if strTypes && props.Str {
			switch newField.Type.Elem().Kind() {
			case reflect.Bool:
				newField.Type = reflect.PtrTo(reflect.TypeOf(BoolS(false)))
			case reflect.Int:
				fallthrough
			case reflect.Int8:
				fallthrough
			case reflect.Int16:
				fallthrough
			case reflect.Int32:
				fallthrough
			case reflect.Int64:
				newField.Type = reflect.PtrTo(reflect.TypeOf(IntS(0)))
			case reflect.Uint:
				fallthrough
			case reflect.Uint8:
				fallthrough
			case reflect.Uint16:
				fallthrough
			case reflect.Uint32:
				fallthrough
			case reflect.Uint64:
				fallthrough
			case reflect.Uintptr:
				newField.Type = reflect.PtrTo(reflect.TypeOf(UIntS(0)))
			case reflect.Float32:
				fallthrough
			case reflect.Float64:
				newField.Type = reflect.PtrTo(reflect.TypeOf(FloatS(0)))
			default:
				// TODO error?
			}
		}

		newField.Tag = field.Tag

		newTypeFields = append(newTypeFields, newField)
	}

	return reflect.StructOf(newTypeFields)
}

func FromUnmarshalObj(val reflect.Value, realObj interface{}) error {
	obj := reflect.ValueOf(realObj)
	if obj.Kind() != reflect.Ptr {
		return InternalError{"object must be a pointer"}
	}
	if obj.IsNil() {
		return InternalError{"object must not be nil"}
	}
	return FromUnmarshalObjVals(val, obj)
}

// FromUnmarshalObj converts an object created with BuildUnmarshalObj, presumably after decoding data into it, into the real object.
// Returns an error if any value fields in the realObj are nil in the val.
func FromUnmarshalObjVals(val reflect.Value, obj reflect.Value) error {
	obj = reflect.Indirect(obj)

	valType := val.Type()
	for i := 0; i < val.NumField(); i++ {
		valTypeField := valType.Field(i)
		fieldName := valTypeField.Name
		valField := val.Field(i)

		objField := obj.FieldByName(valTypeField.Name)
		if objField == (reflect.Value{}) {
			return InternalError{"object missing field in val '" + valTypeField.Name + "'"}
		}

		if valField.IsNil() && objField.Type().Kind() != reflect.Ptr {
			// TODO remove/abstract json tag, so this func is fully encoder-agnostic?
			if jsonTag := valTypeField.Tag.Get("json"); jsonTag != "" {
				if jsonTagParts := strings.Split(jsonTag, ","); len(jsonTagParts) > 0 {
					fieldName = jsonTagParts[0]
				}
			}
			return UserError{"missing required field: " + fieldName} // TODO make missing-required-field an err type?
		}

		if objField.Type().Kind() == reflect.Ptr {
			if objField.IsNil() {
				objField.Set(reflect.New(objField.Type().Elem()))
			}
			objField = reflect.Indirect(objField)
		}
		if valField.Type().Kind() == reflect.Ptr {
			valField = reflect.Indirect(valField)
		}

		if valField.Type().Kind() == reflect.Struct {
			if objField.Type().Kind() != reflect.Struct {
				return InternalError{"val field '" + fieldName + "' type '" + objField.Type().String() + "' does not match struct type '" + valField.Type().String() + "'"}
			}
			if err := FromUnmarshalObjVals(valField, objField); err != nil {
				return errors.New("field '" + fieldName + "': " + err.Error()) // TODO use encoding/json errs? Somehow handle user vs system errs
			}
			continue
		}

		if valField.Type() == reflect.TypeOf(IntS(0)) || valField.Type() == reflect.TypeOf(UIntS(0)) || valField.Type() == reflect.TypeOf(FloatS(0)) || valField.Type() == reflect.TypeOf(BoolS(false)) {
			if !objField.Type().ConvertibleTo(valField.Type()) {
				return InternalError{"val field '" + fieldName + "' type '" + objField.Type().String() + "' is not convertible to object field type '" + valField.Type().String() + "'"}
			}
			valField = valField.Convert(objField.Type())
		}
		if !objField.Type().AssignableTo(valField.Type()) {
			return InternalError{"val field '" + fieldName + "' type '" + objField.Type().String() + "' is not assignable to object field type '" + valField.Type().String() + "'"}
		}
		if !objField.CanSet() {
			return InternalError{"can't set object field '" + fieldName + "'"}
		}

		objField.Set(valField)
	}

	return nil
}

func MarshalJSON(realObj interface{}, version float64) ([]byte, error) {
	obj, err := BuildMarshalObj(realObj, version)
	if err != nil {
		return nil, err
	}
	return json.Marshal(obj)
}

func MarshalJSONIndent(realObj interface{}, prefix, indent string, version float64) ([]byte, error) {
	obj, err := BuildMarshalObj(realObj, version)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(obj, indent, prefix)
}

func BuildMarshalObj(realObj interface{}, version float64) (interface{}, error) {
	// TODO add option to reject any bts with fields not in realObj - https://golang.org/pkg/encoding/json/#Decoder.DisallowUnknownFields
	if realObj == nil {
		return realObj, nil
	}

	obj := reflect.ValueOf(realObj)

	if obj.Kind() == reflect.Ptr && obj.IsNil() {
		return realObj, nil
	}
	obj = reflect.Indirect(obj)

	if obj.Type().Kind() != reflect.Struct {
		return nil, InternalError{"object must be a pointer to a struct"} // TODO handle slices of structs?
	}

	fakeVal := BuildUnmarshalObj(obj, version, false)
	if err := CopyIntoMarshalObj(fakeVal, obj); err != nil {
		return nil, err
	}

	fakeValI := fakeVal.Addr().Interface()
	return fakeValI, nil
}

// CopyIntoMarshalObj copies the data from realVal into the newVal built with BuildUnmarshalObj or BuildJSONType.
// This is only necessary to marshal/encode the real val, not to decode. Decoding with UnmarshalJSON and its compatibility wrappers will preserve existing values in the real value, just like encoding/json.Unmarshal, without this.
func CopyIntoMarshalObj(fakeVal reflect.Value, realVal reflect.Value) error {
	fakeVal = reflect.Indirect(fakeVal) // TODO supported multiple pointers?
	realVal = reflect.Indirect(realVal) // TODO supported multiple pointers?

	if fakeVal == (reflect.Value{}) {
		return errors.New("fakeVal must not be an empty value") // should never happen
	}

	if realVal == (reflect.Value{}) {
		return errors.New("realVal must not be an empty value") // should never happen
	}

	if fakeVal.Type().Kind() != reflect.Struct {
		return errors.New("fakeVal must be a struct")
	}
	if realVal.Type().Kind() != reflect.Struct {
		return errors.New("realVal must be a struct") // should never happen
	}

	fakeValType := fakeVal.Type()
	for i := 0; i < fakeVal.NumField(); i++ {
		fakeValField := fakeVal.Field(i)
		fakeValTypeField := fakeValType.Field(i)
		fieldName := fakeValTypeField.Name

		realValField := realVal.FieldByName(fieldName) // must get by name, because the field order will be different if any newer versions were omitted.
		if realValField == (reflect.Value{}) {
			return errors.New("fakeVal field '" + fieldName + "' not in realVal") // should never happen
		}

		if realValField.Type().Kind() != reflect.Ptr || !realValField.IsNil() {
			// TODO handle all nilable types
			if fakeValField.Type().Kind() != reflect.Ptr {
				return errors.New("realVal field '" + fieldName + "' not a pointer") // should never happen
			}
			if !fakeValField.CanSet() {
				return InternalError{"can't set fakeVal field '" + fieldName + "' to new pointer"} // should never happen
			}
			fakeValField.Set(reflect.New(fakeValField.Type().Elem()))
		}

		fakeValField = reflect.Indirect(fakeValField) // TODO supported multiple pointers?
		realValField = reflect.Indirect(realValField) // TODO supported multiple pointers?

		if fakeValField.Type().Kind() == reflect.Struct {
			if realValField.Type().Kind() != reflect.Struct {
				return InternalError{"fakeVal field '" + fieldName + "' struct type '" + fakeVal.Type().String() + "' does not match realVal field type '" + realValField.Type().String() + "'"} // should never happen
			}
			if err := CopyIntoMarshalObj(fakeValField, realValField); err != nil {
				return errors.New("field '" + fieldName + "': " + err.Error()) // TODO use encoding/json errs? Somehow handle user vs system errs.
			}
			continue
		}

		if fakeValField.Type() == reflect.TypeOf(IntS(0)) || fakeValField.Type() == reflect.TypeOf(UIntS(0)) || fakeValField.Type() == reflect.TypeOf(FloatS(0)) || fakeValField.Type() == reflect.TypeOf(BoolS(false)) {
			if !realValField.Type().ConvertibleTo(fakeValField.Type()) {
				return InternalError{"realVal field '" + fieldName + "' type '" + realValField.Type().String() + "' is not convertible to fakeVal field type '" + fakeValField.Type().String() + "'"} // should never happen
			}
			realValField = realValField.Convert(fakeValField.Type())
		}

		if !fakeValField.Type().AssignableTo(realValField.Type()) {
			return InternalError{"realVal field '" + fieldName + "' type '" + fakeValField.Type().String() + "' is not assignable to fakeVal field type '" + realValField.Type().String() + "'"} // should never happen
		}

		if !fakeValField.CanSet() {
			return InternalError{"can't set fakeVal field '" + fieldName + "'"} // should never happen
		}

		fakeValField.Set(realValField)
	}

	return nil
}
