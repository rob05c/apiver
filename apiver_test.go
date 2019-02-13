package apiver

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func TestUnmarshalJSONBasic(t *testing.T) {
	type Obj struct {
		Foo int     `json:"foo" api:"1.1,str"`
		A   int     `json:"a" api:"1.4,str"`
		F   float32 `json:"f" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": 42}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
}

func TestUnmarshalJSONPtr(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.1,str"`
		A   *int `json:"a" api:"1.2,str"`
	}

	obj := Obj{}
	objJ := `{"foo": 42, "a": 49}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.Foo != 42 {
		t.Errorf("UnmarshalJSON obj.Foo expected: %+v, actual: %+v", 42, obj.Foo)
	}
	if obj.A == nil || *obj.A != 49 {
		t.Errorf("UnmarshalJSON obj.Foo expected: %+v, actual: %+v", 49, obj.A)
	}
}

func TestUnmarshalJSONMissingVal(t *testing.T) {
	type Obj struct {
		Foo int     `json:"foo" api:"1.1,str"`
		A   int     `json:"a" api:"1.4,str"`
		F   float32 `json:"f" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": 42}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected 'missing required field a', actual nil", objJ)
	} else if !strings.Contains(err.Error(), "missing required field") {
		t.Errorf("UnmarshalJSON %+v error expected 'missing required field a', actual %+v", objJ, err)
	}
}

func TestUnmarshalJSONIntStrStr(t *testing.T) {
	type Obj struct {
		Foo int     `json:"foo" api:"1.1,str"`
		A   int     `json:"a" api:"1.4,str"`
		F   float32 `json:"f" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.1)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.Foo != 42 {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, 42, obj.Foo)
	}
}

func TestUnmarshalJSONIntStrInt(t *testing.T) {
	type Obj struct {
		Foo int     `json:"foo" api:"1.1,str"`
		A   int     `json:"a" api:"1.4,str"`
		F   float32 `json:"f" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": 42}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.1)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.Foo != 42 {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, 42, obj.Foo)
	}
}

func TestUnmarshalJSONIntStrStrNegative(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.4,str"`
		I   int `json:"i" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "i": "-99"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.I != -99 {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, -99, obj.I)
	}
}

func TestUnmarshalJSONIntStrIntNegative(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.4,str"`
		I   int `json:"i" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "i": -99}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.I != -99 {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, -99, obj.I)
	}
}

func TestUnmarshalJSONIntStrStrInvalidFloat(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.4,str"`
		I   int `json:"i" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "i": "99.4"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected non-nil, actual nil", objJ)
	}
}

func TestUnmarshalJSONIntStrIntInvalidFloat(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.4,str"`
		I   int `json:"i" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "i": 99.4}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected non-nil, actual nil", objJ)
	}
}

func TestUnmarshalJSONFloatStrFloat(t *testing.T) {
	type Obj struct {
		Foo int     `json:"foo" api:"1.1,str"`
		A   int     `json:"a" api:"1.4,str"`
		F   float32 `json:"f" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "a": "1", "f": 1.92}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.F != 1.92 {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, 1.92, obj.F)
	}
}

func TestUnmarshalJSONFloatStrStr(t *testing.T) {
	type Obj struct {
		Foo int     `json:"foo" api:"1.1,str"`
		A   int     `json:"a" api:"1.4,str"`
		F   float32 `json:"f" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "a": "1", "f": "1.92"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.F != 1.92 {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, 1.92, obj.F)
	}
}

func TestUnmarshalJSONBoolStrBool(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		B   bool `json:"b" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "b": "true"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.B != true {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, true, obj.B)
	}
}

func TestUnmarshalJSONBoolStrStr(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		B   bool `json:"b" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "b": true}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.B != true {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, true, obj.B)
	}
}

func TestUnmarshalJSONBoolStrEmptyFalse(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		B   bool `json:"b" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "b": ""}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.B != false {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, false, obj.B)
	}
}

func TestUnmarshalJSONBoolStrZeroFalse(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		B   bool `json:"b" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "b": "0"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.B != false {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, false, obj.B)
	}
}

func TestUnmarshalJSONBoolStrTrue(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		B   bool `json:"b" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "b": "asdf"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.B != true {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, true, obj.B)
	}
}

func TestUnmarshalJSONBoolStrZeroPointZeroTrue(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		B   bool `json:"b" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "b": "0.0"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.B != true {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, true, obj.B)
	}
}

func TestUnmarshalJSONBoolNumTrue(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		B   bool `json:"b" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "b": 42.1}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.B != true {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, true, obj.B)
	}
}

func TestUnmarshalJSONBoolZeroFalse(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		B   bool `json:"b" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "b": 0}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.B != false {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, false, obj.B)
	}
}

func TestUnmarshalJSONBoolZeroPointZeroFalse(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		B   bool `json:"b" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "b": 0.0}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.B != false {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, false, obj.B)
	}
}

func TestUnmarshalJSONBoolArrTrue(t *testing.T) {
	// TODO verify and change this to be false, for BoolS to emulate Perl.
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		B   bool `json:"b" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "b": []}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.B != true {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, true, obj.B)
	}
}

func TestUnmarshalJSONBoolObjTrue(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		B   bool `json:"b" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "b": {"foo":42}}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.B != true {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, true, obj.B)
	}
}

func TestUnmarshalJSONUintStrUint(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		U   uint `json:"u" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "u": 99}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.U != 99 {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, true, obj.U)
	}
}

func TestUnmarshalJSONUintStrStr(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		U   uint `json:"u" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "u": "99"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.U != 99 {
		t.Errorf("UnmarshalJSON %+v obj.foo expected: %v, actual: %+v", objJ, true, obj.U)
	}
}

func TestUnmarshalJSONUintStrStrInvalidNonNumber(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		U   uint `json:"u" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "u": "99a"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected non-nil, actual nil", objJ)
	}
}

func TestUnmarshalJSONUintStrStrInvalidFloat(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		U   uint `json:"u" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "u": "99.7"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected non-nil, actual nil", objJ)
	}
}

func TestUnmarshalJSONUintStrStrInvalidNegative(t *testing.T) {
	type Obj struct {
		Foo int  `json:"foo" api:"1.4,str"`
		U   uint `json:"u" api:"1.4,str"`
	}

	obj := Obj{}
	objJ := `{"foo": "42", "u": "-99"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected non-nil, actual nil", objJ)
	}
}

func TestUnmarshalJSONIntNonStr(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1"`
	}

	obj := Obj{}
	objJ := `{"foo": "42"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected 'cannot unmarshal string', actual nil", objJ)
	}
}

func TestUnmarshalJSONUIntNonStr(t *testing.T) {
	type Obj struct {
		Foo uint `json:"foo" api:"1.1"`
	}

	obj := Obj{}
	objJ := `{"foo": "42"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected 'cannot unmarshal string', actual nil", objJ)
	}
}

func TestUnmarshalJSONFloatNonStr(t *testing.T) {
	type Obj struct {
		Foo float32 `json:"foo" api:"1.1"`
	}

	obj := Obj{}
	objJ := `{"foo": "42"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected 'cannot unmarshal string', actual nil", objJ)
	}
}

func TestUnmarshalJSONBoolNonStr(t *testing.T) {
	type Obj struct {
		Foo bool `json:"foo" api:"1.1"`
	}

	obj := Obj{}
	objJ := `{"foo": "true"}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected 'cannot unmarshal string', actual nil", objJ)
	}
}

func TestUnmarshalJSONNonPtr(t *testing.T) {
	type Obj struct {
		Foo bool `json:"foo" api:"1.1"`
	}

	obj := Obj{}
	objJ := `{"foo": "true"}`
	err := UnmarshalJSON([]byte(objJ), obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected 'must be a pointer', actual nil", objJ)
	}
}

func TestUnmarshalJSONNewerVersionUnparsed(t *testing.T) {
	type Obj struct {
		A int  `json:"a" api:"1.1"`
		B *int `json:"b" api:"1.5"`
	}

	obj := Obj{}
	objJ := `{"a": 101, "b": 102}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.A != 101 {
		t.Errorf("UnmarshalJSON %+v obj.a expected: %v, actual: %+v", objJ, 101, obj.A)
	}
	if obj.B != nil {
		t.Errorf("UnmarshalJSON %+v obj.b expected: %+v, actual: %+v", objJ, nil, obj.B)
	}
}

func TestUnmarshalJSONNewerVersionUnparsedIntStr(t *testing.T) {
	type Obj struct {
		A int  `json:"a" api:"1.1"`
		B *int `json:"b" api:"1.5,str"`
	}

	obj := Obj{}
	objJ := `{"a": 101, "b": 102}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Errorf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.A != 101 {
		t.Errorf("UnmarshalJSON %+v obj.a expected: %v, actual: %+v", objJ, 101, obj.A)
	}
	if obj.B != nil {
		t.Errorf("UnmarshalJSON %+v obj.b expected: %+v, actual: %+v", objJ, nil, obj.B)
	}
}

func TestUnmarshalJSONNil(t *testing.T) {
	type Obj struct {
		Foo bool `json:"foo" api:"1.1"`
	}

	obj := (*Obj)(nil)
	objJ := `{"foo": "true"}`
	err := UnmarshalJSON([]byte(objJ), obj, 1.4)
	if err == nil {
		t.Errorf("UnmarshalJSON %+v error expected 'must not be nil', actual nil", objJ)
	}
}

func TestUnmarshalJSONNestedObjs(t *testing.T) {
	type B struct {
		Val       int  `json:"val" api:"1.1,str"`
		NewVal    *int `json:"new_val" api:"1.5,str"`
		NonStrVal int  `json:"non_str_val" api:"1.1"`
	}

	type A struct {
		B B `json:"b" api:"1.2"`
	}

	type Obj struct {
		A A `json:"a" api:"1.2"`
	}

	obj := Obj{}
	objJ := `{"a": {
             "b": {
               "val": "42",
               "new_val": "99",
               "non_str_val": 8
              }
            }
          }`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.4)
	if err != nil {
		t.Fatalf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj.A.B.Val != 42 {
		t.Errorf("UnmarshalJSON %+v obj.a.b.val expected: %v, actual: %+v", objJ, 42, obj.A.B.Val)
	}
	if obj.A.B.NewVal != nil {
		t.Errorf("UnmarshalJSON %+v obj.a.b.newval expected: %v, actual: %+v", objJ, nil, obj.A.B.NewVal)
	}
	if obj.A.B.NonStrVal != 8 {
		t.Errorf("UnmarshalJSON %+v obj.a.b.nonstrval expected: %v, actual: %+v", objJ, 8, obj.A.B.NonStrVal)
	}
}

func TestMarshalJSONBasic(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1"`
	}
	obj := Obj{Foo: 42}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":42}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONStr(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1,str"`
	}
	obj := Obj{Foo: 42}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":42}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONPtr(t *testing.T) {
	type Obj struct {
		Foo *int `json:"foo" api:"1.1"`
	}
	v := 42
	obj := Obj{Foo: &v}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":42}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONStrPtr(t *testing.T) {
	type Obj struct {
		Foo *int `json:"foo" api:"1.1,str"`
	}
	v := 42
	obj := Obj{Foo: &v}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":42}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONNewVers(t *testing.T) {
	type Obj struct {
		Foo    int  `json:"foo" api:"1.1,str"`
		NewFoo *int `json:"new_foo" api:"1.5"`
	}
	v := 19
	obj := Obj{Foo: 9, NewFoo: &v}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":9}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONTwoVers(t *testing.T) {
	type Obj struct {
		Foo    int  `json:"foo" api:"1.1"`
		FooTwo *int `json:"foo2" api:"1.2,str"`
	}
	v := 43
	obj := Obj{Foo: 42, FooTwo: &v}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":42,"foo2":43}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONNested(t *testing.T) {
	type B struct {
		Val       int  `json:"val" api:"1.1,str"`
		NewVal    *int `json:"new_val" api:"1.5,str"`
		NonStrVal int  `json:"non_str_val" api:"1.1"`
	}

	type A struct {
		B B `json:"b" api:"1.2"`
	}

	type Obj struct {
		A A `json:"a" api:"1.2"`
	}

	obj := Obj{}
	obj.A.B.Val = 42
	v := 31337
	obj.A.B.NewVal = &v
	obj.A.B.NonStrVal = 8

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"a":{"b":{"val":42,"non_str_val":8}}}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONBasicPtr(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1"`
	}
	obj := Obj{Foo: 42}

	actual, err := MarshalJSON(&obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":42}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONNestedPtr(t *testing.T) {
	type B struct {
		Val       int  `json:"val" api:"1.1,str"`
		NewVal    *int `json:"new_val" api:"1.5,str"`
		NonStrVal int  `json:"non_str_val" api:"1.1"`
	}

	type A struct {
		B B `json:"b" api:"1.2"`
	}

	type Obj struct {
		A A `json:"a" api:"1.2"`
	}

	obj := Obj{}
	obj.A.B.Val = 42
	v := 31337
	obj.A.B.NewVal = &v
	obj.A.B.NonStrVal = 8

	actual, err := MarshalJSON(&obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"a":{"b":{"val":42,"non_str_val":8}}}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONNilObj(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1"`
	}
	obj := (*Obj)(nil)

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `null`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONNilInterface(t *testing.T) {
	actual, err := MarshalJSON(nil, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `null`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONUIntStr(t *testing.T) {
	type Obj struct {
		Foo uint `json:"foo" api:"1.1,str"`
	}
	obj := Obj{Foo: 42}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":42}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONFloatStr(t *testing.T) {
	type Obj struct {
		Foo float64 `json:"foo" api:"1.1,str"`
	}
	obj := Obj{Foo: 42.8}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":42.8}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONBool(t *testing.T) {
	type Obj struct {
		Foo bool `json:"foo" api:"1.1"`
	}
	obj := Obj{Foo: true}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":true}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONBoolStr(t *testing.T) {
	type Obj struct {
		Foo bool `json:"foo" api:"1.1,str"`
	}
	obj := Obj{Foo: true}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":true}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONPrecisionLoss(t *testing.T) {
	// tests that a float32 doesn't lose precision when marshalled.
	type Obj struct {
		Foo float32 `json:"foo" api:"1.1,str"`
	}
	obj := Obj{Foo: 42.8}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"foo":42.8}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestCopyIntoMarshalObjBadInputDifferentStructs(t *testing.T) {
	type A struct {
		A float32 `json:"a" api:"1.1,str"`
	}
	type B struct {
		B float32 `json:"b" api:"1.1,str"`
	}

	a := A{9}
	b := B{1}

	ra := reflect.ValueOf(a)
	rb := reflect.ValueOf(b)
	if err := CopyIntoMarshalObj(ra, rb); err == nil {
		// this test exists as much to ensure it doesn't panic, as to ensure an error
		t.Fatalf("CopyIntoMarshalObj with heterogeneous structs, error expected: not nil, actual: nil")
	}
}

func TestCopyIntoMarshalObjBadInputDifferentBasicTypes(t *testing.T) {
	type A struct {
		A float32 `json:"a" api:"1.1,str"`
	}
	type B int

	a := A{9}
	b := B(42)

	ra := reflect.ValueOf(a)
	rb := reflect.ValueOf(b)

	if err := CopyIntoMarshalObj(ra, rb); err == nil {
		// this test exists as much to ensure it doesn't panic, as to ensure an error
		t.Fatalf("CopyIntoMarshalObj with heterogeneous types, error expected: not nil, actual: nil")
	}

	if err := CopyIntoMarshalObj(rb, ra); err == nil {
		// this test exists as much to ensure it doesn't panic, as to ensure an error
		t.Fatalf("CopyIntoMarshalObj with heterogeneous types, error expected: not nil, actual: nil")
	}
}

func TestCopyIntoMarshalObjBadInputEmptyVals(t *testing.T) {
	type A struct {
		A float32 `json:"a" api:"1.1,str"`
	}
	type B int

	a := A{9}
	b := B(42)

	ra := reflect.ValueOf(a)
	rb := reflect.ValueOf(b)
	re := reflect.Value{}
	if err := CopyIntoMarshalObj(ra, re); err == nil {
		// this test exists as much to ensure it doesn't panic, as to ensure an error
		t.Fatalf("CopyIntoMarshalObj with empty value, error expected: not nil, actual: nil")
	}
	if err := CopyIntoMarshalObj(re, rb); err == nil {
		// this test exists as much to ensure it doesn't panic, as to ensure an error
		t.Fatalf("CopyIntoMarshalObj with empty value, error expected: not nil, actual: nil")
	}
}

func TestCopyIntoMarshalObjBadInputFakeValNonPtr(t *testing.T) {
	type Fake struct {
		A float32 `json:"a" api:"1.1,str"`
	}
	type Real struct {
		A float32 `json:"a" api:"1.1,str"`
	}

	fake := Fake{9}
	real := Real{91}

	rFake := reflect.ValueOf(fake)
	rReal := reflect.ValueOf(real)
	if err := CopyIntoMarshalObj(rFake, rReal); err == nil {
		// this test exists as much to ensure it doesn't panic, as to ensure an error
		t.Fatalf("CopyIntoMarshalObj with empty value, error expected: not nil, actual: nil")
	}
}

func TestCopyIntoMarshalObjBadInputHeterogeneousFieldStructs(t *testing.T) {
	type A struct {
		Foo *float32 `json:"a" api:"1.1,str"`
	}

	type Fake struct {
		A *A `json:"a" api:"1.1,str"`
	}

	type Real struct {
		A float32 `json:"a" api:"1.1,str"`
	}

	f := float32(8.8)
	fake := Fake{A: &A{Foo: &f}}
	real := Real{91}

	rFake := reflect.ValueOf(fake)
	rReal := reflect.ValueOf(real)
	if err := CopyIntoMarshalObj(rFake, rReal); err == nil {
		// this test exists as much to ensure it doesn't panic, as to ensure an error
		t.Fatalf("CopyIntoMarshalObj with empty value, error expected: not nil, actual: nil")
	}
}

func TestUnmarshalJSONSlice(t *testing.T) {
	type Obj struct {
		Foo int     `json:"foo" api:"1.1,str"`
		A   int     `json:"a" api:"1.4,str"`
		F   float32 `json:"f" api:"1.4,str"`
	}

	obj := []Obj{}
	objJ := `[{"foo": 42}, {"foo": "99"}]`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Fatalf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}

	if len(obj) != 2 {
		t.Fatalf("UnmarshalJSON len(obj) expected %v, actual %+v", 2, len(obj))
	}
	if obj[0].Foo != 42 {
		t.Errorf("UnmarshalJSON obj[0].Foo expected %v, actual %+v", 42, obj[0].Foo)
	}
	if obj[1].Foo != 99 {
		t.Errorf("UnmarshalJSON obj[1].Foo expected %v, actual %+v", 99, obj[1].Foo)
	}
}

func TestUnmarshalJSONMap(t *testing.T) {
	type Obj struct {
		Foo int     `json:"foo" api:"1.1,str"`
		A   int     `json:"a" api:"1.4,str"`
		F   float32 `json:"f" api:"1.4,str"`
	}

	obj := map[string]Obj{}
	objJ := `{"a": {"foo": 42}, "b": {"foo": "99"}}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Fatalf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}

	if len(obj) != 2 {
		t.Fatalf("UnmarshalJSON len(obj) expected %v, actual %+v", 2, len(obj))
	}
	if obj["a"].Foo != 42 {
		t.Errorf("UnmarshalJSON obj[0].Foo expected %v, actual %+v", 42, obj["a"].Foo)
	}
	if obj["b"].Foo != 99 {
		t.Errorf("UnmarshalJSON obj[1].Foo expected %v, actual %+v", 99, obj["b"].Foo)
	}
}

func TestUnmarshalJSONMapNil(t *testing.T) {
	type Obj struct {
		Foo int     `json:"foo" api:"1.1,str"`
		A   int     `json:"a" api:"1.4,str"`
		F   float32 `json:"f" api:"1.4,str"`
	}

	obj := (map[string]Obj)(nil)
	objJ := `null`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Fatalf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}

	if obj != nil {
		t.Fatalf("UnmarshalJSON obj expected %v, actual %+v", nil, obj)
	}
}

func TestUnmarshalJSONMapEmpty(t *testing.T) {
	type Obj struct {
		Foo int     `json:"foo" api:"1.1,str"`
		A   int     `json:"a" api:"1.4,str"`
		F   float32 `json:"f" api:"1.4,str"`
	}

	obj := map[string]Obj{}
	objJ := `{}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Fatalf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}

	if len(obj) != 0 {
		t.Fatalf("UnmarshalJSON len(obj) expected %v, actual %+v", 0, len(obj))
	}
}

func TestUnmarshalJSONObjMap(t *testing.T) {
	type B struct {
		C int `json:"c" api:"1.1,str"`
	}
	type A struct {
		BS map[string]B `json:"bs" api:"1.1"`
	}

	// obj := A{BS: map[string]B{"x": {C: 42}, "y": {C: 44}}}
	obj := A{}
	objJ := `{"bs":{"x":{"c":42},"y":{"c":44}}}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Fatalf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}

	if len(obj.BS) != 2 {
		t.Fatalf("UnmarshalJSON len(obj) expected %v, actual %+v", 2, len(obj.BS))
	}
	if obj.BS["x"].C != 42 {
		t.Errorf("UnmarshalJSON obj.BS['x'].C expected %v, actual %+v", 42, obj.BS["x"].C)
	}
	if obj.BS["y"].C != 44 {
		t.Errorf("UnmarshalJSON obj.BS['y'].C expected %v, actual %+v", 44, obj.BS["y"].C)
	}
}
func TestUnmarshalJSONObjMapNonNil(t *testing.T) {
	type B struct {
		C int `json:"c" api:"1.1,str"`
	}
	type A struct {
		BS map[string]B `json:"bs" api:"1.1"`
	}

	// obj := A{BS: map[string]B{"x": {C: 42}, "y": {C: 44}}}
	obj := A{BS: map[string]B{}}
	objJ := `{"bs":{"x":{"c":42},"y":{"c":44}}}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Fatalf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}

	if len(obj.BS) != 2 {
		t.Fatalf("UnmarshalJSON len(obj) expected %v, actual %+v", 2, len(obj.BS))
	}
	if obj.BS["x"].C != 42 {
		t.Errorf("UnmarshalJSON obj.BS['x'].C expected %v, actual %+v", 42, obj.BS["x"].C)
	}
	if obj.BS["y"].C != 44 {
		t.Errorf("UnmarshalJSON obj.BS['y'].C expected %v, actual %+v", 44, obj.BS["y"].C)
	}
}

func TestUnmarshalJSONObjSlice(t *testing.T) {
	type B struct {
		C int `json:"c" api:"1.1,str"`
	}
	type A struct {
		BS []B `json:"bs" api:"1.1"`
	}

	obj := A{}
	objJ := `{"bs":[{"c":42},{"c":44}]}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Fatalf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}

	if len(obj.BS) != 2 {
		t.Fatalf("UnmarshalJSON len(obj) expected %v, actual %+v", 2, len(obj.BS))
	}
	if obj.BS[0].C != 42 {
		t.Errorf("UnmarshalJSON obj.BS[0].C expected %v, actual %+v", 42, obj.BS[0].C)
	}
	if obj.BS[1].C != 44 {
		t.Errorf("UnmarshalJSON obj.BS['y'].C expected %v, actual %+v", 44, obj.BS[1].C)
	}
}

func TestUnmarshalJSONObjSliceNonNil(t *testing.T) {
	type B struct {
		C int `json:"c" api:"1.1,str"`
	}
	type A struct {
		BS []B `json:"bs" api:"1.1"`
	}

	obj := A{BS: []B{}}
	objJ := `{"bs":[{"c":42},{"c":44}]}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Fatalf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}

	if len(obj.BS) != 2 {
		t.Fatalf("UnmarshalJSON len(obj) expected %v, actual %+v", 2, len(obj.BS))
	}
	if obj.BS[0].C != 42 {
		t.Errorf("UnmarshalJSON obj.BS[0].C expected %v, actual %+v", 42, obj.BS[0].C)
	}
	if obj.BS[1].C != 44 {
		t.Errorf("UnmarshalJSON obj.BS['y'].C expected %v, actual %+v", 44, obj.BS[1].C)
	}
}

func TestUnmarshalJSONStr(t *testing.T) {
	obj := "foo"
	objJ := `"bar"`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Fatalf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj != "bar" {
		t.Errorf("UnmarshalJSON obj[0].Foo expected %v, actual %+v", "bar", obj)
	}
}

func TestUnmarshalJSONInt(t *testing.T) {
	obj := 42
	objJ := `24`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Fatalf("UnmarshalJSON %+v error expected nil, actual %+v", objJ, err)
	}
	if obj != 24 {
		t.Errorf("UnmarshalJSON obj[0].Foo expected %v, actual %+v", 24, obj)
	}
}

func TestMarshalSlice(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1"`
	}
	obj := []Obj{{Foo: 42}, {Foo: 44}}

	actual, err := MarshalJSON(&obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `[{"foo":42},{"foo":44}]`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalSliceNil(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1"`
	}
	obj := ([]Obj)(nil)

	actual, err := MarshalJSON(&obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `null`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalSliceEmpty(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1"`
	}
	obj := []Obj{}

	actual, err := MarshalJSON(&obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `[]`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalMap(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1"`
	}
	obj := map[string]Obj{"a": {Foo: 42}, "b": {Foo: 44}}

	actual, err := MarshalJSON(&obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"a":{"foo":42},"b":{"foo":44}}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalMapNil(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1"`
	}
	obj := (map[string]Obj)(nil)

	actual, err := MarshalJSON(&obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `null`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalMapEmpty(t *testing.T) {
	type Obj struct {
		Foo int `json:"foo" api:"1.1"`
	}
	obj := map[string]Obj{}

	actual, err := MarshalJSON(&obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalObjMap(t *testing.T) {
	type B struct {
		C int `json:"c" api:"1.1,str"`
	}
	type A struct {
		BS map[string]B `json:"bs" api:"1.1"`
	}
	obj := A{BS: map[string]B{"x": {C: 42}, "y": {C: 44}}}

	actual, err := MarshalJSON(&obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"bs":{"x":{"c":42},"y":{"c":44}}}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalObjSlice(t *testing.T) {
	type B struct {
		C int `json:"c" api:"1.1,str"`
	}
	type A struct {
		BS []B `json:"bs" api:"1.1"`
	}
	obj := A{BS: []B{{C: 42}, {C: 44}}}

	actual, err := MarshalJSON(&obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"bs":[{"c":42},{"c":44}]}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalObjSlicePtrs(t *testing.T) {
	type B struct {
		C int `json:"c" api:"1.1,str"`
	}
	type A struct {
		BS []*B `json:"bs" api:"1.1"`
	}

	obj := A{BS: []*B{}}
	obj.BS = append(obj.BS, &B{C: 42})
	obj.BS = append(obj.BS, &B{C: 44})

	actual, err := MarshalJSON(&obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"bs":[{"c":42},{"c":44}]}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestMarshalJSONTagDash(t *testing.T) {
	type Obj struct {
		A int `json:"-" api:"1.1,str"`
		B int `json:"-" api:"1.1,str"`
		C int `json:"c" api:"1.1,str"`
	}
	obj := Obj{A: 4, B: 8, C: 16}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"c":16}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON error expected ''%+v'', actual ''%+v''", expected, string(actual))
	}
}

func TestUnmarshalJSONTagDashPtr(t *testing.T) {
	// TODO it never makes sense to include an "api" tag with a json:"-" tag. Here we prove it's ignored correctly; but, should we return a system error instead?

	type Obj struct {
		A int  `json:"a" api:"1.1,str"`
		B *int `json:"-" api:"1.1,str"`
	}

	obj := Obj{}
	objJ := `{"a":4, "B": 42}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err != nil {
		t.Errorf("UnmarshalJSON error expected: '%v', actual '%v'", nil, err)
	}
	if obj.A != 4 {
		t.Errorf("UnmarshalJSON obj.A expected: '%v', actual '%v'", 4, obj.A)
	}
	if obj.B != nil {
		t.Errorf("UnmarshalJSON obj.B expected: '%v', actual '%v'", nil, *obj.B)
	}
}

func TestUnmarshalJSONTagDashRequired(t *testing.T) {
	type Obj struct {
		A int `json:"-" api:"1.1,str"`
		B int `json:"-" api:"1.1,str"`
		C int `json:"c" api:"1.1,str"`
	}

	obj := Obj{}
	objJ := `{"c":16}`
	err := UnmarshalJSON([]byte(objJ), &obj, 1.3)
	if err == nil {
		t.Errorf("UnmarshalJSON error expected: '%v', actual '%v'", "missing required field", err)
	}
}

func TestMarshalJSONTime(t *testing.T) {
	type Obj struct {
		T time.Time `json:"t" api:"1.1,str"`
	}

	now := time.Now()
	obj := Obj{T: now}

	actual, err := MarshalJSON(obj, 1.2)
	if err != nil {
		t.Fatalf("MarshalJSON error expected: nil, actual: %+v", err)
	}

	expected := `{"t":"` + now.Format(time.RFC3339Nano) + `"}`
	if string(actual) != expected {
		t.Errorf("MarshalJSON time\nexpected: %+v\nactual:   %+v", expected, string(actual))
	}
}

// TODO test slice-of-pointers

// TODO test pointers

// TODO test different cases (to test encoding/json behavior of falling back to case-insensitive

// TODO test decoding into an object with newer version fields in the middle, to verify there aren't any "num field" issues, with the "fake object" having different reflect.Type.Field(i) indexes.
