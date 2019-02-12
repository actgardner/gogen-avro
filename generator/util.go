package generator

import (
	"strings"
	"unicode"
)

// ToPublicName returns a go-idiomatic public name. The Avro spec specifies names must start with [A-Za-z_] and contain [A-Za-z0-9_].
// The golang spec says valid identifiers start with [A-Za-z_] and contain [A-Za-z0-9], but the first character must be [A-Z] for the field to be public.
// Additional treatment is applied to the name, in order to remove possible "." from it.
func ToPublicName(name string) string {
	return strings.Title(strings.Trim(strings.Replace(name, ".", "_", -1), "_"))
}

// ToSnake makes filenames snake-case, taken from https://gist.github.com/elwinar/14e1e897fdbe4d3432e1
func ToSnake(in string) string {
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
