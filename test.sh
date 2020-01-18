#!/bin/bash -x

set -e

run_test() {
	echo "Running test $t"
	go generate -v github.com/actgardner/gogen-avro/v7/$1
	go get -t -v github.com/actgardner/gogen-avro/v7/$1
	go test -v github.com/actgardner/gogen-avro/v7/$1
}

go install github.com/actgardner/gogen-avro/v7/cmd/...

if [ $# -eq 0 ]; then
	for t in test/*/; do
		run_test $t
	done
else
	run_test test/$1
fi
