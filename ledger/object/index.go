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

package object

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexModifier -o ./ -s _mock.go

// IndexModifier provides methods for modifying buckets of index.
// Index contains buckets with pn->objID->Bucket hierarchy.
// With using of IndexModifier there is a possibility to set buckets from outside of an index.
type IndexModifier interface {
	CreateIndex(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) *LockedIndex
	// SetIndex adds a bucket with provided pulseNumber and ID
	SetIndex(ctx context.Context, pn insolar.PulseNumber, bucket FilamentIndex) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexAccessor -o ./ -s _mock.go

// IndexAccessor provides an interface for fetching buckets from an index.
type IndexAccessor interface {
	Index(pn insolar.PulseNumber, objID insolar.ID) *LockedIndex
	// ForPNAndJet returns a collection of buckets for a provided pn and jetID
	ForPNAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) []FilamentIndex
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexCleaner -o ./ -s _mock.go

// IndexCleaner provides an interface for removing backets from a storage.
type IndexCleaner interface {
	// DeleteForPN method removes indexes from a storage for a provided
	DeleteForPN(ctx context.Context, pn insolar.PulseNumber)
}
