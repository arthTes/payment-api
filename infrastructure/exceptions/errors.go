package exceptions

import "errors"

var (
	EntityNotFoundError       = errors.New("entity not found")
	PersistenceError          = errors.New("cannot persist error")
	InvalidAmountError        = errors.New("invalid amount value")
	InvalidOperationTypeError = errors.New("invalid operation type value")
	InvalidParameterError     = errors.New("invalid parameter value")
)
