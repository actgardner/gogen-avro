package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v8/container"
	"github.com/actgardner/gogen-avro/v8/test"
	evolution "github.com/actgardner/gogen-avro/v8/test/alias-fixed/evolution"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTripExactBytes(t,
		func() container.AvroRecord { return &UnionIp_addressEvent{} },
		func(r io.Reader) (container.AvroRecord, error) {
			return DeserializeUnionIp_addressEvent(r)
		})
}

func TestEvolution(t *testing.T) {
	test.RoundTripEvolution(t,
		func() container.AvroRecord { return &UnionIp_addressEvent{} },
		func() container.AvroRecord { return &evolution.UnionIPAddressEvent{} },
		func(r io.Reader, schema string) (container.AvroRecord, error) {
			return evolution.DeserializeUnionIPAddressEventFromSchema(r, schema)
		})
}
