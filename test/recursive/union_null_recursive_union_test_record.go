package avro

type UnionNullRecursiveUnionTestRecord struct {
	Null                     interface{}
	RecursiveUnionTestRecord *RecursiveUnionTestRecord
	UnionType                UnionNullRecursiveUnionTestRecordTypeEnum
}

type UnionNullRecursiveUnionTestRecordTypeEnum int

const (
	UnionNullRecursiveUnionTestRecordTypeEnumNull                     UnionNullRecursiveUnionTestRecordTypeEnum = 0
	UnionNullRecursiveUnionTestRecordTypeEnumRecursiveUnionTestRecord UnionNullRecursiveUnionTestRecordTypeEnum = 1
)
