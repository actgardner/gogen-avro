#/bin/bash

go install github.com/alanctgardner/gogen-avro

# Generate structs for the Avro schema
$GOPATH/bin/gogen-avro test/primitive/primitives.avsc > test/primitive/schema.go

# Run the unit tests to round-trip data with this schema
go test github.com/alanctgardner/gogen-avro/test/primitive
