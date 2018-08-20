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
{"Header": {"UnionType": 0}},
{"Header": {"UnionType": 1, "CoreHeader": {"UUID": {"UnionType": 0}, "Hostname": {"UnionType": 0}, "Trace": {"UnionType": 0}}}},
{"Header": {"UnionType": 1, "CoreHeader": {"UUID": {"UnionType": 0}, "Hostname": {"UnionType": 1, "String": "HostnameString"}, "Trace": {"UnionType": 0}}}}
]
`

func TestRoundTrip(t *testing.T) {
	fixtures := make([]Test, 0)
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
		datum, err := DeserializeTest(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, *datum, f)
	}
}
