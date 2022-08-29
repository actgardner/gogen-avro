package avro

import (
	"bytes"
	"testing"

	"github.com/actgardner/gogen-avro/v10/compiler"
	evolution "github.com/actgardner/gogen-avro/v10/test/default-union/evolution"
	"github.com/actgardner/gogen-avro/v10/vm"

	"github.com/stretchr/testify/assert"
)

func TestEvolution(t *testing.T) {
	oldUnionRecord := NewUnionRecord()
	var buf bytes.Buffer
	err := oldUnionRecord.Serialize(&buf)
	assert.Nil(t, err)

	newUnionRecord := evolution.NewUnionRecord()

	deser, err := compiler.CompileSchemaBytes([]byte(oldUnionRecord.Schema()), []byte(newUnionRecord.Schema()))
	assert.Nil(t, err)

	err = vm.Eval(bytes.NewReader(buf.Bytes()), deser, &newUnionRecord)
	assert.Nil(t, err)

	assert.Nil(t, nil, newUnionRecord.UnionNull)
	assert.Equal(t, evolution.UnionStringIntTypeEnumString, newUnionRecord.UnionString.UnionType)
	assert.Equal(t, "hello", newUnionRecord.UnionString.String)
	assert.Equal(t, evolution.UnionUnionRecStringTypeEnumUnionRec, newUnionRecord.UnionRecord.UnionType)
	assert.Equal(t, int32(1), newUnionRecord.UnionRecord.UnionRec.A)
}
