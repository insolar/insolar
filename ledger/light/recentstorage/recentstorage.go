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

//go:generate minimock -i github.com/insolar/insolar/ledger/light/recentstorage.Provider -o ./ -s _mock.go

// Provider provides different types of storages for a specific jet
type Provider interface {
	GetPendingStorage(ctx context.Context, jetID insolar.ID) PendingStorage

	Count() int

	ClonePendingStorage(ctx context.Context, fromJetID, toJetID insolar.ID)

	RemovePendingStorage(ctx context.Context, id insolar.ID)
}

//go:generate minimock -i github.com/insolar/insolar/ledger/light/recentstorage.PendingStorage -o ./ -s _mock.go
type PendingStorage interface {
	AddPendingRequest(ctx context.Context, obj, req insolar.ID)
	SetContextToObject(ctx context.Context, obj insolar.ID, objContext PendingObjectContext)

	GetRequests() map[insolar.ID]PendingObjectContext
	GetRequestsForObject(obj insolar.ID) []insolar.ID

	RemovePendingRequest(ctx context.Context, obj, req insolar.ID)
}
