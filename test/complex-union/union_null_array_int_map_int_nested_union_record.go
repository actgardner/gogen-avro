package avro

type UnionNullArrayIntMapIntNestedUnionRecord struct {
	Null              interface{}
	ArrayInt          []int32
	MapInt            map[string]int32
	NestedUnionRecord *NestedUnionRecord
	UnionType         UnionNullArrayIntMapIntNestedUnionRecordTypeEnum
}

type UnionNullArrayIntMapIntNestedUnionRecordTypeEnum int

const (
	UnionNullArrayIntMapIntNestedUnionRecordTypeEnumNull              UnionNullArrayIntMapIntNestedUnionRecordTypeEnum = 0
	UnionNullArrayIntMapIntNestedUnionRecordTypeEnumArrayInt          UnionNullArrayIntMapIntNestedUnionRecordTypeEnum = 1
	UnionNullArrayIntMapIntNestedUnionRecordTypeEnumMapInt            UnionNullArrayIntMapIntNestedUnionRecordTypeEnum = 2
	UnionNullArrayIntMapIntNestedUnionRecordTypeEnumNestedUnionRecord UnionNullArrayIntMapIntNestedUnionRecordTypeEnum = 3
)
