#/bin/bash

go install github.com/alanctgardner/gogen-avro

# Generate structs for the Avro schema
$GOPATH/bin/gogen-avro test/primitive/primitives.avsc test/primitive
$GOPATH/bin/gogen-avro test/maps/maps.avsc test/maps
$GOPATH/bin/gogen-avro test/arrays/arrays.avsc test/arrays
$GOPATH/bin/gogen-avro test/union/union.avsc test/union
$GOPATH/bin/gogen-avro test/nested/nested.avsc test/nested
$GOPATH/bin/gogen-avro test/enum/enum.avsc test/enum
$GOPATH/bin/gogen-avro test/fixed/fixed.avsc test/fixed
$GOPATH/bin/gogen-avro test/complex-union/union.avsc test/complex-union

# Run the unit tests to round-trip data with this schema
go test  github.com/alanctgardner/gogen-avro/test/primitive
go test  github.com/alanctgardner/gogen-avro/test/maps
go test  github.com/alanctgardner/gogen-avro/test/arrays
go test  github.com/alanctgardner/gogen-avro/test/union
go test  github.com/alanctgardner/gogen-avro/test/nested
go test  github.com/alanctgardner/gogen-avro/test/enum
go test  github.com/alanctgardner/gogen-avro/test/fixed
go test  -v github.com/alanctgardner/gogen-avro/test/complex-union
