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
	"github.com/insolar/insolar/instrumentation/inslogger"
	"go.opencensus.io/stats"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineIndex -o ./ -s _mock.go

// LifelineIndex is a base storage for lifelines.
type LifelineIndex interface {
	// LifelineAccessor provides methods for fetching lifelines.
	LifelineAccessor
	// LifelineModifier provides methods for modifying lifelines.
	LifelineModifier
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineAccessor -o ./ -s _mock.go

// LifelineAccessor provides methods for fetching lifelines.
type LifelineAccessor interface {
	// ForID returns a lifeline from a bucket with provided PN and ObjID
	ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (Lifeline, error)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineModifier -o ./ -s _mock.go

// LifelineModifier provides methods for modifying lifelines.
type LifelineModifier interface {
	// Set set a lifeline to a bucket with provided pulseNumber and ID
	Set(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error
}

//go:generate minimock -i github.com/insolar/insolar/ledger/object.LifelineStateModifier -o ./ -s _mock.go

// LifelineStateModifier provides an interface for changing a state of lifeline.
type LifelineStateModifier interface {
	// SetLifelineUsage updates a last usage fields of a bucket for a provided pulseNumber and an object id
	SetLifelineUsage(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error
}

type LifelineStorage struct {
	idxAccessor IndexAccessor
	idxModifier IndexModifier
}

func NewLifelineStorage(idxAccessor IndexAccessor, idxModifier IndexModifier) *LifelineStorage {
	return &LifelineStorage{idxAccessor: idxAccessor, idxModifier: idxModifier}
}

// ForID returns a lifeline from a bucket with provided PN and ObjID
func (i *LifelineStorage) ForID(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) (Lifeline, error) {
	b := i.idxAccessor.Index(pn, objID)
	if b == nil {
		return Lifeline{}, ErrLifelineNotFound
	}
	b.RLock()
	b.RUnlock()

	return CloneLifeline(b.objectMeta.Lifeline), nil
}

// Set sets a lifeline to a bucket with provided pulseNumber and ID
func (i *LifelineStorage) Set(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID, lifeline Lifeline) error {
	b := i.idxAccessor.Index(pn, objID)
	if b == nil {
		b = i.idxModifier.CreateIndex(ctx, pn, objID)
	}

	b.Lock()
	defer b.Unlock()

	b.objectMeta.Lifeline = lifeline
	b.objectMeta.LifelineLastUsed = pn

	stats.Record(ctx,
		statBucketAddedCount.M(1),
	)

	inslogger.FromContext(ctx).Debugf("[Set] lifeline for obj - %v was set successfully", objID.DebugString())
	return nil
}

// SetLifelineUsage updates a last usage fields of a bucket for a provided pulseNumber and an object id
func (i *LifelineStorage) SetLifelineUsage(ctx context.Context, pn insolar.PulseNumber, objID insolar.ID) error {
	b := i.idxAccessor.Index(pn, objID)
	if b == nil {
		return ErrLifelineNotFound
	}

	b.Lock()
	defer b.Unlock()

	b.objectMeta.LifelineLastUsed = pn

	return nil
}

// EncodeLifeline converts lifeline index into binary format.
func EncodeLifeline(index Lifeline) []byte {
	res, err := index.Marshal()
	if err != nil {
		panic(err)
	}

	return res
}

// MustDecodeLifeline converts byte array into lifeline index struct.
func MustDecodeLifeline(buff []byte) (index Lifeline) {
	idx, err := DecodeLifeline(buff)
	if err != nil {
		panic(err)
	}

	return idx
}

// DecodeLifeline converts byte array into lifeline index struct.
func DecodeLifeline(buff []byte) (Lifeline, error) {
	lfl := Lifeline{}
	err := lfl.Unmarshal(buff)
	return lfl, err
}

// CloneLifeline returns copy of argument idx value.
func CloneLifeline(idx Lifeline) Lifeline {
	if idx.LatestState != nil {
		tmp := *idx.LatestState
		idx.LatestState = &tmp
	}

	if idx.LatestStateApproved != nil {
		tmp := *idx.LatestStateApproved
		idx.LatestStateApproved = &tmp
	}

	if idx.ChildPointer != nil {
		tmp := *idx.ChildPointer
		idx.ChildPointer = &tmp
	}

	if idx.Delegates != nil {
		cp := make([]LifelineDelegate, len(idx.Delegates))
		copy(cp, idx.Delegates)
		idx.Delegates = cp
	} else {
		idx.Delegates = []LifelineDelegate{}
	}

	if idx.EarliestOpenRequest != nil {
		tmp := *idx.EarliestOpenRequest
		idx.EarliestOpenRequest = &tmp
	}

	if idx.PendingPointer != nil {
		tmp := *idx.PendingPointer
		idx.PendingPointer = &tmp
	}

	return idx
}

func (m *Lifeline) SetDelegate(key insolar.Reference, value insolar.Reference) {
	for _, d := range m.Delegates {
		if d.Key == key {
			d.Value = value
			return
		}
	}

	m.Delegates = append(m.Delegates, LifelineDelegate{Key: key, Value: value})
}

func (m *Lifeline) DelegateByKey(key insolar.Reference) (insolar.Reference, bool) {
	for _, d := range m.Delegates {
		if d.Key == key {
			return d.Value, true
		}
	}

	return [64]byte{}, false
}
