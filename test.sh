#!/bin/bash -x

set -e

go install github.com/alanctgardner/gogen-avro

for t in test/*/; do
	echo "Running test $t"
	go generate -v github.com/alanctgardner/gogen-avro/$t
	go get -t -v github.com/alanctgardner/gogen-avro/$t
	go test  github.com/alanctgardner/gogen-avro/$t
done
