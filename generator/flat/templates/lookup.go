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
	var templateDef string
	switch t.(type) {
	case *avro.ArrayField:
		templateDef = ArrayTemplate
	case *avro.MapField:
		templateDef = MapTemplate
	case *avro.UnionField:
		templateDef = UnionTemplate
	case *avro.EnumDefinition:
		templateDef = EnumTemplate
	case *avro.FixedDefinition:
		templateDef = FixedTemplate
	case *avro.RecordDefinition:
		templateDef = RecordTemplate
	default:
		return "", NoTemplateForType
	}
	return Evaluate(templateDef, t)
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

			return convertByteToInitForm(fingerprint), err
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

func convertByteToInitForm(b []byte) string {
	if len(b) == 0 {
		return "[]byte{}"
	}
	s := fmt.Sprintf("[]byte{0x%x", b[0])
	for _, v := range b[1:] {
		s = fmt.Sprintf("%s,0x%x", s, v)
	}
	s = fmt.Sprintf("%s}", s)
	return s
}
