#!/bin/bash -x

set -e

go get gopkg.in/alanctgardner/gogen-avro.v2
go install gopkg.in/alanctgardner/gogen-avro.v2

for t in test/*/; do
	echo "Running test $t"
	go generate -v gopkg.in/alanctgardner/gogen-avro.v2/$t
	go get -t -v gopkg.in/alanctgardner/gogen-avro.v2/$t
	go test gopkg.in/alanctgardner/gogen-avro.v2/$t
done
