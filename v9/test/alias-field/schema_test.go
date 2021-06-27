package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v9/container"
	"github.com/actgardner/gogen-avro/v9/test"
	evolution "github.com/actgardner/gogen-avro/v9/test/alias-field/evolution"
)

func TestEvolution(t *testing.T) {
	test.RoundTripEvolution(t,
		func() container.AvroRecord { return &AliasRecord{} },
		func() container.AvroRecord { return &evolution.AliasRecord{} },
		func(r io.Reader, schema string) (container.AvroRecord, error) {
			record, err := evolution.DeserializeAliasRecordFromSchema(r, schema)
			return &record, err
		})
}
