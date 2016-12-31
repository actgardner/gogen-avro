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
	// Add the Avro definitions for the user-supplied schema
	r, err := addDefinitionForSchema(schemaJson, pkg)
	if err != nil {
		return err
	}

	// Add the generic Avro container schema
	_, err = addDefinitionForSchema([]byte(AVRO_BLOCK_SCHEMA), pkg)
	if err != nil {
		return err
	}

	_, err = addDefinitionForSchema([]byte(AVRO_HEADER_SCHEMA), pkg)
	if err != nil {
		return err
	}

	// Add the container definitions for this file
	containerWriter := &avroContainerWriter{schemaJson, r}
	containerWriter.AddAvroContainerWriter(pkg)
	return nil
}

func addDefinitionForSchema(schemaJson []byte, pkg *Package) (*recordDefinition, error) {
	r, err := deserializeRecordDefinition(schemaJson)
	if err != nil {
		return nil, err
	}

	// Add the Avro type definitions
	r.AddStruct(pkg)
	r.AddSerializer(pkg)
	r.AddDeserializer(pkg)
	return r, nil
}

/* Given a JSON record definition as a JSON encoded string, deserialize the JSON and build the record definition structs */
func deserializeRecordDefinition(schemaJson []byte) (*recordDefinition, error) {
	var schema interface{}
	if err := json.Unmarshal(schemaJson, &schema); err != nil {
		return nil, err
	}
	schemaMap, ok := schema.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("Invalid or unsupported schema - expected map as root JSON object")
	}
	return decodeRecordDefinition(schemaMap)
}

/* Given a map representing a record definition, validate the definition and build the recordDefinition struct */
func decodeRecordDefinition(schemaMap map[string]interface{}) (*recordDefinition, error) {
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

	decodedFields := make([]field, 0)
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

	return &recordDefinition{
		name:   name,
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
	return nil, NewSchemaError(nameStr, NewWrongMapValueTypeError("type", "array, string, map", t))
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
		fieldType, err := decodeFieldDefinitionType("", items, nil, false)
		if err != nil {
			return nil, NewSchemaError(nameStr, err)
		}
		return &arrayField{nameStr, fieldType}, nil
	case "map":
		values, ok := typeMap["values"]
		if !ok {
			return nil, NewSchemaError(nameStr, NewRequiredMapKeyError("values"))
		}
		fieldType, err := decodeFieldDefinitionType("", values, nil, false)
		if err != nil {
			return nil, NewSchemaError(nameStr, err)
		}
		return &mapField{nameStr, fieldType}, nil
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
		return &recordField{nameStr, def.GoType(), def}, nil
	default:
		return nil, NewSchemaError(nameStr, fmt.Errorf("Unknown type name %v", typeStr))
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
