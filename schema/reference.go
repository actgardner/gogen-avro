package schema

import "fmt"

/*
  A named Reference to a user-defined type (fixed, enum, record). Just a wrapper with a name around a qnamed type.
  If a reference is created without a defined type, all parsing-time values are returned as invalid values.
  Once the reference gets assigned a type, it triggers its registered resolvers' Resolve(ref) method.
  This is commonly transparent for almost all types, but some other types depend on their children fields for
  defining their internal data, like Name and GoType.
  Examples of this are:
	- Arrays of base types defined as references: the resolver renames the array to its final name.
	- Union types with named-types' members: once all its children refs are triggered, the union gets renamed.
*/
type Reference struct {
	refName   QualifiedName
	refType   QualifiedAvroType
	resolvers []referenceResolver
}

// referenceResolver is an interface with a function that triggers once this
// reference gets informed about the type it refers to, so no afterwards resolving
// phase is required.
type referenceResolver interface {
	Resolve(ref Reference)
}

const unresolvedFmtString = "unresolved-ref-%s-%s"

var (
	// Ensure interface implementations
	_ QualifiedAvroType = &Reference{}
	_ CompositeType     = &Reference{}
)

func NewReference(name QualifiedName, t QualifiedAvroType) *Reference {
	return &Reference{refName: name, refType: t}
}

func (r Reference) RefName() QualifiedName {
	return r.refName
}

func (r Reference) RefType() QualifiedAvroType {
	return r.refType
}

func (r Reference) IsUntyped() bool {
	return r.refType == nil
}

func (r *Reference) SetType(t QualifiedAvroType) {
	if r.refType != nil {
		panic("Cannot reassign reference type")
	}
	r.refType = t

	for _, resolver := range r.resolvers {
		resolver.Resolve(*r)
	}
}

func (r *Reference) AddResolver(resolver referenceResolver) {
	if resolver == nil {
		panic("Cannot add a nil resolver")
	}
	r.resolvers = append(r.resolvers, resolver)
}

func (r Reference) Name() string {
	if r.refType == nil {
		return fmt.Sprintf(unresolvedFmtString, "name", r.refName)
	}
	return r.refType.Name()
}

func (r Reference) SimpleName() string {
	if r.refType == nil {
		return fmt.Sprintf(unresolvedFmtString, "simplename", r.refName)
	}
	return r.refType.SimpleName()
}

func (r Reference) GoType() string {
	if r.refType == nil {
		return fmt.Sprintf(unresolvedFmtString, "gotype", r.refName)
	}
	return r.refType.GoType()
}

func (r Reference) IsOptional() bool {
	return r.refType.IsOptional()
}

// The name of the method which writes this field onto the wire
func (r Reference) SerializerMethod() string {
	if r.refType == nil {
		return fmt.Sprintf(unresolvedFmtString, "serializermethod", r.refName)
	}
	return r.refType.SerializerMethod()
}

func (r Reference) Definition(scope map[QualifiedName]interface{}) (interface{}, error) {
	return r.refType.Definition(scope)
}

func (r Reference) QualifiedName() QualifiedName {
	return r.refType.QualifiedName()
}

func (r Reference) Aliases() []QualifiedName {
	return r.refType.Aliases()
}

func (r Reference) Children() []AvroType {
	if ct, ok := r.refType.(CompositeType); ok {
		return ct.Children()
	}
	// References can only point to Definitions and thus have no children
	// unless they're pointing to composite types
	return []AvroType{}
}

func (r Reference) DefaultValue(lvalue string, rvalue interface{}) (string, error) {
	return r.refType.DefaultValue(lvalue, rvalue)
}

func (r Reference) WrapperType() string {
	if r.refType == nil {
		return fmt.Sprintf(unresolvedFmtString, "wrappertype", r.refName)
	}
	return r.refType.WrapperType()
}

func (r Reference) IsReadableBy(f AvroType, visited map[QualifiedName]interface{}) bool {
	if ref, ok := f.(*Reference); ok {
		f = ref.refType
	}
	return r.refType.IsReadableBy(f, visited)
}
