package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v7/container"
	"github.com/actgardner/gogen-avro/v7/test"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTripExactBytes(t,
		func() container.AvroRecord { return &UnionIp_addressEvent{} },
		func(r io.Reader) (container.AvroRecord, error) {
			return DeserializeUnionIp_addressEvent(r)
		})
}
