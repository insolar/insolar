package blockexplorer

import (
	"errors"
)

// Custom errors possibly useful to check by block explorer callers.
var (
	ErrInvalidRef        = errors.New("invalid reference")
	ErrObjectDeactivated = errors.New("object is deactivated")
	ErrNotFound          = errors.New("object not found")
	ErrUnexpectedReply   = errors.New("unexpected reply")
	ErrStateNotAvailable = errors.New("object state is not available")
)
