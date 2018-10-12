package core

import "github.com/pkg/errors"

var (
	// ErrUnknown returned when error type cannot be defined.
	ErrUnknown = errors.New("unknown error")
	// ErrDeactivated returned when requested object is deactivated.
	ErrDeactivated = errors.New("object is deactivated")
)
