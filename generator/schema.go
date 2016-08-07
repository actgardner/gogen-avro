package generator

import (
	"encoding/json"
	"fmt"
)

/*
 Deserialize the JSON definiton of a record and generate structs, deserializer and serializer methods.
 This function only supports JSON maps at the moment, where "type" -> "record". Avro also allows for schemas
 which are JSON arrays or JSON strings, but we don't currently support those as the root JSON type.
*/
func DeserializeRecordSchema(packageName string, schemaJson []byte, pkg *Package) error {
	r, err := deserializeRecordDefinition(schemaJson)
	if err != nil {
		return fmt.Errorf("Error decoding schema JSON: %v", err)
	}
	r.AddStruct(pkg)
	r.AddSerializer(pkg)
	r.AddDeserializer(pkg)
	return nil
}

/* Given a JSON record definition as a JSON encoded string, deserialize the JSON and build the record definition structs */
func deserializeRecordDefinition(schemaJson []byte) (*recordDefinition, error) {
	var schema interface{}
	if err := json.Unmarshal(schemaJson, &schema); err != nil {
		return nil, err
	}
	schemaMap, ok := schema.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid or unsupported schema - expected map")
	}
	return decodeRecordDefinition(schemaMap)
}

/* Given a map representing a record definition, validate the definition and build the recordDefinition struct */
func decodeRecordDefinition(schemaMap map[string]interface{}) (*recordDefinition, error) {
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
		name, ok := field["name"]
		if !ok {
			return nil, fmt.Errorf("Field is missing requird 'name' field")
		}
		nameStr, ok := name.(string)
		if !ok {
			return nil, fmt.Errorf("Field 'name' must be string type")
		}
		t, ok := field["type"]
		if !ok {
			return nil, fmt.Errorf("Field %q is missing required 'type' field", nameStr)
		}
		def, hasDef := field["default"]
		fieldStruct, err := decodeFieldDefinitionType(nameStr, t, def, hasDef)
		if err != nil {
			return nil, err
		}

		decodedFields = append(decodedFields, fieldStruct)
	}
	return &recordDefinition{
		name:   nameStr,
		fields: decodedFields,
	}, nil
}

func decodeFieldDefinitionType(nameStr string, t, def interface{}, hasDef bool) (field, error) {
	switch t.(type) {
	case string:
		typeStr := t.(string)
		return createFieldStruct(nameStr, typeStr, def, hasDef)
	case []interface{}:
		return decodeUnionDefinition(nameStr, def, hasDef, t.([]interface{}))
	case map[string]interface{}:
		return decodeComplexDefinition(nameStr, t.(map[string]interface{}))
	}
	return nil, fmt.Errorf("Field %v has invalid type - must be a JSON string, array or map", nameStr)
}

func decodeUnionDefinition(nameStr string, def interface{}, hasDef bool, fieldList []interface{}) (field, error) {
	unionFields := make([]field, 0)
	for i, f := range fieldList {
		var fieldDef field
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
		fieldType, err := decodeFieldDefinitionType("", items, nil, false)
		if err != nil {
			return nil, fmt.Errorf("Array %v item definition is invalid - %v", err)
		}
		return &arrayField{nameStr, fieldType}, nil
	case "map":
		values, ok := typeMap["values"]
		if !ok {
			return nil, fmt.Errorf("Field %q must have an 'values' field", nameStr)
		}
		fieldType, err := decodeFieldDefinitionType("", values, nil, false)
		if err != nil {
			return nil, fmt.Errorf("Array %v item definition is invalid - %v", err)
		}
		return &mapField{nameStr, fieldType}, nil
	case "enum":
		symbols, ok := typeMap["symbols"]
		if !ok {
			return nil, fmt.Errorf("Field %q must have a 'symbols' field", nameStr)
		}
		symbolSlice, ok := symbols.([]interface{})
		if !ok {
			return nil, fmt.Errorf("Field %q 'symbols' must be an array", nameStr)
		}
		symbolStr, ok := interfaceSliceToStringSlice(symbolSlice)
		if !ok {
			return nil, fmt.Errorf("Field %q 'symbols' must be an array of strings", nameStr)
		}
		typeName, ok := typeMap["name"]
		if !ok {
			return nil, fmt.Errorf("Field %q must have a 'name' field", nameStr)
		}
		typeNameStr, ok := typeName.(string)
		if !ok {
			return nil, fmt.Errorf("Field %q 'name' must be a string", nameStr)
		}
		return &enumField{nameStr, typeNameStr, "", false, symbolStr}, nil
	case "fixed":
		size, ok := typeMap["size"]
		if !ok {
			return nil, fmt.Errorf("Field %q must have a 'size' field", nameStr)
		}
		sizeInt, ok := size.(float64)
		if !ok {
			return nil, fmt.Errorf("Field %q 'size' must be an integer", nameStr)
		}
		typeName, ok := typeMap["name"]
		if !ok {
			return nil, fmt.Errorf("Field %q must have a 'name' field", nameStr)
		}
		typeNameStr, ok := typeName.(string)
		if !ok {
			return nil, fmt.Errorf("Field %q 'name' must be a string", nameStr)
		}
		return &fixedField{nameStr, typeNameStr, nil, false, int(sizeInt)}, nil
	case "record":
		def, err := decodeRecordDefinition(typeMap)
		if err != nil {
			return nil, err
		}
		return &recordField{nameStr, def.GoType(), def}, nil
	default:
		return nil, fmt.Errorf("Unknown complex type %v", typeStr)
	}
}

func createFieldStruct(nameStr, typeStr string, def interface{}, hasDef bool) (field, error) {
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
		return &recordField{nameStr, typeStr, nil}, nil
	}
}
