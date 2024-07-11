package timezone

import "errors"

var (
	ErrLocalTimezoneAlreadySet = errors.New("local timezone already set")
	ErrInvalidTimezone         = errors.New("invalid timezone")
)
