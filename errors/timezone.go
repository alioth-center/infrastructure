package errors

type LocalTimezoneAlreadySetError struct{}

func (e LocalTimezoneAlreadySetError) Error() string {
	return "local timezone already set"
}

func NewLocalTimezoneAlreadySetError() LocalTimezoneAlreadySetError {
	return LocalTimezoneAlreadySetError{}
}
