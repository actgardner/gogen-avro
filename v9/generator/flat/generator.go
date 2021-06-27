package flat

import (
	"github.com/actgardner/gogen-avro/v9/generator"
	"github.com/actgardner/gogen-avro/v9/generator/flat/templates"
	avro "github.com/actgardner/gogen-avro/v9/schema"
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

func (f *FlatPackageGenerator) Add(def avro.Node) error {
	file, err := templates.Template(def)
	if err == nil {
		// If there's a template for this definition, add it to the package
		filename := generator.ToSnake(def.Name()) + ".go"
		f.files.AddFile(filename, file)
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

	for _, child := range def.Children() {
		if err := f.Add(child); err != nil {
			return err
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
