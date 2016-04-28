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
	testBytes  = &bytesField{"BytesField", []byte{}, true}
	testFixed  = &fixedField{"FixedField", "FixedType", []byte{}, false, 16}
	testEnum   = &enumField{"EnumField", "EnumType", "", false, []string{"a", "b"}}
	testRecord = &recordField{"NestedRecordField", "NestedRecord", nil}
)

/* For each field type, ensure we add the correct functions (including dependencies), structs and imports to each file */
func TestIntSerialize(t *testing.T) {
	p := NewPackage("avro")
	testInt.AddSerializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "encodeInt"}, {"", "writeInt"}})
}

func TestLongSerialize(t *testing.T) {
	p := NewPackage("avro")
	testLong.AddSerializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "encodeInt"}, {"", "writeLong"}})
}

func TestFloatSerialize(t *testing.T) {
	p := NewPackage("avro")
	testFloat.AddSerializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io", "math"})
	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "encodeFloat"}, {"", "writeFloat"}})
}

func TestDoubleSerialize(t *testing.T) {
	p := NewPackage("avro")
	testDouble.AddSerializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io", "math"})
	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "encodeFloat"}, {"", "writeDouble"}})
}

func TestBytesSerialize(t *testing.T) {
	p := NewPackage("avro")
	testBytes.AddSerializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "encodeInt"}, {"", "writeBytes"}, {"", "writeLong"}})
}

func TestStringSerialize(t *testing.T) {
	p := NewPackage("avro")
	testString.AddSerializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter", "StringWriter"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "encodeInt"}, {"", "writeLong"}, {"", "writeString"}})
}

func TestBoolSerialize(t *testing.T) {
	p := NewPackage("avro")
	testBool.AddSerializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "writeBool"}})
}

func TestFixedSerialize(t *testing.T) {
	p := NewPackage("avro")
	testFixed.AddSerializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Structs(), []string{})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "writeFixedType"}})
}

func TestFixedType(t *testing.T) {
	p := NewPackage("avro")
	testFixed.AddStruct(p)

	assert.Equal(t, p.Files(), []string{"fixed_type.go"})

	utilFile, _ := p.File("fixed_type.go")
	assert.Equal(t, utilFile.Imports(), []string{})
	assert.Equal(t, utilFile.Structs(), []string{"FixedType"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{})
}

func TestEnumSerialize(t *testing.T) {
	p := NewPackage("avro")
	testEnum.AddSerializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "encodeInt"}, {"", "writeEnumType"}, {"", "writeInt"}})
}

func TestEnumType(t *testing.T) {
	p := NewPackage("avro")
	testEnum.AddStruct(p)

	assert.Equal(t, p.Files(), []string{"enum_type.go"})

	utilFile, _ := p.File("enum_type.go")
	assert.Equal(t, utilFile.Imports(), []string{})
	assert.Equal(t, utilFile.Structs(), []string{"EnumType"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"EnumType", "String"}})
}

func TestIntDeserialize(t *testing.T) {
	p := NewPackage("avro")
	testInt.AddDeserializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "readInt"}})
}

func TestLongDeserialize(t *testing.T) {
	p := NewPackage("avro")
	testLong.AddDeserializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "readLong"}})
}

func TestFloatDeserialize(t *testing.T) {
	p := NewPackage("avro")
	testFloat.AddDeserializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"encoding/binary", "io", "math"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "readFloat"}})
}

func TestDoubleDeserialize(t *testing.T) {
	p := NewPackage("avro")
	testDouble.AddDeserializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"encoding/binary", "io", "math"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "readDouble"}})
}

func TestBytesDeserialize(t *testing.T) {
	p := NewPackage("avro")
	testBytes.AddDeserializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "readBytes"}, {"", "readLong"}})
}

func TestStringDeserialize(t *testing.T) {
	p := NewPackage("avro")
	testString.AddDeserializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "readLong"}, {"", "readString"}})
}

func TestBoolDeserialize(t *testing.T) {
	p := NewPackage("avro")
	testBool.AddDeserializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "readBool"}})
}

func TestFixedDeserialize(t *testing.T) {
	p := NewPackage("avro")
	testFixed.AddDeserializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "readFixedType"}})
}

func TestEnumDeserialize(t *testing.T) {
	p := NewPackage("avro")
	testEnum.AddDeserializer(p)

	assert.Equal(t, p.Files(), []string{UTIL_FILE})

	utilFile, _ := p.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"io"})
	assert.Equal(t, utilFile.Functions(), []FunctionName{{"", "readEnumType"}, {"", "readInt"}})
}

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

func TestRecordDeserializerMethod(t *testing.T) {
	primitiveRecord := &recordDefinition{
		name:   "PrimitiveStruct",
		fields: []field{testInt, testString, testLong, testFloat, testDouble, testBool, testRecord},
	}

	pkg := NewPackage("avro")
	primitiveRecord.AddDeserializer(pkg)

	assert.Equal(t, pkg.Files(), []string{UTIL_FILE, "primitive_struct.go"})
	utilFile, _ := pkg.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"encoding/binary", "io", "math"})
	expectedFunctions := []FunctionName{{"", "readBool"}, {"", "readDouble"}, {"", "readLong"}, {"", "readInt"}, {"", "readString"}, {"", "readFloat"}, {"", "readPrimitiveStruct"}}
	sort.Sort(FunctionNameList(expectedFunctions))
	assert.Equal(t, utilFile.Functions(), expectedFunctions)
	assert.Equal(t, utilFile.Structs(), []string{})

	structFile, _ := pkg.File("primitive_struct.go")
	assert.Equal(t, structFile.Imports(), []string{"io"})
	assert.Equal(t, structFile.Functions(), []FunctionName{{"", "DeserializePrimitiveStruct"}})
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

func TestArrayStructDeserializer(t *testing.T) {
	arrayRecord := &recordDefinition{
		name:   "ArrayStruct",
		fields: []field{&arrayField{"IntArray", testInt}, &arrayField{"StringArray", testString}, &arrayField{"FloatArray", testFloat}, &arrayField{"DoubleArray", testDouble}, &arrayField{"LongArray", testLong}, &arrayField{"BoolArray", testBool}, &arrayField{"RecordArray", testRecord}},
	}

	pkg := NewPackage("avro")
	arrayRecord.AddDeserializer(pkg)

	assert.Equal(t, pkg.Files(), []string{"array_struct.go", UTIL_FILE})
	utilFile, _ := pkg.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"encoding/binary", "io", "math"})
	expectedFunctions := []FunctionName{{"", "readBool"}, {"", "readDouble"}, {"", "readLong"}, {"", "readInt"}, {"", "readString"}, {"", "readFloat"}, {"", "readArrayBool"}, {"", "readArrayDouble"}, {"", "readArrayLong"}, {"", "readArrayInt"}, {"", "readArrayString"}, {"", "readArrayFloat"}, {"", "readArrayNestedRecord"}, {"", "readArrayStruct"}}
	sort.Sort(FunctionNameList(expectedFunctions))
	assert.Equal(t, utilFile.Functions(), expectedFunctions)

	assert.Equal(t, utilFile.Structs(), []string{})
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

func TestMapDeserializer(t *testing.T) {
	mapRecord := &recordDefinition{
		name:   "MapStruct",
		fields: []field{&mapField{"IntMap", testInt}, &mapField{"StringMap", testString}, &mapField{"FloatMap", testFloat}, &mapField{"DoubleMap", testDouble}, &mapField{"LongMap", testLong}, &mapField{"BoolMap", testBool}, &mapField{"RecordMap", testRecord}},
	}

	pkg := NewPackage("avro")
	mapRecord.AddDeserializer(pkg)

	assert.Equal(t, pkg.Files(), []string{"map_struct.go", UTIL_FILE})
	utilFile, _ := pkg.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"encoding/binary", "fmt", "io", "math"})
	expectedFunctions := []FunctionName{{"", "readBool"}, {"", "readDouble"}, {"", "readLong"}, {"", "readInt"}, {"", "readString"}, {"", "readFloat"}, {"", "readMapBool"}, {"", "readMapDouble"}, {"", "readMapLong"}, {"", "readMapInt"}, {"", "readMapString"}, {"", "readMapFloat"}, {"", "readMapNestedRecord"}, {"", "readMapStruct"}}
	sort.Sort(FunctionNameList(expectedFunctions))
	assert.Equal(t, utilFile.Functions(), expectedFunctions)

	assert.Equal(t, utilFile.Structs(), []string{})

	structFile, _ := pkg.File("map_struct.go")
	assert.Equal(t, structFile.Imports(), []string{"io"})
	assert.Equal(t, structFile.Functions(), []FunctionName{{"", "DeserializeMapStruct"}})
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

	assert.Equal(t, pkg.Files(), []string{UTIL_FILE, "union_struct.go"})
	utilFile, _ := pkg.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"fmt", "io", "math"})
	expectedFunctions := []FunctionName{{"", "writeBool"}, {"", "writeDouble"}, {"", "writeLong"}, {"", "writeInt"}, {"", "writeString"}, {"", "writeFloat"}, {"", "writeNull"}, {"", "writeUnionIntStringFloatDoubleLongBoolNestedRecordNull"}, {"", "writeUnionStruct"}, {"", "encodeInt"}, {"", "encodeFloat"}}
	sort.Sort(FunctionNameList(expectedFunctions))
	assert.Equal(t, utilFile.Functions(), expectedFunctions)

	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter", "StringWriter"})

	structFile, _ := pkg.File("union_struct.go")
	assert.Equal(t, structFile.Imports(), []string{"io"})
	assert.Equal(t, structFile.Functions(), []FunctionName{{"UnionStruct", "Serialize"}})
}

func TestPrimitiveUnionDeserializer(t *testing.T) {
	record := &recordDefinition{
		name:   "UnionStruct",
		fields: []field{&unionField{"UnionField", false, []field{testInt, testString, testFloat, testDouble, testLong, testBool, testRecord, &nullField{}}}},
	}

	pkg := NewPackage("avro")
	record.AddDeserializer(pkg)

	assert.Equal(t, pkg.Files(), []string{UTIL_FILE, "union_struct.go"})
	utilFile, _ := pkg.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"encoding/binary", "fmt", "io", "math"})
	expectedFunctions := []FunctionName{{"", "readBool"}, {"", "readDouble"}, {"", "readLong"}, {"", "readInt"}, {"", "readString"}, {"", "readFloat"}, {"", "readNull"}, {"", "readUnionIntStringFloatDoubleLongBoolNestedRecordNull"}, {"", "readUnionStruct"}}
	sort.Sort(FunctionNameList(expectedFunctions))
	assert.Equal(t, utilFile.Functions(), expectedFunctions)

	assert.Equal(t, utilFile.Structs(), []string{})

	structFile, _ := pkg.File("union_struct.go")
	assert.Equal(t, structFile.Imports(), []string{"io"})
	assert.Equal(t, structFile.Functions(), []FunctionName{{"", "DeserializeUnionStruct"}})
}

func TestRecursiveUnionStructDef(t *testing.T) {
	record := &recordDefinition{
		name:   "RecursiveStruct",
		fields: []field{&unionField{"RecursiveField", false, []field{&nullField{}, &recordField{typeName: "RecursiveStruct"}}}},
	}
	expectedStructDef := `type RecursiveStruct struct {
	RecursiveField UnionNullRecursiveStruct
}
`
	fmtSrc, err := format.Source([]byte(record.structDefinition()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedStructDef)
}

func TestRecursiveUnionSerializer(t *testing.T) {
	record := &recordDefinition{
		name:   "RecursiveStruct",
		fields: []field{&unionField{"RecursiveField", false, []field{&nullField{}, &recordField{typeName: "RecursiveStruct"}}}},
	}

	pkg := NewPackage("avro")
	record.AddSerializer(pkg)

	assert.Equal(t, pkg.Files(), []string{UTIL_FILE, "recursive_struct.go"})
	utilFile, _ := pkg.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"fmt", "io"})
	expectedFunctions := []FunctionName{{"", "writeLong"}, {"", "writeNull"}, {"", "writeUnionNullRecursiveStruct"}, {"", "writeRecursiveStruct"}, {"", "encodeInt"}}
	sort.Sort(FunctionNameList(expectedFunctions))
	assert.Equal(t, utilFile.Functions(), expectedFunctions)

	assert.Equal(t, utilFile.Structs(), []string{"ByteWriter"})

	structFile, _ := pkg.File("recursive_struct.go")
	assert.Equal(t, structFile.Imports(), []string{"io"})
	assert.Equal(t, structFile.Functions(), []FunctionName{{"RecursiveStruct", "Serialize"}})
}

func TestRecursiveUnionDeserializer(t *testing.T) {
	record := &recordDefinition{
		name:   "RecursiveStruct",
		fields: []field{&unionField{"RecursiveField", false, []field{&nullField{}, &recordField{typeName: "RecursiveStruct"}}}},
	}

	pkg := NewPackage("avro")
	record.AddDeserializer(pkg)

	assert.Equal(t, pkg.Files(), []string{UTIL_FILE, "recursive_struct.go"})
	utilFile, _ := pkg.File(UTIL_FILE)
	assert.Equal(t, utilFile.Imports(), []string{"fmt", "io"})
	expectedFunctions := []FunctionName{{"", "readLong"}, {"", "readNull"}, {"", "readUnionNullRecursiveStruct"}, {"", "readRecursiveStruct"}}
	sort.Sort(FunctionNameList(expectedFunctions))
	assert.Equal(t, utilFile.Functions(), expectedFunctions)

	assert.Equal(t, utilFile.Structs(), []string{})

	structFile, _ := pkg.File("recursive_struct.go")
	assert.Equal(t, structFile.Imports(), []string{"io"})
	assert.Equal(t, structFile.Functions(), []FunctionName{{"", "DeserializeRecursiveStruct"}})
}
