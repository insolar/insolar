/*
 *    Copyright 2018 INS Ecosystem
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

package networkcoordinator

import (
	"context"

	"github.com/insolar/insolar/core"
)

// Coordinator interface contains NetworkState dependent methods
type Coordinator interface {
	// GetCert returns certificate object by node reference, using discovery nodes for signing
	GetCert(context.Context, *core.RecordRef) (core.Certificate, error)

	// SetPulse uses PulseManager component for saving pulse info
	SetPulse(ctx context.Context, pulse core.Pulse) error

	// signCertHandler is used by MsgBus handler for signing certificate
	signCertHandler(ctx context.Context, p core.Parcel) (core.Reply, error)
}
