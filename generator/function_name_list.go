package generator

// Implement the Sortable interface for FunctionNames
type FunctionNameList []FunctionName

func (f FunctionNameList) Len() int {
	return len(f)
}

func (f FunctionNameList) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// Sort functions by the struct to which they're attached first, then the name of the method itself. If the function isn't attached to a struct, put it at the bottom
func (f FunctionNameList) Less(i, j int) bool {
	if f[i].Str == "" && f[j].Str != "" {
		return true
	}
	if f[i].Str != "" && f[j].Str == "" {
		return false
	}
	if f[i].Str != "" && f[j].Str != "" {
		if f[i].Str > f[j].Str {
			return true
		} else if f[i].Str < f[j].Str {
			return false
		}
	}
	return f[i].Name < f[j].Name
}
