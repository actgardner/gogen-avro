package schema

import "fmt"

// Common methods for all types
type abstractField struct {
	name       string
	simpleName string
	goType     string
	optional   bool
	definition interface{}
}

var (
	// Ensure interface implementation
	_ AbstractType = &abstractField{}
)

func (p *abstractField) Name() string {
	return p.name
}

func (p *abstractField) setName(name string) {
	p.name = name
}

func (p *abstractField) SimpleName() string {
	return p.simpleName
}

func (p *abstractField) setSimpleName(simpleName string) {
	p.simpleName = simpleName
}

func (p *abstractField) GoType() string {
	if p.IsOptional() {
		return fmt.Sprintf("*%s", p.goType)
	}
	return p.goType
}

func (p *abstractField) setGoType(goType string) {
	p.goType = goType
}

func (p *abstractField) IsOptional() bool {
	return p.optional
}

func (p *abstractField) setOptional() {
	p.optional = true
}

func (p *abstractField) Definition(_ map[QualifiedName]interface{}) (interface{}, error) {
	return p.definition, nil
}

func (p *abstractField) definitionAsStringMap() map[string]interface{} {
	return p.definition.(map[string]interface{})
}
