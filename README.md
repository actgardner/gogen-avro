gogen-avro
===

[![Build Status](https://travis-ci.org/alanctgardner/gogen-avro.svg?branch=master)](https://travis-ci.org/alanctgardner/gogen-avro)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/alanctgardner/gogen-avro/master/LICENSE)
[![Version 4.0.1](https://img.shields.io/badge/version-4.0.1-lightgrey.svg)](https://gopkg.in/alanctgardner/gogen-avro.v4)

Generate Go structures and serializer / deserializer methods from Avro schemas. Generated serializers/deserializers are 2-8x faster than goavro, and you get compile-time safety for getting and setting fields.

### Installation

gogen-avro is a tool which you install on your system (usually on your GOPATH), and run as part of your build process. To install gogen-avro to `$GOPATH/bin/`, run:

```
go get gopkg.in/alanctgardner/gogen-avro.v4/...
```

### Usage

To generate Go source files from one or more Avro schema files, run:

```
gogen-avro.v4 [--package=<package name>] <output directory> <avro schema files>
```

You can also use a `go:generate` directive in a source file ([example](https://github.com/alanctgardner/gogen-avro/blob/master/test/primitive/schema_test.go)):

```
//go:generate $GOPATH/bin/gogen-avro.v4 . primitives.avsc
```

The generated source files contain structs for each schema, plus a function `Serialize(io.Writer)` to encode the contents into the given `io.Writer`, and `Deserialize<RecordType>(io.Reader)` to read a struct from the given `io.Reader`.

### Container File Support

gogen-avro generates a struct for each record type defined in the supplied schemas. Container file support is implemented in a generic way for all generated structs. The package `container` has a `Writer` which wraps an `io.Writer` and accepts some arguments for block size (in records) and codec (for compression). 

The `WriteRecord` method in `container.Writer` accepts an `AvroRecord`, which is an interface implemented by every generated record struct. 

An example of how to write a container file can be found in `example/container/example.go`.

[Godocs for the container package](https://godoc.org/github.com/alanctgardner/gogen-avro/container)

### Example

The `example` directory contains simple example projects with an Avro schema. Once you've installed gogen-avro on your GOPATH, you can install the example projects:

```
# Build the Go source files from the Avro schema using the generate directive
go generate github.com/alanctgardner/gogen-avro/example

# Install the example projects on the gopath
go install github.com/alanctgardner/gogen-avro/example/record
go install github.com/alanctgardner/gogen-avro/example/container
```

### Naming

Gogen-avro converts field and type names to be valid, public Go names using snaker (github.com/serenize/snaker). 
As a result the generated structs may have names that are slightly different from the names in the Avro schema - 
illegal characters will be removed and the first character will be upper-case. 

Gogen-avro respects namespaces and aliases when resolving type names. However, generated files will all be placed directly
into the package specified by the user. This may cause issues in rare cases where two types have different namespaces but the
same name.

### Type Conversion

Gogen-avro produces a Go struct which reflects the structure of your Avro schema. Most Go types map neatly onto Avro types:

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

`union` is more complicated than primitive types. We generate a struct and enum whose name is uniquely determined by the types in the union. For a field whose type is `["null", "int"]` we generate the following:

```
type UnionNullInt struct {
	// All the possible types the union could take on
	Null               interface{}
	Int                int32
	// Which field actually has data in it - defaults to the first type in the list, "null"
	UnionType          UnionNullTypeEnum
}

type UnionNullIntTypeEnum int

const (
	UnionNullIntTypeEnumNull            UnionNullIntTypeEnum = 0
	UnionNullIntTypeEnumInt             UnionNullIntTypeEnum = 1
)
```

### Versioning

This tool is versioned using [gopkg.in](http://labix.org/gopkg.in).
The API is guaranteed to be stable within a release. This guarantee applies to:
- the public members of generated structures
- the public methods attached to generated structures
- the command-line arguments of the tool itself

Only bugfixes will be backported to existing major releases.
This means that source files generated with the same major release may differ, but they will never break your build.

4.0
---
- Support for writing object container files is no longer experimental
- `container` package now works with the generated code for any record type
- Aliases and namespaces are now used properly to resolve types
- Record structs expose a `Schema` method which includes metadata from the schema definition 

3.0
---
- Experimental support for writing object container files
- Improved variable and type names
- Support for custom package names as a command line argument


2.0
---
- Bug fixes for arrays and maps with record members
- Refactored internals significantly

1.0
---
- Initial release
- No longer supported - no more bugfixes are being backported

### Thanks

Thanks to LinkedIn's [goavro library](https://github.com/linkedin/goavro), for providing the encoders for primitives.
