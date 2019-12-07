package templates

type PrimitiveContext struct {
	goType           string
	wrapperType      string
	serializerMethod string
}

func (p *PrimitiveContext) GoType() string {
	return p.goType
}

func (p *PrimitiveContext) ConstructorMethod() string {
	return ""
}

func (p *PrimitiveContext) SerializerMethod() string {
	return ""
}

func (p *PrimitiveContext) Template() string {
	return ""
}

func (p *PrimitiveContext) WrapperType() string {
	return wrapperType
}
