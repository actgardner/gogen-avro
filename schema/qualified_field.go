package schema

import (
	"fmt"

	"github.com/actgardner/gogen-avro/generator"
)

// qualifiedField is a primitive field with a qualified name and aliases
type qualifiedField struct {
	primitiveField
	qname   QualifiedName
	aliases []QualifiedName
}

func newQualifiedField(qname QualifiedName, aliases []QualifiedName, definition interface{}) qualifiedField {
	q := qualifiedField{aliases: aliases}
	q.definition = definition
	// Force set QName for creating all other inner fields from it
	q.setQualifiedName(qname)
	return q
}

func (q *qualifiedField) QualifiedName() QualifiedName {
	return q.qname
}

// Public name and GoType are remade if qname changes
func (q *qualifiedField) setQualifiedName(qname QualifiedName) {
	q.qname = qname
	publicName := generator.ToPublicName(qname.String())
	simpleName := generator.ToPublicSimpleName(qname.String())
	q.setName(publicName)
	q.setSimpleName(simpleName)
	q.setGoType(publicName)
	q.setSerializerMethod(fmt.Sprintf("write%s", publicName))
}

func (q *qualifiedField) Aliases() []QualifiedName {
	return q.aliases
}

func (q *qualifiedField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	if _, ok := scope[q.qname]; ok {
		return q.qname.String(), nil
	}
	scope[q.qname] = 1
	return q.primitiveField.Definition(scope)
}
