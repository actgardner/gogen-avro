package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v8/container"
	"github.com/actgardner/gogen-avro/v8/test"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTrip(t, func() container.AvroRecord { return &ArrayTestRecord{} },
		func(r io.Reader) (container.AvroRecord, error) {
			record, err := DeserializeArrayTestRecord(r)
			return &record, err
		})
}
