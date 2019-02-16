package types

import (
	"fmt"
	"strings"

	"github.com/actgardner/gogen-avro/generator"
)

const enumTypeDef = `
%v
type %v int32

const (
%v
)
`

const enumTypeStringer = `
func (e %v) String() string {
	switch e {
%v
	}
	return "unknown"
}
`

const enumSerializerDef = `
func %v(r %v, w io.Writer) error {
	return writeInt(int32(r), w)
}
`

const enumDeserializerDef = `
func %v(r io.Reader) (%v, error) {
	val, err := readInt(r)
	return %v(val), err
}
`

type EnumDefinition struct {
	name       QualifiedName
	aliases    []QualifiedName
	symbols    []string
	doc        string
	definition map[string]interface{}
}

func NewEnumDefinition(name QualifiedName, aliases []QualifiedName, symbols []string, doc string, definition map[string]interface{}) *EnumDefinition {
	return &EnumDefinition{
		name:       name,
		aliases:    aliases,
		symbols:    symbols,
		doc:        doc,
		definition: definition,
	}
}

func (e *EnumDefinition) Name() string {
	return e.GoType()
}

func (e *EnumDefinition) SimpleName() string {
	return e.name.Name
}

func (e *EnumDefinition) AvroName() QualifiedName {
	return e.name
}

func (e *EnumDefinition) Aliases() []QualifiedName {
	return e.aliases
}

func (e *EnumDefinition) GoType() string {
	return generator.ToPublicName(e.name.Name)
}

func (e *EnumDefinition) typeList() string {
	typeStr := ""
	for i, t := range e.symbols {
		typeStr += fmt.Sprintf("%v %v = %v\n", generator.ToPublicName(e.GoType()+strings.Title(t)), e.GoType(), i)
	}
	return typeStr
}

func (e *EnumDefinition) stringerList() string {
	stringerStr := ""
	for _, t := range e.symbols {
		stringerStr += fmt.Sprintf("case %v:\n return %q\n", generator.ToPublicName(e.GoType()+strings.Title(t)), t)
	}
	return stringerStr
}

func (e *EnumDefinition) structDef() string {
	var doc string
	if e.doc != "" {
		doc = fmt.Sprintf("// %v", e.doc)
	}
	return fmt.Sprintf(enumTypeDef, doc, e.GoType(), e.typeList())
}

func (e *EnumDefinition) stringerDef() string {
	return fmt.Sprintf(enumTypeStringer, e.GoType(), e.stringerList())
}

func (e *EnumDefinition) serializerMethodDef() string {
	return fmt.Sprintf(enumSerializerDef, e.SerializerMethod(), e.GoType())
}

func (e *EnumDefinition) SerializerMethod() string {
	return "write" + e.GoType()
}

func (e *EnumDefinition) deserializerMethodDef() string {
	return fmt.Sprintf(enumDeserializerDef, e.DeserializerMethod(), e.GoType(), e.GoType())
}

func (e *EnumDefinition) DeserializerMethod() string {
	return "read" + e.GoType()
}

func (e *EnumDefinition) filename() string {
	return generator.ToSnake(e.GoType()) + ".go"
}

func (e *EnumDefinition) AddStruct(p *generator.Package, _ bool) error {
	p.AddStruct(e.filename(), e.GoType(), e.structDef())
	p.AddFunction(e.filename(), e.GoType(), "String", e.stringerDef())
	return nil
}

func (e *EnumDefinition) AddSerializer(p *generator.Package) {
	p.AddStruct(UTIL_FILE, "ByteWriter", byteWriterInterface)
	p.AddFunction(UTIL_FILE, "", "writeInt", writeIntMethod)
	p.AddFunction(UTIL_FILE, "", "encodeInt", encodeIntMethod)
	p.AddFunction(UTIL_FILE, "", e.SerializerMethod(), e.serializerMethodDef())
	p.AddImport(UTIL_FILE, "io")
}

func (e *EnumDefinition) AddDeserializer(p *generator.Package) {
	p.AddFunction(UTIL_FILE, "", "readInt", readIntMethod)
	p.AddFunction(UTIL_FILE, "", e.DeserializerMethod(), e.deserializerMethodDef())
	p.AddImport(UTIL_FILE, "io")
}

func (s *EnumDefinition) ResolveReferences(n *Namespace) error {
	return nil
}

func (s *EnumDefinition) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	if _, ok := scope[s.name]; ok {
		return s.name.String(), nil
	}
	scope[s.name] = 1
	return s.definition, nil
}

func (s *EnumDefinition) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	if _, ok := rvalue.(string); !ok {
		return "", fmt.Errorf("Expected string as default for field %v, got %q", lvalue, rvalue)
	}

	return fmt.Sprintf("%v = %v", lvalue, generator.ToPublicName(s.GoType()+strings.Title(rvalue.(string)))), nil
}
