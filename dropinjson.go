package apiver

import (
	"encoding/json"
	"io"
	"reflect"
)

// NewJSON returns a encoding/json compatible object for marshalling and unmarshalling.
// This is designed to simplify dropping in an apiver in place of encoding/json calls.
// Example:
//
//  json := apiver.NewJSON(1.4)
//  if err := json.Unmarshal(bts, &obj); err != nil {
//    return err
//  }
//
func NewJSON(version float64) EncodingJSONDropIn {
	return EncodingJSONDropIn{Version: version}
}

type EncodingJSONDropIn struct {
	Version float64
}

func (j EncodingJSONDropIn) Marshal(v interface{}) ([]byte, error) {
	return MarshalJSON(v, j.Version)
}

func (j EncodingJSONDropIn) MarshalIndent(v interface{}, prefix, indent string) ([]byte, error) {
	return MarshalJSONIndent(v, prefix, indent, j.Version)
}

func (j EncodingJSONDropIn) Unmarshal(data []byte, v interface{}) error {
	return UnmarshalJSON(data, v, j.Version)
}

type JSONDecoder struct {
	Version float64
	D       *json.Decoder
}

func (j EncodingJSONDropIn) NewDecoder(r io.Reader) *JSONDecoder {
	return &JSONDecoder{Version: j.Version, D: json.NewDecoder(r)}
}
func (d *JSONDecoder) Buffered() io.Reader        { return d.D.Buffered() }
func (d *JSONDecoder) DisallowUnknownFields()     { d.D.DisallowUnknownFields() }
func (d *JSONDecoder) More() bool                 { return d.D.More() }
func (d *JSONDecoder) Token() (json.Token, error) { return d.D.Token() }
func (d *JSONDecoder) UseNumber()                 { d.D.UseNumber() }
func (d *JSONDecoder) Decode(realObj interface{}) error {
	// TODO reduce duplication with UnmarshalJSON
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

	newVal := BuildUnmarshalObj(obj, d.Version, true)
	newValI := newVal.Addr().Interface()

	if err := d.D.Decode(newValI); err != nil {
		return err
	}

	if err := FromUnmarshalObj(newVal, realObj); err != nil {
		return err
	}
	return nil
}

type JSONEncoder struct {
	Version float64
	E       *json.Encoder
}

func (j EncodingJSONDropIn) NewEncoder(w io.Writer) *JSONEncoder {
	return &JSONEncoder{Version: j.Version, E: json.NewEncoder(w)}
}
func (e *JSONEncoder) SetEscapeHTML(on bool)           { e.E.SetEscapeHTML(on) }
func (e *JSONEncoder) SetIndent(prefix, indent string) { e.E.SetIndent(prefix, indent) }
func (e *JSONEncoder) Encode(v interface{}) error {
	obj, err := BuildMarshalObj(v, e.Version)
	if err != nil {
		return err
	}
	return e.E.Encode(obj)
}
