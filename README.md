gogen-avro
===

[![Build Status](https://travis-ci.org/alanctgardner/gogen-avro.svg?branch=master)](https://travis-ci.org/alanctgardner/gogen-avro)

Generate Go structures and serializer / deserializer methods from Avro schemas. Based on LinkedIn's goavro package.

### Installation

```
go get github.com/alanctgardner/gogen-avro
go install github.com/alanctgardner/gogen-avro
```

### Usage

```
gogen-avro <output directory> <avro schema files>
```

Gogen-avro will write a Go file in the output directory for each `record`, `fixed`, and `enum` type defined in the schema files. Multiple schema files are supported.

Generated structs have a function `Serialize(io.Writer)` to encode the contents into the given `io.Writer`, and `Deserialize<RecordType>(io.Reader)` to read a struct from the given `io.Reader`. See `test/primitive/schema_test.go` for examples of encoding and decoding.

### Type Conversion

Go types mostly map neatly onto Avro types:

| Avro Type     | Go Type           | Notes                                                                                                                |
|---------------|-------------------|----------------------------------------------------------------------------------------------------------------------|
| null          | interface{}       | This is just a placeholder, nothing is encoded/decoded                                                               |
| boolean       | bool              |                                                                                                                      |
| int, long     | int32,int64       |                                                                                                                      |
| float, double | float32, float64  |                                                                                                                      |
| bytes         | []byte            |                                                                                                                      |
| string        | string            |                                                                                                                      |
| enum          | custom type       | Generates a type with a constant for each symbol                                                                     |
| array<type>   | []<type>          |                                                                                                                      |
| map<type>     | map[string]<type> |                                                                                                                      |
| fixed         | [<n>]byte         | Fixed fields are given a custom type, which is an alias for an appropriately sized byte array                        |
| union         | custom type       | Unions are handled as a struct with one field per possible type, and an enum field to dictate which field to read    |

`union` is the exception. To avoid a round-trip through `interface{}`, we generate a struct and enumeration whose name is uniquely determined by the types in the union. This can get pretty hairy - for a field whose type is `["int", "string", "float", "double", "long", "bool", "null"]` we generate the following:

```
type UnionIntStringFloatDoubleLongBoolNull struct {
	// All the possible types the union could take on
	Int                int32
	String             string
	Float              float32
	Double             float64
	Long               int64
	Bool               bool
	Null               interface{}
	// Which field actually has data in it
	UnionType          UnionIntStringFloatDoubleLongBoolNullTypeEnum
}

// These names are obscenely long to guarantee uniqueness
type UnionIntStringFloatDoubleLongBoolNullTypeEnum int

const (
	UnionIntStringFloatDoubleLongBoolNullTypeEnumInt                UnionIntStringFloatDoubleLongBoolNullTypeEnum = 0
	UnionIntStringFloatDoubleLongBoolNullTypeEnumString             UnionIntStringFloatDoubleLongBoolNullTypeEnum = 1
	UnionIntStringFloatDoubleLongBoolNullTypeEnumFloat              UnionIntStringFloatDoubleLongBoolNullTypeEnum = 2
	UnionIntStringFloatDoubleLongBoolNullTypeEnumDouble             UnionIntStringFloatDoubleLongBoolNullTypeEnum = 3
	UnionIntStringFloatDoubleLongBoolNullTypeEnumLong               UnionIntStringFloatDoubleLongBoolNullTypeEnum = 4
	UnionIntStringFloatDoubleLongBoolNullTypeEnumBool               UnionIntStringFloatDoubleLongBoolNullTypeEnum = 5
	UnionIntStringFloatDoubleLongBoolNullTypeEnumNull               UnionIntStringFloatDoubleLongBoolNullTypeEnum = 6
)
``` 

### TODO / Caveats

This package doesn't implement the entire Avro 1.7.7 specification, specifically:

- Schema resolution
- Framing - generate RPCs and container format readers/writers

### Thanks

Thanks to LinkedIn's [goavro library](https://github.com/linkedin/goavro), for providing the encoders for primitives.
