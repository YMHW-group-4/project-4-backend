package errors

import (
	"errors"
	"fmt"
)

// errInvalidArgument is the base error for ErrInvalidArgument.
var errInvalidArgument = errors.New("argument is invalid")

// errInvalidData the base error for ErrInvalidData.
var errInvalidData = errors.New("invalid data")

// ErrInvalidArgument can be used when an argument is invalid. A message can be specified
// using the standard golang formatting rules.
func ErrInvalidArgument(format string, a ...interface{}) error {
	return fmt.Errorf("%w: %s", errInvalidArgument, fmt.Sprintf(format, a...))
}

// ErrInvalidData can be used when the data is invalid. A message can be specified
// using the standard golang formatting rules.
func ErrInvalidData(format string, a ...interface{}) error {
	return fmt.Errorf("%w: %s", errInvalidData, fmt.Sprintf(format, a...))
}
