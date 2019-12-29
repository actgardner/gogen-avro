package resolver

import (
	"fmt"

	avro "github.com/actgardner/gogen-avro/schema"
)

// ResolveTypes ensures that all avro.References in a tree have an avro.Definition associated with them
func ResolveTypes(rootType avro.AvroType, defs map[avro.QualifiedName]avro.Definition) error {
	return resolveReferences(rootType, defs)
}

// ResolveDefinitions resolves the References in each Definition one level deep
func ResolveDefinitions(defs map[avro.QualifiedName]avro.Definition) error {
	for _, def := range defs {
		for _, child := range def.Children() {
			err := resolveReferences(child, defs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func resolveReferences(t avro.AvroType, defs map[avro.QualifiedName]avro.Definition) error {
	// If we reach a Reference to a Definition, resolve it and stop recursing
	if ref, ok := t.(*avro.Reference); ok {
		ref.Def, ok = defs[ref.TypeName]
		if !ok {
			return fmt.Errorf("Unable to resolve type reference %v", ref.TypeName)
		}
		return nil
	}

	// Otherwise recursively search for further References
	for _, child := range t.Children() {
		err := resolveReferences(child, defs)
		if err != nil {
			return err
		}
	}

	return nil
}
