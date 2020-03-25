package avro

import (
	"bytes"
	"github.com/actgardner/gogen-avro/soe"
	"testing"

	"github.com/actgardner/gogen-avro/compiler"
	evolution "github.com/actgardner/gogen-avro/test/default-union/evolution"
	"github.com/actgardner/gogen-avro/vm"

	"github.com/stretchr/testify/assert"
)

func TestEvolution(t *testing.T) {
	oldUnionRecord := NewUnionRecord()
	var buf bytes.Buffer
	writer := soe.NewWriter(&buf,  oldUnionRecord.AvroCRC64Fingerprint())
	err := oldUnionRecord.Serialize(writer)
	assert.Nil(t, err)

	newUnionRecord := evolution.NewUnionRecord()

	deser, err := compiler.CompileSchemaBytes([]byte(oldUnionRecord.Schema()), []byte(newUnionRecord.Schema()))
	assert.Nil(t, err)

	err = vm.Eval(soe.NewReader(&buf), deser, newUnionRecord)
	assert.Nil(t, err)

	assert.Equal(t, evolution.UnionNullStringTypeEnumNull, newUnionRecord.UnionNull.UnionType)
	assert.Equal(t, evolution.UnionStringIntTypeEnumString, newUnionRecord.UnionString.UnionType)
	assert.Equal(t, "hello", newUnionRecord.UnionString.String)
	assert.Equal(t, evolution.UnionUnionRecStringTypeEnumUnionRec, newUnionRecord.UnionRecord.UnionType)
	assert.Equal(t, int32(1), newUnionRecord.UnionRecord.UnionRec.A)
}
