/*
 * CODE GENERATED AUTOMATICALLY WITH github.com/alanctgardner/gogen-avro
 * THIS FILE SHOULD NOT BE EDITED BY HAND
 *
 * SOURCE:
 *     nested.avsc
 */
package avro

import (
	"io"
)

type NumberRecord struct {
	IntField    int32
	LongField   int64
	FloatField  float32
	DoubleField float64
}

func DeserializeNumberRecord(r io.Reader) (*NumberRecord, error) {
	return readNumberRecord(r)
}

func (r *NumberRecord) Schema() string {
	return "{\"fields\":[{\"name\":\"IntField\",\"type\":\"int\"},{\"name\":\"LongField\",\"type\":\"long\"},{\"name\":\"FloatField\",\"type\":\"float\"},{\"name\":\"DoubleField\",\"type\":\"double\"}],\"name\":\"NumberRecord\",\"type\":\"record\"}"
}

func (r *NumberRecord) Serialize(w io.Writer) error {
	return writeNumberRecord(r, w)
}
