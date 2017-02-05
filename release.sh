#!/bin/bash

# Rewrite references from github.com/alanctgardner/gogen-avro to gopkg.in/alanctgardner/gogen-avro.<version>

if [ "$#" -ne 1 ]; then
  echo "Usage: $0 <version>"
  exit 1
fi
 
GITHUB_REPO="github.com/alanctgardner/gogen-avro"
VERSION="$1"
GOPKG_REPO="gopkg.in/alanctgardner/gogen-avro.$VERSION"

sed -i "s|$GITHUB_REPO|$GOPKG_REPO|" container/*.go generator/*.go types/*.go main.go example/*/*.go test.sh 
