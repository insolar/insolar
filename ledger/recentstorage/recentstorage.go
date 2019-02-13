/*
 *    Copyright 2019 Insolar
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

package recentstorage

import (
	"context"

	"github.com/insolar/insolar/core"
)

// Provider provides different types of storages for a specific jet
//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.Provider -o ./ -s _mock.go
type Provider interface {
	GetIndexStorage(ctx context.Context, jetID core.RecordID) RecentIndexStorage
	GetPendingStorage(ctx context.Context, jetID core.RecordID) PendingStorage

	CloneIndexStorage(ctx context.Context, fromJetID, toJetID core.RecordID)
	ClonePendingStorage(ctx context.Context, fromJetID, toJetID core.RecordID)

	DecreaseIndexesTTL(ctx context.Context) map[core.RecordID][]core.RecordID

	RemovePendingStorage(ctx context.Context, id core.RecordID)
}

// RecentIndexStorage is a struct which contains `recent indexes` for a specific jet
// `recent index` is a index which was called between TTL-border
// If index is put to a recent storage, it'll be there for TTL-pulses at least
//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.RecentIndexStorage -o ./ -s _mock.go
type RecentIndexStorage interface {
	AddObject(ctx context.Context, id core.RecordID)
	AddObjectWithTLL(ctx context.Context, id core.RecordID, ttl int)

	GetObjects() map[core.RecordID]int

	DecreaseIndexTTL(ctx context.Context) []core.RecordID

	FilterNotExistWithLock(ctx context.Context, candidates []core.RecordID, fn func(filtered []core.RecordID))
}

//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.PendingStorage -o ./ -s _mock.go
type PendingStorage interface {
	AddPendingRequest(ctx context.Context, obj, req core.RecordID)
	SetContextToObject(ctx context.Context, obj core.RecordID, objContext PendingObjectContext)

	GetRequests() map[core.RecordID]PendingObjectContext
	GetRequestsForObject(obj core.RecordID) []core.RecordID

	RemovePendingRequest(ctx context.Context, obj, req core.RecordID)
}
