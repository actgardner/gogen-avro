package templates

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"

	avro "github.com/actgardner/gogen-avro/schema"
	"github.com/actgardner/gogen-avro/schema/canonical"
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
	t, err := template.New("").Funcs(template.FuncMap{
		"definitionFingerprint": func(def avro.Definition) (string, error) {
			cf := canonical.DefinitionCanonicalForm(def)
			encoded, err := json.Marshal(cf)
			if err != nil {
				return "", err
			}
			fingerprint := canonical.AvroCRC64Fingerprint(encoded)
			return fingerprint, err
		},
	}).Parse(templateStr)
	if err != nil {
		return "", err
	}

	err = t.Execute(buf, obj)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
