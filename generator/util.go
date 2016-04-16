package generator

import (
	"sort"
	"strings"
)

func toPublicName(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToUpper(name[0:1]) + name[1:]
}

func concatSortedMap(m map[string]string, sep string) string {
	keys := make([]string, 0)
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	s := ""
	for _, k := range keys {
		s += m[k] + sep
	}
	return s
}

func interfaceSliceToStringSlice(iSlice []interface{}) ([]string, bool) {
	var ok bool
	stringSlice := make([]string, len(iSlice))
	for i, v := range iSlice {
		stringSlice[i], ok = v.(string)
		if !ok {
			return nil, false
		}
	}
	return stringSlice, true
}
