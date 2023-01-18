package blockchain

import (
	"errors"
	"fmt"
)

// ErrInvalidTransaction can be used when the transaction is invalid. A message can be specified
// using the standard golang formatting rules.
func ErrInvalidTransaction(format string, a ...interface{}) error {
	return fmt.Errorf("%w: %s", errors.New("invalid transaction"), fmt.Sprintf(format, a...))
}

// ErrInvalidBLock can be used when the block is invalid. A message can be specified
// using the standard golang formatting rules.
func ErrInvalidBLock(format string, a ...interface{}) error {
	return fmt.Errorf("%w: %s", errors.New("invalid block"), fmt.Sprintf(format, a...))
}
