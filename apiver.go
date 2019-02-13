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

	// if obj.Type().Kind() != reflect.Struct {
	// 	return InternalError{"object must be a pointer to a struct"}
	// }

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

	if typ.Kind() == reflect.Slice {
		return reflect.SliceOf(BuildUnmarshalType(typ.Elem(), version, strTypes))
	}
	if typ.Kind() == reflect.Map {
		return reflect.MapOf(BuildUnmarshalType(typ.Key(), version, strTypes), BuildUnmarshalType(typ.Elem(), version, strTypes))
	}
	if typ.Kind() != reflect.Struct {
		return typ // if it's not a slice, map, or struct, return the type as-is
	}

	newTypeFields := []reflect.StructField{}

	// changedAnyFields is whether any fields were changed (pointers, versions omitted, str types, etc)
	// If we don't change anything, we want to use the original struct verbatim.
	// This is important for external packages with custom marshal/unmarshal funcs that set unexported fields,
	// Because Go reflect.StructOf currently panics on unexported fields.
	// See https://github.com/golang/go/issues/25401
	changedAnyFields := false
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// if isExported := field.PkgPath == ""; !isExported {
		// 	continue // TODO copy unexported fields verbatim?
		// }

		// fmt.Println("DEBUG but type '" + typ.String() + "' field '" + field.Name + "' PkgPath '" + field.PkgPath + "'")

		props := GetTagProperties(field.Tag.Get(TagName))
		if props.Version > version {
			changedAnyFields = true // we skipped a field, structs are different
			continue
		}

		newField := reflect.StructField{}
		newField.Name = field.Name
		newField.Type = field.Type
		newField.PkgPath = field.PkgPath
		if newField.Type.Kind() == reflect.Struct {
			newType := BuildUnmarshalType(newField.Type, version, strTypes)
			if newType != newField.Type {
				changedAnyFields = true // we changed a field that was a struct, structs are different
			}
			newField.Type = newType
		}

		if newField.Type.Kind() != reflect.Ptr && props.Version != 0.0 {
			// no need to pointer-ify fields with no "api:version" tag
			// TODO verify this is correct

			// convert all versioned fields to pointers
			// this lets us later verify value=required fields exist, and return an error if any value field is nil.
			// Without this, we can't distinguish empty from missing values.
			newField.Type = reflect.PtrTo(newField.Type)
			changedAnyFields = true // we changed a field into a pointer, structs are different
		}

		if strTypes && props.Str {
			switch newField.Type.Elem().Kind() {
			case reflect.Bool:
				newField.Type = reflect.PtrTo(reflect.TypeOf(BoolS(false)))
				changedAnyFields = true // we changed a field str type, structs are different
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
				changedAnyFields = true // we changed a field str type, structs are different
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
				changedAnyFields = true // we changed a field str type, structs are different
			case reflect.Float32:
				fallthrough
			case reflect.Float64:
				newField.Type = reflect.PtrTo(reflect.TypeOf(FloatS(0)))
				changedAnyFields = true // we changed a field str type, structs are different
			default:
				// TODO error?
			}
		}

		newField.Tag = field.Tag

		newTypeFields = append(newTypeFields, newField)
	}

	if !changedAnyFields {
		// fmt.Println("DEBUG type '" + typ.String() + "' changed no fields, using verbatim!")
		return typ
	}
	return reflect.StructOf(newTypeFields)
}

// FromUnmarshalObj converts an object created with BuildUnmarshalObj, presumably after decoding data into it, into the real object.
// Returns an error if any value fields in the realObj are nil in the val.
func FromUnmarshalObj(fakeVal reflect.Value, realObj interface{}) error {
	realVal := reflect.ValueOf(realObj)
	if realVal.Kind() != reflect.Ptr {
		return InternalError{"object must be a pointer"}
	}
	if realVal.IsNil() {
		return InternalError{"object must not be nil"}
	}
	return SetUnmarshalObj(fakeVal, realVal)
}

func SetUnmarshalObj(fakeVal reflect.Value, realVal reflect.Value) error {
	fakeVal = reflect.Indirect(fakeVal)

	if fakeVal == (reflect.Value{}) {
		return nil
	}

	for realVal.Type().Kind() == reflect.Ptr {
		if realVal.IsNil() {
			realVal.Set(reflect.New(realVal.Type().Elem()))
		}
		realVal = reflect.Indirect(realVal)
	}

	if fakeVal.Type().Kind() == reflect.Slice {
		if realVal.Type().Kind() != reflect.Slice {
			return InternalError{"realVal '" + realVal.Type().String() + "' slice type does not match fakeVal type '" + fakeVal.Type().String() + "'"}
		}
		for i := 0; i < fakeVal.Len(); i++ {
			newRealValElem := reflect.New(realVal.Type().Elem())
			if err := SetUnmarshalObj(fakeVal.Index(i), newRealValElem); err != nil {
				return errors.New("setting slice type '" + realVal.Type().String() + "': " + err.Error())
			}
			newRealValElem = reflect.Indirect(newRealValElem)
			realVal.Set(reflect.Append(realVal, newRealValElem))
		}
		return nil
	} else if fakeVal.Type().Kind() == reflect.Map {
		if realVal.Type().Kind() != reflect.Map {
			return InternalError{"realVal '" + realVal.Type().String() + "' map type does not match fakeVal type '" + fakeVal.Type().String() + "'"}
		}

		if fakeVal.IsNil() {
			return nil
		}

		if realVal.IsNil() {
			realVal.Set(reflect.MakeMapWithSize(realVal.Type(), fakeVal.Len()))
		}

		for _, fakeValKey := range fakeVal.MapKeys() {
			fakeValVal := fakeVal.MapIndex(fakeValKey)

			realValKey := reflect.New(realVal.Type().Key())
			if err := SetUnmarshalObj(fakeValKey, realValKey); err != nil {
				return errors.New("setting map type '" + realVal.Type().String() + "' key: " + err.Error())
			}

			realValVal := reflect.New(realVal.Type().Elem())
			if err := SetUnmarshalObj(fakeValVal, realValVal); err != nil {
				return errors.New("setting map type '" + realVal.Type().String() + "' val: " + err.Error())
			}

			realValKey = reflect.Indirect(realValKey)
			realValVal = reflect.Indirect(realValVal)
			realVal.SetMapIndex(realValKey, realValVal)
		}
		return nil
	} else if fakeVal.Type().Kind() == reflect.Struct {
		if realVal.Type().Kind() != reflect.Struct {
			return InternalError{"realVal '" + realVal.Type().String() + "' struct type does not match fakeVal type '" + fakeVal.Type().String() + "'"}
		}

		for i := 0; i < fakeVal.NumField(); i++ {
			fakeValTypeField := fakeVal.Type().Field(i)
			fakeValField := fakeVal.Field(i)
			realValField := realVal.FieldByName(fakeValTypeField.Name)

			if isExported := fakeValTypeField.PkgPath == ""; !isExported {
				continue // TODO copy unexported fields verbatim?
			}

			// fmt.Println("DEBUG type '" + fakeVal.Type().String() + "' field '" + fakeValTypeField.Name + "' PkgPath '" + fakeValTypeField.PkgPath + "'")

			// fieldTagName is the user-facing field name: json tag if it exists, else the struct field name.
			// TODO remove/abstract json tag, so this func is fully encoder-agnostic?
			fieldTagName := fakeValTypeField.Name
			if jsonTag := fakeValTypeField.Tag.Get("json"); jsonTag != "" {
				if jsonTagParts := strings.Split(jsonTag, ","); len(jsonTagParts) > 0 {
					fieldTagName = jsonTagParts[0]
				}
			}

			if realValField == (reflect.Value{}) {
				return InternalError{"object missing field in val '" + fieldTagName + "'"}
			}

			if fakeValField.Type().Kind() == reflect.Ptr && fakeValField.IsNil() && realValField.Type().Kind() != reflect.Ptr {
				// If the "fake" val unmarshaled from json is nil (that is, was not in the JSON), and the real val type isn't a pointer, and thus "required," return an error: missing required field.
				return UserError{"missing required field: " + fieldTagName} // TODO make missing-required-field an err type?
			}

			if err := SetUnmarshalObj(fakeValField, realValField); err != nil {
				// TODO better error objects, to set the "user field name" if SetUnmarshalObj returns a user err, and the struct name if it returns a system err.
				return errors.New("field '" + fieldTagName + "':" + err.Error())
			}
		}
		return nil
	} else { // not struct, slice, or map
		if realVal.Type().Kind() == reflect.Ptr {
			if realVal.IsNil() {
				realVal.Set(reflect.New(realVal.Type().Elem()))
			}
			realVal = reflect.Indirect(realVal)
		}
		if fakeVal.Type().Kind() == reflect.Ptr {
			// TODO return success if nil?
			fakeVal = reflect.Indirect(fakeVal)
		}

		if fakeVal.Type() == reflect.TypeOf(IntS(0)) || fakeVal.Type() == reflect.TypeOf(UIntS(0)) || fakeVal.Type() == reflect.TypeOf(FloatS(0)) || fakeVal.Type() == reflect.TypeOf(BoolS(false)) {
			if !realVal.Type().ConvertibleTo(fakeVal.Type()) {
				// fmt.Println("DEBUG SetUnmarshalObj returning err not convertible")
				return InternalError{"realVal type '" + realVal.Type().String() + "' is not convertible to fakeVal type '" + fakeVal.Type().String() + "'"}
			}
			fakeVal = fakeVal.Convert(realVal.Type())
		}
		if !realVal.Type().AssignableTo(fakeVal.Type()) {
			// fmt.Println("DEBUG SetUnmarshalObj returning err not assignable")
			return InternalError{"realVal type '" + realVal.Type().String() + "' is not assignable to fakeVal type '" + fakeVal.Type().String() + "'"}
		}
		if !realVal.CanSet() {
			// fmt.Println("DEBUG SetUnmarshalObj returning err not settable")
			return InternalError{"can't set realVal type '" + realVal.Type().String() + "'"}
		}

		realVal.Set(fakeVal)
		// fmt.Println("DEBUG SetUnmarshalObj returning nil success non-struct")
		return nil
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

	// if obj.Type().Kind() != reflect.Struct {
	// 	return nil, InternalError{"object must be a pointer to a struct"} // TODO handle slices of structs?
	// }

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
	if fakeVal == (reflect.Value{}) {
		return errors.New("fakeVal must not be an empty value") // should never happen
	}
	return doCopyIntoMarshalObj(fakeVal, realVal, fakeVal.Type().String())
}

func doCopyIntoMarshalObj(fakeVal reflect.Value, realVal reflect.Value, fieldName string) error {
	if fakeVal == (reflect.Value{}) {
		return errors.New("fakeVal must not be an empty value") // should never happen
	}

	if realVal == (reflect.Value{}) {
		return errors.New("realVal must not be an empty value") // should never happen
	}

	if realVal.Type().Kind() == reflect.Ptr && realVal.IsNil() {
		return nil // TODO verify? test?
	}

	// TODO handle all nilable types
	for fakeVal.Type().Kind() == reflect.Ptr {
		if fakeVal.IsNil() {
			if !fakeVal.CanSet() {
				return InternalError{"can't set fakeVal field '" + fieldName + "' to new pointer"} // should never happen
			}
			fakeVal.Set(reflect.New(fakeVal.Type().Elem()))
		}
		fakeVal = reflect.Indirect(fakeVal) // TODO supported multiple pointers?
	}
	realVal = reflect.Indirect(realVal) // TODO supported multiple pointers?

	if fakeVal.Type() == realVal.Type() {
		// if the types are identical, set directly
		if !fakeVal.CanSet() {
			return InternalError{"can't set fakeVal '" + fakeVal.Type().String() + "'"} // should never happen
		}
		fakeVal.Set(realVal)
		return nil
	}

	if fakeVal.Type().Kind() == reflect.Struct {
		if realVal.Type().Kind() != reflect.Struct {
			return errors.New("fakeVal is a struct, realVal must also be a struct") // should never happen
		}

		fakeValType := fakeVal.Type()
		for i := 0; i < fakeVal.NumField(); i++ {
			fakeValField := fakeVal.Field(i)

			if isExported := fakeValType.Field(i).PkgPath == ""; !isExported {
				continue // TODO copy unexported fields verbatim?
			}

			fieldName := fakeValType.Field(i).Name
			realValField := realVal.FieldByName(fieldName) // must get by name, because the field order will be different if any newer versions were omitted.
			if realValField == (reflect.Value{}) {
				return errors.New("fakeVal field '" + fieldName + "' not in realVal") // should never happen
			}
			if err := doCopyIntoMarshalObj(fakeValField, realValField, fieldName); err != nil {
				return errors.New("struct field '" + fieldName + "' error: " + err.Error())
			}
		}

		return nil
	}

	if fakeVal.Type().Kind() == reflect.Slice {
		if realVal.Type().Kind() != reflect.Slice {
			return InternalError{"fakeVal '" + fakeVal.Type().String() + "' slice types do not match, realval type '" + realVal.Type().String() + "'"}
		}

		if realVal.IsNil() {
			return nil
		}

		if !fakeVal.CanSet() {
			return InternalError{"can't set fakeVal field '" + fieldName + "' to appended slice"} // should never happen
		}

		if fakeVal.IsNil() {
			fakeVal.Set(reflect.MakeSlice(fakeVal.Type(), 0, realVal.Len()))
		}

		for i := 0; i < realVal.Len(); i++ {
			fakeValElem := reflect.New(fakeVal.Type().Elem())
			if err := doCopyIntoMarshalObj(fakeValElem, realVal.Index(i), fakeValElem.Type().String()); err != nil {
				return errors.New("setting slice type '" + fakeVal.Type().String() + "': " + err.Error())
			}
			fakeValElem = reflect.Indirect(fakeValElem)
			fakeVal.Set(reflect.Append(fakeVal, fakeValElem))
		}
		return nil
	}

	if fakeVal.Type().Kind() == reflect.Map {
		if realVal.Type().Kind() != reflect.Map {
			return InternalError{"fakeVal '" + fakeVal.Type().String() + "' map types do not match, realval type '" + realVal.Type().String() + "'"}
		}

		if realVal.IsNil() {
			return nil
		}

		if !fakeVal.CanSet() {
			return InternalError{"can't set fakeVal field '" + fieldName + "' to new map"} // should never happen
		}

		if fakeVal.IsNil() {
			fakeVal.Set(reflect.MakeMapWithSize(fakeVal.Type(), realVal.Len()))
		}

		for _, realValKey := range realVal.MapKeys() {
			realValVal := realVal.MapIndex(realValKey)

			fakeValKey := reflect.New(fakeVal.Type().Key())
			if err := doCopyIntoMarshalObj(fakeValKey, realValKey, fakeValKey.Type().String()); err != nil {
				return errors.New("copying map type '" + fakeVal.Type().String() + "' key: " + err.Error())
			}

			fakeValVal := reflect.New(fakeVal.Type().Elem())
			if err := doCopyIntoMarshalObj(fakeValVal, realValVal, fakeValVal.Type().String()); err != nil {
				return errors.New("copying map type '" + fakeVal.Type().String() + "' val: " + err.Error())
			}
			fakeValKey = reflect.Indirect(fakeValKey)
			fakeValVal = reflect.Indirect(fakeValVal)
			fakeVal.SetMapIndex(fakeValKey, fakeValVal)
		}
		return nil
	}

	if fakeVal.Type() == reflect.TypeOf(IntS(0)) || fakeVal.Type() == reflect.TypeOf(UIntS(0)) || fakeVal.Type() == reflect.TypeOf(FloatS(0)) || fakeVal.Type() == reflect.TypeOf(BoolS(false)) {
		if !realVal.Type().ConvertibleTo(fakeVal.Type()) {
			return InternalError{"realVal field '" + fieldName + "' type '" + realVal.Type().String() + "' is not convertible to fakeVal field type '" + fakeVal.Type().String() + "'"} // should never happen
		}
		realVal = realVal.Convert(fakeVal.Type())
	}

	if !fakeVal.Type().AssignableTo(realVal.Type()) {
		return InternalError{"realVal field '" + fieldName + "' type '" + realVal.Type().String() + "' is not assignable to fakeVal field type '" + fakeVal.Type().String() + "'"} // should never happen
	}

	if !fakeVal.CanSet() {
		return InternalError{"can't set fakeVal field '" + fieldName + "'"} // should never happen
	}

	fakeVal.Set(realVal)
	return nil
}
