## gogen-avro

[![CircleCI](https://circleci.com/gh/actgardner/gogen-avro/tree/master.svg?style=svg)](https://circleci.com/gh/actgardner/gogen-avro/tree/master)
[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/actgardner/gogen-avro/master/LICENSE)
[![Version 10.0.0](https://img.shields.io/badge/version-10.0.0-lightgrey.svg)](https://github.com/actgardner/gogen-avro/releases)

Generates type-safe Go code based on your Avro schemas, including serializers and deserializers that support Avro's schema evolution rules.
Also supports deserializing generic Avro data (in beta).

### Table of contents

<!--ts-->
   * [Table of contents](#table-of-contents)
   * [Installation](#installation)
   * [Generic Data](#generic-data)
   * [Usage](#usage)
   * [Generated Methods](#generated-methods)
   * [Working with Object Container Files (OCF)](#working-with-object-container-files-ocf)
   * [Single Object Encoding](#single-object-encoding)
   * [Examples](#examples)
   * [Naming](#naming)
   * [Type Conversion](#type-conversion)
   * [Versioning](#Versioning)
   * [Companies/Projects Using gogen-avro](#companiesprojects-using-gogen-avro)
   * [Reporting Issues](#reporting-issues)
   * [Alternatives](#alternatives)
<!--te-->


### Installation

gogen-avro has two parts: a tool which you install on your system (usually on your GOPATH) to generate code, and a runtime library that gets imported.

To generate structs, install the command-line tool (this isn't necessary for the generic package):

```
go install github.com/actgardner/gogen-avro/v10/cmd/...@latest
```

This will put the `gogen-avro` binary in `$GOPATH/bin`, which should be part of your PATH.

### Generic Data

_Note: Generic data is support is currently beta. Please report any issues or feature requests!_

To deserialize generic Avro data into Go structs without generating code, instantiate a new `generic.Codec` with the reader and writer schemas:

```
import "github.com/actgardner/gogen-avro/v10/generic"

codec, err := NewCodecFromSchema(writerSchema, readerSchema) // writerSchema and readerSchema may be the same schema if you're not dealing with evolution
datum, err := codec.Deserialize(r) // r is an io.Reader with the Avro-encoded bytes
```

### Usage

To generate Go source files from one or more Avro schema files, run:

```
gogen-avro [--containers=false] [--sources-comment=false] [--short-unions=false] [--package=<package name>] <output directory> <avro schema files>
```

You can also use a `go:generate` directive in a source file ([example](https://github.com/actgardner/gogen-avro/blob/master/v10/test/primitive/generate.go#L3)):

```
//go:generate $GOPATH/bin/gogen-avro . primitives.avsc
```

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

An example of how to write a container file can be found in [example/container/example.go](https://github.com/actgardner/gogen-avro/blob/master/v10/example/container/example.go).

[Godocs for the container package](https://godoc.org/github.com/actgardner/gogen-avro/v10/container)

### Single-Object Encoding

An example of how to read and write Single-Object encoded messages (for use with Kafka, for instance) can be found in [example/single_object/example.go](https://github.com/actgardner/gogen-avro/blob/master/v10/example/single_object/example.go).

[Godocs for the soe package](https://godoc.org/github.com/actgardner/gogen-avro/v10/soe)

### Examples

The `example` directory contains simple example projects with an Avro schema. Once you've installed the CLI tool you can install the example projects:

```
# Build the Go source files from the Avro schema using the generate directive
go generate github.com/actgardner/gogen-avro/v10/example

# Install the example projects on the GOPATH
go install github.com/actgardner/gogen-avro/v10/example/record
go install github.com/actgardner/gogen-avro/v10/example/container
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
| `map<type>`     | `map[string]<type>` |                                                  |
| `fixed`         | `[<n>]byte`         | Fixed fields are given a custom type, which is an alias for an appropriately sized byte array                        |
| `union`         | custom struct       | Unions are handled as a struct with one field per possible type, and an enum field to dictate which field to read    |

`union` is more complicated than primitive types. We generate a struct and enum whose name is uniquely determined by the types in the union. For a field whose type is `["string", "int"]` we generate the following:

```
type UnionStringInt struct {
	// All the possible types the union could take on, except `null`
	String             string
	Int                int32
	// Which field actually has data in it - defaults to the first type in the list
	UnionType          UnionStringIntTypeEnum
}

type UnionStringIntTypeEnum int

const (
	UnionStringIntTypeEnumInt             UnionStringIntTypeEnum = 1
)
```

`null` unions are unique - to set a union to `null`, simply the set the pointer the union to `nil`, This lines up with [Avro's JSON encoding](https://avro.apache.org/docs/current/spec.html#json_encoding).

### Versioning

As of v7, this project uses go modules. Imports should be written as `github.com/actgardner/gogen-avro/v7`. Generated code will reference the same major version as the CLI tool.

v6.x releases used semver tags (ex. `v6.0.0`) which are compatible with dep and modules. 

Before version 6.0 this project used gopkg.in for versioning of both the code generation tool and library. Older versions are still available on gopkg.in.

See [Releases](https://github.com/actgardner/gogen-avro/releases) for the changelogs.

### Companies/Projects Using gogen-avro

If you're using gogen-avro in production, let us know! File an issue or open a PR to add your company or project here.

### Reporting Issues

When reporting issues, please include your reader and writer schemas, and the output from the compiler logs by adding this to one of your source files:

```
import (
	"github.com/actgardner/gogen-avro/v10/compiler"
)

func init() {
	compiler.LoggingEnabled = true
}
```

The logs will be printed on stdout.

### Alternatives

This project is designed to generate type-safe code from Avro schemas, and to handle schema evolution according to the Avro specification. In cases where the application logic is tied to the Avro schema this makes developing code easier and less error-prone by removing run-time type casting.

There are cases where code generation may not be appropriate - if you need to handle many different schemas in a generic way, or if the schema data isn't available at compile-time. 

[goavro library](https://github.com/linkedin/goavro) is an alternative library that lacks support for evolution, but provides a comprehensive interface for generic data.

[hamba/avro](https://github.com/hamba/avro) is another library oriented towards high performance for generic data.

[nrwiersma/avro-benchmarks](https://github.com/nrwiersma/avro-benchmarks) provides benchmarks of several alternative Avro libraries.
