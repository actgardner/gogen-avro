package avro

import (
	"bytes"
	"github.com/alanctgardner/gogen-avro/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	moreFields          = &MoreFields{}
	lessFields          = &LessFields{}
	moreFieldsSchema, _ = types.AvroTypeFromString(moreFields.Schema())
	lessFieldsSchema, _ = types.AvroTypeFromString(lessFields.Schema())
)

func TestAddRecords(t *testing.T) {
	oldRecord := &LessFields{2, 12.34, "oldstring"}
	var buf bytes.Buffer
	err := oldRecord.Serialize(&buf)
	if err != nil {
		t.Fatal(err)
	}

	datum, err := DeserializeMoreFields(lessFieldsSchema, &buf)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, datum.IntField, int32(2))
	assert.Equal(t, datum.LongField, int64(2))
	assert.Equal(t, datum.FloatField, float32(12.34))
	assert.Equal(t, datum.DoubleField, float64(5.6))
	assert.Equal(t, datum.StringField, "oldstring")
	assert.Equal(t, datum.BoolField, true)
}
