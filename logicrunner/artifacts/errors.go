package artifacts

import (
	"errors"
)

// Custom errors possibly useful to check by artifact manager callers.
var (
	ErrObjectDeactivated = errors.New("object is deactivated")
	ErrNotFound          = errors.New("object not found")
	ErrNoReply           = errors.New("timeout while awaiting reply from watermill")
)
