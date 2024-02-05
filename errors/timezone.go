package errors

import "github.com/alioth-center/infrastructure/utils/values"

type LocalTimezoneAlreadySetError struct{}

func (e LocalTimezoneAlreadySetError) Error() string {
	return "local timezone already set"
}

func NewLocalTimezoneAlreadySetError() LocalTimezoneAlreadySetError {
	return LocalTimezoneAlreadySetError{}
}

type InvalidTimezoneError struct {
	Timezone string `json:"timezone"`
}

func (e InvalidTimezoneError) Error() string {
	return values.BuildStrings("invalid timezone: ", e.Timezone)
}

func NewInvalidTimezoneError(timezone string) InvalidTimezoneError {
	return InvalidTimezoneError{Timezone: timezone}
}
