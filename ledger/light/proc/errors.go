package proc

import (
	"errors"
)

var (
	ErrNotExecutor      = errors.New("not executor for jet")
	ErrExecutorMismatch = errors.New("sender isn't executor for object")
	ErrNotActivated     = errors.New("object should be activated")
)
