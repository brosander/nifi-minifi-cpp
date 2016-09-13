package processor

type PropertyValueImpl struct {
	val string
}

func NewPropertyValue(val string) *PropertyValueImpl {
	return &PropertyValueImpl{val: val}
}

func (p *PropertyValueImpl) AsString() string {
	return p.val
}
