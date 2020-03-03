package schema

import (
	"fmt"
	"strings"
)

type UnionField struct {
	qualifiedField
	children       []AvroType
	generatedTypes []AvroType
	optTypeIdx     int
	unresolvedRefs int
	namingFunc     func()
}

func NewUnionField(name string, children []AvroType, definition []interface{}) *UnionField {
	u := &UnionField{
		children:   children,
		optTypeIdx: -1,
	}
	u.definition = definition
	u.namingFunc = func() {
		// Use constructor base name for name regeneration func
		u.regenerateName(name)
	}

	// Process fields. Processing implies detecting the optional field, if any,
	// as well as registering reference resolvers, if any. If there are no reference
	// fields, then this union can be named straight ahead. Otherwise, the naming occurs
	// after all references have been resolved.
	for i, f := range u.children {
		switch t := f.(type) {
		case *NullField:
			// Optional field detected in position 'i'
			u.optTypeIdx = i
		case *Reference:
			// Add reference resolver for this union's field
			t.AddResolver(u)
			u.unresolvedRefs++
		}
	}

	// Cache the generated types slice, optionally removing the null type elem.
	if u.optTypeIdx == -1 {
		u.generatedTypes = u.children
	} else {
		u.setOptional()
		u.generatedTypes = make([]AvroType, len(u.children)-1)
		i := 0
		for j, child := range u.children {
			if j != u.optTypeIdx {
				u.generatedTypes[i] = child
				i++
			}
		}
	}

	// No references, name the union just now
	if u.unresolvedRefs == 0 {
		u.namingFunc()
	}
	return u
}

// Resolve runtime data when item type is a reference, since its data
// cannot be always known at this array's creation time.
func (u *UnionField) Resolve(ref Reference) {
	u.unresolvedRefs--
	switch {
	case u.unresolvedRefs < 0:
		panic("Union resolved too much times")
	case u.unresolvedRefs == 0:
		u.namingFunc()
	}
}

func (u *UnionField) regenerateName(basename string) {
	// If no base name is given, generate one from its fields
	if basename == "" {
		var str strings.Builder
		str.WriteString("Union")
		for _, f := range u.children {
			str.WriteString(f.Name())
		}
		basename = str.String()
	}
	u.setQualifiedName(QualifiedName{Name: basename})
}

// Create the generated types cache slice for the template,
// removing the null type elem from it, if there's one.
func (u *UnionField) cacheGeneratedTypes() {
	if u.optTypeIdx == -1 {
		u.generatedTypes = u.children
		return
	}
	u.generatedTypes = make([]AvroType, len(u.children)-1)
	i := 0
	for j, child := range u.children {
		if j != u.optTypeIdx {
			u.generatedTypes[i] = child
			i++
		}
	}
}

func (s *UnionField) Children() []AvroType {
	return s.children
}

func (s *UnionField) UnionEnumType() string {
	return fmt.Sprintf("%vTypeEnum", s.Name())
}

func (s *UnionField) ItemName(item AvroType) string {
	return s.UnionEnumType() + item.Name()
}

func (s *UnionField) ItemTypes() []AvroType {
	return s.generatedTypes
}

func (s *UnionField) OptionalIndex() int {
	return s.optTypeIdx
}

func (s *UnionField) ItemConstructor(f AvroType) string {
	if constructor, ok := getConstructableForType(f); ok {
		return constructor.ConstructorMethod()
	}
	return ""
}

func (s *UnionField) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	def := make([]interface{}, len(s.children))
	var err error
	for i, item := range s.children {
		def[i], err = item.Definition(scope)
		if err != nil {
			return nil, err
		}
	}
	return def, nil
}

func (s *UnionField) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	defaultType := s.children[0]
	init := fmt.Sprintf("%v = %v\n", lvalue, s.ConstructorMethod())
	lvalue = fmt.Sprintf("%v.%v", lvalue, defaultType.Name())
	constructorCall := ""
	if constructor, ok := getConstructableForType(defaultType); ok {
		constructorCall = fmt.Sprintf("%v = %v\n", lvalue, constructor.ConstructorMethod())
	}
	assignment, err := defaultType.DefaultValue(lvalue, rvalue)
	return init + constructorCall + assignment, err
}

func (s *UnionField) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
	// Report if *any* writer type could be deserialized by the reader
	for _, t := range s.children {
		if t.IsReadableBy(f, visited) {
			return true
		}
	}
	return false
}

func (s *UnionField) ConstructorMethod() string {
	if s.IsOptional() {
		return fmt.Sprintf("&%s{}", s.Name())
	}
	return fmt.Sprintf("%s{}", s.Name())
}

func (s *UnionField) Equals(reader *UnionField) bool {
	if len(reader.children) != len(s.children) {
		return false
	}

	for i, t := range s.children {
		readerType := reader.children[i]
		if writerRef, ok := t.(QualifiedAvroType); ok {
			if readerRef, ok := readerType.(QualifiedAvroType); ok {
				if readerRef.QualifiedName() != writerRef.QualifiedName() {
					return false
				}
			} else {
				return false
			}
		} else if t != readerType {
			return false
		}
	}
	return true
}
