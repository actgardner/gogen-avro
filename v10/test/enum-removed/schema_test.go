package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v10/container"
	"github.com/actgardner/gogen-avro/v10/test"
	evolution "github.com/actgardner/gogen-avro/v10/test/enum-removed/evolution"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTripExactBytes(t,
		func() container.AvroRecord { return &EnumTestRecord{} },
		func(r io.Reader) (container.AvroRecord, error) {
			record, err := DeserializeEnumTestRecord(r)
			return &record, err
		})
}

func TestEvolution(t *testing.T) {
	test.RoundTripEvolution(t,
		func() container.AvroRecord { return &EnumTestRecord{} },
		func() container.AvroRecord { return &evolution.EnumTestRecord{} },
		func(r io.Reader, schema string) (container.AvroRecord, error) {
			record, err := evolution.DeserializeEnumTestRecordFromSchema(r, schema)
			return &record, err
		})
}
