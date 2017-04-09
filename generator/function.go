package generator

import (
	"fmt"
)

type FunctionName struct {
	// The target struct type, if applicable
	Struct string
	// The function name
	Name string
}

type Function struct {
	File string
	Name *FunctionName
	Arguments []*StructField
	ReturnTypes []string
	Body string

	Dependencies []Block
}

func (f *Function) Signature() string {
	argList := make([]string, 0)
	for _, arg := range f.Arguments {
		argList = append(argList, arg.String())
	}

	args := strings.Join(argList, ", ")
	returns := strings.Join(f.ReturnTypes, ", ")

	return fmt.Sprintf("%v (%v) (%v)", f.Name.Name, args, returns)
}

func (f *Function) String() string {
	argList := make([]string, 0)
	for _, arg := range f.Arguments {
		argList = append(argList, arg.String())
	}

	args := strings.Join(argList, ", ")
	returns := strings.Join(f.ReturnTypes, ", ")

	if f.Name.Struct == "" {
		return fmt.Sprintf("func (s *%v) %v {\n%v}", f.Name.Struct, f.Signature(), f.Body)
	}

	return fmt.Sprintf("func %v {\n%v}", f.Signature(), f.Body)
}

func (f *Function) AddToPackage(p *Package) {
	p.getOrAddFile(f.File).functions[f.Name] = f.Body
	p.AddBlocks(f.Dependencies)
}
