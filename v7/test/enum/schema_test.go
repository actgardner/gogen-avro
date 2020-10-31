package avro

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/actgardner/gogen-avro/v7/container"
	"github.com/actgardner/gogen-avro/v7/test"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTripExactBytes(t, func() container.AvroRecord { return &EnumTestRecord{} }, func(r io.Reader) (interface{}, error) {
		return DeserializeEnumTestRecord(r)
	})
}

func TestInvalidStringConversion(t *testing.T) {
	enumified, err := NewTestEnumTypeValue("bogus")
	assert.Error(t, err)
	assert.Equal(t, TestEnumType(-1), enumified)
}
