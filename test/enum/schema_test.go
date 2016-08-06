package avro

import (
	"bytes"
	"encoding/json"
	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

//go:generate $GOPATH/bin/gogen-avro . enum.avsc

/* Round-trip some primitive values through our serializer and goavro to verify */
const fixtureJson = `
[
	{
		"EnumField": 0
	},
	{
		"EnumField": 2
	}
]
`

func TestEnumFixture(t *testing.T) {
	fixtures := make([]EnumTestRecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	if err != nil {
		t.Fatal(err)
	}

	schemaJson, err := ioutil.ReadFile("enum.avsc")
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
		datum, err := codec.Decode(&buf)
		if err != nil {
			t.Fatal(err)
		}
		record := datum.(*goavro.Record)
		recordVal, err := record.Get("EnumField")
		if err != nil {
			t.Fatal(err)
		}
		if recordVal.(goavro.Enum).Value != f.EnumField.String() {
			t.Fatalf("EnumField %v is not equal to %v", recordVal, f.EnumField)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]EnumTestRecord, 0)
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
		datum, err := DeserializeEnumTestRecord(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *datum, f)
	}
}
