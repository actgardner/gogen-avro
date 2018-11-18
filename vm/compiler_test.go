package vm

import (
	"testing"

	"github.com/actgardner/gogen-avro/types"
	"github.com/stretchr/testify/assert"
)

func TestCompilePrimitive(t *testing.T) {
	reader := `
{
  "type": "record",
  "name": "test",
  "fields": [
    {"name": "one", "type":"string"},
    {"name": "two", "type":"int"}
  ]
}
`

	writer := `
{
  "type": "record",
  "name": "test",
  "fields": [
    {"name": "two", "type":"int"},
    {"name": "one", "type":"string"}
  ]
}
`

	readerNs := types.NewNamespace(false)
	readerType, err := readerNs.TypeForSchema([]byte(reader))
	assert.Nil(t, err)

	err = readerType.ResolveReferences(readerNs)
	assert.Nil(t, err)

	writerNs := types.NewNamespace(false)
	writerType, err := writerNs.TypeForSchema([]byte(writer))
	assert.Nil(t, err)

	err = writerType.ResolveReferences(writerNs)
	assert.Nil(t, err)

	program, err := Compile(writerType, readerType)

	expectedProgram := []Instruction{
		Instruction{Op: Read, Type: Int, Field: 0},
		Instruction{Op: Set, Type: 3, Field: 1},
		Instruction{Op: Read, Type: String, Field: 0},
		Instruction{Op: Set, Type: 8, Field: 0},
	}
	assert.Equal(t, expectedProgram, program)
	assert.Nil(t, err)
}
