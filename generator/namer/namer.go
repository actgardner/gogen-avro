package namer

import (
	"fmt"
	"strings"

	avro "github.com/actgardner/gogen-avro/schema"
)

type Namer struct {
	// typeNamer generates the names of Go types
	typeNamer NameFormatter

	// fieldNamer generates names for struct fields
	fieldNamer NameFormatter
}

func NewNamer(typeNamer NameFormatter, fieldNamer NameFormatter) *Namer {
	return &Namer{
		typeNamer:  typeNamer,
		fieldNamer: fieldNamer,
	}
}

// Apply generates metadata about Go struct names, file paths and method names for generated code
func (n *Namer) Apply(node avro.Node) error {
	// Because the names of parent types depend on child types, names are generated depth-first
	for _, child := range node.Children() {
		err := n.Apply(child)
		if err != nil {
			return err
		}
	}

	switch t := node.(type) {
	// Traverse References if the underlying definition hasn't been resolved yet
	case *avro.Reference:
		t.SetGeneratorMetadata(MetadataKey, &TypeMetadata{})
		if !t.HasGeneratorMetadata(MetadataKey) {
			n.Apply(t)
		}

	// User-defined types
	case *avro.EnumDefinition:
		metadata := &EnumMetadata{}
		metadata.Name = n.typeNamer.Format(t.AvroName().Name)
		metadata.GoType = metadata.Name
		metadata.SerializerMethod = fmt.Sprintf("write%v", metadata.Name)
		metadata.ConstructorMethod = fmt.Sprintf("make(%v, 0)", metadata.GoType)
		metadata.WrapperType = "types.Int"
		metadata.SerializerMethod = "write" + metadata.GoType
		metadata.FromStringMethod = "New" + metadata.GoType + "Value"
		metadata.SymbolNames = make([]string, len(t.Symbols()))
		for i, s := range t.Symbols() {
			metadata.SymbolNames[i] = metadata.GoType + strings.Title(s)
		}
		t.SetGeneratorMetadata(MetadataKey, metadata)

	case *avro.FixedDefinition:
		metadata := &FixedMetadata{}
		metadata.Name = n.typeNamer.Format(t.AvroName().Name)
		metadata.GoType = fmt.Sprintf("*%v", metadata.Name)
		metadata.SerializerMethod = fmt.Sprintf("write%v", metadata.Name)
		metadata.ConstructorMethod = fmt.Sprintf("make(%v, 0)", metadata.GoType)
		metadata.WrapperType = fmt.Sprintf("%vWrapper", metadata.Name)
		metadata.SerializerMethod = "write" + metadata.GoType
		metadata.WrapperType = fmt.Sprintf("%vWrapper", metadata.GoType)
		t.SetGeneratorMetadata(MetadataKey, metadata)

	case *avro.RecordDefinition:
		metadata := &RecordMetadata{}
		metadata.Name = n.typeNamer.Format(t.AvroName().Name)
		metadata.GoType = "*" + metadata.Name
		metadata.SerializerMethod = fmt.Sprintf("write%v", metadata.Name)
		metadata.ConstructorMethod = fmt.Sprintf("New%v()", metadata.GoType)
		metadata.WrapperType = ""
		metadata.SerializerMethod = fmt.Sprintf("write%v", metadata.Name)
		metadata.NewWriterMethod = fmt.Sprintf("New%vWriter", metadata.Name)
		metadata.ConstructorMethod = fmt.Sprintf("New%v()", metadata.Name)
		metadata.RecordReaderTypeName = metadata.Name + "Reader"

		for _, f := range t.Fields() {
			f.SetGeneratorMetadata(MetadataKey, &FieldMetadata{n.fieldNamer.Format(f.Name())})
		}
		t.SetGeneratorMetadata(MetadataKey, metadata)

	// Complex types
	case *avro.ArrayField:
		metadata := &TypeMetadata{}
		itemMetadata := t.ItemType().GetGeneratorMetadata(MetadataKey).(*TypeMetadata)
		metadata.Name = "Array" + itemMetadata.Name
		metadata.GoType = fmt.Sprintf("[]%v", itemMetadata.GoType)
		metadata.SerializerMethod = fmt.Sprintf("write%v", metadata.Name)
		metadata.ConstructorMethod = fmt.Sprintf("make(%v, 0)", metadata.GoType)
		metadata.WrapperType = fmt.Sprintf("%vWrapper", metadata.Name)
		t.SetGeneratorMetadata(MetadataKey, metadata)

	case *avro.MapField:
		metadata := &TypeMetadata{}
		itemMetadata := t.ItemType().GetGeneratorMetadata(MetadataKey).(*TypeMetadata)
		metadata.Name = "Map" + itemMetadata.Name
		metadata.GoType = fmt.Sprintf("*%v", itemMetadata.GoType)
		metadata.SerializerMethod = fmt.Sprintf("write%v", metadata.Name)
		metadata.ConstructorMethod = fmt.Sprintf("New%v()", metadata.Name)
		metadata.WrapperType = ""
		t.SetGeneratorMetadata(MetadataKey, metadata)

	case *avro.UnionField:
		metadata := &TypeMetadata{}
		itemMetadata := make([]*TypeMetadata, len(t.AvroTypes()))
		for i, at := range t.AvroTypes() {
			fmt.Printf("%T\n", at)
			itemMetadata[i] = at.GetGeneratorMetadata(MetadataKey).(*TypeMetadata)
		}
		compositeUnionName := "Union"
		for _, item := range itemMetadata {
			compositeUnionName += item.Name
		}
		metadata.Name = n.typeNamer.Format(compositeUnionName)
		metadata.GoType = "*" + metadata.Name
		metadata.SerializerMethod = fmt.Sprintf("write%v", metadata.Name)
		metadata.ConstructorMethod = fmt.Sprintf("New%v()", metadata.Name)
		metadata.WrapperType = ""
		t.SetGeneratorMetadata(MetadataKey, metadata)

	// Primitive types
	case *avro.BoolField:
		metadata := &TypeMetadata{}
		metadata.GoType = "bool"
		metadata.SerializerMethod = "vm.WriteBool"
		metadata.Name = "Bool"
		metadata.WrapperType = "types.Boolean"
		t.SetGeneratorMetadata(MetadataKey, metadata)

	case *avro.BytesField:
		metadata := &TypeMetadata{}
		metadata.GoType = "[]byte"
		metadata.SerializerMethod = "vm.WriteBytes"
		metadata.Name = "Bytes"
		metadata.WrapperType = "types.Bytes"
		t.SetGeneratorMetadata(MetadataKey, metadata)

	case *avro.StringField:
		metadata := &TypeMetadata{}
		metadata.GoType = "string"
		metadata.SerializerMethod = "vm.WriteString"
		metadata.Name = "String"
		metadata.WrapperType = "types.String"
		t.SetGeneratorMetadata(MetadataKey, metadata)

	case *avro.FloatField:
		metadata := &TypeMetadata{}
		metadata.GoType = "float32"
		metadata.SerializerMethod = "vm.WriteFloat"
		metadata.Name = "Float"
		metadata.WrapperType = "types.Float"
		t.SetGeneratorMetadata(MetadataKey, metadata)

	case *avro.DoubleField:
		metadata := &TypeMetadata{}
		metadata.GoType = "float64"
		metadata.SerializerMethod = "vm.WriteDouble"
		metadata.Name = "Double"
		metadata.WrapperType = "types.Double"
		t.SetGeneratorMetadata(MetadataKey, metadata)

	case *avro.IntField:
		metadata := &TypeMetadata{}
		metadata.GoType = "int32"
		metadata.SerializerMethod = "vm.WriteInt"
		metadata.Name = "Int"
		metadata.WrapperType = "types.Int"
		t.SetGeneratorMetadata(MetadataKey, metadata)

	case *avro.LongField:
		metadata := &TypeMetadata{}
		metadata.GoType = "int64"
		metadata.SerializerMethod = "vm.WriteLong"
		metadata.Name = "Long"
		metadata.WrapperType = "types.Long"
		t.SetGeneratorMetadata(MetadataKey, metadata)

	case *avro.NullField:
		metadata := &TypeMetadata{}
		metadata.GoType = "*types.NullVal"
		metadata.SerializerMethod = "vm.WriteNull"
		metadata.Name = "Null"
		metadata.WrapperType = ""
		t.SetGeneratorMetadata(MetadataKey, metadata)
	}
	return nil
}
