package avro

import (
	"bytes"
	"testing"

	"github.com/actgardner/gogen-avro/container"
	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our container file writer and goavro to verify

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

	reader, err := goavro.NewOCFReader(bytes.NewReader(buf.Bytes()))
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

	reader, err := NewEventReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	for i := range fixtures {
		datum, err := reader.Read()
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, datum, fixtures[i])
	}
}
