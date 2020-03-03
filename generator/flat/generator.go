package flat

import (
	"github.com/actgardner/gogen-avro/generator"
	"github.com/actgardner/gogen-avro/generator/flat/templates"
	avro "github.com/actgardner/gogen-avro/schema"
)

// FlatPackageGenerator emits a file per generated type, all in a single Go package without handling namespacing
type FlatPackageGenerator struct {
	containers bool
	files      *generator.Package
}

func NewFlatPackageGenerator(files *generator.Package, containers bool) *FlatPackageGenerator {
	return &FlatPackageGenerator{
		containers: containers,
		files:      files,
	}
}

type namedDefinition interface {
	Name() string
}

func (f *FlatPackageGenerator) Add(def namedDefinition) error {
	contents, err := templates.Template(def)
	if err == nil {
		// If there's a template for this definition, add it to the package
		filename := generator.ToSnake(def.Name()) + ".go"
		f.files.AddFile(filename, contents)
	} else {
		if err != templates.NoTemplateForType {
			return err
		}
	}

	if r, ok := def.(*avro.RecordDefinition); ok && f.containers {
		if err := f.addRecordContainer(r); err != nil {
			return err
		}
	}

	if ct, ok := def.(avro.CompositeType); ok {
		for _, child := range ct.Children() {
			// Avoid references
			if _, ok := child.(*avro.Reference); !ok {
				if err := f.Add(child); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

func (f *FlatPackageGenerator) addRecordContainer(def *avro.RecordDefinition) error {
	containerFilename := generator.ToSnake(def.Name()) + "_container.go"
	file, err := templates.Evaluate(templates.RecordContainerTemplate, def)
	if err != nil {
		return err
	}
	f.files.AddFile(containerFilename, file)
	return nil
}
