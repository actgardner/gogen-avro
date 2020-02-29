package avro

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/actgardner/gogen-avro/container"
	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our container file writer and goavro to verify

func TestNullEncoding(t *testing.T) {
	roundTripWithCodec(container.Null, t)
}

func TestDeflateEncoding(t *testing.T) {
	roundTripWithCodec(container.Deflate, t)
}

func TestSnappyEncoding(t *testing.T) {
	roundTripWithCodec(container.Snappy, t)
}

// Round-trip some primitive values through our container file writer and reader
func TestGogenNullEncoding(t *testing.T) {
	roundTripGogenWithCodec(container.Null, t)
}

func TestGogenDeflateEncoding(t *testing.T) {
	roundTripGogenWithCodec(container.Deflate, t)
}

func TestGogenSnappyEncoding(t *testing.T) {
	roundTripGogenWithCodec(container.Snappy, t)
}

func roundTripWithCodec(codec container.Codec, t *testing.T) {
	fixtures := make([]PrimitiveTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	// Write the container file contents to the buffer
	var containerWriter *container.Writer
	containerWriter, err = NewPrimitiveTestRecordWriter(&buf, codec, 2)
	assert.Nil(t, err)

	for _, f := range fixtures {
		// Write the record to the container file
		err = containerWriter.WriteRecord(&f)
		assert.Nil(t, err)
	}

	// Flush the buffers to ensure the last block has been written
	err = containerWriter.Flush()
	assert.Nil(t, err)

	reader, err := goavro.NewOCFReader(bytes.NewReader(buf.Bytes()))
	assert.Nil(t, err)

	var i int
	for reader.Scan() {
		datum, err := reader.Read()
		assert.Nil(t, err)

		compareFixtureGoAvro(t, datum, fixtures[i])
		i = i + 1
	}
	assert.Equal(t, i, len(fixtures))
}

func roundTripGogenWithCodec(codec container.Codec, t *testing.T) {
	fixtures := make([]PrimitiveTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	// Write the container file contents to the buffer
	var containerWriter *container.Writer
	containerWriter, err = NewPrimitiveTestRecordWriter(&buf, codec, 2)
	assert.Nil(t, err)

	for _, f := range fixtures {
		// Write the record to the container file
		err = containerWriter.WriteRecord(&f)
		assert.Nil(t, err)
	}

	// Flush the buffers to ensure the last block has been written
	err = containerWriter.Flush()
	if err != nil {
		t.Fatal(err)
	}

	reader, err := NewPrimitiveTestRecordReader(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatal(err)
	}

	for i := range fixtures {
		record, err := reader.Read()
		assert.Nil(t, err)
		assert.Equal(t, record, fixtures[i])
	}
}
