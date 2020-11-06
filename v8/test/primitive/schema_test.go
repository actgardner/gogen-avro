package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v8/container"
	"github.com/actgardner/gogen-avro/v8/test"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTrip(t,
		func() container.AvroRecord { return &PrimitiveTestRecord{} },
		func(r io.Reader) (container.AvroRecord, error) {
			return DeserializePrimitiveTestRecord(r)
		})
}
