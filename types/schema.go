package types

import (
	"encoding/json"
	"fmt"
)

const UTIL_FILE = "primitive.go"

/* 
  Given an Avro schema as a JSON string, decode it and return the Field defined at the top level:
    - a single record definition (JSON map)
    - a union of multiple types (JSON array)  
    - an already-defined type (JSON string)
 */
func FieldDefinitionForSchema(schemaJson []byte) (Field, error) {
	var schema interface{}
	if err := json.Unmarshal(schemaJson, &schema); err != nil {
		return nil, err
	}
	return decodeFieldDefinitionType("", schema, nil, false)
}

func decodeFieldDefinitionType(nameStr string, t, def interface{}, hasDef bool) (Field, error) {
	switch t.(type) {
	case string:
		typeStr := t.(string)
		return createFieldStruct(nameStr, typeStr, def, hasDef)
	case []interface{}:
		return decodeUnionDefinition(nameStr, def, hasDef, t.([]interface{}))
	case map[string]interface{}:
		return decodeComplexDefinition(nameStr, t.(map[string]interface{}))
	}
	return nil, NewSchemaError(nameStr, NewWrongMapValueTypeError("type", "array, string, map", t))
}

/* Given a map representing a record definition, validate the definition and build the recordDefinition struct */
func decodeRecordDefinition(schemaMap map[string]interface{}) (*RecordDefinition, error) {
	typeStr, err := getMapString(schemaMap, "type")
	if err != nil {
		return nil, err
	}

	if typeStr != "record" {
		return nil, fmt.Errorf("Type of record must be 'record'")
	}

	name, err := getMapString(schemaMap, "name")
	if err != nil {
		return nil, err
	}

	fieldList, err := getMapArray(schemaMap, "fields")
	if err != nil {
		return nil, err
	}

	decodedFields := make([]Field, 0)
	for _, f := range fieldList {
		field, ok := f.(map[string]interface{})
		if !ok {
			return nil, NewWrongMapValueTypeError("fields", "map[]", field)
		}
		fieldName, err := getMapString(field, "name")
		if err != nil {
			return nil, err
		}
		t, ok := field["type"]
		if !ok {
			return nil, NewRequiredMapKeyError("type")
		}
		def, hasDef := field["default"]
		fieldStruct, err := decodeFieldDefinitionType(fieldName, t, def, hasDef)
		if err != nil {
			return nil, err
		}

		decodedFields = append(decodedFields, fieldStruct)
	}

	return &RecordDefinition{
		name:   name,
		fields: decodedFields,
	}, nil
}

func decodeUnionDefinition(nameStr string, def interface{}, hasDef bool, FieldList []interface{}) (Field, error) {
	unionFields := make([]Field, 0)
	for i, f := range FieldList {
		var fieldDef Field
		var err error
		if i == 0 {
			fieldDef, err = decodeFieldDefinitionType("", f, def, hasDef)
		} else {
			fieldDef, err = decodeFieldDefinitionType("", f, nil, false)
		}
		if err != nil {
			return nil, err
		}
		unionFields = append(unionFields, fieldDef)
	}
	return &unionField{nameStr, hasDef, unionFields}, nil
}

func decodeComplexDefinition(nameStr string, typeMap map[string]interface{}) (Field, error) {
	typeStr, err := getMapString(typeMap, "type")
	if err != nil {
		return nil, NewSchemaError(nameStr, err)
	}
	switch typeStr {
	case "array":
		items, ok := typeMap["items"]
		if !ok {
			return nil, NewSchemaError(nameStr, NewRequiredMapKeyError("items"))
		}
		FieldType, err := decodeFieldDefinitionType("", items, nil, false)
		if err != nil {
			return nil, NewSchemaError(nameStr, err)
		}
		return &arrayField{nameStr, FieldType}, nil
	case "map":
		values, ok := typeMap["values"]
		if !ok {
			return nil, NewSchemaError(nameStr, NewRequiredMapKeyError("values"))
		}
		FieldType, err := decodeFieldDefinitionType("", values, nil, false)
		if err != nil {
			return nil, NewSchemaError(nameStr, err)
		}
		return &mapField{nameStr, FieldType}, nil
	case "enum":
		symbolSlice, err := getMapArray(typeMap, "symbols")
		if err != nil {
			return nil, NewSchemaError(nameStr, err)
		}
		symbolStr, ok := interfaceSliceToStringSlice(symbolSlice)
		if !ok {
			return nil, NewSchemaError(nameStr, fmt.Errorf("'symbols' must be an array of strings"))
		}
		typeNameStr, err := getMapString(typeMap, "name")
		if err != nil {
			return nil, NewSchemaError(nameStr, err)
		}
		return &enumField{nameStr, typeNameStr, "", false, symbolStr}, nil
	case "fixed":
		size, err := getMapFloat(typeMap, "size")
		if err != nil {
			return nil, NewSchemaError(nameStr, err)
		}
		typeNameStr, err := getMapString(typeMap, "name")
		if err != nil {
			return nil, NewSchemaError(nameStr, err)
		}
		return &fixedField{nameStr, typeNameStr, nil, false, int(size)}, nil
	case "record":
		def, err := decodeRecordDefinition(typeMap)
		if err != nil {
			return nil, NewSchemaError(nameStr, err)
		}
		return &recordField{nameStr, def.FieldType(), def}, nil
	default:
		return nil, NewSchemaError(nameStr, fmt.Errorf("Unknown type name %v", typeStr))
	}
}

func createFieldStruct(nameStr, typeStr string, def interface{}, hasDef bool) (Field, error) {
	switch typeStr {
	case "string":
		var defStr string
		var ok bool
		if hasDef {
			defStr, ok = def.(string)
			if !ok {
				return nil, fmt.Errorf("Default value must be string type")
			}

		}
		return &stringField{nameStr, defStr, hasDef}, nil
	case "int":
		var defInt int32
		if hasDef {
			defFloat, ok := def.(float64)
			if !ok {
				return nil, fmt.Errorf("Default must be float type")
			}
			defInt = int32(defFloat)

		}
		return &intField{nameStr, defInt, hasDef}, nil
	case "long":
		var defInt int64
		if hasDef {
			defFloat, ok := def.(float64)
			if !ok {
				return nil, fmt.Errorf("Field %q default must be float type", nameStr)
			}
			defInt = int64(defFloat)
		}
		return &longField{nameStr, defInt, hasDef}, nil
	case "float":
		var defFloat float64
		var ok bool
		if hasDef {
			defFloat, ok = def.(float64)
			if !ok {
				return nil, fmt.Errorf("Field %q default must be float type", nameStr)
			}
		}
		return &floatField{nameStr, float32(defFloat), hasDef}, nil
	case "double":
		var defFloat float64
		var ok bool
		if hasDef {
			defFloat, ok = def.(float64)
			if !ok {
				return nil, fmt.Errorf("Field %q default must be float type", nameStr)
			}
		}
		return &doubleField{nameStr, defFloat, hasDef}, nil
	case "boolean":
		var defBool bool
		var ok bool
		if hasDef {
			defBool, ok = def.(bool)
			if !ok {
				return nil, fmt.Errorf("Field %q default must be bool type", nameStr)
			}

		}
		return &boolField{nameStr, defBool, hasDef}, nil
	case "bytes":
		var defBytes []byte
		if hasDef {
			defString, ok := def.(string)
			if !ok {
				return nil, fmt.Errorf("Field %q default must be string type", nameStr)
			}
			defBytes = []byte(defString)
		}
		return &bytesField{nameStr, defBytes, hasDef}, nil
	case "null":
		return &nullField{nameStr, hasDef}, nil
	default:
		return &recordField{nameStr, typeStr, nil}, nil
	}
}
