package avro

import (
	"bytes"
	"testing"

	"github.com/actgardner/gogen-avro/v7/compiler"
	evolution "github.com/actgardner/gogen-avro/v7/test/alias-field/evolution"
	"github.com/actgardner/gogen-avro/v7/vm"

	"github.com/stretchr/testify/assert"
)

func TestEvolution(t *testing.T) {
	oldAliasRecord := NewAliasRecord()
	oldAliasRecord.A = "hi"
	oldAliasRecord.C = "bye"

	var buf bytes.Buffer
	err := oldAliasRecord.Serialize(&buf)
	assert.Nil(t, err)

	newAliasRecord := evolution.NewAliasRecord()

	deser, err := compiler.CompileSchemaBytes([]byte(oldAliasRecord.Schema()), []byte(newAliasRecord.Schema()))
	assert.Nil(t, err)

	err = vm.Eval(bytes.NewReader(buf.Bytes()), deser, newAliasRecord)
	assert.Nil(t, err)

	assert.Equal(t, "hi", newAliasRecord.B)
	assert.Equal(t, "bye", newAliasRecord.D)
}
