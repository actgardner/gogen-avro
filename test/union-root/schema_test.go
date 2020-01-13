package avro

import (
	"bytes"
	"testing"

	"github.com/linkedin/goavro/v2"
	"github.com/stretchr/testify/assert"
)

// Round-trip some primitive values through our serializer and goavro to verify
var fixtures = []Event{
	{
		Id:       "id1",
		Start_ip: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		End_ip:   [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		Id:       "differentid",
		Start_ip: [16]byte{0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255},
		End_ip:   [16]byte{0, 1, 3, 7, 15, 31, 63, 127, 254, 2, 4, 6},
	},
}

func compareFixtureGoAvro(t *testing.T, actual interface{}, expected interface{}) {
	record := actual.(map[string]interface{})
	fixture := expected.(Event)
	id, ok := record["id"]
	assert.Equal(t, ok, true)
	assert.Equal(t, id, fixture.Id)

	startIp, ok := record["start_ip"]
	assert.Equal(t, ok, true)
	assert.Equal(t, startIp, fixture.Start_ip[:])

	endIp, ok := record["end_ip"]
	assert.Equal(t, ok, true)
	assert.Equal(t, endIp, fixture.End_ip[:])
}

func TestRootUnionFixture(t *testing.T) {
	codec, err := goavro.NewCodec(fixtures[0].Schema())
	assert.Nil(t, err)

	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = writeEvent(&f, &buf)
		assert.Nil(t, err)

		datum, remaining, err := codec.NativeFromBinary(buf.Bytes())
		assert.Nil(t, err)

		assert.Equal(t, 0, len(remaining))
		compareFixtureGoAvro(t, datum, f)
	}
}

func TestRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err := writeEvent(&f, &buf)
		assert.Nil(t, err)

		datum, err := DeserializeEvent(&buf)
		assert.Nil(t, err)

		assert.Equal(t, datum, &f)
	}
}
