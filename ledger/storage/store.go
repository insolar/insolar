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
	"github.com/insolar/insolar/ledger/storage/object"
)

// Store is used by context unaware clients who can work inside transactions as well as outside.
type Store interface {
	GetObjectIndex(ctx context.Context, jetID insolar.ID, ref *insolar.ID) (*object.Lifeline, error)
	SetObjectIndex(ctx context.Context, jetID insolar.ID, ref *insolar.ID, idx *object.Lifeline) error

	// Deprecated: use insolar.PulseStorage.Current() instead
	GetLatestPulse(ctx context.Context) (*Pulse, error)
	GetPulse(ctx context.Context, num insolar.PulseNumber) (*Pulse, error)
}
