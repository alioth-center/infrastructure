package errors

type PromiseCompletedError struct{}

func (e PromiseCompletedError) Error() string {
	return "promise completed"
}

func NewPromiseCompletedError() PromiseCompletedError {
	return PromiseCompletedError{}
}
