package avro

import (
	"bytes"
	"testing"

	evolution "github.com/actgardner/gogen-avro/v10/test/evolve-union/evolution"

	"github.com/stretchr/testify/assert"
)

func TestEvolution(t *testing.T) {
	oldUnionRecord := NewUnionRecord()
	var id int32 = 1
	oldUnionRecord.Id = &id
	oldUnionRecord.A = "hi"
	var name = "abcd"
	oldUnionRecord.Name = &name

	var buf bytes.Buffer
	err := oldUnionRecord.Serialize(&buf)
	assert.Nil(t, err)

	newUnionRecord, err := evolution.DeserializeUnionRecordFromSchema(&buf, oldUnionRecord.Schema())
	assert.Nil(t, err)
	assert.Equal(t, "hi", *newUnionRecord.A)

	assert.Equal(t, "abcd", newUnionRecord.Name)
}
