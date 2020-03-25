package avro

import (
	"bytes"
	"testing"

	evolution "github.com/actgardner/gogen-avro/test/evolve-union/evolution"

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

	newUnionRecord, err := evolution.DeserializeUnionRecordFromSchema(&buf, NewUnionRecord().Schema())
	assert.Nil(t, err)
	assert.Equal(t, evolution.UnionNullStringTypeEnumString, newUnionRecord.A.UnionType)
	assert.Equal(t, "hi", newUnionRecord.A.String)

	assert.Equal(t, "abcd", newUnionRecord.Name)
}
