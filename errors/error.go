package errors

import (
	"errors"
	"fmt"
)

// errInvalidArgument is the base error for ErrInvalidArgument.
var errInvalidArgument = errors.New("invalid argument")

// errInvalidOperation is the base error for ErrInvalidOperation.
var errInvalidOperation = errors.New("invalid operation")

// ErrInvalidArgument can be used when an argument is invalid. A message can be specified
// using the standard golang formatting rules.
func ErrInvalidArgument(format string, a ...interface{}) error {
	return fmt.Errorf("%w: %s", errInvalidArgument, fmt.Sprintf(format, a...))
}

// ErrInvalidOperation can be used when the operation is invalid. A message can be specified
// using the standard golang formatting rules.
func ErrInvalidOperation(format string, a ...interface{}) error {
	return fmt.Errorf("%w: %s", errInvalidOperation, fmt.Sprintf(format, a...))
}
