// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package reply

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/flow"
)

// OK is a generic reply for signaling a positive result.
type OK struct {
}

// Type implementation of Reply interface.
func (e *OK) Type() insolar.ReplyType {
	return TypeOK
}

// NotOK is a generic reply for signaling a negative result.
type NotOK struct {
}

// Type implementation of Reply interface.
func (e *NotOK) Type() insolar.ReplyType {
	return TypeNotOK
}

// Error is common error reaction.
type Error struct {
	ErrType ErrType
}

// Type implementation of Reply interface.
func (e *Error) Type() insolar.ReplyType {
	return TypeError
}

// Error returns concrete error for stored type.
func (e *Error) Error() error {
	switch e.ErrType {
	case ErrDeactivated:
		return insolar.ErrDeactivated
	case ErrStateNotAvailable:
		return insolar.ErrStateNotAvailable
	case ErrHotDataTimeout:
		return insolar.ErrHotDataTimeout
	case ErrNoPendingRequests:
		return insolar.ErrNoPendingRequest
	case FlowCancelled:
		return flow.ErrCancelled
	}

	return insolar.ErrUnknown
}
