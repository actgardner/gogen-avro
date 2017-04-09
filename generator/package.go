package generator

import (
	"path/filepath"
	"sort"
)

// Represents the output package
type Package struct {
	name  string
	files map[string]*File
}

func NewPackage(name string) *Package {
	return &Package{name: name, files: make(map[string]*File)}
}

func (p *Package) WriteFiles(targetDir string) error {
	for _, f := range p.files {
		err := f.WriteFile(p.name, filepath.Join(targetDir, f.name))
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Package) AddBlocks(blocks Block...) {
	for _i, block := range blocks {
		block.AddToPackage(p)
	}
}

func (p *Package) Files() []string {
	files := make([]string, 0)
	for file, _ := range p.files {
		files = append(files, file)
	}
	sort.Strings(files)
	return files
}

func (p *Package) File(name string) (*File, bool) {
	file, ok := p.files[name]
	return file, ok
}

func (p *Package) AddHeader(file, header string) {
	f, ok := p.files[file]
	if !ok {
		f = NewFile(file)
		p.files[file] = f
	}

	f.headers = append(f.headers, header)
}

func (p *Package) getOrAddFile(file string) *File {
	f, ok := p.files[file]
	if !ok {
		f = NewFile(file)
		p.files[file] = f
	}
	return f
}

func (p *Package) HasStruct(file, name string) bool {
	f, ok := p.files[file]
	if !ok {
		return false
	}
	_, ok = f.structs[name]
	return ok
}

func (p *Package) HasFunction(file, str, name string) bool {
	f, ok := p.files[file]
	if !ok {
		return false
	}
	_, ok = f.functions[FunctionName{str, name}]
	return ok
}

func (p *Package) HasImport(file, name string) bool {
	f, ok := p.files[file]
	if !ok {
		return false
	}
	_, ok = f.imports[name]
	return ok
}
