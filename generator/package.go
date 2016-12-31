package generator

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

const UTIL_FILE = "primitive.go"

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
		err := f.writeFile(p.name, filepath.Join(targetDir, f.name))
		if err != nil {
			return err
		}
	}
	return nil
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

func (p *Package) addFunction(file, str, name, def string) {
	f, ok := p.files[file]
	if !ok {
		f = NewFile(file)
		p.files[file] = f
	}
	f.functions[FunctionName{str, name}] = def
}

func (p *Package) addStruct(file, name, def string) {
	f, ok := p.files[file]
	if !ok {
		f = NewFile(file)
		p.files[file] = f
	}
	f.structs[name] = def
}

func (p *Package) addImport(file, name string) {
	f, ok := p.files[file]
	if !ok {
		f = NewFile(file)
		p.files[file] = f
	}
	f.imports[name] = 1
}

func (p *Package) addConstant(file, name string, value interface{}) {
	f, ok := p.files[file]
	if !ok {
		f = NewFile(file)
		p.files[file] = f
	}
	f.constants[name] = value
}

func (p *Package) hasStruct(file, name string) bool {
	f, ok := p.files[file]
	if !ok {
		return false
	}
	_, ok = f.structs[name]
	return ok
}

func (p *Package) hasFunction(file, str, name string) bool {
	f, ok := p.files[file]
	if !ok {
		return false
	}
	_, ok = f.functions[FunctionName{str, name}]
	return ok
}

func (p *Package) hasImport(file, name string) bool {
	f, ok := p.files[file]
	if !ok {
		return false
	}
	_, ok = f.imports[name]
	return ok
}

// Represents a Go source file in the output package
type File struct {
	name      string
	functions map[FunctionName]string
	structs   map[string]string
	imports   map[string]interface{}
	constants map[string]interface{}
}

func NewFile(name string) *File {
	return &File{
		name: name,
		functions: make(map[FunctionName]string),
		structs: make(map[string]string),
		imports: make(map[string]interface{}),
		constants: make(map[string]interface{}),
	}
}

type FunctionName struct {
	// The target struct type, if there is one
	str string
	// The function name
	name string
}

func (f *File) Imports() []string {
	imports := make([]string, 0)
	for i, _ := range f.imports {
		imports = append(imports, i)
	}
	sort.Strings(imports)
	return imports
}

func (f *File) Structs() []string {
	structs := make([]string, 0)
	for s, _ := range f.structs {
		structs = append(structs, s)
	}
	sort.Strings(structs)
	return structs
}

func (f *File) Functions() []FunctionName {
	funcs := make([]FunctionName, 0)
	for f, _ := range f.functions {
		funcs = append(funcs, f)
	}
	sort.Sort(FunctionNameList(funcs))
	return funcs
}

func (f *File) importString() string {
	if len(f.imports) == 0 {
		return ""
	}
	imports := "import (\n"
	for i, _ := range f.imports {
		imports += fmt.Sprintf("%q\n", i)
	}
	imports += ")"
	return imports
}

func (f *File) constantString() string {
	if len(f.constants) == 0 {
		return ""
	}
	constants := "const (\n"
	for name, value := range f.constants {
		// For strings, quote the right-hand side
		if valueString, ok := value.(string); ok {
			constants += fmt.Sprintf("%s = %q\n", name, valueString)
		} else {
			constants += fmt.Sprintf("%s = %s\n", name, value)
		}
	}
	constants += ")"
	return constants
}

func (f *File) structString() string {
	structs := ""
	keys := make([]string, 0)
	for k := range f.structs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		structs += f.structs[k] + "\n"
	}
	return structs
}

func (f *File) functionString() string {
	funcs := ""
	keys := make(FunctionNameList, 0)
	for k := range f.functions {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	for _, k := range keys {
		funcs += f.functions[k] + "\n"
	}
	return funcs
}

/* Write the contents of the file:
- imports (all in one block)
- struct definitions (alphabetically)
- functions (sorted alphabetically by struct to which they're attached, then unattached funcs)
TODO: It'd be better to group funcs attached to a struct with the struct definition
*/
func (f *File) writeFile(pkgName, targetFile string) error {
	src := fmt.Sprintf("package %v\n%v\n%v\n%v\n%v\n", pkgName, f.importString(), f.constantString(), f.structString(), f.functionString())
	fileContent, err := format.Source([]byte(src))
	if err != nil {
		return fmt.Errorf("Error formatting file %v - %v\n\nContents: %v", f.name, err, src)
	}
	err = ioutil.WriteFile(targetFile, fileContent, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Error writing file %v - %v", f.name, err)
	}
	return nil
}

/* Implement the Sortable interface for FunctionNames */
type FunctionNameList []FunctionName

func (f FunctionNameList) Len() int {
	return len(f)
}

func (f FunctionNameList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

/* Sort functions by the struct to which they're attached first, then the name of the method itself. If the function isn't attached to a struct, put it at the bottom */
func (f FunctionNameList) Less(i, j int) bool {
	if f[i].str == "" && f[j].str != "" {
		return true
	}
	if f[i].str != "" && f[j].str == "" {
		return false
	}
	if f[i].str != "" && f[j].str != "" {
		if f[i].str > f[j].str {
			return true
		} else if f[i].str < f[j].str {
			return false
		}
	}
	return f[i].name < f[j].name
}
