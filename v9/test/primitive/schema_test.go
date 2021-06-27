package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v9/container"
	"github.com/actgardner/gogen-avro/v9/test"
	evolution "github.com/actgardner/gogen-avro/v9/test/primitive/evolution"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTrip(t,
		func() container.AvroRecord { return &PrimitiveTestRecord{} },
		func(r io.Reader) (container.AvroRecord, error) {
			record, err := DeserializePrimitiveTestRecord(r)
			return &record, err
		})
}

func TestEvolution(t *testing.T) {
	test.RoundTripEvolution(t,
		func() container.AvroRecord { return &PrimitiveTestRecord{} },
		func() container.AvroRecord { return &evolution.PrimitiveTestRecord{} },
		func(r io.Reader, schema string) (container.AvroRecord, error) {
			record, err := evolution.DeserializePrimitiveTestRecordFromSchema(r, schema)
			return &record, err
		})
}
