package reply

import "github.com/insolar/insolar/core"

// Error is common reaction for methods returning id to lifeline states.
type Error struct {
	ErrType ErrType
}

// Type implementation of Reply interface.
func (e *Error) Type() core.ReplyType {
	return TypeError
}

func (e *Error) Error() error {
	switch e.ErrType {
	case ErrDeactivated:
		return core.ErrDeactivated
	}
	return core.ErrUnknown
}
