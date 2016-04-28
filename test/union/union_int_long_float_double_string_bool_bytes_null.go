package avro

type UnionIntLongFloatDoubleStringBoolBytesNull struct {
	Int       int32
	Long      int64
	Float     float32
	Double    float64
	String    string
	Bool      bool
	Bytes     []byte
	Null      interface{}
	UnionType UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum
}

type UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum int

const (
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumInt    UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 0
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumLong   UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 1
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumFloat  UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 2
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumDouble UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 3
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumString UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 4
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumBool   UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 5
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumBytes  UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 6
	UnionIntLongFloatDoubleStringBoolBytesNullTypeEnumNull   UnionIntLongFloatDoubleStringBoolBytesNullTypeEnum = 7
)
