// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package errors

import (
	"fmt"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/pulse"
)

func NewPulseRoundMismatchError(pn pulse.Number, msg string) error {
	return &nextPulseRoundError{pn: pn, s: msg}
}

func NewPulseRoundMismatchErrorDef(pn pulse.Number, filterPN pulse.Number, localID insolar.ShortNodeID, from interface{}, details string) error {
	msg := fmt.Sprintf("packet pulse number mismatched: expected=%v, actual=%v, local=%d, from=%v, details=%v",
		filterPN, pn, localID, from, details)
	return NewPulseRoundMismatchError(pn, msg)
}

func IsMismatchPulseError(err error) (bool, pulse.Number) {
	pr, ok := err.(*nextPulseRoundError)
	if !ok {
		return false, pulse.Unknown
	}
	return !pr.pn.IsUnknown(), pr.pn
}

type nextPulseRoundError struct {
	pn pulse.Number
	s  string
}

func (e *nextPulseRoundError) Error() string {
	return e.s
}
