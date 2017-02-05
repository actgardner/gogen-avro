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

type NestedRecord struct {
	StringField string
	BoolField   bool
	BytesField  []byte
}

func DeserializeNestedRecord(r io.Reader) (*NestedRecord, error) {
	return readNestedRecord(r)
}

func (r *NestedRecord) Schema() string {
	return "{\"fields\":[{\"name\":\"StringField\",\"type\":\"string\"},{\"name\":\"BoolField\",\"type\":\"boolean\"},{\"name\":\"BytesField\",\"type\":\"bytes\"}],\"name\":\"NestedRecord\",\"type\":\"record\"}"
}

func (r *NestedRecord) Serialize(w io.Writer) error {
	return writeNestedRecord(r, w)
}
