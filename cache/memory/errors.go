package memory

import "strings"

type ValueTypeNotMatchError struct {
	Want     Type
	Actually Type
}

func (e ValueTypeNotMatchError) Error() string {
	builder := strings.Builder{}
	builder.WriteString("value type not match, want ")
	builder.WriteString(string(e.Want))
	builder.WriteString(" but actually ")
	builder.WriteString(string(e.Actually))
	return builder.String()
}

func NewValueTypeNotMatchError(want, actually Type) ValueTypeNotMatchError {
	return ValueTypeNotMatchError{Want: want, Actually: actually}
}

type ReceiverTypeIncorrectError struct {
	Type string
	Nil  bool
}

func (e ReceiverTypeIncorrectError) Error() string {
	builder := strings.Builder{}
	builder.WriteString("receiver must be a non nil pointer but actually ")
	builder.WriteString(e.Type)
	if e.Nil {
		builder.WriteString(" is nil")
	}
	return builder.String()
}

func NewReceiverTypeIncorrectError(t string, isNil bool) ReceiverTypeIncorrectError {
	return ReceiverTypeIncorrectError{Type: t, Nil: isNil}
}
