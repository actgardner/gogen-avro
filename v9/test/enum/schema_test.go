package avro

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/actgardner/gogen-avro/v9/container"
	"github.com/actgardner/gogen-avro/v9/test"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTripExactBytes(t,
		func() container.AvroRecord { return &EnumTestRecord{} },
		func(r io.Reader) (container.AvroRecord, error) {
			record, err := DeserializeEnumTestRecord(r)
			return &record, err
		})
}

func TestInvalidStringConversion(t *testing.T) {
	enumified, err := NewTestEnumTypeValue("bogus")
	assert.Error(t, err)
	assert.Equal(t, TestEnumType(-1), enumified)
}
