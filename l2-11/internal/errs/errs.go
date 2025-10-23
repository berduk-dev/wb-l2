package errs

import "errors"

var (
	ErrEventNotFound     = errors.New("event not found")
	ErrNoEvents          = errors.New("no events found for this user")
	ErrInvalidDateFormat = errors.New("invalid date format")
)
