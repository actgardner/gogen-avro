## gogen-avro


[![Build Status](https://travis-ci.org/actgardner/gogen-avro.svg?branch=master)](https://travis-ci.org/actgardner/gogen-avro)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/actgardner/gogen-avro/master/LICENSE)
[![Version 6.2.0](https://img.shields.io/badge/version-6.2.0-lightgrey.svg)](https://github.com/actgardner/gogen-avro/releases)

Generates type-safe Go code based on your Avro schemas, including serializers and deserializers that support Avro's schema evolution rules.

### Table of contents

<!--ts-->
   * [Table of contents](#table-of-contents)
   * [Installation](#installation)
   * [Usage](#usage)
   * [Generated Methods](#generated-methods)
   * [Working with Object Container Files (OCF)](#working-with-object-container-files-ocf)
   * [Example](#example)
   * [Naming](#naming)
   * [Type Conversion](#type-conversion)
   * [Versioning](#Versioning)
   * [Reporting Issues](#reporting-issues)
   * [Thanks](#thanks)
<!--te-->


### Installation

gogen-avro has two parts: a tool which you install on your system (usually into your GOPATH) to generate code, and a runtime library that gets imported.

To install the latest version of the gogen-avro executable to `$GOPATH/bin/` and generate structs,
change to any directory that is both outside of your GOPATH and outside of a module (a temp directory is fine);
then run:

```
go get github.com/actgardner/gogen-avro/v7/cmd/gogen-avro@latest
```

We recommend pinning a specific SHA of the gogen-avro tool when you compile your schemas with a tool like [retool](https://github.com/twitchtv/retool). This will ensure your builds are repeatable.

For the library imports, you should manage the dependency on this repo using Go modules.

### Usage

To generate Go source files from one or more Avro schema files, run:

```
gogen-avro [--package=<package name>] <output directory> <avro schema files>
```

You can also use a `go:generate` directive in a source file ([example](https://github.com/actgardner/gogen-avro/blob/master/test/primitive/generate.go#L3)):

```
//go:generate $GOPATH/bin/gogen-avro . primitives.avsc
```

Note: If you want to parse multiple `.avsc` files into a single Go package (a single folder), make sure you put them all in one line. gogen-avro produces a file, `primitive.go`, that will be overwritten if you run it multiple times with different `.avsc` files and the same output folder.


### Generated Methods

For each record in the provided schemas, gogen-avro will produce a struct, and the following methods:

#### `New<RecordType>()`
A constructor to create a new record struct, with no values set.

#### `New<RecordType>Writer(writer io.Writer, codec container.Codec, recordsPerBlock int64) (*container.Writer, error)`
Creates a new `container.Writer` which writes generated structs to `writer` with Avro OCF format. This is the method you want if you're writing Avro to files. `codec` supports `Identity`, `Deflate` and `Snappy` encodings per the Avro spec.

#### `New<RecordType>Reader(reader io.Reader) (<RecordTypeReader>, error)`
Creates a new `<RecordTypeReader>` which reads data in the Avro OCF format into generated structs. This is the method you want if you're reading Avro data from files. It will handle the codec and schema evolution for you based on the OCF headers and the reader schema used to generate the structs.

#### `<RecordType>.Serialize(io.Writer) error`
Write the contents of the struct into the given `io.Writer` in the Avro binary format, with no Avro Object Container File (OCF) framing.

#### `Deserialize<RecordType>(io.Reader) (<RecordType>, error)`
Read Avro data from the given `io.Reader` and deserialize it into the generated struct. This assumes the schema used to write the data is identical to the schema used to generate the struct. This method assumes there's no OCF framing. This method is also slow because it re-compiles the bytecode for your type every time - if you need performance you should call `compiler.Compile` once and then `vm.Eval` for each record.

### Working with Object Container Files (OCF)

An example of how to write a container file can be found in [example/container/example.go](https://github.com/actgardner/gogen-avro/blob/master/example/container/example.go).

[Godocs for the container package](https://godoc.org/github.com/actgardner/gogen-avro/container)

### Example

The `example` directory contains simple example projects with an Avro schema. Once you've installed gogen-avro on your GOPATH, you can install the example projects:

```
# Build the Go source files from the Avro schema using the generate directive
go generate github.com/actgardner/gogen-avro/v7/example

# Install the example projects on the GOPATH
go install github.com/actgardner/gogen-avro/v7/example/record
go install github.com/actgardner/gogen-avro/v7/example/container
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
| `null`          | `interface{}`       | This is just a placeholder, nothing is encoded/decoded                                                               |
| `boolean`       | `bool`              |                                                                                                                      |
| `int, long`     | `int32, int64`      |                                                                                                                      |
| `float, double` | `float32, float64`  |                                                                                                                      |
| `bytes`         | `[]byte`            |                                                                                                                      |
| `string`        | `string`            |                                                                                                                      |
| `enum`          | custom type       | Generates a type with a constant for each symbol                                                                     |
| `array<type>`   | `[]<type>`          |                                                                                                                      |
| `map<type>`     | custom struct | Generates a struct with a field `M`, `M` has the type `map[string]<type>`                                                  |
| `fixed`         | `[<n>]byte`         | Fixed fields are given a custom type, which is an alias for an appropriately sized byte array                        |
| `union`         | custom struct     | Unions are handled as a struct with one field per possible type, and an enum field to dictate which field to read    |

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

Until version 6.0 this project used gopkg.in for versioning of both the code generation tool and library. Older versions are still available on gopkg.in.

Releases from 6.0 onward use semver tags (ex. `v6.0.0`) which are compatible with dep and modules. See [Releases](https://github.com/actgardner/gogen-avro/releases).

### Reporting Issues

When reporting issues, please include your reader and writer schemas, and the output from the compiler logs by adding this to one of your source files:

```
import (
	"github.com/actgardner/gogen-avro/v7/compiler"
)

func init() {
	compiler.LoggingEnabled = true
}
```

The logs will be printed on stdout.

### Thanks

Thanks to LinkedIn's [goavro library](https://github.com/linkedin/goavro), for providing the encoders for primitives.
