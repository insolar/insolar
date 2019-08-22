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
	"github.com/insolar/insolar/insolar/record"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexModifier -o ./ -s _mock.go -g

// IndexModifier provides methods for modifying buckets of index.
// Lifeline contains buckets with pn->objID->Bucket hierarchy.
// With using of IndexModifier there is a possibility to set buckets from outside of an index.
type IndexModifier interface {
	// SetIndex adds a bucket with provided pulseNumber and ID
	SetIndex(ctx context.Context, pn insolar.PulseNumber, index record.Index) error
	// UpdateLastKnownPulse updates last know pulse to given one for all objects from this pulse
	UpdateLastKnownPulse(ctx context.Context, pn insolar.PulseNumber) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.MemoryIndexModifier -o ./ -s _mock.go

// MemoryIndexModifier writes index to in-memory storage.
type MemoryIndexModifier interface {
	Set(ctx context.Context, pn insolar.PulseNumber, index record.Index)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexAccessor -o ./ -s _mock.go -g

// IndexAccessor provides an interface for fetching buckets from an index.
type IndexAccessor interface {
	ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (record.Index, error)
	// ForPulse returns a collection of buckets for a provided pulse number.
	ForPulse(ctx context.Context, pn insolar.PulseNumber) ([]record.Index, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexStorage -o ./ -s _mock.go -g

type IndexStorage interface {
	IndexAccessor
	IndexModifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.MemoryIndexStorage -o ./ -s _mock.go -g

type MemoryIndexStorage interface {
	IndexAccessor
	MemoryIndexModifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexCleaner -o ./ -s _mock.go -g

// IndexCleaner provides an interface for removing backets from a storage.
type IndexCleaner interface {
	// DeleteForPN method removes indexes from a storage for a provided
	DeleteForPN(ctx context.Context, pn insolar.PulseNumber)
}
