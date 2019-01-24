/*
 *    Copyright 2018 Insolar
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

// Provider provides a recent storage for jet
//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.Provider -o ./ -s _mock.go
type Provider interface {
	GetStorage(ctx context.Context, jetID core.RecordID) RecentStorage
	CloneStorage(ctx context.Context, fromJetID, toJetID core.RecordID)
	RemoveStorage(ctx context.Context, id core.RecordID)
}

// RecentStorage is a base interface for the storage of recent objects and indexes
//go:generate minimock -i github.com/insolar/insolar/ledger/recentstorage.RecentStorage -o ./ -s _mock.go
type RecentStorage interface {
	AddObject(ctx context.Context, id core.RecordID)
	AddObjectWithTLL(ctx context.Context, id core.RecordID, ttl int)

	AddPendingRequest(ctx context.Context, obj, req core.RecordID)
	RemovePendingRequest(ctx context.Context, obj, req core.RecordID)

	GetObjects() map[core.RecordID]int
	GetRequests() map[core.RecordID]map[core.RecordID]struct{}
	GetRequestsForObject(obj core.RecordID) []core.RecordID

	IsRecordIDCached(obj core.RecordID) bool

	DecreaseTTL(ctx context.Context)
}
