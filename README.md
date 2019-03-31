gogen-avro
===

[![Build Status](https://travis-ci.org/actgardner/gogen-avro.svg?branch=master)](https://travis-ci.org/actgardner/gogen-avro)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/actgardner/gogen-avro/master/LICENSE)
[![Version 5.3.0](https://img.shields.io/badge/version-5.3.0-lightgrey.svg)](https://gopkg.in/actgardner/gogen-avro.v5)

Generate Go structures and serializer / deserializer methods from Avro schemas. Generated serializers/deserializers are 2-8x faster than goavro, and you get compile-time safety for getting and setting fields.

### Installation

gogen-avro is a tool which you install on your system (usually on your GOPATH), and run as part of your build process. To install gogen-avro to `$GOPATH/bin/`, first download the repository:

```
go get -d github.com/actgardner/gogen-avro
```

Then run:

```
go install github.com/actgardner/gogen-avro/gogen-avro
```

Or download and install a fixed release from gopkg.in:

```
go get -d gopkg.in/actgardner/gogen-avro.v5
go install gopkg.in/actgardner/gogen-avro.v5/gogen-avro
```

### Usage

To generate Go source files from one or more Avro schema files, run:

```
gogen-avro [--package=<package name>] [--containers] <output directory> <avro schema files>
```

You can also use a `go:generate` directive in a source file ([example](https://github.com/actgardner/gogen-avro/blob/master/test/primitive/schema_test.go)):

```
//go:generate $GOPATH/bin/gogen-avro . primitives.avsc
```

Note: If you want to parse multiple `.avsc` files into a single Go package (a single folder), make sure you put them all in one line. gogen-avro produces a file, `primitive.go`, that will be overwritten if you run it multiple times with different `.avsc` files and the same output folder.

For each record in the provided schemas, gogen-avro will produce a struct, and the following methods:

- `New<RecordType>()` - a constructor to create a new record struct with the default values from the Avro schema
- `<RecordType>.Serialize(io.Writer)` - a method to encode the contents of the struct into the given `io.Writer`
- `Deserialize<RecordType>(io.Reader)` - a method to read a struct from the given `io.Reader`

Passing the `--containers` flag also generates a method `New<RecordType>Writer(w io.Writer, codec Codec, batchSize int)` for each record type.
This is a convenience method to generate a new container writer.

The containers flag is disabled by default, because the generated files have to import the containers package. 

### Container File Support

gogen-avro generates a struct definition for each record type defined in the supplied schemas. 
The `WriteRecord` method in `container.Writer` accepts an `AvroRecord`, an interface implemented by every generated record struct.

To create a new `container.Writer`, you can specify the schema manually in `container.NewWriter`, or you can use the `--containers` flag to generate methods for each record type. 

An example of how to write a container file can be found in `example/container/example.go`.

**Experimental:** gogen-avro now supports unpacking Object Container Files. There's a `container.Reader` which will unpack the OCF framing and feed the records into a generated struct deserializer. This should only be used when you're 100% sure the reader and writer schemas are identical - you may see panics, corrupt or incomplete data when reading with a different schema than the writer.

[Godocs for the container package](https://godoc.org/github.com/actgardner/gogen-avro/container)

### Example

The `example` directory contains simple example projects with an Avro schema. Once you've installed gogen-avro on your GOPATH, you can install the example projects:

```
# Build the Go source files from the Avro schema using the generate directive
go generate github.com/actgardner/gogen-avro/example

# Install the example projects on the gopath
go install github.com/actgardner/gogen-avro/example/record
go install github.com/actgardner/gogen-avro/example/container
```

### Naming

Gogen-avro converts field and type names to be valid, public Go names by following a few simple steps:

- removing leading underscore characters (`_`)
- upper-casing the first letter of the name

This minimizes the risk that two fields with different Avro names will have the same Go name.

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
