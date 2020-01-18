package templates

import (
	"bytes"
	"fmt"
	"text/template"

	avro "github.com/actgardner/gogen-avro/v7/schema"
)

var NoTemplateForType = fmt.Errorf("No template exists for supplied type")

func Template(t avro.Node) (string, error) {
	var template string
	switch t.(type) {
	case *avro.ArrayField:
		template = ArrayTemplate
	case *avro.MapField:
		template = MapTemplate
	case *avro.UnionField:
		template = UnionTemplate
	case *avro.EnumDefinition:
		template = EnumTemplate
	case *avro.FixedDefinition:
		template = FixedTemplate
	case *avro.RecordDefinition:
		template = RecordTemplate
	default:
		return "", NoTemplateForType
	}
	return Evaluate(template, t)
}

func Evaluate(templateStr string, obj interface{}) (string, error) {
	buf := &bytes.Buffer{}
	t, err := template.New("").Parse(templateStr)
	if err != nil {
		return "", err
	}

	err = t.Execute(buf, obj)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
