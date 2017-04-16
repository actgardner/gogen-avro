package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

const UTIL_FILE = "primitive.go"

/*
  An Avro qualified name, which includes an optional namespace and the type name.
*/
type QualifiedName struct {
	Namespace string
	Name      string
}

func (q QualifiedName) String() string {
	if q.Namespace == "" {
		return q.Name
	}
	return q.Namespace + "." + q.Name
}

type Schema struct {
	Root       AvroType
	JSONSchema []byte
}

/*
  Namespace is a mapping of QualifiedNames to their Definitions, used to resolve
  type lookups within a schema.
*/
type Namespace struct {
	Definitions map[QualifiedName]Definition
	Schemas     []Schema
}

func NewNamespace() *Namespace {
	return &Namespace{
		Definitions: make(map[QualifiedName]Definition),
		Schemas:     make([]Schema, 0),
	}
}

/*
  Add a new type definition to the namespace. Returns an error if the type is already defined.
*/
func (n *Namespace) RegisterDefinition(d Definition) error {
	if _, ok := n.Definitions[d.AvroName()]; ok {
		return fmt.Errorf("Conflicting definitions for %v", d.AvroName())
	}
	n.Definitions[d.AvroName()] = d

	for _, alias := range d.Aliases() {
		if _, ok := n.Definitions[alias]; ok {
			return fmt.Errorf("Conflicting alias for %v - %v", d.AvroName(), alias)
		}
		n.Definitions[alias] = d
	}
	return nil
}

/*
  Parse a name according to the Avro spec:
  - If the name contains a dot ('.'), the last part is the name and the rest is the namespace
  - Otherwise, the enclosing namespace is used
*/
func ParseAvroName(enclosing, name string) QualifiedName {
	lastIndex := strings.LastIndex(name, ".")
	if lastIndex != -1 {
		return QualifiedName{name[:lastIndex], name[lastIndex+1:]}
	}
	return QualifiedName{enclosing, name}
}

/*
  Given an Avro schema as a JSON string, decode it and return the AvroType defined at the top level:
    - a single record definition (JSON map)
    - a union of multiple types (JSON array)
    - an already-defined type (JSON string)

The Avro type defined at the top level and all the type definitions beneath it will also be added to this Namespace.
*/
func (n *Namespace) TypeForSchema(schemaJson []byte) (AvroType, error) {
	var schema interface{}
	if err := json.Unmarshal(schemaJson, &schema); err != nil {
		return nil, err
	}

	field, err := n.decodeTypeDefinition("", schema)
	if err != nil {
		return nil, err
	}

	n.Schemas = append(n.Schemas, Schema{field, schemaJson})
	return field, nil
}

func (n *Namespace) decodeTypeDefinition(namespace string, schema interface{}) (AvroType, error) {
	switch schema.(type) {
	case string:
		typeStr := schema.(string)
		return n.getTypeByName(namespace, typeStr, schema), nil

	case []interface{}:
		return n.decodeUnionDefinition(namespace, schema.([]interface{}))

	case map[string]interface{}:
		return n.decodeComplexDefinition(namespace, schema.(map[string]interface{}))

	}

	return nil, NewWrongMapValueTypeError("type", "array, string, map", schema)
}

/*
   Given a map representing a record definition, validate the definition and build the RecordDefinition struct.
*/
func (n *Namespace) decodeRecordDefinition(namespace string, schemaMap map[string]interface{}) (Definition, error) {
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

	decodedFields := make([]*Field, 0)
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

		fieldType, err := n.decodeTypeDefinition(namespace, t)
		if err != nil {
			return nil, err
		}

		def, hasDef := field["default"]
		fieldStruct := NewField(fieldName, fieldType, def, hasDef, field)

		decodedFields = append(decodedFields, fieldStruct)
	}

	aliases, err := parseAliases(schemaMap, namespace)
	if err != nil {
		return nil, err
	}

	return NewRecordDefinition(ParseAvroName(namespace, name), aliases, decodedFields, schemaMap), nil
}

/* Given a map representing an enum definition, validate the definition and build the EnumDefinition struct.
 */
func (n *Namespace) decodeEnumDefinition(namespace string, schemaMap map[string]interface{}) (Definition, error) {
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

	return NewEnumDefinition(ParseAvroName(namespace, name), aliases, symbolStr, schemaMap), nil
}

/* Given a map representing a fixed definition, validate the definition and build the FixedDefinition struct. */
func (n *Namespace) decodeFixedDefinition(namespace string, schemaMap map[string]interface{}) (Definition, error) {
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

	return NewFixedDefinition(ParseAvroName(namespace, name), aliases, int(sizeBytes), schemaMap), nil
}

func (n *Namespace) decodeUnionDefinition(namespace string, fieldList []interface{}) (AvroType, error) {
	unionFields := make([]AvroType, 0)
	for _, f := range fieldList {
		fieldDef, err := n.decodeTypeDefinition(namespace, f)
		if err != nil {
			return nil, err
		}

		unionFields = append(unionFields, fieldDef)
	}
	return NewUnionField(unionFields, fieldList), nil
}

func (n *Namespace) decodeComplexDefinition(namespace string, typeMap map[string]interface{}) (AvroType, error) {
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

		fieldType, err := n.decodeTypeDefinition(namespace, items)
		if err != nil {
			return nil, err
		}

		return NewArrayField(fieldType, typeMap), nil

	case "map":
		values, ok := typeMap["values"]
		if !ok {
			return nil, NewRequiredMapKeyError("values")
		}

		fieldType, err := n.decodeTypeDefinition(namespace, values)
		if err != nil {
			return nil,  err
		}

		return NewMapField(fieldType, typeMap), nil

	case "enum":
		definition, err := n.decodeEnumDefinition(namespace, typeMap)
		if err != nil {
			return nil, err
		}

		err = n.RegisterDefinition(definition)
		if err != nil {
			return nil, err
		}
		return NewReference(definition.AvroName()), nil

	case "fixed":
		definition, err := n.decodeFixedDefinition(namespace, typeMap)
		if err != nil {
			return nil, err
		}

		err = n.RegisterDefinition(definition)
		if err != nil {
			return nil, err
		}

		return NewReference(definition.AvroName()), nil

	case "record":
		definition, err := n.decodeRecordDefinition(namespace, typeMap)
		if err != nil {
			return nil, err
		}

		err = n.RegisterDefinition(definition)
		if err != nil {
			return nil, err
		}

		return NewReference(definition.AvroName()), nil

	default:
		// If the type isn't a special case, it's a primitive or a reference to an existing type
		return n.getTypeByName(namespace, typeStr, typeMap), nil
	}
}

func (n *Namespace) getTypeByName(namespace string, typeStr string, definition interface{}) AvroType {
	switch typeStr {
	case "int":
		return &intField{definition}

	case "long":
		return &longField{definition}

	case "float":
		return &floatField{definition}

	case "double":
		return &doubleField{definition}

	case "boolean":
		return &boolField{definition}

	case "bytes":
		return &bytesField{definition}

	case "string":
		return &stringField{definition}

	case "null":
		return &nullField{definition}
	}

	return  NewReference(ParseAvroName(namespace, typeStr))
}

/*
  Parse out all the aliases from a definition map - returns an empty slice if no aliases exist.
  Returns an error if the aliases key exists but the value isn't a list of strings.
*/
func parseAliases(objectMap map[string]interface{}, namespace string) ([]QualifiedName, error) {
	aliases, ok := objectMap["aliases"]
	if !ok {
		return make([]QualifiedName, 0), nil
	}

	aliasList, ok := aliases.([]interface{})
	if !ok {
		return nil, fmt.Errorf("Field aliases expected to be array, got %v", aliases)
	}

	qualifiedAliases := make([]QualifiedName, 0, len(aliasList))

	for _, alias := range aliasList {
		aliasString, ok := alias.(string)
		if !ok {
			return nil, fmt.Errorf("Field aliases expected to be array of strings, got %v", aliases)
		}
		qualifiedAliases = append(qualifiedAliases, ParseAvroName(namespace, aliasString))
	}
	return qualifiedAliases, nil
}
