package schema

type AbstractType interface {
	Name() string
	SimpleName() string
	GoType() string
	Definition(scope map[QualifiedName]interface{}) (interface{}, error)
	IsOptional() bool
}

func copyDefinition(x map[string]interface{}) map[string]interface{} {
	if x == nil {
		return x
	}
	x1 := make(map[string]interface{})
	for name, val := range x {
		x1[name] = val
	}
	return x1
}
