# apiver
A Go API versioning package, for encoding and decoding structs per Semantic Versioning.

This package lets you use Semantic Versioning (https://semver.org/), specifically minor versioning, on a Go API, via struct tags.

To use it, add struct tags for the minor version each field was added in, then call `UnmarshalJSON(bytes, object, version)` or `MarshalJSON(obj, version)` with the current version. For example, if an HTTP client requests `https://example.com/api/1.1/foo`, you would load the `Foo` object, then call `MarshalJSON(foo, 1.1)`.

The `encoding/json` functions `Marshal`, `MarshalIndent`, `Unmarshal`, and `NewDecoder` are also implemented as object methods, to make it easier to use this package as a drop-in replacement for `encoding/json`. For example:

```go
json := apiver.NewJSON(1.1)
bts, err := json.Marshal(obj)
```

Accepting weakly-typed strings for primitives is also supported, via a `str` tag field. This is different than the `encoding/json` tag field `,string`, which will only decode strings. This `str` tag field accepts both strings and the native primitive type (integers, unsigned integers, floats, or bools). The `encoding/json` `,string` tag field will also encode the value as a string, whereas `apiver` `,str` will encode the value as a JSON number or boolean.

Example struct:

```go
	type Obj struct {
		Foo int      `json:"foo" api:"1.1,str"`
		A   int      `json:"a" api:"1.1,str"`
		F   *float32 `json:"f" api:"1.4,str"`
	}
```

For more examples, see the tests.

# Other Encodings

Currently, only JSON marshal and unmarshal functions are provided. But the package is structured such that adding additional encodings would be relatively easy. The functions takes real objects, parse their tags, dynamically create new objects with the appropriate fields and tags, then pass them to `encoding/json` `Marshal` and `Unmarshal`.

Hence, to add another encoding, the new functions need only call the object construction functions and pass the result to their encoder and decoder functions. See `BuildUnmarshalObj`, `FromUnmarshalObj`, and `BuildMarshalObj`.

# Tests

The package currently has 130 tests in 1104 lines, and `go test -cover` reports 90.7%.

Comparatively, the functions which make up the core reflection logic comprise 238 lines, and the total project including `encoding/json` wrappers is 392 lines of code.
