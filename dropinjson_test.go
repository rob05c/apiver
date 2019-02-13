package apiver

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

func TestNewJSONMarshal(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	json := NewJSON(1.3)

	a := 24
	obj := Obj{Foo: 42, A: &a}

	actual, err := json.Marshal(obj)
	if err != nil {
		t.Fatalf("json.Marshal error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":42}`
	if string(actual) != expected {
		t.Errorf("json.Marshal expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestNewJSONMarshalIndent(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	json := NewJSON(1.3)

	a := 24
	obj := Obj{Foo: 42, A: &a}

	actual, err := json.MarshalIndent(obj, " ", "\t")
	if err != nil {
		t.Fatalf("json.Marshal error expected: nil, actual: %+v", err)
	}

	expected := `{
	 "foo": 42
	}`
	if string(actual) != expected {
		t.Errorf("json.Marshal expected ''\n%+v\n'', actual ''\n%+v\n''", expected, string(actual))
	}
}

func TestNewJSONUnmarshal(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	jObj := `{"foo":42}`

	json := NewJSON(1.3)

	obj := Obj{}
	if err := json.Unmarshal([]byte(jObj), &obj); err != nil {
		t.Fatalf("json.Marshal error expected: nil, actual: %+v", err)
	}

	if obj.Foo != 42 {
		t.Errorf("json.Marshal expected obj.Foo %+v, actual %+v", 42, obj.Foo)
	}
}

func TestNewJSONDecoder(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	jObj := `{"foo":42}`

	json := NewJSON(1.3)

	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jObj)))

	obj := Obj{}

	if err := decoder.Decode(&obj); err != nil {
		t.Fatalf("json.Decoder error expected: nil, actual: %+v", err)
	}

	if obj.Foo != 42 {
		t.Errorf("json.Marshal expected obj.Foo %+v, actual %+v", 42, obj.Foo)
	}
}

func TestNewJSONDecoderBuffered(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	jObj := `{"foo":42}`

	json := NewJSON(1.3)

	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jObj)))

	buf := decoder.Buffered()
	if buf == nil {
		t.Errorf("decoder.Buffered expected non-nil, actual nil")
	}

	obj := Obj{}

	if err := decoder.Decode(&obj); err != nil {
		t.Fatalf("json.Decoder error expected: nil, actual: %+v", err)
	}

	if obj.Foo != 42 {
		t.Errorf("json.Marshal expected obj.Foo %+v, actual %+v", 42, obj.Foo)
	}
}

func TestNewJSONDecoderDisallowUnknownFields(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	jObj := `{"foo":42, "unknown-field":96}`

	json := NewJSON(1.3)

	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jObj)))

	decoder.DisallowUnknownFields()

	obj := Obj{}

	if err := decoder.Decode(&obj); err == nil {
		t.Fatalf("json.Decoder error expected: 'unknown field', actual: nil")
	}
}

func TestNewJSONDecoderToken(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	jObj := `{"foo":42, "unknown-field":96}`

	jsonn := NewJSON(1.3)

	decoder := jsonn.NewDecoder(bytes.NewBuffer([]byte(jObj)))

	token, err := decoder.Token() // (json.Token, error) { return d.D.Token() }
	if err != nil {
		t.Fatalf("json.Decoder.Token error expected: nil, actual: %v", err)
	}
	if token != json.Delim('{') {
		t.Fatalf("json.Decoder.Token expected: Delim, actual: %v", token)
	}

	obj := Obj{}

	if err := decoder.Decode(&obj); err == nil {
		t.Fatalf("json.Decoder error expected: nil, actual: %+v", err)
	}
}

func TestNewJSONDecoderMore(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	jObj := `{"foo":42, "unknown-field":96}`

	jsonn := NewJSON(1.3)

	decoder := jsonn.NewDecoder(bytes.NewBuffer([]byte(jObj)))

	if !decoder.More() {
		t.Errorf("json.Decoder.More expected: true, actual: false")
	}

	obj := Obj{}

	if err := decoder.Decode(&obj); err != nil {
		t.Fatalf("json.Decoder error expected: nil, actual: %+v", err)
	}

	if decoder.More() {
		t.Errorf("json.Decoder.More expected: false, actual: true")
	}
}

func TestNewJSONDecoderUseNumber(t *testing.T) {
	type Obj struct {
		Foo interface{} `json:"foo" api:"1.1,str"`
	}

	jObj := `{"foo":42.9}`

	jsonn := NewJSON(1.3)

	decoder := jsonn.NewDecoder(bytes.NewBuffer([]byte(jObj)))

	decoder.UseNumber()

	obj := Obj{}

	if err := decoder.Decode(&obj); err != nil {
		t.Fatalf("json.Decoder error expected: nil, actual: %+v", err)
	}

	if jNum, ok := obj.Foo.(json.Number); !ok {
		t.Errorf("json.Decoder.UseNumber obj.Foo expected: json.Number, actual: %T", obj.Foo)
	} else if string(jNum) != "42.9" {
		t.Errorf("json.Decoder.UseNumber obj.Foo expected: '42.9', actual: %v", jNum)
	}
}

func TestNewJSONDecoderBadInputNonPtr(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	jObj := `{"foo":42}`

	json := NewJSON(1.3)

	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jObj)))

	obj := Obj{}

	if err := decoder.Decode(obj); err == nil {
		t.Errorf("json.Decoder error expected: 'must be pointer', actual: %+v", err)
	}
}

func TestNewJSONDecoderBadInputNilPtr(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	jObj := `{"foo":42}`

	json := NewJSON(1.3)

	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jObj)))

	obj := (*Obj)(nil)

	if err := decoder.Decode(obj); err == nil {
		t.Errorf("json.Decoder error expected: 'must be non-nil', actual: %+v", err)
	}
}

func TestNewJSONDecoderBadInputNonStruct(t *testing.T) {
	jObj := `{"foo":42}`

	json := NewJSON(1.3)

	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jObj)))

	i := 0
	obj := &i

	if err := decoder.Decode(obj); err == nil {
		t.Errorf("json.Decoder error expected: 'must be struct', actual: %+v", err)
	}
}

func TestNewJSONDecoderBadInputMalformedJSON(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	jObj := `{"foo":42`

	json := NewJSON(1.3)

	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jObj)))

	obj := Obj{}

	if err := decoder.Decode(&obj); err == nil {
		t.Fatalf("json.Decoder error expected: 'malformed json', actual: %+v", err)
	}
}

func TestNewJSONDecoderBadInputMissingRequiredField(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1,str"`
		A   int `json:"a" api:"1.1,str"`
	}

	jObj := `{"foo":42}`

	json := NewJSON(1.3)

	decoder := json.NewDecoder(bytes.NewBuffer([]byte(jObj)))

	obj := Obj{}

	if err := decoder.Decode(&obj); err == nil {
		t.Fatalf("json.Decoder error expected: 'missing required field', actual: %+v", err)
	}
}

func TestNewJSONEncoder(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.4,str"`
	}

	json := NewJSON(1.3)

	obj := Obj{Foo: 42}
	expected := `{"foo":42}`

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	if err := enc.Encode(&obj); err != nil {
		t.Fatalf("json.Decoder error expected: nil, actual: %+v", err)
	}

	if actual := strings.TrimSpace(string(buf.Bytes())); expected != actual {
		t.Errorf("json.Marshal expected: '%+v', actual: '%+v'", expected, actual)
	}
}

func TestNewJSONEncoderSetEscapeHTMLOn(t *testing.T) {
	type Obj struct {
		Foo string `json:"foo" api:"1.1,str"`
		A   *int   `json:"a" api:"1.4,str"`
	}

	json := NewJSON(1.3)

	obj := Obj{Foo: "&"}
	expected := `{"foo":"\u0026"}`

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(&obj); err != nil {
		t.Fatalf("json.Decoder error expected: nil, actual: %+v", err)
	}

	if actual := strings.TrimSpace(string(buf.Bytes())); expected != actual {
		t.Errorf("json.Marshal expected: '%+v', actual: '%+v'", expected, actual)
	}
}

func TestNewJSONEncoderSetEscapeHTMLOff(t *testing.T) {
	type Obj struct {
		Foo string `json:"foo" api:"1.1,str"`
		A   *int   `json:"a" api:"1.4,str"`
	}

	json := NewJSON(1.3)

	obj := Obj{Foo: "&"}
	expected := `{"foo":"&"}`

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(&obj); err != nil {
		t.Fatalf("json.Decoder error expected: nil, actual: %+v", err)
	}

	if actual := strings.TrimSpace(string(buf.Bytes())); expected != actual {
		t.Errorf("json.Marshal expected: '%+v', actual: '%+v'", expected, actual)
	}
}

func TestNewJSONEncoderSetIndent(t *testing.T) {
	type Obj struct {
		Foo string `json:"foo" api:"1.1,str"`
		A   *int   `json:"a" api:"1.4,str"`
	}

	json := NewJSON(1.3)

	obj := Obj{Foo: "bar"}
	expected := `{
	  "foo": "bar"
	}`

	buf := bytes.Buffer{}
	enc := json.NewEncoder(&buf)
	enc.SetIndent("\t", "  ")
	if err := enc.Encode(&obj); err != nil {
		t.Fatalf("json.Decoder error expected: nil, actual: %+v", err)
	}

	if actual := strings.TrimSpace(string(buf.Bytes())); expected != actual {
		t.Errorf("json.Marshal expected: '%+v', actual: '%+v'", expected, actual)
	}
}
