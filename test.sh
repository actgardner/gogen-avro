#!/bin/bash -x

set -e

run_test() {
	echo "Running test $t"
	go generate -v gopkg.in/alanctgardner/gogen-avro.v5/$1
	go get -t -v gopkg.in/alanctgardner/gogen-avro.v5/$1
	go test -v gopkg.in/alanctgardner/gogen-avro.v5/$1
}

go install gopkg.in/alanctgardner/gogen-avro.v5/gogen-avro

if [ $# -eq 0 ]; then
	for t in test/*/; do
		run_test $t	
	done
else
	run_test test/$1
fi
