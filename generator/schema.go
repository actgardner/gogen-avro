package generator

import (
	"encoding/json"
	"fmt"
	"go/format"
)

/*
	Given a JSON Avro schema, produce a struct and serializer/deserializer pair
	TODO: Figure out how this should handle multiple record definitions
*/
func GenerateForSchema(schemaJson []byte) (string, error) {
	r, err := decodeSchema(schemaJson)
	if err != nil {
		return "", fmt.Errorf("Error decoding schema JSON: %v", err)
	}
	imports := make(map[string]string)
	ns := make(map[string]string)
	structDef := r.structDefinition()
	r.namespaceMap(imports, ns)
	src := fmt.Sprintf("package avro\n%v\n%v\n%v", concatSortedMap(imports, "\n"), structDef, concatSortedMap(ns, "\n"))
	fmtSrc, err := format.Source([]byte(src))
	return string(fmtSrc), err
}

/* Decode the schema for a single Record */
func decodeSchema(schemaJson []byte) (*recordDefinition, error) {
	var schema interface{}
	if err := json.Unmarshal(schemaJson, &schema); err != nil {
		return nil, err
	}
	schemaMap, ok := schema.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid or unsupported schema - expected map")
	}
	t, ok := schemaMap["type"]
	if !ok {
		return nil, fmt.Errorf("Schema is missing required field 'type'")
	}
	typeStr, ok := t.(string)
	if !ok || typeStr != "record" {
		return nil, fmt.Errorf("Schema type must be 'record'")
	}
	name, ok := schemaMap["name"]
	if !ok {
		return nil, fmt.Errorf("Record schema missing required field 'name'")
	}
	nameStr, ok := name.(string)
	if !ok {
		return nil, fmt.Errorf("Record schema field 'name' must be string")
	}
	fields, ok := schemaMap["fields"]
	if !ok {
		return nil, fmt.Errorf("Record schema missing required field 'fields'")
	}
	fieldList, ok := fields.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Record schema field 'fields' must be an array")
	}
	decodedFields := make([]field, 0)
	for _, f := range fieldList {
		field, ok := f.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("Record schema field 'fields' elements must be maps")
		}
		fieldStruct, err := decodeField(field)
		if err != nil {
			fmt.Printf("Decoded field err %v - %v\n", f, err)
			return nil, err
		}
		decodedFields = append(decodedFields, fieldStruct)
	}
	return &recordDefinition{
		name:   nameStr,
		fields: decodedFields,
	}, nil
}

func decodeField(fieldMap map[string]interface{}) (field, error) {
	name, ok := fieldMap["name"]
	if !ok {
		return nil, fmt.Errorf("Field is missing requird 'name' field")
	}
	nameStr, ok := name.(string)
	if !ok {
		return nil, fmt.Errorf("Field 'name' must be string type")
	}
	return decodeFieldDefinition(nameStr, fieldMap)
}

func decodeFieldDefinition(nameStr string, fieldMap map[string]interface{}) (field, error) {
	t, ok := fieldMap["type"]
	if !ok {
		return nil, fmt.Errorf("Field %q is missing required 'type' field", nameStr)
	}
	def, hasDef := fieldMap["default"]
	switch t.(type) {
	case string:
		typeStr := t.(string)
		return createFieldStruct(nameStr, typeStr, def, hasDef, fieldMap)
	case []interface{}:
		return decodeUnionDefinition(nameStr, def, hasDef, t.([]interface{}))
	case map[string]interface{}:
		return decodeComplexDefinition(nameStr, t.(map[string]interface{}))
	}
	return nil, fmt.Errorf("")
}

func decodeUnionDefinition(nameStr string, def interface{}, hasDef bool, fieldList []interface{}) (field, error) {
	unionFields := make([]field, 0)
	for i, f := range fieldList {
		typeStr, ok := f.(string)
		if !ok {
			return nil, fmt.Errorf("Union members for %v is not of type strings", nameStr)
		}
		var fieldDef field
		var err error
		if i == 0 {
			fieldDef, err = createFieldStruct("", typeStr, def, hasDef, nil)
		} else {
			fieldDef, err = createFieldStruct("", typeStr, nil, false, nil)
		}
		if err != nil {
			return nil, err
		}
		unionFields = append(unionFields, fieldDef)
	}
	return &unionField{nameStr, hasDef, unionFields}, nil
}

func decodeComplexDefinition(nameStr string, typeMap map[string]interface{}) (field, error) {
	t, ok := typeMap["type"]
	if !ok {
		return nil, fmt.Errorf("Field %q is missing required complex 'type' field", nameStr)
	}
	typeStr, ok := t.(string)
	if !ok {
		return nil, fmt.Errorf("Field %q complex 'type' field must be string", nameStr)
	}
	switch typeStr {
	case "array":
		items, ok := typeMap["items"]
		if !ok {
			return nil, fmt.Errorf("Field %q must have an 'items' field", nameStr)
		}
		var fieldType field
		var err error
		switch items.(type) {
		case string:
			fieldType, err = createFieldStruct("", items.(string), nil, false, nil)
		case map[string]interface{}:
			fieldType, err = decodeFieldDefinition("", items.(map[string]interface{}))
		case []interface{}:
			fieldType, err = decodeUnionDefinition("", nil, false, items.([]interface{}))

		default:
			return nil, fmt.Errorf("Array %v items type must be a string or map", nameStr)
		}
		if err != nil {
			return nil, fmt.Errorf("Array %v item definition is invalid - %v", err)
		}
		return &arrayField{nameStr, fieldType}, nil
	case "map":
		items, ok := typeMap["values"]
		if !ok {
			return nil, fmt.Errorf("Field %q must have an 'values' field", nameStr)
		}
		var fieldType field
		var err error
		switch items.(type) {
		case string:
			fieldType, err = createFieldStruct("", items.(string), nil, false, nil)
		case map[string]interface{}:
			fieldType, err = decodeFieldDefinition("", items.(map[string]interface{}))
		case []interface{}:
			fieldType, err = decodeUnionDefinition("", nil, false, items.([]interface{}))
		default:
			return nil, fmt.Errorf("Array %v items type must be a string or map", nameStr)
		}
		if err != nil {
			return nil, fmt.Errorf("Array %v item definition is invalid - %v", err)
		}
		return &mapField{nameStr, fieldType}, nil
	default:
		return nil, fmt.Errorf("Unknown complex type %v", typeStr)
	}
}

func createFieldStruct(nameStr, typeStr string, def interface{}, hasDef bool, fieldMap map[string]interface{}) (field, error) {
	switch typeStr {
	case "string":
		var defStr string
		var ok bool
		if hasDef {
			defStr, ok = def.(string)
			if !ok {
				return nil, fmt.Errorf("Field %q default must be string type", nameStr)
			}

		}
		return &stringField{nameStr, defStr, hasDef}, nil
	case "int":
		var defInt int32
		if hasDef {
			defFloat, ok := def.(float64)
			if !ok {
				return nil, fmt.Errorf("Field %q default must be float type", nameStr)
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
		return &recordField{nameStr, typeStr}, nil
	}
}
