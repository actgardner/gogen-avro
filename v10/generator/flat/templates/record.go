package templates

const RecordTemplate = `
import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/actgardner/gogen-avro/v10/vm/types"
	"github.com/actgardner/gogen-avro/v10/vm"
	"github.com/actgardner/gogen-avro/v10/compiler"
)

var _ = fmt.Printf

{{ if ne .Doc "" }}// {{ .Doc}}{{ end }}
type {{ .Name }} struct {
{{ range $i, $field := .Fields -}}
	{{ if ne $field.Doc "" }}// {{ $field.Doc }}{{ end }}
	{{ if ne $field.Tags "" -}}
		{{ $field.GoName }} {{ $field.Type.GoType }} ` + "`{{ $field.Tags }}`" + `
	{{ else -}}
		{{ $field.GoName }} {{ $field.Type.GoType }}
	{{ end -}}
{{ end }}
}

const {{ .Name }}AvroCRC64Fingerprint = {{ definitionFingerprint . }}

func {{ .ConstructorMethod }} ({{ .GoType}}) {
	r := {{ .Name }}{}
	{{ range $i, $field := .Fields -}}
	{{ if $.ConstructableForField $field | ne "" -}}
		{{ if not (hasNullDefault $field.Type) -}}
			{{ $.ConstructableForField $field }}
		{{ end -}}
	{{ end -}}
	{{ if and .HasDefault (not (hasNullDefault $field.Type)) -}}
		{{ $.DefaultForField $field }}
	{{ end -}}
	{{ end -}}
	return r
}

func Deserialize{{ .Name }}(r io.Reader) ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}
	deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func Deserialize{{ .Name }}FromSchema(r io.Reader, schema string) ({{ .GoType }}, error) {
	t := {{ .ConstructorMethod }}

	deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	if err != nil {
		return t, err
	}

	err = vm.Eval(r, deser, &t)
	return t, err
}

func {{ .SerializerMethod }}(r {{ .GoType }}, w io.Writer) error {
	var err error
	{{ range $i, $field := .Fields -}}
    {{ if $.IsSimpleNullUnion $field -}}
	if r.{{ .GoName }} == nil {
		err = vm.WriteLong({{ $.SimpleNullUnionNullIndex $field }}, w)
		if err != nil {
			return err
		}
	} else {
		err = vm.WriteLong(int64({{ $.SimpleNullUnionNonNullIndex $field }}), w)
		if err != nil {
			return err
		}

		{{ if $.IsSimpleNullUnionOfPrimitive $field -}}
			err = vm.Write{{ $.SimpleNullUnionItemType $field }}( *r.{{ .GoName }}, w)
		{{ else -}}
			err = write{{ $.SimpleNullUnionItemType $field }}( *r.{{ .GoName }}, w)
		{{ end -}}
	}
	{{ else if  $.IsArrayOfSimpleNullUnion $field -}}
	err = vm.WriteLong(int64(len(r.{{ .GoName }})), w)
	if err != nil {
		return err
	}
	if (len(r.{{ .GoName }}) != 0) {
		for _, e := range r.{{ .GoName }} {
			if e == nil {
				return nil
			}
			err = vm.Write{{ $.ArraySimpleNullUnionItemType $field }}(*e, w)
			if err != nil {
				return err
			}
		}
		err = vm.WriteLong(0, w)
	}
	{{ else -}}
	err = {{ .Type.SerializerMethod }}( r.{{ .GoName }}, w)
	if err != nil {
		return err
	}
	{{ end -}}
	{{ end -}}
	return err
}

func (r {{ .GoType }}) Serialize(w io.Writer) error {
	return {{ .SerializerMethod }}(r, w)
}

func (r {{ .GoType }}) Schema() string {
	return {{ printf "%q" .Schema }}
}

func (r {{ .GoType }}) SchemaName() string {
	return {{ printf "%q" .AvroName.String }}
}

func (_ {{ .GoType }}) SetBoolean(v bool) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetInt(v int32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetLong(v int64) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetFloat(v float32) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetDouble(v float64) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetBytes(v []byte) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetString(v string) { panic("Unsupported operation") }
func (_ {{ .GoType }}) SetUnionElem(v int64) { panic("Unsupported operation") }

func (r *{{ .GoType }}) Get(i int) types.Field {
	switch (i) {
	{{ range $i, $field := .Fields -}}
	case {{ $i }}:
		{{ if $.ConstructableForField $field | ne "" -}}
			{{ $.ConstructableForField $field }}
		{{ end -}}
		{{ if ne $field.Type.WrapperType "" -}}
			{{ if $.IsSimpleNullUnion $field -}}
			if r.{{ $field.GoName }} == nil {
				var {{ $field.GoName }} = new({{slice .Type.GoType 1 }})
				r.{{ $field.GoName }} = {{ $field.GoName }}
			}
			{{ if $.IsSimpleNullUnionOfPrimitive $field -}}
			w := {{ $field.Type.WrapperType }}{Target: r.{{ $field.GoName }}}
			{{ else -}}
			w := r.{{ $field.GoName }}
			{{ end -}}
			{{ else if  $.IsArrayOfSimpleNullUnion $field -}}
			w := types.ArrayOfNullable{{ $.ArraySimpleNullUnionItemType $field }}Union{Target: &r.{{ $field.GoName }}}
			{{ else -}}
			w := {{ $field.Type.WrapperType }}{Target: &r.{{ $field.GoName }}}
			{{ end -}}
			{{ if $field.Type.WrapperPointer }}
			return &w
			{{ else }}
			return w
			{{ end }}
		{{ else -}}
			return r.{{ $field.GoName }}
		{{ end -}}
	{{ end -}}
	}
	panic("Unknown field index")
}

func (r *{{ .GoType }}) SetDefault(i int) {
	switch (i) {
	{{ range $i, $field := .Fields -}}
        {{ if .HasDefault -}}
	case {{ $i }}:
       	 	{{ $.DefaultForField $field }}
		return
	{{ end -}}
	{{ end -}}
	}
	panic("Unknown field index")
}

func (r *{{ .GoType }}) NullField(i int) { 
	switch (i) {
	{{ range $i, $field := .Fields -}}
        {{ if isNullable $field.Type -}}
	case {{ $i }}:
		r.{{ $field.GoName }} = nil
		return
	{{ end -}}
	{{ end -}}
	}
	panic("Not a nullable field index")
}

func (_ {{ .GoType }}) AppendMap(key string) types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) AppendArray() types.Field { panic("Unsupported operation") }
func (_ {{ .GoType }}) HintSize(int) { panic("Unsupported operation") }
func (_ {{ .GoType }}) Finalize() { }


func (_ {{ .GoType}}) AvroCRC64Fingerprint() []byte {
  return []byte({{ .Name }}AvroCRC64Fingerprint)
}

func (r {{ .GoType }}) MarshalJSON() ([]byte, error) {
	{{ if .Fields -}}
	var err error
	{{ end -}}
	output := make(map[string]json.RawMessage)
	{{ range $i, $field := .Fields -}}
	{{ if $.IsSimpleNullUnion $field -}}
		if r.{{ $field.GoName }} == nil {
			output[{{ printf "%q" $field.Name }}], err = []byte("null"), nil
		} else {
			output[{{ printf "%q" $field.Name }}], err = json.Marshal(map[string]interface{}{
				{{ if $.IsSimpleNullUnionOfPrimitive $field -}}
					"{{slice .Type.GoType 1 }}": *r.{{ .GoName }},
				{{ else -}}
					"{{ $.SimpleNullUnionKey $field }}": r.{{ .GoName }},
				{{ end -}}
			})
		}
    {{ else if  $.IsArrayOfSimpleNullUnion $field -}}
		if r.{{ $field.GoName }} != nil {
			y := make([]*map[string]{{ slice .Type.GoType 3 }}, len(r.{{ $field.GoName }}))
			for i, e := range r.{{ $field.GoName }} {
			if e == nil {
					y[i] = nil
				} else {
					tmp := map[string]{{ slice .Type.GoType 3 }}{"{{ $.ArraySimpleNullUnionNonNullUnionKey $field }}" : *e}
					y[i] = &tmp
				}
			}
			output[{{ printf "%q" $field.Name }}], err = json.Marshal(y)
		}
	{{ else if  $.IsMapOfSimpleNullUnion $field -}}
		if r.{{ $field.GoName }} != nil {
			y := make(map[string]*map[string]{{ slice .Type.GoType 12 }}, len(r.{{ $field.GoName }}))
			for k, v := range r.IntField {
				if v == nil {
					y[k] = nil
				} else {
					tmp := map[string]{{ slice .Type.GoType 12 }}{"{{ $.MapSimpleNullUnionNonNullUnionKey $field }}" : *v}
					y[k] = &tmp
				}
			}
			output["IntField"], err = json.Marshal(y)
		}
	{{ else -}}
	output[{{ printf "%q" $field.Name }}], err = json.Marshal(r.{{ $field.GoName}})
    {{ end -}}
	if err != nil {
		return nil, err
	}
	{{ end -}}
	return json.Marshal(output)	
}

{{ range $i, $field := .Fields -}}
{{ if  $.IsArrayOfSimpleNullUnion $field -}}
func (r *{{ $.GoType }}) Unmarshal{{ .Name -}}JSON(data []byte) (error) {
	y := make([]*map[string]{{ slice .Type.GoType 3 }}, 0)

	if err := json.Unmarshal(data, &y); err != nil {
		return nil
	}

	length := len(y)
	r.{{ .Name -}} = make([]{{ slice .Type.GoType 2 }}, length)

	for i, e := range y {
		if (e == nil) {
			r.{{ .Name -}}[i] = nil
		} else {
			tmp := (*e)["{{ $.ArraySimpleNullUnionNonNullUnionKey $field }}"]
			r.{{ .Name -}}[i] = &tmp
		}
	}

	return nil
}
{{ else if $.IsSimpleNullUnion $field -}}
func (r *{{ $.GoType }}) Unmarshal{{ .Name -}}JSON(data []byte) (error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	if len(fields) > 1 {
		return fmt.Errorf("more than one type supplied for union")
	}

	if v, ok := fields["{{ $.SimpleNullUnionKey $field }}"]; ok {
		r.{{ .GoName -}} = new({{slice .Type.GoType 1 }})
		json.Unmarshal(v, r.{{ .GoName -}})
	}

	return nil
}
{{ else if $.IsMapOfSimpleNullUnion $field -}}
func (r *{{ $.GoType }}) Unmarshal{{ .Name -}}JSON(data []byte) (error) {
  y := make(map[string]*{{ .Type.GoType }}, 0)

	if err := json.Unmarshal(data, &y); err != nil {
		return nil
	}

	length := len(y)
	r.{{ .Name -}} = make({{ .Type.GoType }}, length)

	for i, e := range y {
		if (e == nil) {
			r.{{ .Name -}}[i] = nil
		} else {
			tmp := (*e)["{{ $.MapSimpleNullUnionNonNullUnionKey $field }}"]
			r.{{ .Name -}}[i] = tmp
		}
	}

	return nil
}
{{ end -}}
{{ end }}

func (r *{{ .GoType }}) UnmarshalJSON(data []byte) (error) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	{{ if .Fields -}}
	
	var val json.RawMessage
	{{ end -}}
	{{ range $i, $field := .Fields -}}
		val = func() json.RawMessage {
			if v, ok := fields[{{ printf "%q" $field.Name }}]; ok {
				return v
			}
			{{ range $j, $alias := $field.Aliases -}}
			if v, ok := fields[{{ printf "%q" $alias }}]; ok {
				return v
			}
			{{ end -}}
			return nil
		}()

		if val != nil {
			{{ if $.HasInlinedCustomUnmarshalMethod $field -}}
			if err := r.Unmarshal{{ .Name -}}JSON(val); err != nil {
				return err
			}
			{{ else -}}
			if err := json.Unmarshal([]byte(val), &r.{{ .GoName}}); err != nil {
				return err
			}
			{{ end -}}
		} else {
				{{ if .HasDefault -}}
			{{ if $.ConstructableForField $field | ne "" -}}
			{{ $.ConstructableForField $field }}
			{{ end -}}
				{{ $.DefaultForField $field }}
			{{ else -}}
			return fmt.Errorf("no value specified for {{ $field.Name }}")
			{{ end -}}
		}
		{{ end -}}
	return nil
}
`
