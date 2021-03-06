package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v9/container"
	"github.com/actgardner/gogen-avro/v9/test"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTrip(t,
		func() container.AvroRecord { return &NestedMap{} },
		func(r io.Reader) (container.AvroRecord, error) {
			record, err := DeserializeNestedMap(r)
			return &record, err
		})
}
