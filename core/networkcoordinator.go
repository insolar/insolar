/*
 *    Copyright 2018 Insolar
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

import (
	"context"
)

// NetworkCoordinator encapsulates logic of network configuration
//go:generate minimock -i github.com/insolar/insolar/core.NetworkCoordinator -o ../testutils -s _mock.go
type NetworkCoordinator interface {
	// GetCert returns certificate object by node reference, using discovery nodes for signing
	GetCert(context.Context, *RecordRef) (Certificate, error)

	// ValidateCert checks certificate signature
	ValidateCert(context.Context, AuthorizationCertificate) (bool, error)

	// TODO: Remove this method, use SetPulse instead
	WriteActiveNodes(ctx context.Context, number PulseNumber, activeNodes []Node) error

	// SetPulse uses PulseManager component for saving pulse info
	SetPulse(ctx context.Context, pulse Pulse) error

	// IsStarted returns true if component was started and false in other way
	IsStarted() bool
}
