// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package serialization

import (
	"github.com/pkg/errors"
)

func ErrPayloadLengthMismatch(expected, actual int64) error {
	return errors.Errorf("payload length mismatch - expected: %d, actual: %d", expected, actual)
}

func ErrMalformedPulseNumber(err error) error {
	return errors.Wrap(err, "malformed pulse number")
}

func ErrMalformedHeader(err error) error {
	return errors.Wrap(err, "malformed header")
}

func ErrMalformedPacketBody(err error) error {
	return errors.Wrap(err, "malformed packet body")
}

func ErrMalformedPacketSignature(err error) error {
	return errors.Wrap(err, "invalid packet signature")
}
