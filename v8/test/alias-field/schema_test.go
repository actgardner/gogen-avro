package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v8/container"
	"github.com/actgardner/gogen-avro/v8/test"
	evolution "github.com/actgardner/gogen-avro/v8/test/alias-field/evolution"
)

func TestEvolution(t *testing.T) {
	test.RoundTripEvolution(t,
		func() container.AvroRecord { return NewAliasRecord() },
		func() container.AvroRecord { return evolution.NewAliasRecord() },
		func(r io.Reader, schema string) (container.AvroRecord, error) {
			return evolution.DeserializeAliasRecordFromSchema(r, schema)
		})
}
