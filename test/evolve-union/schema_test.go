package avro

import (
	"bytes"
	"testing"

	"github.com/actgardner/gogen-avro/compiler"
	evolution "github.com/actgardner/gogen-avro/test/evolve-union/evolution"
	"github.com/actgardner/gogen-avro/vm"

	"github.com/stretchr/testify/assert"
)

func TestEvolution(t *testing.T) {
	oldUnionRecord := NewUnionRecord()
	oldUnionRecord.Id = &UnionNullInt{UnionType: UnionNullIntTypeEnumInt, Int: 1}
	oldUnionRecord.A = "hi"
	oldUnionRecord.Name = &UnionNullString{UnionType: UnionNullStringTypeEnumString, String: "abcd"}

	var buf bytes.Buffer
	err := oldUnionRecord.Serialize(&buf)
	assert.Nil(t, err)

	newUnionRecord := evolution.NewUnionRecord()

	deser, err := compiler.CompileSchemaBytes([]byte(oldUnionRecord.Schema()), []byte(newUnionRecord.Schema()))
	assert.Nil(t, err)

	err = vm.Eval(bytes.NewReader(buf.Bytes()), deser, newUnionRecord)
	assert.Nil(t, err)

	assert.Equal(t, evolution.UnionNullStringTypeEnumString, newUnionRecord.A.UnionType)
	assert.Equal(t, "hi", newUnionRecord.A.String)

	assert.Equal(t, "abcd", newUnionRecord.Name)
}
