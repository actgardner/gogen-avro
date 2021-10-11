package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v10/container"
	"github.com/actgardner/gogen-avro/v10/test"
	evolution "github.com/actgardner/gogen-avro/v10/test/alias-record/evolution"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTripExactBytes(t,
		func() container.AvroRecord { return &UnionNestedRecordNestedTestRecord{} },
		func(r io.Reader) (container.AvroRecord, error) {
			record, err := DeserializeUnionNestedRecordNestedTestRecord(r)
			return &record, err
		})
}

func TestEvolution(t *testing.T) {
	test.RoundTripEvolution(t,
		func() container.AvroRecord { return &UnionNestedRecordNestedTestRecord{} },
		func() container.AvroRecord { return &evolution.UnionAliasedRecordNestedTestRecord{} },
		func(r io.Reader, schema string) (container.AvroRecord, error) {
			record, err := evolution.DeserializeUnionAliasedRecordNestedTestRecordFromSchema(r, schema)
			return &record, err
		})
}
