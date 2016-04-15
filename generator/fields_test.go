package generator

import (
	"github.com/stretchr/testify/assert"
	"go/format"
	"testing"
)

var (
	testInt    = &intField{"IntField", 1, true}
	testString = &stringField{"StringField", "", false}
	testFloat  = &floatField{"FloatField", 1, true}
	testDouble = &doubleField{"DoubleField", 1, true}
	testLong   = &longField{"LongField", 1, true}
	testBool   = &boolField{"BoolField", true, true}
	testRecord = &recordField{"NestedRecordField", "NestedRecord", nil}
)

func TestRecordStructDef(t *testing.T) {
	primitiveRecord := &recordDefinition{
		name:   "PrimitiveStruct",
		fields: []field{testInt, testString, testLong, testFloat, testDouble, testBool, testRecord},
	}

	expectedStructDef := `type PrimitiveStruct struct {
	IntField          int32
	StringField       string
	LongField         int64
	FloatField        float32
	DoubleField       float64
	BoolField         bool
	NestedRecordField *NestedRecord
}

func (r PrimitiveStruct) Serialize(w io.Writer) error {
	return writePrimitiveStruct(&r, w)
}
`
	fmtSrc, err := format.Source([]byte(primitiveRecord.structDefinition()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedStructDef)
}

func TestRecordSerializerMethod(t *testing.T) {
	primitiveRecord := &recordDefinition{
		name:   "PrimitiveStruct",
		fields: []field{testInt, testString, testLong, testFloat, testDouble, testBool, testRecord},
	}

	expectedSerializer := `func writePrimitiveStruct(r *PrimitiveStruct, w io.Writer) error {
	var err error
	err = writeInt(r.IntField, w)
	if err != nil {
		return err
	}
	err = writeString(r.StringField, w)
	if err != nil {
		return err
	}
	err = writeLong(r.LongField, w)
	if err != nil {
		return err
	}
	err = writeFloat(r.FloatField, w)
	if err != nil {
		return err
	}
	err = writeDouble(r.DoubleField, w)
	if err != nil {
		return err
	}
	err = writeBool(r.BoolField, w)
	if err != nil {
		return err
	}
	err = writeNestedRecord(r.NestedRecordField, w)
	if err != nil {
		return err
	}

	return nil
}`
	imports := make(map[string]string)
	ns := make(map[string]string)

	primitiveRecord.namespaceMap(imports, ns)
	// Verify we import the correct packages
	_, ok := imports["io"]
	assert.Equal(t, ok, true)
	_, ok = imports["math"]
	assert.Equal(t, ok, true)
	assert.Equal(t, len(imports), 2)

	assert.Equal(t, len(ns), 11)
	_, ok = ns["writeBool"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeDouble"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeLong"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeInt"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeString"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeFloat"]
	assert.Equal(t, ok, true)
	_, ok = ns["encodeInt"]
	assert.Equal(t, ok, true)
	_, ok = ns["encodeFloat"]
	assert.Equal(t, ok, true)
	_, ok = ns["StringWriter"]
	assert.Equal(t, ok, true)
	_, ok = ns["ByteWriter"]
	assert.Equal(t, ok, true)

	fmtSrc, err := format.Source([]byte(ns["writePrimitiveStruct"]))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedSerializer)
}

func TestArrayStructDef(t *testing.T) {
	arrayRecord := &recordDefinition{
		name:   "ArrayStruct",
		fields: []field{&arrayField{"IntArray", testInt}, &arrayField{"StringArray", testString}, &arrayField{"FloatArray", testFloat}, &arrayField{"DoubleArray", testDouble}, &arrayField{"LongArray", testLong}, &arrayField{"BoolArray", testBool}, &arrayField{"RecordArray", testRecord}},
	}
	expectedStructDef := `type ArrayStruct struct {
	IntArray    []int32
	StringArray []string
	FloatArray  []float32
	DoubleArray []float64
	LongArray   []int64
	BoolArray   []bool
	RecordArray []*NestedRecord
}

func (r ArrayStruct) Serialize(w io.Writer) error {
	return writeArrayStruct(&r, w)
}
`
	fmtSrc, err := format.Source([]byte(arrayRecord.structDefinition()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedStructDef)
}

func TestArrayStructSerializer(t *testing.T) {
	arrayRecord := &recordDefinition{
		name:   "ArrayStruct",
		fields: []field{&arrayField{"IntArray", testInt}, &arrayField{"StringArray", testString}, &arrayField{"FloatArray", testFloat}, &arrayField{"DoubleArray", testDouble}, &arrayField{"LongArray", testLong}, &arrayField{"BoolArray", testBool}, &arrayField{"RecordArray", testRecord}},
	}
	expectedSerializer := `func writeArrayStruct(r *ArrayStruct, w io.Writer) error {
	var err error
	err = writeArrayInt(r.IntArray, w)
	if err != nil {
		return err
	}
	err = writeArrayString(r.StringArray, w)
	if err != nil {
		return err
	}
	err = writeArrayFloat(r.FloatArray, w)
	if err != nil {
		return err
	}
	err = writeArrayDouble(r.DoubleArray, w)
	if err != nil {
		return err
	}
	err = writeArrayLong(r.LongArray, w)
	if err != nil {
		return err
	}
	err = writeArrayBool(r.BoolArray, w)
	if err != nil {
		return err
	}
	err = writeArrayNestedRecord(r.RecordArray, w)
	if err != nil {
		return err
	}

	return nil
}`
	imports := make(map[string]string)
	ns := make(map[string]string)

	arrayRecord.namespaceMap(imports, ns)
	// Verify we import the correct packages
	_, ok := imports["io"]
	assert.Equal(t, ok, true)
	_, ok = imports["math"]
	assert.Equal(t, ok, true)
	assert.Equal(t, len(imports), 2)

	assert.Equal(t, len(ns), 18)
	_, ok = ns["writeBool"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeDouble"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeLong"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeInt"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeString"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeFloat"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeArrayBool"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeArrayDouble"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeArrayLong"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeArrayInt"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeArrayString"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeArrayFloat"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeArrayNestedRecord"]
	assert.Equal(t, ok, true)

	_, ok = ns["encodeInt"]
	assert.Equal(t, ok, true)
	_, ok = ns["encodeFloat"]
	assert.Equal(t, ok, true)
	_, ok = ns["StringWriter"]
	assert.Equal(t, ok, true)
	_, ok = ns["ByteWriter"]
	assert.Equal(t, ok, true)

	fmtSrc, err := format.Source([]byte(ns["writeArrayStruct"]))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedSerializer)
}

func TestMapStructDef(t *testing.T) {
	mapRecord := &recordDefinition{
		name:   "MapStruct",
		fields: []field{&mapField{"IntMap", testInt}, &mapField{"StringMap", testString}, &mapField{"FloatMap", testFloat}, &mapField{"DoubleMap", testDouble}, &mapField{"LongMap", testLong}, &mapField{"BoolMap", testBool}, &mapField{"RecordMap", testRecord}},
	}
	expectedStructDef := `type MapStruct struct {
	IntMap    map[string]int32
	StringMap map[string]string
	FloatMap  map[string]float32
	DoubleMap map[string]float64
	LongMap   map[string]int64
	BoolMap   map[string]bool
	RecordMap map[string]*NestedRecord
}

func (r MapStruct) Serialize(w io.Writer) error {
	return writeMapStruct(&r, w)
}
`
	fmtSrc, err := format.Source([]byte(mapRecord.structDefinition()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedStructDef)
}

func TestMapSerializer(t *testing.T) {
	mapRecord := &recordDefinition{
		name:   "MapStruct",
		fields: []field{&mapField{"IntMap", testInt}, &mapField{"StringMap", testString}, &mapField{"FloatMap", testFloat}, &mapField{"DoubleMap", testDouble}, &mapField{"LongMap", testLong}, &mapField{"BoolMap", testBool}, &mapField{"RecordMap", testRecord}},
	}
	expectedSerializer := `func writeMapStruct(r *MapStruct, w io.Writer) error {
	var err error
	err = writeMapInt(r.IntMap, w)
	if err != nil {
		return err
	}
	err = writeMapString(r.StringMap, w)
	if err != nil {
		return err
	}
	err = writeMapFloat(r.FloatMap, w)
	if err != nil {
		return err
	}
	err = writeMapDouble(r.DoubleMap, w)
	if err != nil {
		return err
	}
	err = writeMapLong(r.LongMap, w)
	if err != nil {
		return err
	}
	err = writeMapBool(r.BoolMap, w)
	if err != nil {
		return err
	}
	err = writeMapNestedRecord(r.RecordMap, w)
	if err != nil {
		return err
	}

	return nil
}`

	imports := make(map[string]string)
	ns := make(map[string]string)

	mapRecord.namespaceMap(imports, ns)
	// Verify we import the correct packages
	_, ok := imports["io"]
	assert.Equal(t, ok, true)
	_, ok = imports["math"]
	assert.Equal(t, ok, true)
	assert.Equal(t, len(imports), 2)

	assert.Equal(t, len(ns), 18)
	_, ok = ns["writeBool"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeDouble"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeLong"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeInt"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeString"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeFloat"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeMapBool"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeMapDouble"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeMapLong"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeMapInt"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeMapString"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeMapFloat"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeMapNestedRecord"]
	assert.Equal(t, ok, true)

	_, ok = ns["encodeInt"]
	assert.Equal(t, ok, true)
	_, ok = ns["encodeFloat"]
	assert.Equal(t, ok, true)
	_, ok = ns["StringWriter"]
	assert.Equal(t, ok, true)
	_, ok = ns["ByteWriter"]
	assert.Equal(t, ok, true)
	fmtSrc, err := format.Source([]byte(ns["writeMapStruct"]))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedSerializer)
}

func TestPrimitiveUnionStructDef(t *testing.T) {
	record := &recordDefinition{
		name:   "UnionStruct",
		fields: []field{&unionField{"UnionField", false, []field{testInt, testString, testFloat, testDouble, testLong, testBool, testRecord, &nullField{}}}},
	}
	expectedStructDef := `type UnionStruct struct {
	UnionField UnionIntStringFloatDoubleLongBoolNestedRecordNull
}

func (r UnionStruct) Serialize(w io.Writer) error {
	return writeUnionStruct(&r, w)
}
`
	fmtSrc, err := format.Source([]byte(record.structDefinition()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedStructDef)

}

func TestPrimitiveUnionSerializer(t *testing.T) {
	record := &recordDefinition{
		name:   "UnionStruct",
		fields: []field{&unionField{"UnionField", false, []field{testInt, testString, testFloat, testDouble, testLong, testBool, testRecord, &nullField{}}}},
	}
	expectedSerializer := `func writeUnionStruct(r *UnionStruct, w io.Writer) error {
	var err error
	err = writeUnionIntStringFloatDoubleLongBoolNestedRecordNull(r.UnionField, w)
	if err != nil {
		return err
	}

	return nil
}`
	imports := make(map[string]string)
	ns := make(map[string]string)

	record.namespaceMap(imports, ns)
	// Verify we import the correct packages
	assert.Equal(t, len(imports), 3)
	_, ok := imports["io"]
	assert.Equal(t, ok, true)
	_, ok = imports["math"]
	assert.Equal(t, ok, true)
	_, ok = imports["fmt"]
	assert.Equal(t, ok, true)

	assert.Equal(t, len(ns), 15)
	_, ok = ns["writeBool"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeDouble"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeLong"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeInt"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeString"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeFloat"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeNull"]
	assert.Equal(t, ok, true)
	_, ok = ns["writeUnionIntStringFloatDoubleLongBoolNestedRecordNull"]
	assert.Equal(t, ok, true)

	_, ok = ns["encodeInt"]
	assert.Equal(t, ok, true)
	_, ok = ns["encodeFloat"]
	assert.Equal(t, ok, true)
	_, ok = ns["StringWriter"]
	assert.Equal(t, ok, true)
	_, ok = ns["ByteWriter"]
	assert.Equal(t, ok, true)
	fmtSrc, err := format.Source([]byte(ns["writeUnionStruct"]))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedSerializer)
}
