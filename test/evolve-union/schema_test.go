package avro

import (
	"bytes"
	"github.com/actgardner/gogen-avro/soe"
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
	writer := soe.NewWriter(&buf, oldUnionRecord.AvroCRC64Fingerprint())
	err := oldUnionRecord.Serialize(writer)
	assert.Nil(t, err)

	newUnionRecord, err := evolution.DeserializeUnionRecordFromSchema(soe.NewReader(&buf), NewUnionRecord().Schema())
	assert.Nil(t, err)
	assert.Equal(t, evolution.UnionNullStringTypeEnumString, newUnionRecord.A.UnionType)
	assert.Equal(t, "hi", newUnionRecord.A.String)

	assert.Equal(t, "abcd", newUnionRecord.Name)
}
