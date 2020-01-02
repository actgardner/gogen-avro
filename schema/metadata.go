package schema

// An embeddable struct that we can use to attach data for code generation to an AvroType or Definition
type generatorMetadata struct {
	data map[interface{}]interface{}
}

func (r *generatorMetadata) HasGeneratorMetadata(key interface{}) bool {
	if r.data == nil {
		return false
	}
	_, ok := r.data[key]
	return ok
}

func (r *generatorMetadata) SetGeneratorMetadata(key, value interface{}) {
	if r.data == nil {
		r.data = make(map[interface{}]interface{})
	}

	r.data[key] = value
}

func (r *generatorMetadata) GetGeneratorMetadata(key interface{}) interface{} {
	return r.data[key]
}
