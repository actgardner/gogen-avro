package schema

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePrimitiveRecordSchema(t *testing.T) {
	schemaString := `
	{
		"type": "record",
		"name": "PrimitiveTest",
		"fields": [
			{"name":"StringField", "type":"string"},
			{"name":"IntField", "type":"int"},
			{"name":"LongField", "type":"long"},
			{"name":"BoolField", "type":"boolean"},
			{"name":"FloatField", "type": "float"},
			{"name":"DoubleField", "type": "double"},
			{"name":"BytesField", "type": "bytes"},
			{"name":"RecordField", "type": "NestedRecord"}
		]
	}
`
	schemaField, err := FieldDefinitionForSchema([]byte(schemaString))
	assert.Nil(t, err)
	schemaRecord, ok := schemaField.(*recordField)
	assert.True(t, ok)
	schema := schemaRecord.def
	assert.Equal(t, schema.name, "PrimitiveTest")
	assert.Equal(t, len(schema.fields), 8)
	assert.Equal(t, schema.fields[0].(*stringField).name, "StringField")
	assert.Equal(t, schema.fields[0].(*stringField).hasDefault, false)
	assert.Equal(t, schema.fields[1].(*intField).name, "IntField")
	assert.Equal(t, schema.fields[1].(*intField).hasDefault, false)
	assert.Equal(t, schema.fields[2].(*longField).name, "LongField")
	assert.Equal(t, schema.fields[2].(*longField).hasDefault, false)
	assert.Equal(t, schema.fields[3].(*boolField).name, "BoolField")
	assert.Equal(t, schema.fields[3].(*boolField).hasDefault, false)
	assert.Equal(t, schema.fields[4].(*floatField).name, "FloatField")
	assert.Equal(t, schema.fields[4].(*floatField).hasDefault, false)
	assert.Equal(t, schema.fields[5].(*doubleField).name, "DoubleField")
	assert.Equal(t, schema.fields[5].(*doubleField).hasDefault, false)
	assert.Equal(t, schema.fields[6].(*bytesField).name, "BytesField")
	assert.Equal(t, schema.fields[6].(*bytesField).hasDefault, false)
	assert.Equal(t, schema.fields[7].(*recordField).name, "RecordField")
	assert.Equal(t, schema.fields[7].(*recordField).typeName, "NestedRecord")
}

func TestParsePrimitiveMapRecordSchema(t *testing.T) {
	schemaString := `
	{
		"type": "record",
		"name": "ComplexRecord",
		"fields": [
			{"name": "StringMapField", "type": {"type": "map", "values": "string"}},
			{"name": "IntMapField", "type": {"type": "map", "values": "int"}},
			{"name": "LongMapField", "type": {"type": "map", "values": "long"}},
			{"name": "BoolMapField", "type": {"type": "map", "values": "boolean"}},
			{"name": "FloatMapField", "type": {"type": "map", "values": "float"}},
			{"name": "DoubleMapField", "type": {"type": "map", "values": "double"}},
			{"name": "BytesMapField", "type": {"type": "map", "values": "bytes"}},
			{"name": "RecordMapField", "type": {"type": "map", "values": "NestedRecord"}}
		]
	}
	`
	schemaField, err := FieldDefinitionForSchema([]byte(schemaString))
	assert.Nil(t, err)
	schemaRecord, ok := schemaField.(*recordField)
	assert.True(t, ok)
	schema := schemaRecord.def

	assert.Equal(t, schema.name, "ComplexRecord")
	assert.Equal(t, len(schema.fields), 8)
	assert.Equal(t, schema.fields[0].(*mapField).name, "StringMapField")
	_, ok = schema.fields[0].(*mapField).itemType.(*stringField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[1].(*mapField).name, "IntMapField")
	_, ok = schema.fields[1].(*mapField).itemType.(*intField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[2].(*mapField).name, "LongMapField")
	_, ok = schema.fields[2].(*mapField).itemType.(*longField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[3].(*mapField).name, "BoolMapField")
	_, ok = schema.fields[3].(*mapField).itemType.(*boolField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[4].(*mapField).name, "FloatMapField")
	_, ok = schema.fields[4].(*mapField).itemType.(*floatField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[5].(*mapField).name, "DoubleMapField")
	_, ok = schema.fields[5].(*mapField).itemType.(*doubleField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[6].(*mapField).name, "BytesMapField")
	_, ok = schema.fields[6].(*mapField).itemType.(*bytesField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[7].(*mapField).name, "RecordMapField")
	_, ok = schema.fields[7].(*mapField).itemType.(*recordField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[7].(*mapField).itemType.(*recordField).typeName, "NestedRecord")
}

func TestParsePrimitiveArrayRecordSchema(t *testing.T) {
	schemaString := `
	{
		"type": "record",
		"name": "ComplexRecord",
		"fields": [
			{"name": "StringArrayField", "type": {"type": "array", "items": "string"}},
			{"name": "IntArrayField", "type": {"type": "array", "items": "int"}},
			{"name": "LongArrayField", "type": {"type": "array", "items": "long"}},
			{"name": "BoolArrayField", "type": {"type": "array", "items": "boolean"}},
			{"name": "FloatArrayField", "type": {"type": "array", "items": "float"}},
			{"name": "DoubleArrayField", "type": {"type": "array", "items": "double"}},
			{"name": "BytesArrayField", "type": {"type": "array", "items": "bytes"}},
			{"name": "RecordArrayField", "type": {"type": "array", "items": "NestedRecord"}}
		]
	}
	`
	schemaField, err := FieldDefinitionForSchema([]byte(schemaString))
	assert.Nil(t, err)
	schemaRecord, ok := schemaField.(*recordField)
	assert.True(t, ok)
	schema := schemaRecord.def

	assert.Equal(t, schema.name, "ComplexRecord")
	assert.Equal(t, len(schema.fields), 8)
	assert.Equal(t, schema.fields[0].(*arrayField).name, "StringArrayField")
	_, ok = schema.fields[0].(*arrayField).itemType.(*stringField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[1].(*arrayField).name, "IntArrayField")
	_, ok = schema.fields[1].(*arrayField).itemType.(*intField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[2].(*arrayField).name, "LongArrayField")
	_, ok = schema.fields[2].(*arrayField).itemType.(*longField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[3].(*arrayField).name, "BoolArrayField")
	_, ok = schema.fields[3].(*arrayField).itemType.(*boolField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[4].(*arrayField).name, "FloatArrayField")
	_, ok = schema.fields[4].(*arrayField).itemType.(*floatField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[5].(*arrayField).name, "DoubleArrayField")
	_, ok = schema.fields[5].(*arrayField).itemType.(*doubleField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[6].(*arrayField).name, "BytesArrayField")
	_, ok = schema.fields[6].(*arrayField).itemType.(*bytesField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[7].(*arrayField).name, "RecordArrayField")
	_, ok = schema.fields[7].(*arrayField).itemType.(*recordField)
	assert.Equal(t, ok, true)
	assert.Equal(t, schema.fields[7].(*arrayField).itemType.(*recordField).typeName, "NestedRecord")
}

func TestParsePrimitiveUnionRecordSchema(t *testing.T) {
	schemaString := `
	{
		"type": "record",
		"name": "UnionRecord",
		"fields": [
			{"name": "UnionField", "type": ["string", "int", "long", "boolean", "float", "double", "bytes", "null", "NestedRecord"]}
		]
	}
`
	schemaField, err := FieldDefinitionForSchema([]byte(schemaString))
	assert.Nil(t, err)
	schemaRecord, ok := schemaField.(*recordField)
	assert.True(t, ok)
	schema := schemaRecord.def

	assert.Equal(t, schema.name, "UnionRecord")
	assert.Equal(t, len(schema.fields), 1)
	unionField, ok := schema.fields[0].(*unionField)
	assert.Equal(t, unionField.name, "UnionField")
	_, ok = unionField.itemType[0].(*stringField)
	assert.Equal(t, ok, true)
	_, ok = unionField.itemType[1].(*intField)
	assert.Equal(t, ok, true)
	_, ok = unionField.itemType[2].(*longField)
	assert.Equal(t, ok, true)
	_, ok = unionField.itemType[3].(*boolField)
	assert.Equal(t, ok, true)
	_, ok = unionField.itemType[4].(*floatField)
	assert.Equal(t, ok, true)
	_, ok = unionField.itemType[5].(*doubleField)
	assert.Equal(t, ok, true)
	_, ok = unionField.itemType[6].(*bytesField)
	assert.Equal(t, ok, true)
	_, ok = unionField.itemType[7].(*nullField)
	assert.Equal(t, ok, true)
	_, ok = unionField.itemType[8].(*recordField)
	assert.Equal(t, ok, true)
	assert.Equal(t, unionField.itemType[8].(*recordField).typeName, "NestedRecord")
}

func TestParsePrimitiveSchema(t *testing.T) {
	schemaString := `"string"`
	schemaField, err := FieldDefinitionForSchema([]byte(schemaString))
	assert.Nil(t, err)
	field, ok := schemaField.(*stringField)
	assert.True(t, ok)
	assert.Equal(t, field.name, "")
	assert.Equal(t, field.hasDefault, false)
}

func TestParsePrimitiveUnionSchema(t *testing.T) {
	schemaString := `
	[
		"string",
		"int",
		"long",
		"boolean",
		"float",
		"double",
		"bytes",
		"NestedRecord"
	]
`
	schemaField, err := FieldDefinitionForSchema([]byte(schemaString))
	assert.Nil(t, err)
	schemaUnion, ok := schemaField.(*unionField)
	assert.True(t, ok)
	assert.Equal(t, schemaUnion.name, "")
	assert.Equal(t, len(schemaUnion.itemType), 8)
	assert.Equal(t, schemaUnion.itemType[0].(*stringField).name, "")
	assert.Equal(t, schemaUnion.itemType[0].(*stringField).hasDefault, false)
	assert.Equal(t, schemaUnion.itemType[1].(*intField).name, "")
	assert.Equal(t, schemaUnion.itemType[1].(*intField).hasDefault, false)
	assert.Equal(t, schemaUnion.itemType[2].(*longField).name, "")
	assert.Equal(t, schemaUnion.itemType[2].(*longField).hasDefault, false)
	assert.Equal(t, schemaUnion.itemType[3].(*boolField).name, "")
	assert.Equal(t, schemaUnion.itemType[3].(*boolField).hasDefault, false)
	assert.Equal(t, schemaUnion.itemType[4].(*floatField).name, "")
	assert.Equal(t, schemaUnion.itemType[4].(*floatField).hasDefault, false)
	assert.Equal(t, schemaUnion.itemType[5].(*doubleField).name, "")
	assert.Equal(t, schemaUnion.itemType[5].(*doubleField).hasDefault, false)
	assert.Equal(t, schemaUnion.itemType[6].(*bytesField).name, "")
	assert.Equal(t, schemaUnion.itemType[6].(*bytesField).hasDefault, false)
	assert.Equal(t, schemaUnion.itemType[7].(*recordField).name, "")
	assert.Equal(t, schemaUnion.itemType[7].(*recordField).typeName, "NestedRecord")
}
