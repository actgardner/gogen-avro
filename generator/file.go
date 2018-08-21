package generator

import (
	"fmt"
	"go/format"
	"io/ioutil"
	"sort"
	"strings"
)

// File represents a Go source file in the output package
type File struct {
	name      string
	headers   []string
	functions map[FunctionName]string
	structs   map[string]string
	imports   map[string]interface{}
	constants map[string]interface{}
}

func NewFile(name string) *File {
	return &File{
		name:      name,
		functions: make(map[FunctionName]string),
		structs:   make(map[string]string),
		imports:   make(map[string]interface{}),
		constants: make(map[string]interface{}),
	}
}

type FunctionName struct {
	// The target struct type, if there is one
	Str string
	// The function name
	Name string
}

// Write the contents of the file:
//   - imports (all in one block)
//   - struct definitions (alphabetically)
//   - functions (sorted alphabetically by struct to which they're attached, then unattached funcs)

// TODO: It'd be better to group funcs attached to a struct with the struct definition
func (f *File) WriteFile(pkgName, targetFile string) error {
	src := fmt.Sprintf("%v\n\npackage %v\n%v\n%v\n%v\n%v\n", f.headerString(), pkgName, f.importString(), f.constantString(), f.structString(), f.functionString())
	fileContent, err := format.Source([]byte(src))
	if err != nil {
		return fmt.Errorf("Error formatting file %v - %v\n\nContents: %v", f.name, err, src)
	}
	err = ioutil.WriteFile(targetFile, fileContent, 0640)
	if err != nil {
		return fmt.Errorf("Error writing file %v - %v", f.name, err)
	}
	return nil
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

func (f *File) headerString() string {
	if len(f.headers) == 0 {
		return ""
	}

	return strings.Join(f.headers, "\n")
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
