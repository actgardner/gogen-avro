package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v10/container"
	"github.com/actgardner/gogen-avro/v10/test"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTripGoGenOnly(t,
		func() container.AvroRecord { return &Parent{} },
		func(r io.Reader) (container.AvroRecord, error) {
			record, err := DeserializeParent(r)
			return &record, err
		})
}
