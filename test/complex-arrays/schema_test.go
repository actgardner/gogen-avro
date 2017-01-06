package avro

import (
	"bytes"
	"encoding/json"
	"github.com/linkedin/goavro"
	avrotest "github.com/alanctgardner/gogen-avro/test"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

/* Round trip some records nested in arrays */
const fixtureJson = `
[
  {"Children": []},
  {"Children": [{"Name": "test-record"}]},
  {"Children": [{"Name": "test-record"}, {"Name": "test-record-2"}]}
]
`

func TestPrimitiveFixture(t *testing.T) {
	fixtures := make([]Parent, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	schemaJson, err := ioutil.ReadFile("arrays.avsc")
	if err != nil {
		t.Fatal(err)
	}

	codec, err := goavro.NewCodec(string(schemaJson))
	if err != nil {
		t.Fatal(err)
	}

	for _, f := range fixtures {
		avrotest.RoundTripGoAvroTest(&f, codec, t)
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]Parent, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
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
		datum, err := DeserializeParent(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *datum, f)
	}
}
