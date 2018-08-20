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

func TestSnappyEncoding(t *testing.T) {
	roundTripWithCodec(container.Deflate, t)
}

func TestDeflateEncoding(t *testing.T) {
	roundTripWithCodec(container.Snappy, t)
}

// Test that extra metadata in the schema is included
func TestEventSchemaMetadata(t *testing.T) {
	event := &Event{}
	var eventJson map[string]interface{}
	assert.Nil(t, json.Unmarshal([]byte(event.Schema()), &eventJson))
	metadata, ok := eventJson["metadata"]
	assert.Equal(t, ok, true)
	metadataMap, ok := metadata.(map[string]interface{})
	assert.Equal(t, ok, true)
	assert.Equal(t, metadataMap["key"], "value")
}

func TestIPSchemaMetadata(t *testing.T) {
	event := &Event{}
	var eventJson map[string]interface{}
	assert.Nil(t, json.Unmarshal([]byte(event.Schema()), &eventJson))
	fields, ok := eventJson["fields"]
	assert.Equal(t, ok, true)
	fieldList, ok := fields.([]interface{})
	assert.Equal(t, ok, true)
	field, ok := fieldList[1].(map[string]interface{})
	assert.Equal(t, ok, true)
	typeField, ok := field["type"]
	assert.Equal(t, ok, true)
	typeMap, ok := typeField.(map[string]interface{})
	assert.Equal(t, ok, true)
	metadata, ok := typeMap["metadata"]
	assert.Equal(t, ok, true)
	metadataMap, ok := metadata.(map[string]interface{})
	assert.Equal(t, ok, true)
	assert.Equal(t, metadataMap["a"], "b")
	assert.Equal(t, metadataMap["c"], float64(123))
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
