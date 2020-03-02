// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insolar

import "github.com/pkg/errors"

var (
	// ErrUnknown is returned when error type cannot be defined.
	ErrUnknown = errors.New("unknown error")
	// ErrDeactivated is returned when requested object is deactivated.
	ErrDeactivated = errors.New("object is deactivated")
	// ErrStateNotAvailable is returned when requested object is deactivated.
	ErrStateNotAvailable = errors.New("object state is not available")
	// ErrHotDataTimeout is returned when no hot data received for a specific jet
	ErrHotDataTimeout = errors.New("requests were abandoned due to hot-data timeout")
	// ErrNoPendingRequest is returned when there are no pending requests on current LME
	ErrNoPendingRequest = errors.New("no pending requests are available")
	// ErrNotFound is returned when something not found
	ErrNotFound = errors.New("not found")
)
