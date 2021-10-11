package generic

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrimitive(t *testing.T) {
	s := []byte(`
{
	"type": "record",
	"name": "primitive",
	"fields": [
		{
			"name": "intfield",
			"type": "int",
			"default": 1
		},
		{
			"name": "bytesfield",
			"type": "bytes",
			"default": "\u00fe"
		}
	]
}
`)

	r := bytes.NewBuffer([]byte{2, 4, byte('h'), byte('i')})

	codec, err := NewCodecFromSchema(s, s)
	assert.NoError(t, err)

	datum, err := codec.Deserialize(r)
	assert.NoError(t, err)

	assert.Equal(t, map[string]interface{}{"intfield": int32(1), "bytesfield": []byte{'h', 'i'}}, datum)
}

func TestComplex(t *testing.T) {
	s := []byte(`
{
	"type": "record",
	"name": "complex",
	"fields": [
		{
			"name": "mapfield",
			"type": {
				"type": "map",
				"values": "int"
			}
		},
		{
			"name": "arrayfield",
			"type": {
				"type": "array",
				"items": "string"
			}
		}
	]
}
`)
	r := bytes.NewBuffer([]byte{2, 6, 'k', 'e', 'y', 7, 0, 4, 2, 'a', 2, 'b', 0})
	codec, err := NewCodecFromSchema(s, s)
	assert.NoError(t, err)

	datum, err := codec.Deserialize(r)
	assert.NoError(t, err)

	assert.Equal(t, map[string]interface{}{"mapfield": map[string]interface{}{"key": int32(-4)}, "arrayfield": []interface{}{"a", "b"}}, datum)
}

func TestUnion(t *testing.T) {
	s := []byte(`
[
	"int",
	"long",
	"string"
]
`)
	codec, err := NewCodecFromSchema(s, s)
	assert.NoError(t, err)

	for _, f := range []struct {
		data     []byte
		expected interface{}
	}{
		{data: []byte{0, 2}, expected: int32(1)},
		{data: []byte{2, 4}, expected: int64(2)},
		{data: []byte{4, 4, 'h', 'i'}, expected: string("hi")},
	} {
		r := bytes.NewBuffer(f.data)
		datum, err := codec.Deserialize(r)
		assert.NoError(t, err)

		assert.Equal(t, f.expected, datum)
	}
}

func TestEnum(t *testing.T) {
	s := []byte(`
{
	"type": "enum",
	"name": "testenum",
	"symbols": [
		"ONE",
		"TWO",
		"THREE"
	]
}
`)
	codec, err := NewCodecFromSchema(s, s)
	assert.NoError(t, err)

	for _, f := range []struct {
		data     []byte
		expected interface{}
	}{
		{data: []byte{2}, expected: "TWO"},
		{data: []byte{4}, expected: "THREE"},
	} {
		r := bytes.NewBuffer(f.data)
		datum, err := codec.Deserialize(r)
		assert.NoError(t, err)

		assert.Equal(t, f.expected, datum)
	}
}

func TestLinkedList(t *testing.T) {
	s := []byte(`
{
	"type": "record",
	"name": "elem",
	"fields": [
		{
			"name": "next",
			"type": ["null", "elem"]
		},
		{
			"name": "val",
			"type": "int"
		}
	]
}
`)
	codec, err := NewCodecFromSchema(s, s)
	assert.NoError(t, err)

	for _, f := range []struct {
		data     []byte
		expected interface{}
	}{
		{data: []byte{0, 1}, expected: map[string]interface{}{"next": nil, "val": int32(-1)}},
		{data: []byte{2, 0, 4, 2}, expected: map[string]interface{}{"next": map[string]interface{}{"next": nil, "val": int32(2)}, "val": int32(1)}},
	} {
		r := bytes.NewBuffer(f.data)
		datum, err := codec.Deserialize(r)
		assert.NoError(t, err)

		assert.Equal(t, f.expected, datum)
	}
}
