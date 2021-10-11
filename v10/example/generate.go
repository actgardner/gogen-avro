package generate

// Use a go:generate directive to build the Go structs for `example.avsc`
// These files are used for all of the example projects
// Source files will be in a package called `example/avro`

//go:generate mkdir -p ./avro
//go:generate $GOPATH/bin/gogen-avro -containers ./avro example.avsc
