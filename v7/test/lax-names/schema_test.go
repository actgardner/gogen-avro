package avro

import (
	"testing"

	"github.com/actgardner/gogen-avro/v7/compiler"

	"github.com/stretchr/testify/assert"
)

func TestLaxNames(t *testing.T) {
	record := NewRecordSchema()

	// Should throw an error by default
	_, err := compiler.CompileSchemaBytes([]byte(record.Schema()), []byte(record.Schema()))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Incompatible types by name")

	// With AllowLaxNames() option it should compile schemas successfully
	program, err := compiler.CompileSchemaBytes(
		[]byte(record.Schema()),
		[]byte(record.Schema()),
		compiler.AllowLaxNames(),
	)
	assert.Nil(t, err)
	assert.NotNil(t, program)
}
