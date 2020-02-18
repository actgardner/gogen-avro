package compiler_test

import (
	"testing"

	"github.com/actgardner/gogen-avro/compiler"
)

const s1 = `{
    "name": "R",
    "type": "record",
    "fields": [
        {
            "name": "F1",
            "type": "int"
        }
    ]
}`

const s2 = `{
    "name": "R",
    "type": "record",
    "fields": [
        {
            "name": "F1",
            "type": "string"
        }
    ]
}
`

func TestIncompatibility(t *testing.T) {
	_, err := compiler.CompileSchemaBytes([]byte(s1), []byte(s2))
	if err == nil {
		t.Fatalf("unexpected success compiling incompatible schemas")
	}
}
