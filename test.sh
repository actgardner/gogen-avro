#/bin/bash

set -e

go install github.com/alanctgardner/gogen-avro

go generate  github.com/alanctgardner/gogen-avro/test/primitive
go test  github.com/alanctgardner/gogen-avro/test/primitive

go generate  github.com/alanctgardner/gogen-avro/test/maps
go test  github.com/alanctgardner/gogen-avro/test/maps

go generate  github.com/alanctgardner/gogen-avro/test/arrays
go test  github.com/alanctgardner/gogen-avro/test/arrays

go generate  github.com/alanctgardner/gogen-avro/test/union
go test  github.com/alanctgardner/gogen-avro/test/union

go generate  github.com/alanctgardner/gogen-avro/test/nested
go test  github.com/alanctgardner/gogen-avro/test/nested

go generate  github.com/alanctgardner/gogen-avro/test/enum
go test  github.com/alanctgardner/gogen-avro/test/enum

go generate  github.com/alanctgardner/gogen-avro/test/fixed
go test  github.com/alanctgardner/gogen-avro/test/fixed

go generate  github.com/alanctgardner/gogen-avro/test/complex-union
go test  github.com/alanctgardner/gogen-avro/test/complex-union

go generate  github.com/alanctgardner/gogen-avro/test/recursive
go test  github.com/alanctgardner/gogen-avro/test/recursive
