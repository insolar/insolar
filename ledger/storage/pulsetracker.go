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

package storage

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

// Pulse is a record containing pulse info.
type Pulse struct {
	Prev         *insolar.PulseNumber
	Next         *insolar.PulseNumber
	SerialNumber int
	Pulse        insolar.Pulse
}

// PulseTracker allows to modify state of the pulse inside storage (internal or external)
// DEPRECATED
//go:generate minimock -i github.com/insolar/insolar/ledger/storage.PulseTracker -o ./ -s _mock.go
type PulseTracker interface {
	GetPulse(ctx context.Context, num insolar.PulseNumber) (*Pulse, error)
	GetPreviousPulse(ctx context.Context, num insolar.PulseNumber) (*Pulse, error)
	GetNthPrevPulse(ctx context.Context, n uint, from insolar.PulseNumber) (*Pulse, error)
	GetLatestPulse(ctx context.Context) (*Pulse, error)

	AddPulse(ctx context.Context, pulse insolar.Pulse) error

	DeletePulse(ctx context.Context, num insolar.PulseNumber) error
}
