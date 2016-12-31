package types

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


