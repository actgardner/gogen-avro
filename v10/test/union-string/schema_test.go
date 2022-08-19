package avro

import (
	"io"
	"testing"

	"github.com/actgardner/gogen-avro/v10/container"
	"github.com/actgardner/gogen-avro/v10/test"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTripExactBytes(t,
		func() container.AvroRecord { return &Record_v1{} },
		func(r io.Reader) (container.AvroRecord, error) {
			record, err := DeserializeRecord_v1(r)
			return &record, err
		})
}

func TestRoundTripString(t *testing.T) {
  test.RoundTripExactBytes(t,
    func() container.AvroRecord {
			return &Record_v1{
				Foo: NewUnionNullStringString("MyFooValue"),
			}
		},
    func(r io.Reader) (container.AvroRecord, error) {
      record, err := DeserializeRecord_v1(r)
      return &record, err
    })
}
