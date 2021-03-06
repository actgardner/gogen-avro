package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v9/container"
	"github.com/actgardner/gogen-avro/v9/test"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTripExactBytes(t,
		func() container.AvroRecord { return &NestedTestRecord{} },
		func(r io.Reader) (container.AvroRecord, error) {
			record, err := DeserializeNestedTestRecord(r)
			return &record, err
		})
}
