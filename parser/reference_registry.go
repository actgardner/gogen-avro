package parser

import (
	"fmt"

	"github.com/actgardner/gogen-avro/schema"
)

// referenceRegistry is a special registry for all qnamed types, which could be possibly referenced
// via aliases or type names elsewhere. The registry not only avoids duplicating references in memory,
// but also triggers the resolution of a reference, if it was created before defining their real type.
// By triggering the resolution, all types whose properties depend on their children or base types can
// be finally set.
type referenceRegistry struct {
	qualifiedTypes map[schema.QualifiedName]*schema.Reference
}

func NewReferenceRegistry() *referenceRegistry {
	return &referenceRegistry{
		qualifiedTypes: make(map[schema.QualifiedName]*schema.Reference),
	}
}

// CreateReference returns a Reference for a given qnamed type. If the reference already exists, its value is returned
// instead of creating a duplicated one. If the type being registered is unknown at registration time, the reference
// gets registered untyped. Once the real, final type is registered, the reference is updated using its SetType method.
// An untyped reference being setted this way triggers its resolution, allowing all its "owner" types to refresh their
// internal data, if needed.
func (reg *referenceRegistry) CreateReference(name schema.QualifiedName, t schema.QualifiedAvroType) (*schema.Reference, error) {
	if t == nil {
		return reg.getOrCreateUntypedReference(name), nil
	}

	ref := reg.getOrCreateUntypedReference(name)
	if !ref.IsUntyped() {
		return nil, fmt.Errorf("Conflicting definitions for %v", name)
	}

	// Trigger type update for reference
	ref.SetType(t)

	// Trigger type update for all of its aliases
	for _, alias := range ref.Aliases() {
		aliasedRef := reg.getOrCreateUntypedReference(alias)
		if !aliasedRef.IsUntyped() {
			return nil, fmt.Errorf("Alias %s from %s is conflicting with definitions for %s", alias, aliasedRef.RefName(), ref.RefName())
		}
		aliasedRef.SetType(t)
	}
	return ref, nil
}

func (reg *referenceRegistry) getOrCreateUntypedReference(name schema.QualifiedName) *schema.Reference {
	ref := reg.qualifiedTypes[name]
	if ref == nil {
		ref = schema.NewReference(name, nil)
		reg.qualifiedTypes[name] = ref
	}
	return ref
}
