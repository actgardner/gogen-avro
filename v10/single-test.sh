#!/bin/bash -ex

go install ./cmd/...

go generate ./test/$1
go get -t -v -u ./...
go test ./test/$1
