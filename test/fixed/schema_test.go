package avro

import (
	"bytes"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
var fixtures = []FixedTestRecord{
	{
		FixedField: [12]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		FixedField: [12]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255},
	},
	{
		FixedField: [12]byte{0, 1, 3, 7, 15, 31, 63, 127, 255, 0, 2, 128},
	},
}

func TestFixedFixture(t *testing.T) {
	schemaJson, err := ioutil.ReadFile("fixed.avsc")
	if err != nil {
		t.Fatal(err)
	}
	codec, err := goavro.NewCodec(string(schemaJson))
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		if err != nil {
			t.Fatal(err)
		}
		datum, _, err := codec.NativeFromBinary(buf.Bytes())
		if err != nil {
			t.Fatal(err)
		}
		record := datum.(map[string]interface{})
		recordVal, ok := record["FixedField"]
		if !ok {
			t.Fatalf("GOT: %#v; WANT: %#v", ok, true)
		}
		if !reflect.DeepEqual(recordVal, ([]byte)((f.FixedField)[:])) {
			t.Fatalf("FixedField %v is not equal to %v", recordVal.([]byte), ([]byte)((f.FixedField)[:]))
		}
	}
}

func TestRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err := f.Serialize(&buf)
		if err != nil {
			t.Fatal(err)
		}
		datum, err := DeserializeFixedTestRecord(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *datum, f)
	}
}
