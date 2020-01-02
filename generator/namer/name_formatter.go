package namer

import (
	"regexp"
	"strings"
	"unicode"
)

var invalidTokensExpr = regexp.MustCompile(`[._\s]+`)

// NameFormatter is the interface defining a function for converting
// a name to a go-idiomatic name.
type NameFormatter interface {
	Format(name string) string
}

// DefaultNameFormatter implements the NameFormatter interface with the
// backwards-compatible public name generator function.
type DefaultNameFormatter struct {
}

func NewDefaultNameFormatter() *DefaultNameFormatter {
	return &DefaultNameFormatter{}
}

// NamespaceNameFormatter is like DefaultNameFormatter but taking into account
// special tokens so namespaced names can be generated safely.
type NamespaceNameFormatter struct {
	shortNames bool
}

// NewNamespaceNameFormatter returns a namespace-aware namer.
func NewNamespaceNameFormatter(shortNames bool) *NamespaceNameFormatter {
	return &NamespaceNameFormatter{shortNames: shortNames}
}

// ToPublicName implements the backwards-compatible name converter in
// DefaultNameFormatter.
func (d *DefaultNameFormatter) Format(name string) string {
	lastIndex := strings.LastIndex(name, ".")
	name = name[lastIndex+1:]
	return strings.Title(strings.Trim(name, "_"))
}

// ToPublicName implements the go-idiomatic public name as in DefaultNameFormatter's
// struct, but with additional treatment applied in order to remove possible
// invalid tokens from it. Final string is then converted to camel-case.
func (n *NamespaceNameFormatter) Format(name string) string {
	if n.shortNames {
		if parts := strings.Split(name, "."); len(parts) > 2 {
			name = strings.Join(parts[len(parts)-2:], ".")
		}
	}
	name = invalidTokensExpr.ReplaceAllString(name, " ")
	return strings.Replace(strings.Title(name), " ", "", -1)
}

// SnakeNameFormatter returns names in snake_case, primarily for idiomatic filenames
type SnakeNameFormatter struct{}

func NewSnakeNameFormatter() *SnakeNameFormatter {
	return &SnakeNameFormatter{}
}

func (_ *SnakeNameFormatter) Format(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
