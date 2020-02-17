// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package foundation

// Error elementary string based error struct satisfying builtin error interface
//    foundation.Error{"some err"}
type Error struct {
	S string
}

// Error returns error in string format
func (e *Error) Error() string {
	return e.S
}
