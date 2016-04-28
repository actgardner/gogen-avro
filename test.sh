#/bin/bash

go install github.com/alanctgardner/gogen-avro

# Generate structs for the Avro schema
$GOPATH/bin/gogen-avro test/primitive test/primitive/primitives.avsc
$GOPATH/bin/gogen-avro test/maps test/maps/maps.avsc
$GOPATH/bin/gogen-avro test/arrays test/arrays/arrays.avsc
$GOPATH/bin/gogen-avro test/union test/union/union.avsc
$GOPATH/bin/gogen-avro test/nested test/nested/nested.avsc
$GOPATH/bin/gogen-avro test/enum test/enum/enum.avsc
$GOPATH/bin/gogen-avro test/fixed test/fixed/fixed.avsc
$GOPATH/bin/gogen-avro test/complex-union test/complex-union/union.avsc

# Run the unit tests to round-trip data with this schema
go test  github.com/alanctgardner/gogen-avro/test/primitive
go test  github.com/alanctgardner/gogen-avro/test/maps
go test  github.com/alanctgardner/gogen-avro/test/arrays
go test  github.com/alanctgardner/gogen-avro/test/union
go test  github.com/alanctgardner/gogen-avro/test/nested
go test  github.com/alanctgardner/gogen-avro/test/enum
go test  github.com/alanctgardner/gogen-avro/test/fixed
go test  github.com/alanctgardner/gogen-avro/test/complex-union
