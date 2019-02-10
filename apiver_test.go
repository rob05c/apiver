package apiver

import (
	"reflect"
	"strings"
	"testing"
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

func TestMarshalJSONNonStruct(t *testing.T) {
	if _, err := MarshalJSON(42, 1.2); err == nil {
		t.Errorf("MarshalJSON error expected: 'must be a struct', actual: nil")
	}
	if _, err := MarshalJSON("foo", 1.2); err == nil {
		t.Errorf("MarshalJSON error expected: 'must be a struct', actual: nil")
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

// TODO test pointers

// TODO test different cases (to test encoding/json behavior of falling back to case-insensitive

// TODO test decoding into an object with newer version fields in the middle, to verify there aren't any "num field" issues, with the "fake object" having different reflect.Type.Field(i) indexes.
