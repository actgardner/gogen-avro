package avro

import (
	"bytes"
	"github.com/linkedin/goavro"
	"github.com/stretchr/testify/assert"
	"testing"
)

/* Round-trip some primitive values through our serializer and goavro to verify */
var fixtures = []Event{
	{
		ID:      "id1",
		StartIP: [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		EndIP:   [16]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	},
	{
		ID:      "differentid",
		StartIP: [16]byte{0, 255, 0, 255, 0, 255, 0, 255, 0, 255, 0, 255},
		EndIP:   [16]byte{0, 1, 3, 7, 15, 31, 63, 127, 254, 2, 4, 6},
	},
}

func compareFixtureGoAvro(t *testing.T, actual interface{}, expected interface{}) {
	record := actual.(*goavro.Record)
	fixture := expected.(Event)
	id, err := record.Get("id")
	assert.Nil(t, err)
	assert.Equal(t, id, fixture.ID)
	startIp, err := record.Get("start_ip")
	assert.Nil(t, err)
	assert.Equal(t, startIp.(goavro.Fixed).Value, fixture.StartIP[:])
	endIp, err := record.Get("end_ip")
	assert.Nil(t, err)
	assert.Equal(t, endIp.(goavro.Fixed).Value, fixture.EndIP[:])
}

func TestRootUnionFixture(t *testing.T) {
	codec, err := goavro.NewCodec(fixtures[0].Schema())
	if err != nil {
		t.Fatal(err)
	}
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err = writeEvent(&f, &buf)
		if err != nil {
			t.Fatal(err)
		}
		datum, err := codec.Decode(&buf)
		if err != nil {
			t.Fatal(err)
		}
		compareFixtureGoAvro(t, datum, f)
	}
}

func TestRoundTrip(t *testing.T) {
	var buf bytes.Buffer
	for _, f := range fixtures {
		buf.Reset()
		err := writeEvent(&f, &buf)
		if err != nil {
			t.Fatal(err)
		}
		datum, err := readEvent(&buf)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, datum, &f)
	}
}
