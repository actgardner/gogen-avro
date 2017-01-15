#!/bin/bash -x

set -e

run_test() {
	echo "Running test $t"
	go generate -v gopkg.in/alanctgardner/gogen-avro.v3/$1
	go get -t -v gopkg.in/alanctgardner/gogen-avro.v3/$1
	go test -v gopkg.in/alanctgardner/gogen-avro.v3/$1
}

go get gopkg.in/alanctgardner/gogen-avro.v3
go install gopkg.in/alanctgardner/gogen-avro.v3

if [ $# -eq 0 ]; then
	for t in test/*/; do
		run_test $t	
	done
else
	run_test test/$1
fi
