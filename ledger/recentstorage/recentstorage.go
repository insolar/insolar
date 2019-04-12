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

package recentstorage

import (
	"context"

	"github.com/insolar/insolar/insolar"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.Provider -o ./ -s _mock.go

// Provider provides different types of storages for a specific jet
type Provider interface {
	GetIndexStorage(ctx context.Context, jetID insolar.ID) RecentIndexStorage
	GetPendingStorage(ctx context.Context, jetID insolar.ID) PendingStorage

	Count() int

	CloneIndexStorage(ctx context.Context, fromJetID, toJetID insolar.ID)
	ClonePendingStorage(ctx context.Context, fromJetID, toJetID insolar.ID)

	DecreaseIndexesTTL(ctx context.Context) map[insolar.ID][]insolar.ID

	RemovePendingStorage(ctx context.Context, id insolar.ID)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.RecentIndexStorage -o ./ -s _mock.go

// RecentIndexStorage is a struct which contains `recent indexes` for a specific jet
// `recent index` is a index which was called between TTL-border
// If index is put to a recent storage, it'll be there for TTL-pulses at least
type RecentIndexStorage interface {
	AddObject(ctx context.Context, id insolar.ID)
	AddObjectWithTLL(ctx context.Context, id insolar.ID, ttl int)

	GetObjects() map[insolar.ID]int

	DecreaseIndexTTL(ctx context.Context) []insolar.ID

	FilterNotExistWithLock(ctx context.Context, candidates []insolar.ID, fn func(filtered []insolar.ID))
}

//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.PendingStorage -o ./ -s _mock.go
type PendingStorage interface {
	AddPendingRequest(ctx context.Context, obj, req insolar.ID)
	SetContextToObject(ctx context.Context, obj insolar.ID, objContext PendingObjectContext)

	GetRequests() map[insolar.ID]PendingObjectContext
	GetRequestsForObject(obj insolar.ID) []insolar.ID

	RemovePendingRequest(ctx context.Context, obj, req insolar.ID)
}
