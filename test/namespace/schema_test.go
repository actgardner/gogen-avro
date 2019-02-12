package avro

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
const fixtureJson = `
[
{"Header": {"UnionType": 0}},
{"Header": {"UnionType": 1, "Headerworks_CoreHeader": {"UUID": {"UnionType": 0}, "Hostname": {"UnionType": 0}, "Trace": {"UnionType": 0}}}},
{"Header": {"UnionType": 1, "Headerworks_CoreHeader": {"UUID": {"UnionType": 0}, "Hostname": {"UnionType": 1, "String": "HostnameString"}, "Trace": {"UnionType": 0}}}}
]
`

func TestRoundTrip(t *testing.T) {
	fixtures := make([]Com_avro_test_testrecord, 0)
	err := json.Unmarshal([]byte(fixtureJson), &fixtures)
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = f.Serialize(&buf)
		assert.Nil(t, err)

		datum, err := DeserializeCom_avro_test_testrecord(&buf)
		assert.Nil(t, err)
		assert.Equal(t, *datum, f)
	}
}
