package avro

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"

	"github.com/actgardner/gogen-avro/v7/compiler"
	"github.com/actgardner/gogen-avro/v7/parser"
	"github.com/actgardner/gogen-avro/v7/resolver"
	"github.com/actgardner/gogen-avro/v7/test"
	"github.com/actgardner/gogen-avro/v7/vm"

	"github.com/linkedin/goavro/v2"
	"github.com/stretchr/testify/assert"
)

func TestRoundTrip(t *testing.T) {
	test.RoundTrip(t, &PrimitiveTestRecord{}, func(r io.Reader) (interface{}, error) {
		return DeserializePrimitiveTestRecord(r)
	})
}

func BenchmarkSerializePrimitiveRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	record := PrimitiveTestRecord{1, 2, 3.4, 5.6, "789", true, []byte{1, 2, 3, 4}}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := record.Serialize(buf)
		assert.Nil(b, err)
	}
}

func BenchmarkSerializePrimitiveGoavro(b *testing.B) {
	schemaJson, err := ioutil.ReadFile("primitives.avsc")
	assert.Nil(b, err)

	codec, err := goavro.NewCodec(string(schemaJson))
	assert.Nil(b, err)

	someRecord := map[string]interface{}{
		"IntField":    int32(1),
		"LongField":   int64(2),
		"FloatField":  float32(3.4),
		"DoubleField": float64(5.6),
		"StringField": "789",
		"BoolField":   true,
		"BytesField":  []byte{1, 2, 3, 4},
	}
	buf := make([]byte, 0, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := codec.BinaryFromNative(buf, someRecord)
		assert.Nil(b, err)
	}
}

func BenchmarkDeserializePrimitiveRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	record := PrimitiveTestRecord{1, 2, 3.4, 5.6, "789", true, []byte{1, 2, 3, 4}}
	err := record.Serialize(buf)
	assert.Nil(b, err)

	recordBytes := buf.Bytes()

	schemaJson, err := ioutil.ReadFile("primitives.avsc")
	assert.Nil(b, err)

	readerNs := parser.NewNamespace(false)
	readerType, err := readerNs.TypeForSchema(schemaJson)
	assert.Nil(b, err)

	for _, def := range readerNs.Roots {
		err = resolver.ResolveDefinition(def, readerNs.Definitions)
		assert.Nil(b, err)
	}

	deser, err := compiler.Compile(readerType, readerType)
	assert.Nil(b, err)

	var target PrimitiveTestRecord

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := vm.Eval(bytes.NewReader(recordBytes), deser, &target)
		assert.Nil(b, err)
	}
}

func BenchmarkDeserializePrimitiveGoavro(b *testing.B) {
	schemaJson, err := ioutil.ReadFile("primitives.avsc")
	assert.Nil(b, err)

	codec, err := goavro.NewCodec(string(schemaJson))
	assert.Nil(b, err)

	someRecord := map[string]interface{}{
		"IntField":    int32(1),
		"LongField":   int64(2),
		"FloatField":  float32(3.4),
		"DoubleField": float64(5.6),
		"StringField": "789",
		"BoolField":   true,
		"BytesField":  []byte{1, 2, 3, 4},
	}

	buf, err := codec.BinaryFromNative(nil, someRecord)
	assert.Nil(b, err)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _, err := codec.NativeFromBinary(buf)
		assert.Nil(b, err)
	}
}
