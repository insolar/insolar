// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexAccessor -o ./ -s _mock.go -g

// IndexAccessor provides an interface for fetching buckets from an index.
type IndexAccessor interface {
	ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (record.Index, error)
	// ForPulse returns a collection of buckets for a provided pulse number.
	ForPulse(ctx context.Context, pn insolar.PulseNumber) ([]record.Index, error)
	// LastKnownForID returns a latest version of an index
	LastKnownForID(ctx context.Context, objID insolar.ID) (record.Index, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.MemoryIndexModifier -o ./ -s _mock.go

// MemoryIndexModifier writes index to in-memory storage.
type MemoryIndexModifier interface {
	Set(ctx context.Context, pn insolar.PulseNumber, index record.Index)
	SetIfNone(ctx context.Context, pn insolar.PulseNumber, index record.Index)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.MemoryIndexAccessor -o ./ -s _mock.go -g

// MemoryIndexAccessor provides an interface for fetching buckets from an index.
type MemoryIndexAccessor interface {
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
	MemoryIndexAccessor
	MemoryIndexModifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.IndexCleaner -o ./ -s _mock.go -g

// IndexCleaner provides an interface for removing backets from a storage.
type IndexCleaner interface {
	// DeleteForPN method removes indexes from a storage for a provided
	DeleteForPN(ctx context.Context, pn insolar.PulseNumber)
}
