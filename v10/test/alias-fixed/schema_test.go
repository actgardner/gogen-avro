package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v10/container"
	"github.com/actgardner/gogen-avro/v10/test"
	evolution "github.com/actgardner/gogen-avro/v10/test/alias-fixed/evolution"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTripExactBytes(t,
		func() container.AvroRecord { return &UnionIp_addressEvent{} },
		func(r io.Reader) (container.AvroRecord, error) {
			record, err := DeserializeUnionIp_addressEvent(r)
			return &record, err
		})
}

func TestEvolution(t *testing.T) {
	test.RoundTripEvolution(t,
		func() container.AvroRecord { return &UnionIp_addressEvent{} },
		func() container.AvroRecord { return &evolution.UnionIPAddressEvent{} },
		func(r io.Reader, schema string) (container.AvroRecord, error) {
			record, err := evolution.DeserializeUnionIPAddressEventFromSchema(r, schema)
			return &record, err
		})
}
