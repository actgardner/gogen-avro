package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v8/container"
	"github.com/actgardner/gogen-avro/v8/test"
	evolution "github.com/actgardner/gogen-avro/v8/test/alias-record/evolution"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTripExactBytes(t,
		func() container.AvroRecord { return &UnionNestedRecordNestedTestRecord{} },
		func(r io.Reader) (container.AvroRecord, error) {
			return DeserializeUnionNestedRecordNestedTestRecord(r)
		})
}

func TestEvolution(t *testing.T) {
	test.RoundTripEvolution(t,
		func() container.AvroRecord { return &UnionNestedRecordNestedTestRecord{} },
		func() container.AvroRecord { return &evolution.UnionAliasedRecordNestedTestRecord{} },
		func(r io.Reader) (container.AvroRecord, error) {
			return evolution.DeserializeUnionAliasedRecordNestedTestRecord(r)
		})
}
