package generator

import (
	"fmt"
)

type Block interface {
	AddToPackage(*Package)
}

type Import {
	File string
	Package string
}

func (i *Import) AddToPackage(p *Package) {
	p.getOrAddFile(c.File).imports[i.Package] = 1
}

func Constant struct {
	File string
	Name string
	Value interface{}

	Dependencies []Block
}

func (c *Constant) AddToPackage(p *Package) {
	p.getOrAddFile(c.File).constants[name] = value
	p.AddBlocks(c.Dependencies)
}

type Interface struct {
	File string
	Name string
	Functions []*Function
}

func (i *Interface) String() string {
	funcList := ""
	for _, f := range(i.Functions) {
		funcList += fmt.Sprintf("%v %v\n", f.Name, field.Type)
	}
	return fmt.Sprintf("type %v struct {\n%v}", s.Name, fieldList)
}

func (i *Interface) AddToPackage(p *Package) {
	p.getOrAddFile(i.File).interfaces[i.Name] = value
	p.AddBlocks(i.Dependencies)
}

