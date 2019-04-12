//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package insolar

import (
	"context"
)

//go:generate minimock -i github.com/insolar/insolar/insolar.NetworkCoordinator -o ../testutils -s _mock.go

// NetworkCoordinator encapsulates logic of network configuration
type NetworkCoordinator interface {
	// GetCert returns certificate object by node reference, using discovery nodes for signing
	GetCert(context.Context, *Reference) (Certificate, error)

	// ValidateCert checks certificate signature
	ValidateCert(context.Context, AuthorizationCertificate) (bool, error)

	// SetPulse uses PulseManager component for saving pulse info
	SetPulse(ctx context.Context, pulse Pulse) error

	// IsStarted returns true if component was started and false in other way
	IsStarted() bool
}
