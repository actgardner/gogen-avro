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
	testRecord = &recordField{"NestedRecordField", "NestedRecord"}
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

	expectedSerializer := `func (r *PrimitiveStruct) Serialize(w io.Writer) error {
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

	fmtSrc, err := format.Source([]byte(primitiveRecord.serializerMethod()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedSerializer)
}

func TestRecordAuxStructs(t *testing.T) {
	primitiveRecord := &recordDefinition{
		name:   "PrimitiveStruct",
		fields: []field{testInt, testString, testLong, testFloat, testDouble, testBool, testRecord},
	}
	expectedSerializer := `import "io"
import "math"

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type StringWriter interface {
	WriteString(string) (int, error)
}

func encodeFloat(w io.Writer, byteCount int, bits uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}
	for i := 0; i < byteCount; i++ {
		if bw != nil {
			err = bw.WriteByte(byte(bits & 255))
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(bits&255))
		}
		bits = bits >> 8
	}
	if bw == nil {
		_, err = w.Write(bb)
		return err
	}
	return nil
}

func encodeInt(w io.Writer, byteCount int, encoded int64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}

	if encoded == 0 {
		if bw != nil {
			err = bw.WriteByte(0)
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(0))
		}
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			if bw != nil {
				err = bw.WriteByte(b)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, b)
			}
		}
	}
	if bw == nil {
		_, err := w.Write(bb)
		return err
	}
	return nil

}

func writeBool(r bool, w io.Writer) error {
	var b byte
	if r {
		b = byte(1)
	}

	var err error
	if bw, ok := w.(ByteWriter); ok {
		err = bw.WriteByte(b)
	} else {
		bb := make([]byte, 1)
		bb[0] = b
		_, err = w.Write(bb)
	}
	if err != nil {
		return err
	}
	return nil
}

func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}

func writeFloat(r float32, w io.Writer) error {
	bits := uint64(math.Float32bits(r))
	const byteCount = 4
	return encodeFloat(w, byteCount, bits)
}

func writeInt(r int32, w io.Writer) error {
	downShift := uint32(31)
	encoded := int64((r << 1) ^ (r >> downShift))
	const maxByteSize = 5
	return encodeInt(w, maxByteSize, encoded)
}

func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := int64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}

func writeString(r string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	if sw, ok := w.(StringWriter); ok {
		_, err = sw.WriteString(r)
	} else {
		_, err = w.Write([]byte(r))
	}
	return err
}

`
	fmtSrc, err := format.Source([]byte(primitiveRecord.auxStructs()))
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
	expectedSerializer := `func (r *ArrayStruct) Serialize(w io.Writer) error {
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
	err = writeArrayNestedRecordRecord(r.RecordArray, w)
	if err != nil {
		return err
	}

	return nil
}`
	fmtSrc, err := format.Source([]byte(arrayRecord.serializerMethod()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedSerializer)
}

func TestArrayStructAuxStructs(t *testing.T) {
	arrayRecord := &recordDefinition{
		name:   "ArrayStruct",
		fields: []field{&arrayField{"IntArray", testInt}, &arrayField{"StringArray", testString}, &arrayField{"FloatArray", testFloat}, &arrayField{"DoubleArray", testDouble}, &arrayField{"LongArray", testLong}, &arrayField{"BoolArray", testBool}, &arrayField{"RecordArray", testRecord}},
	}
	expectedStructDef := `import "io"
import "math"

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type StringWriter interface {
	WriteString(string) (int, error)
}

func encodeFloat(w io.Writer, byteCount int, bits uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}
	for i := 0; i < byteCount; i++ {
		if bw != nil {
			err = bw.WriteByte(byte(bits & 255))
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(bits&255))
		}
		bits = bits >> 8
	}
	if bw == nil {
		_, err = w.Write(bb)
		return err
	}
	return nil
}

func encodeInt(w io.Writer, byteCount int, encoded int64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}

	if encoded == 0 {
		if bw != nil {
			err = bw.WriteByte(0)
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(0))
		}
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			if bw != nil {
				err = bw.WriteByte(b)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, b)
			}
		}
	}
	if bw == nil {
		_, err := w.Write(bb)
		return err
	}
	return nil

}

func writeArrayBool(r []bool, w io.Writer) error {
	err := writeLong(len(r), w)
	if err != nil {
		return err
	}
	for _, e := range r {
		err = writeBool(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeArrayDouble(r []float64, w io.Writer) error {
	err := writeLong(len(r), w)
	if err != nil {
		return err
	}
	for _, e := range r {
		err = writeDouble(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeArrayFloat(r []float32, w io.Writer) error {
	err := writeLong(len(r), w)
	if err != nil {
		return err
	}
	for _, e := range r {
		err = writeFloat(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeArrayInt(r []int32, w io.Writer) error {
	err := writeLong(len(r), w)
	if err != nil {
		return err
	}
	for _, e := range r {
		err = writeInt(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeArrayLong(r []int64, w io.Writer) error {
	err := writeLong(len(r), w)
	if err != nil {
		return err
	}
	for _, e := range r {
		err = writeLong(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeArrayNestedRecordRecord(r []*NestedRecord, w io.Writer) error {
	err := writeLong(len(r), w)
	if err != nil {
		return err
	}
	for _, e := range r {
		err = writeNestedRecord(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeArrayString(r []string, w io.Writer) error {
	err := writeLong(len(r), w)
	if err != nil {
		return err
	}
	for _, e := range r {
		err = writeString(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeBool(r bool, w io.Writer) error {
	var b byte
	if r {
		b = byte(1)
	}

	var err error
	if bw, ok := w.(ByteWriter); ok {
		err = bw.WriteByte(b)
	} else {
		bb := make([]byte, 1)
		bb[0] = b
		_, err = w.Write(bb)
	}
	if err != nil {
		return err
	}
	return nil
}

func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}

func writeFloat(r float32, w io.Writer) error {
	bits := uint64(math.Float32bits(r))
	const byteCount = 4
	return encodeFloat(w, byteCount, bits)
}

func writeInt(r int32, w io.Writer) error {
	downShift := uint32(31)
	encoded := int64((r << 1) ^ (r >> downShift))
	const maxByteSize = 5
	return encodeInt(w, maxByteSize, encoded)
}

func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := int64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}

func writeString(r string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	if sw, ok := w.(StringWriter); ok {
		_, err = sw.WriteString(r)
	} else {
		_, err = w.Write([]byte(r))
	}
	return err
}

`
	fmtSrc, err := format.Source([]byte(arrayRecord.auxStructs()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedStructDef)
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
	expectedSerializer := `func (r *MapStruct) Serialize(w io.Writer) error {
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
	err = writeMapNestedRecordRecord(r.RecordMap, w)
	if err != nil {
		return err
	}

	return nil
}`
	fmtSrc, err := format.Source([]byte(mapRecord.serializerMethod()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedSerializer)
}

func TestMapAuxStructs(t *testing.T) {
	mapRecord := &recordDefinition{
		name:   "MapStruct",
		fields: []field{&mapField{"IntMap", testInt}, &mapField{"StringMap", testString}, &mapField{"FloatMap", testFloat}, &mapField{"DoubleMap", testDouble}, &mapField{"LongMap", testLong}, &mapField{"BoolMap", testBool}, &mapField{"RecordMap", testRecord}},
	}
	expectedStructDef := `import "io"
import "math"

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type StringWriter interface {
	WriteString(string) (int, error)
}

func encodeFloat(w io.Writer, byteCount int, bits uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}
	for i := 0; i < byteCount; i++ {
		if bw != nil {
			err = bw.WriteByte(byte(bits & 255))
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(bits&255))
		}
		bits = bits >> 8
	}
	if bw == nil {
		_, err = w.Write(bb)
		return err
	}
	return nil
}

func encodeInt(w io.Writer, byteCount int, encoded int64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}

	if encoded == 0 {
		if bw != nil {
			err = bw.WriteByte(0)
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(0))
		}
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			if bw != nil {
				err = bw.WriteByte(b)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, b)
			}
		}
	}
	if bw == nil {
		_, err := w.Write(bb)
		return err
	}
	return nil

}

func writeBool(r bool, w io.Writer) error {
	var b byte
	if r {
		b = byte(1)
	}

	var err error
	if bw, ok := w.(ByteWriter); ok {
		err = bw.WriteByte(b)
	} else {
		bb := make([]byte, 1)
		bb[0] = b
		_, err = w.Write(bb)
	}
	if err != nil {
		return err
	}
	return nil
}

func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}

func writeFloat(r float32, w io.Writer) error {
	bits := uint64(math.Float32bits(r))
	const byteCount = 4
	return encodeFloat(w, byteCount, bits)
}

func writeInt(r int32, w io.Writer) error {
	downShift := uint32(31)
	encoded := int64((r << 1) ^ (r >> downShift))
	const maxByteSize = 5
	return encodeInt(w, maxByteSize, encoded)
}

func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := int64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}

func writeMapBool(r map[string]bool, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeBool(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeMapDouble(r map[string]float64, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeDouble(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeMapFloat(r map[string]float32, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeFloat(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeMapInt(r map[string]int32, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeInt(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeMapLong(r map[string]int64, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeLong(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeMapNestedRecordRecord(r map[string]*NestedRecord, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeNestedRecord(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeMapString(r map[string]string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	for k, e := range r {
		err = writeString(k, w)
		if err != nil {
			return err
		}
		err = writeString(e, w)
		if err != nil {
			return err
		}
	}
	return writeLong(0, w)
}

func writeString(r string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	if sw, ok := w.(StringWriter); ok {
		_, err = sw.WriteString(r)
	} else {
		_, err = w.Write([]byte(r))
	}
	return err
}

`
	fmtSrc, err := format.Source([]byte(mapRecord.auxStructs()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedStructDef)
}

func TestPrimitiveUnionStructDef(t *testing.T) {
	record := &recordDefinition{
		name:   "UnionStruct",
		fields: []field{&unionField{"UnionField", false, []field{testInt, testString, testFloat, testDouble, testLong, testBool, testRecord, &nullField{}}}},
	}
	expectedStructDef := `type UnionStruct struct {
	UnionField UnionIntStringFloatDoubleLongBoolNestedRecordRecordNull
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
	expectedSerializer := `func (r *UnionStruct) Serialize(w io.Writer) error {
	var err error
	err = writeUnionIntStringFloatDoubleLongBoolNestedRecordRecordNull(r.UnionField, w)
	if err != nil {
		return err
	}

	return nil
}`
	fmtSrc, err := format.Source([]byte(record.serializerMethod()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedSerializer)

}

func TestPrimitiveUnionAuxStructs(t *testing.T) {
	record := &recordDefinition{
		name:   "UnionStruct",
		fields: []field{&unionField{"UnionField", false, []field{testInt, testString, testFloat, testDouble, testLong, testBool, testRecord, &nullField{}}}},
	}
	expectedStructDef := `import "fmt"
import "io"
import "math"

type ByteWriter interface {
	Grow(int)
	WriteByte(byte) error
}

type StringWriter interface {
	WriteString(string) (int, error)
}

type UnionIntStringFloatDoubleLongBoolNestedRecordRecordNull struct {
	Int                int32
	String             string
	Float              float32
	Double             float64
	Long               int64
	Bool               bool
	NestedRecordRecord *NestedRecord
	Null               interface{}
	UnionType          UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnum
}

type UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnum int

const (
	UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumInt                UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnum = 0
	UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumString             UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnum = 1
	UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumFloat              UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnum = 2
	UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumDouble             UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnum = 3
	UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumLong               UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnum = 4
	UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumBool               UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnum = 5
	UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumNestedRecordRecord UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnum = 6
	UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumNull               UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnum = 7
)

func encodeFloat(w io.Writer, byteCount int, bits uint64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}
	for i := 0; i < byteCount; i++ {
		if bw != nil {
			err = bw.WriteByte(byte(bits & 255))
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(bits&255))
		}
		bits = bits >> 8
	}
	if bw == nil {
		_, err = w.Write(bb)
		return err
	}
	return nil
}

func encodeInt(w io.Writer, byteCount int, encoded int64) error {
	var err error
	var bb []byte
	bw, ok := w.(ByteWriter)
	// To avoid reallocations, grow capacity to the largest possible size
	// for this integer
	if ok {
		bw.Grow(byteCount)
	} else {
		bb = make([]byte, 0, byteCount)
	}

	if encoded == 0 {
		if bw != nil {
			err = bw.WriteByte(0)
			if err != nil {
				return err
			}
		} else {
			bb = append(bb, byte(0))
		}
	} else {
		for encoded > 0 {
			b := byte(encoded & 127)
			encoded = encoded >> 7
			if !(encoded == 0) {
				b |= 128
			}
			if bw != nil {
				err = bw.WriteByte(b)
				if err != nil {
					return err
				}
			} else {
				bb = append(bb, b)
			}
		}
	}
	if bw == nil {
		_, err := w.Write(bb)
		return err
	}
	return nil

}

func writeBool(r bool, w io.Writer) error {
	var b byte
	if r {
		b = byte(1)
	}

	var err error
	if bw, ok := w.(ByteWriter); ok {
		err = bw.WriteByte(b)
	} else {
		bb := make([]byte, 1)
		bb[0] = b
		_, err = w.Write(bb)
	}
	if err != nil {
		return err
	}
	return nil
}

func writeDouble(r float64, w io.Writer) error {
	bits := uint64(math.Float64bits(r))
	const byteCount = 8
	return encodeFloat(w, byteCount, bits)
}

func writeFloat(r float32, w io.Writer) error {
	bits := uint64(math.Float32bits(r))
	const byteCount = 4
	return encodeFloat(w, byteCount, bits)
}

func writeInt(r int32, w io.Writer) error {
	downShift := uint32(31)
	encoded := int64((r << 1) ^ (r >> downShift))
	const maxByteSize = 5
	return encodeInt(w, maxByteSize, encoded)
}

func writeLong(r int64, w io.Writer) error {
	downShift := uint64(63)
	encoded := int64((r << 1) ^ (r >> downShift))
	const maxByteSize = 10
	return encodeInt(w, maxByteSize, encoded)
}

func writeNull(_ interface{}, _ io.Writer) error {
	return nil
}

func writeString(r string, w io.Writer) error {
	err := writeLong(int64(len(r)), w)
	if err != nil {
		return err
	}
	if sw, ok := w.(StringWriter); ok {
		_, err = sw.WriteString(r)
	} else {
		_, err = w.Write([]byte(r))
	}
	return err
}

func writeUnionIntStringFloatDoubleLongBoolNestedRecordRecordNull(r UnionIntStringFloatDoubleLongBoolNestedRecordRecordNull, w io.Writer) error {
	err := writeLong(int64(r.UnionType), w)
	if err != nil {
		return err
	}
	switch r.UnionType {
	case UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumInt:
		return writeInt(r.Int, w)
	case UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumString:
		return writeString(r.String, w)
	case UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumFloat:
		return writeFloat(r.Float, w)
	case UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumDouble:
		return writeDouble(r.Double, w)
	case UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumLong:
		return writeLong(r.Long, w)
	case UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumBool:
		return writeBool(r.Bool, w)
	case UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumNestedRecordRecord:
		return writeNestedRecord(r.NestedRecordRecord, w)
	case UnionIntStringFloatDoubleLongBoolNestedRecordRecordNullTypeEnumNull:
		return writeNull(r.Null, w)

	}
	return fmt.Errorf("Invalid value for UnionIntStringFloatDoubleLongBoolNestedRecordRecordNull")
}

`
	fmtSrc, err := format.Source([]byte(record.auxStructs()))
	assert.Nil(t, err)
	assert.Equal(t, string(fmtSrc), expectedStructDef)

}
