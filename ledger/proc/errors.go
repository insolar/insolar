package proc

import "github.com/pkg/errors"

var (
	ErrWaiterNotLocked = errors.New("unlocked waiter unlock attempt")
)
