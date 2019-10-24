package generator

import (
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
)

const (
	CaseTitle = "title"
	CaseCamel = "camel"

	invalidTokensExpr = `[._\s]+`
)

// Namer is the interface defining a function for converting
// a name to a go-idiomatic public name.
type Namer interface {
	// ToPublicName returns a go-idiomatic public name. The Avro spec
	// specifies names must start with [A-Za-z_] and contain [A-Za-z0-9_].
	// The golang spec says valid identifiers start with [A-Za-z_] and contain
	// [A-Za-z0-9], but the first character must be [A-Z] for the field to be
	// public.
	ToPublicName(name string) string
}

// DefaultNamer implements the Namer interface with the
// backwards-compatible public name generator function.
type DefaultNamer struct {
	nameCase string
}

// NamespaceNamer is like DefaultNamer but taking into account
// special tokens so namespaced names can be generated safely.
type NamespaceNamer struct {
	shortNames bool
	nameCase   string
	re         *regexp.Regexp
}

var (
	namer Namer = &DefaultNamer{nameCase: "title"}
)

// NewDefaultNamer returns a default namer.
func NewDefaultNamer(nameCase string) *DefaultNamer {
	return &DefaultNamer{nameCase: nameCase}
}

// NewNamespaceNamer returns a namespace-aware namer.
func NewNamespaceNamer(shortNames bool, nameCase string) *NamespaceNamer {
	return &NamespaceNamer{shortNames: shortNames, nameCase: nameCase, re: regexp.MustCompile(invalidTokensExpr)}
}

// SetNamer sets the generator's global namer
func SetNamer(n Namer) {
	namer = n
}

// ToPublicName implements the backwards-compatible name converter in
// DefaultNamer.
// Returns a go-idiomatic public name. The Avro spec
// specifies names must start with [A-Za-z_] and contain [A-Za-z0-9_].
// The golang spec says valid identifiers start with [A-Za-z_] and contain
// [A-Za-z0-9], but the first character must be [A-Z] for the field to be
// public.
func (d *DefaultNamer) ToPublicName(name string) string {
	lastIndex := strings.LastIndex(name, ".")
	name = strings.Trim(name[lastIndex+1:], "_")
	switch d.nameCase {
	case "title":
		name = strings.Title(name)
	case "camel":
		name = strcase.ToCamel(name)
	}
	return name
}

// ToPublicName implements the go-idiomatic public name as in DefaultNamer's
// struct, but with additional treatment applied in order to remove possible
// invalid tokens from it. Final string is then converted to camel-case.
func (n *NamespaceNamer) ToPublicName(name string) string {
	if n.shortNames {
		if parts := strings.Split(name, "."); len(parts) > 2 {
			name = strings.Join(parts[len(parts)-2:], ".")
		}
	}
	name = n.re.ReplaceAllString(name, " ")

	switch n.nameCase {
	case "title":
		name = strings.Replace(strings.Title(name), " ", "", -1)
	case "camel":
		name = strcase.ToCamel(name)
	}
	return name
}
