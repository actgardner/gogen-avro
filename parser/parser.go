package parser

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	avro_schema "github.com/actgardner/gogen-avro/schema"
)

// ParseAvroName parses a name according to the Avro spec:
//   - If the name contains a dot ('.'), the last part is the name and the rest is the namespace
//   - Otherwise, the enclosing namespace is used
func ParseAvroName(enclosing, name string) avro_schema.QualifiedName {
	lastIndex := strings.LastIndex(name, ".")
	if lastIndex != -1 {
		enclosing = name[:lastIndex]
	}
	return avro_schema.QualifiedName{enclosing, name[lastIndex+1:]}
}

// Parser decodes Avro type definitions, validates them and accumulates them so the references can be resolved
type Parser struct {
	Definitions map[avro_schema.QualifiedName]avro_schema.Definition
}

func NewParser() *Parser {
	return &Parser{
		Definitions: make(map[avro_schema.QualifiedName]avro_schema.Definition),
	}
}

// RegisterDefinition adds a new type definition to the namespace. Returns an error if the type is already defined.
func (n *Parser) RegisterDefinition(d avro_schema.Definition) error {
	if curDef, ok := n.Definitions[d.AvroName()]; ok {
		if !reflect.DeepEqual(curDef, d) {
			return fmt.Errorf("Conflicting definitions for %v", d.AvroName())
		}
		return nil
	}
	n.Definitions[d.AvroName()] = d

	for _, alias := range d.Aliases() {
		if existing, ok := n.Definitions[alias]; ok {
			return fmt.Errorf("Alias for %q is %q, but %q is already aliased with that name", d.AvroName(), alias, existing.AvroName())
		}
		n.Definitions[alias] = d
	}
	return nil
}

// Parse accepts an Avro schema as a JSON string and adds any type definitions to the given Parser
func (n *Parser) Parse(schemaJson []byte) (avro_schema.AvroType, error) {
	var schema interface{}
	if err := json.Unmarshal(schemaJson, &schema); err != nil {
		return nil, err
	}

	field, err := n.decodeTypeDefinition("topLevel", "", schema)
	if err != nil {
		return nil, err
	}

	return field, nil
}

func (n *Parser) decodeTypeDefinition(name, namespace string, schema interface{}) (avro_schema.AvroType, error) {
	switch schema.(type) {
	case string:
		typeStr := schema.(string)
		return n.getTypeByName(namespace, typeStr, schema), nil

	case []interface{}:
		return n.decodeUnionDefinition(name, namespace, schema.([]interface{}))

	case map[string]interface{}:
		return n.decodeComplexDefinition(name, namespace, schema.(map[string]interface{}))

	}

	return nil, NewWrongMapValueTypeError("type", "array, string, map", schema)
}

// Given a map representing a record definition, validate the definition and build the RecordDefinition struct.
func (n *Parser) decodeRecordDefinition(namespace string, schemaMap map[string]interface{}) (avro_schema.Definition, error) {
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

	var rDocString string
	if rDoc, ok := schemaMap["doc"]; ok {
		rDocString, ok = rDoc.(string)
		if !ok {
			return nil, NewWrongMapValueTypeError("doc", "string", rDoc)
		}
	}

	if _, ok := schemaMap["namespace"]; ok {
		namespace, err = getMapString(schemaMap, "namespace")
		if err != nil {
			return nil, err
		}
	}

	fieldList, err := getMapArray(schemaMap, "fields")
	if err != nil {
		return nil, err
	}

	decodedFields := make([]*avro_schema.Field, 0)
	for i, f := range fieldList {
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

		fieldType, err := n.decodeTypeDefinition(fieldName, namespace, t)
		if err != nil {
			return nil, err
		}

		var docString string
		if doc, ok := field["doc"]; ok {
			docString, ok = doc.(string)
			if !ok {
				return nil, NewWrongMapValueTypeError("doc", "string", doc)
			}
		}

		var fieldTags string
		if tags, ok := field["golang.tags"]; ok {
			fieldTags, ok = tags.(string)
			if !ok {
				return nil, NewWrongMapValueTypeError("golang.tags", "string", tags)
			}
		}

		var fieldAliases []string
		if aliases, ok := field["aliases"]; ok {
			aliasList, ok := aliases.([]interface{})
			if !ok {
				return nil, NewWrongMapValueTypeError("aliases", "[]string", aliases)
			}

			for _, aliasVal := range aliasList {
				aliasStr, ok := aliasVal.(string)
				if !ok {
					return nil, NewWrongMapValueTypeError("aliases", "[]string", aliases)
				}
				fieldAliases = append(fieldAliases, aliasStr)
			}
		}

		def, hasDef := field["default"]
		fieldStruct := avro_schema.NewField(fieldName, fieldType, def, hasDef, fieldAliases, docString, field, i, fieldTags)

		decodedFields = append(decodedFields, fieldStruct)
	}

	aliases, err := parseAliases(schemaMap, namespace)
	if err != nil {
		return nil, err
	}

	return avro_schema.NewRecordDefinition(ParseAvroName(namespace, name), aliases, decodedFields, rDocString, schemaMap), nil
}

// decodeEnumDefinition accepts a namespace and a map representing an enum definition,
// it validates the definition and build the EnumDefinition struct.
func (n *Parser) decodeEnumDefinition(namespace string, schemaMap map[string]interface{}) (avro_schema.Definition, error) {
	typeStr, err := getMapString(schemaMap, "type")
	if err != nil {
		return nil, err
	}

	if typeStr != "enum" {
		return nil, fmt.Errorf("Type of enum must be 'enum'")
	}

	if _, ok := schemaMap["namespace"]; ok {
		namespace, err = getMapString(schemaMap, "namespace")
		if err != nil {
			return nil, err
		}
	}

	name, err := getMapString(schemaMap, "name")
	if err != nil {
		return nil, err
	}

	symbolSlice, err := getMapArray(schemaMap, "symbols")
	if err != nil {
		return nil, err
	}

	symbolStr, ok := interfaceSliceToStringSlice(symbolSlice)
	if !ok {
		return nil, fmt.Errorf("'symbols' must be an array of strings")
	}

	aliases, err := parseAliases(schemaMap, namespace)
	if err != nil {
		return nil, err
	}

	var docString string
	if doc, ok := schemaMap["doc"]; ok {
		if docString, ok = doc.(string); !ok {
			return nil, fmt.Errorf("'doc' must be a string")
		}
	}

	return avro_schema.NewEnumDefinition(ParseAvroName(namespace, name), aliases, symbolStr, docString, schemaMap), nil
}

// decodeFixedDefinition accepts a namespace and a map representing a fixed definition,
// it validates the definition and build the Fixedschema.Definition struct.
func (n *Parser) decodeFixedDefinition(namespace string, schemaMap map[string]interface{}) (avro_schema.Definition, error) {
	typeStr, err := getMapString(schemaMap, "type")
	if err != nil {
		return nil, err
	}

	if typeStr != "fixed" {
		return nil, fmt.Errorf("Type of fixed must be 'fixed'")
	}

	if _, ok := schemaMap["namespace"]; ok {
		namespace, err = getMapString(schemaMap, "namespace")
		if err != nil {
			return nil, err
		}
	}

	name, err := getMapString(schemaMap, "name")
	if err != nil {
		return nil, err
	}

	sizeBytes, err := getMapFloat(schemaMap, "size")
	if err != nil {
		return nil, err
	}

	aliases, err := parseAliases(schemaMap, namespace)
	if err != nil {
		return nil, err
	}

	return avro_schema.NewFixedDefinition(ParseAvroName(namespace, name), aliases, int(sizeBytes), schemaMap), nil
}

func (n *Parser) decodeUnionDefinition(name, namespace string, fieldList []interface{}) (avro_schema.AvroType, error) {
	unionFields := make([]avro_schema.AvroType, 0)
	for _, f := range fieldList {
		fieldDef, err := n.decodeTypeDefinition(name, namespace, f)
		if err != nil {
			return nil, err
		}

		unionFields = append(unionFields, fieldDef)
	}

	return avro_schema.NewUnionField(unionFields, fieldList), nil
}

func (n *Parser) decodeComplexDefinition(name, namespace string, typeMap map[string]interface{}) (avro_schema.AvroType, error) {
	typeStr, err := getMapString(typeMap, "type")
	if err != nil {
		return nil, err
	}
	switch typeStr {
	case "array":
		items, ok := typeMap["items"]
		if !ok {
			return nil, NewRequiredMapKeyError("items")
		}

		fieldType, err := n.decodeTypeDefinition(name, namespace, items)
		if err != nil {
			return nil, err
		}

		return avro_schema.NewArrayField(fieldType, typeMap), nil

	case "map":
		values, ok := typeMap["values"]
		if !ok {
			return nil, NewRequiredMapKeyError("values")
		}

		fieldType, err := n.decodeTypeDefinition(name, namespace, values)
		if err != nil {
			return nil, err
		}

		return avro_schema.NewMapField(fieldType, typeMap), nil

	case "enum":
		definition, err := n.decodeEnumDefinition(namespace, typeMap)
		if err != nil {
			return nil, err
		}

		err = n.RegisterDefinition(definition)
		if err != nil {
			return nil, err
		}
		return avro_schema.NewReference(definition.AvroName()), nil

	case "fixed":
		definition, err := n.decodeFixedDefinition(namespace, typeMap)
		if err != nil {
			return nil, err
		}

		err = n.RegisterDefinition(definition)
		if err != nil {
			return nil, err
		}

		return avro_schema.NewReference(definition.AvroName()), nil

	case "record":
		definition, err := n.decodeRecordDefinition(namespace, typeMap)
		if err != nil {
			return nil, err
		}

		err = n.RegisterDefinition(definition)
		if err != nil {
			return nil, err
		}

		return avro_schema.NewReference(definition.AvroName()), nil

	default:
		// If the type isn't a special case, it's a primitive or a reference to an existing type
		return n.getTypeByName(namespace, typeStr, typeMap), nil
	}
}

func (n *Parser) getTypeByName(namespace string, typeStr string, definition interface{}) avro_schema.AvroType {
	switch typeStr {
	case "int":
		return avro_schema.NewIntField(definition)

	case "long":
		return avro_schema.NewLongField(definition)

	case "float":
		return avro_schema.NewFloatField(definition)

	case "double":
		return avro_schema.NewDoubleField(definition)

	case "boolean":
		return avro_schema.NewBoolField(definition)

	case "bytes":
		return avro_schema.NewBytesField(definition)

	case "string":
		return avro_schema.NewStringField(definition)

	case "null":
		return avro_schema.NewNullField(definition)
	}

	return avro_schema.NewReference(ParseAvroName(namespace, typeStr))
}

// parseAliases parses out all the aliases from a definition map - returns an empty slice if no aliases exist.
// Returns an error if the aliases key exists but the value isn't a list of strings.
func parseAliases(objectMap map[string]interface{}, namespace string) ([]avro_schema.QualifiedName, error) {
	aliases, ok := objectMap["aliases"]
	if !ok {
		return make([]avro_schema.QualifiedName, 0), nil
	}

	aliasList, ok := aliases.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Field aliases expected to be array, got %v", aliases)
	}

	qualifiedAliases := make([]avro_schema.QualifiedName, 0, len(aliasList))

	for _, alias := range aliasList {
		aliasString, ok := alias.(string)
		if !ok {
			return nil, fmt.Errorf("Field aliases expected to be array of strings, got %v", aliases)
		}
		qualifiedAliases = append(qualifiedAliases, ParseAvroName(namespace, aliasString))
	}
	return qualifiedAliases, nil
}
