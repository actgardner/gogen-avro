package avro

import (
	"bytes"
	"github.com/alanctgardner/gogen-avro/container"
	"github.com/linkedin/goavro"
	"testing"
)

/* Round-trip some primitive values through our container file writer and goavro to verify */

func TestNullEncoding(t *testing.T) {
	roundTripWithCodec(container.Null, t)
}

func TestSnappyEncoding(t *testing.T) {
	roundTripWithCodec(container.Deflate, t)
}

func TestDeflateEncoding(t *testing.T) {
	roundTripWithCodec(container.Snappy, t)
}

func roundTripWithCodec(codec container.Codec, t *testing.T) {
	var buf bytes.Buffer
	// Write the container file contents to the buffer
	containerWriter, err := container.NewWriter(&buf, codec, 2)
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

	reader, err := goavro.NewReader(goavro.FromReader(&buf))
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
