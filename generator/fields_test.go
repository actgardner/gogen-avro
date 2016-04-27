package generator

import (
	"github.com/stretchr/testify/assert"
	"go/format"
	"sort"
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

	pkg := NewPackage("avro")
	primitiveRecord.AddSerializer(pkg)

	assert.Equal(t, pkg.Files(), []string{UTIL_FILE, "primitive_struct.go"})
	utilFile, _ := pkg.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io", "math"})
	expectedFunctions := []FunctionName{{"", "writeBool"}, {"", "writeDouble"}, {"", "writeLong"}, {"", "writeInt"}, {"", "writeString"}, {"", "writeFloat"}, {"", "encodeInt"}, {"", "encodeFloat"}, {"", "writePrimitiveStruct"}}
	sort.Sort(FunctionNameList(expectedFunctions))
	assert.Equal(t, utilFile.Functions(), expectedFunctions)
	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter", "StringWriter"})

	structFile, _ := pkg.File("primitive_struct.go")
	assert.Equal(t, structFile.Imports(), []string{"io"})
	assert.Equal(t, structFile.Functions(), []FunctionName{{"PrimitiveStruct", "Serialize"}})
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

	pkg := NewPackage("avro")
	arrayRecord.AddSerializer(pkg)

	assert.Equal(t, pkg.Files(), []string{"array_struct.go", UTIL_FILE})
	utilFile, _ := pkg.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io", "math"})
	expectedFunctions := []FunctionName{{"", "writeBool"}, {"", "writeDouble"}, {"", "writeLong"}, {"", "writeInt"}, {"", "writeString"}, {"", "writeFloat"}, {"", "writeArrayBool"}, {"", "writeArrayDouble"}, {"", "writeArrayLong"}, {"", "writeArrayInt"}, {"", "writeArrayString"}, {"", "writeArrayFloat"}, {"", "writeArrayNestedRecord"}, {"", "encodeInt"}, {"", "encodeFloat"}, {"", "writeArrayStruct"}}
	sort.Sort(FunctionNameList(expectedFunctions))
	assert.Equal(t, utilFile.Functions(), expectedFunctions)

	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter", "StringWriter"})
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

	pkg := NewPackage("avro")
	mapRecord.AddSerializer(pkg)

	assert.Equal(t, pkg.Files(), []string{"map_struct.go", UTIL_FILE})
	utilFile, _ := pkg.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io", "math"})
	expectedFunctions := []FunctionName{{"", "writeBool"}, {"", "writeDouble"}, {"", "writeLong"}, {"", "writeInt"}, {"", "writeString"}, {"", "writeFloat"}, {"", "writeMapBool"}, {"", "writeMapDouble"}, {"", "writeMapLong"}, {"", "writeMapInt"}, {"", "writeMapString"}, {"", "writeMapFloat"}, {"", "writeMapNestedRecord"}, {"", "encodeInt"}, {"", "encodeFloat"}, {"", "writeMapStruct"}}
	sort.Sort(FunctionNameList(expectedFunctions))
	assert.Equal(t, utilFile.Functions(), expectedFunctions)

	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter", "StringWriter"})

	structFile, _ := pkg.File("map_struct.go")
	assert.Equal(t, structFile.Imports(), []string{"io"})
	assert.Equal(t, structFile.Functions(), []FunctionName{{"MapStruct", "Serialize"}})
}

func TestPrimitiveUnionStructDef(t *testing.T) {
	record := &recordDefinition{
		name:   "UnionStruct",
		fields: []field{&unionField{"UnionField", false, []field{testInt, testString, testFloat, testDouble, testLong, testBool, testRecord, &nullField{}}}},
	}
	expectedStructDef := `type UnionStruct struct {
	UnionField UnionIntStringFloatDoubleLongBoolNestedRecordNull
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

	pkg := NewPackage("avro")
	record.AddSerializer(pkg)

	assert.Equal(t, pkg.Files(), []string{UTIL_FILE, "union_int_string_float_double_long_bool_nested_record_null.go", "union_struct.go"})
	utilFile, _ := pkg.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io", "math"})
	expectedFunctions := []FunctionName{{"", "writeBool"}, {"", "writeDouble"}, {"", "writeLong"}, {"", "writeInt"}, {"", "writeString"}, {"", "writeFloat"}, {"", "writeNull"}, {"", "writeUnionIntStringFloatDoubleLongBoolNestedRecordNull"}, {"", "writeUnionStruct"}, {"", "encodeInt"}, {"", "encodeFloat"}}
	sort.Sort(FunctionNameList(expectedFunctions))
	assert.Equal(t, utilFile.Functions(), expectedFunctions)

	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter", "StringWriter"})

	structFile, _ := pkg.File("union_int_string_float_double_long_bool_nested_record_null.go")
	assert.Equal(t, structFile.Imports(), []string{"fmt"})
	assert.Equal(t, structFile.Functions(), []FunctionName{})
}
