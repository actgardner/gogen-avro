package generator

import (
	"fmt"
)

type StructField struct {
	Name string
	Type string
}

type (s *StructField) String() {
	return fmt.Sprintf("%v %v", s.Name, s.Type)
}

type Struct struct {
	File string
	Name string
	Imports []string
	Fields []*StructField

	Dependencies []Block
}

func (s *Struct) String() string {
	fieldList := ""
	for _, field := range(s.Fields) {
		fieldList += fmt.Sprintf("%v %v\n", field.Name, field.Type)
	}
	return fmt.Sprintf("type %v struct {\n%v}", s.Name, fieldList)
}

func (s *Struct) AddToPackage(p *Package) {
	f, ok := p.files[s.File]
	if !ok {
		f = NewFile(s.File)
		p.files[s.File] = f
	}
	f.structs[s.Name] = s

	for _, i := range(s.Imports) {
		p.AddImport(s.File, i)
	}
}
