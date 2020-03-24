package avro

import (
	"bytes"
	"github.com/actgardner/gogen-avro/soe"
	"testing"

	"github.com/actgardner/gogen-avro/container"
	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

func TestNullEncoding(t *testing.T) {
	roundTripWithCodec(container.Null, t)
}

func TestSnappyEncoding(t *testing.T) {
	roundTripWithCodec(container.Deflate, t)
}

func TestDeflateEncoding(t *testing.T) {
	roundTripWithCodec(container.Snappy, t)
}

func TestGogenNullEncoding(t *testing.T) {
	roundTripGogenWithCodec(container.Null, t)
}

func TestGogenSnappyEncoding(t *testing.T) {
	roundTripGogenWithCodec(container.Deflate, t)
}

func TestGogenDeflateEncoding(t *testing.T) {
	roundTripGogenWithCodec(container.Snappy, t)
}

func roundTripWithCodec(codec container.Codec, t *testing.T) {
	var buf bytes.Buffer
	UID := make([]byte, 8)
	writer := soe.NewWriter(&buf, UID)
	// Write the container file contents to the buffer
	containerWriter, err := NewEventWriter(writer, codec, 2)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range fixtures {
		// Write the record to the container file
		err = containerWriter.WriteRecord(&f)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Flush the buffers to ensure the last block has been written
	err = containerWriter.Flush()
	if err != nil {
		t.Fatal(err)
	}

	reader, err := goavro.NewOCFReader(soe.NewReader(&buf))
	if err != nil {
		t.Fatal(err)
	}

	var i int
	for reader.Scan() {
		datum, err := reader.Read()
		if err != nil {
			t.Fatal(err)
		}
		compareFixtureGoAvro(t, datum, fixtures[i])
		i = i + 1
	}
}

func roundTripGogenWithCodec(codec container.Codec, t *testing.T) {
	var buf bytes.Buffer
	// Write the container file contents to the buffer
	containerWriter, err := NewEventWriter(&buf, codec, 2)
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range fixtures {
		// Write the record to the container file
		err = containerWriter.WriteRecord(&f)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Flush the buffers to ensure the last block has been written
	err = containerWriter.Flush()
	if err != nil {
		t.Fatal(err)
	}

	reader, err := NewEventReader(&buf)
	if err != nil {
		t.Fatal(err)
	}

	for i := range fixtures {
		datum, err := reader.Read()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, datum, &fixtures[i])
	}
}
