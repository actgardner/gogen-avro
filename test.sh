#/bin/bash

go install github.com/alanctgardner/gogen-avro

# Generate structs for the Avro schema
$GOPATH/bin/gogen-avro test/primitive/primitives.avsc > test/primitive/schema.go
$GOPATH/bin/gogen-avro test/maps/maps.avsc > test/maps/schema.go
$GOPATH/bin/gogen-avro test/arrays/arrays.avsc > test/arrays/schema.go

# Run the unit tests to round-trip data with this schema
go test github.com/alanctgardner/gogen-avro/test/primitive
go test github.com/alanctgardner/gogen-avro/test/maps
go test github.com/alanctgardner/gogen-avro/test/arrays
