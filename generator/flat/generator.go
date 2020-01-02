package flat

import (
	"github.com/actgardner/gogen-avro/generator"
	"github.com/actgardner/gogen-avro/generator/namer"
	"github.com/actgardner/gogen-avro/generator/templates"
	avro "github.com/actgardner/gogen-avro/schema"
)

type Named interface {
	AvroName() avro.QualifiedName
}

// FlatPackageGenerator emits a file per generated type, all in a single Go package without handling namespacing
type FlatPackageGenerator struct {
	containers bool
	fileNamer  namer.NameFormatter
	files      *generator.Package
}

func NewFlatPackageGenerator(files *generator.Package, fileNamer namer.NameFormatter, containers bool) *FlatPackageGenerator {
	return &FlatPackageGenerator{
		containers: containers,
		fileNamer:  fileNamer,
		files:      files,
	}
}

func (f *FlatPackageGenerator) Add(def avro.Node) error {
	file, err := templates.Template(def)
	if err == nil {
		name := def.(Named).AvroName().Name
		// If there's a template for this definition, add it to the package
		filename := f.fileNamer.Format(name) + ".go"
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
	name := Named(def).AvroName().Name
	containerFilename := f.fileNamer.Format(name + "Container.go")
	file, err := templates.Evaluate(templates.RecordContainerTemplate, def)
	if err != nil {
		return err
	}
	f.files.AddFile(containerFilename, file)
	return nil
}
