/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package core

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
	// ErrTooManyPendingRequests is returned when a limit of pending requests has been reached on a current LME
	ErrTooManyPendingRequests = errors.New("the limit of pending requests count has been reached")
	// ErrNoNodes is returned if no matching nodes found
	ErrNoNodes = errors.New("no matching nodes")
)
