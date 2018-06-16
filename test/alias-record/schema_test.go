package avro

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
  {
    "OtherField": {
      "StringField": "789", 
      "BoolField": true, 
      "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"
    }
  },
  {
    "OtherField": {
      "StringField": "abcdghejw", 
      "BoolField": true, 
      "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"
    }
  },
  {
    "OtherField": {
      "StringField": "jdnwjkendwedddedee", 
      "BoolField": true, 
      "BytesField": "VGhpcyBpcyBhIHRlc3Qgc3RyaW5n"
    }
  }
]
`

func TestRoundTrip(t *testing.T) {
	fixtures := make([]NestedTestRecord, 0)
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
		datum, err := DeserializeNestedTestRecord(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *datum, f)
	}
}
