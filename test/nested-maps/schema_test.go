package avro

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
  {
	"MapOfMaps": {
      "m": {
        "NestedMap": {
          "m": {
            "NestedArray": ["some","strings"],
            "NestedArray1": ["other","values"]
          }
        }
      }
	}
  },
  {
	"MapOfMaps": {
      "m": {
        "NestedMap": {
          "m": {
            "NestedArray": ["some","strings"],
            "NestedArray1": ["other","values"]
          }
        }
      }
	}
  },
  {
	"MapOfMaps": {
      "m": {
        "NestedMap": {
          "m": {
            "NestedArray": ["some","strings"],
            "NestedArray1": ["other","values"]
          }
        }
      }
	}
  }
]
`

func BenchmarkMapOfMapsRecord(b *testing.B) {
	buf := new(bytes.Buffer)
	record := NestedMap{
		&MapMapArrayString{
			M: map[string]*MapArrayString{
				"key1": {
					M: map[string][]string{
						"array1": {"value1", "value2"},
						"array2": {"value3", "value4"},
					},
				},
				"key2": {
					M: map[string][]string{
						"array3": {"value5"},
						"array4": {},
					},
				},
			},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := record.Serialize(buf)
		assert.Nil(b, err)
	}
}

func BenchmarkMapOfMapsGoavro(b *testing.B) {
	schemaJson, err := ioutil.ReadFile("nested-maps.avsc")
	assert.Nil(b, err)

	codec, err := goavro.NewCodec(string(schemaJson))
	assert.Nil(b, err)

	someRecord := map[string]interface{}{
		"MapOfMaps": map[string]interface{}{
			"key1": map[string][]string{
				"array1": {"value1", "value2"},
				"array2": {"value3", "value4"},
			},
			"key2": map[string][]string{
				"array3": {"value5"},
				"array4": {},
			},
		},
	}
	buf := make([]byte, 0, 1024)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := codec.BinaryFromNative(buf, someRecord)
		assert.Nil(b, err)
	}
}

func TestRoundTrip(t *testing.T) {
	fixtures := make([]NestedMap, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializeNestedMap(&buf)
		assert.Nil(t, err)
		assert.Equal(t, datum.MapOfMaps, f.MapOfMaps)
	}
}
