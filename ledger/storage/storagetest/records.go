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

package storagetest

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/testutils"
)

// AddRandIndex adds random index.
func AddRandIndex(
	ctx context.Context,
	// t *testing.T,
	objectStorage storage.ObjectStorage,
	jetID insolar.ID,
	pulsenum insolar.PulseNumber,
) (*insolar.ID, error) {
	parentID := testutils.RandomID()
	err := objectStorage.SetObjectIndex(ctx, jetID, &parentID, &object.Lifeline{
		LatestState: &parentID,
	})
	return &parentID, err
}

// AddRandBlob adds random blob.
func AddRandBlob(
	ctx context.Context,
	objectStorage storage.ObjectStorage,
	jetID insolar.ID,
	pulsenum insolar.PulseNumber,
) (*insolar.ID, error) {
	randID := testutils.RandomID()
	return objectStorage.SetBlob(ctx, jetID, pulsenum, randID[:])
}

// AddRandRecord adds random record.
func AddRandRecord(
	ctx context.Context,
	objectStorage storage.ObjectStorage,
	jetID insolar.ID,
	pulsenum insolar.PulseNumber,
) (*insolar.ID, error) {

	randID := testutils.RandomID()
	record := object.CodeRecord{
		Code: &randID,
	}
	return objectStorage.SetRecord(
		ctx,
		jetID,
		pulsenum,
		&record,
	)
}

// AddRandDrop adds random drop.
func AddRandDrop(
	ctx context.Context,
	modifier drop.Modifier,
	accessor drop.Accessor,
	jetID insolar.ID,
	pulsenum insolar.PulseNumber,
) (*drop.Drop, error) {

	hash1 := testutils.RandomID()
	hash2 := testutils.RandomID()
	drop := drop.Drop{
		Pulse:    pulsenum,
		PrevHash: hash1[:],
		Hash:     hash2[:],
		JetID:    insolar.JetID(jetID),
	}
	err := modifier.Set(ctx, drop)
	if err != nil {
		return nil, err
	}
	resDrop, err := accessor.ForPulse(ctx, insolar.JetID(jetID), pulsenum)
	return &resDrop, err
}
