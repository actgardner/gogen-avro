package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v7/container"
	"github.com/actgardner/gogen-avro/v7/test"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTrip(t, func() container.AvroRecord { return &Parent{} }, func(r io.Reader) (interface{}, error) {
		return DeserializeParent(r)
	})
}
