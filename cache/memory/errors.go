package memory

type ValueTypeNotMatchError struct {
	Want     Type
	Actually Type
}

func (e ValueTypeNotMatchError) Error() string {
	return "value type not match, want " + string(e.Want) + " but actually " + string(e.Actually)
}

func NewValueTypeNotMatchError(want, actually Type) ValueTypeNotMatchError {
	return ValueTypeNotMatchError{Want: want, Actually: actually}
}
