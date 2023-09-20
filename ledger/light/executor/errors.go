package executor

import (
	"github.com/pkg/errors"
)

var (
	ErrWaiterNotLocked = errors.New("unlocked waiter unlock attempt")
	ErrWriteClosed     = errors.New("requested pulse is closed for writing")
)
