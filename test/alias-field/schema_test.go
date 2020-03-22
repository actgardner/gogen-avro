package avro

import (
	"bytes"
	"github.com/actgardner/gogen-avro/singleobject"
	"testing"

	"github.com/actgardner/gogen-avro/compiler"
	evolution "github.com/actgardner/gogen-avro/test/alias-field/evolution"
	"github.com/actgardner/gogen-avro/vm"

	"github.com/stretchr/testify/assert"
)

func TestEvolution(t *testing.T) {
	oldAliasRecord := NewAliasRecord()
	oldAliasRecord.A = "hi"
	oldAliasRecord.C = "bye"

	var buf bytes.Buffer
	writer := singleobject.NewWriter(&buf, AliasRecordUID)
	err := oldAliasRecord.Serialize(writer)
	assert.Nil(t, err)

	newAliasRecord := evolution.NewAliasRecord()

	deser, err := compiler.CompileSchemaBytes([]byte(oldAliasRecord.Schema()), []byte(newAliasRecord.Schema()))
	assert.Nil(t, err)

	err = vm.Eval(singleobject.NewReader(&buf), deser, newAliasRecord)
	assert.Nil(t, err)

	assert.Equal(t, "hi", newAliasRecord.B)
	assert.Equal(t, "bye", newAliasRecord.D)
}
