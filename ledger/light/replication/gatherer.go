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

package replication

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

//go:generate minimock -i github.com/insolar/insolar/ledger/light/replication.DataGatherer -o ./ -s _mock.go

// DataGatherer is an inteface, that provides methods for gathering a heavy payload for a provided pulse
type DataGatherer interface {
	// ForPulseAndJet returns HeavyPayload message for a provided pulse and a jetID
	ForPulseAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) (*message.HeavyPayload, error)
}

// LightDataGatherer is a realisation of DataGatherer
type LightDataGatherer struct {
	dropAccessor         drop.Accessor
	blobsAccessor        blob.CollectionAccessor
	recsAccessor         object.RecordCollectionAccessor
	indexReplicaAccessor object.IndexBucketAccessor
}

// NewDataGatherer creates a new instance of LightDataGatherer
func NewDataGatherer(
	dropAccessor drop.Accessor,
	blobsAccessor blob.CollectionAccessor,
	recsAccessor object.RecordCollectionAccessor,
	indexReplicaAccessor object.IndexBucketAccessor,
) *LightDataGatherer {
	return &LightDataGatherer{
		dropAccessor:         dropAccessor,
		blobsAccessor:        blobsAccessor,
		recsAccessor:         recsAccessor,
		indexReplicaAccessor: indexReplicaAccessor,
	}
}

// ForPulseAndJet returns HeavyPayload message for a provided pulse and a jetID
func (d *LightDataGatherer) ForPulseAndJet(
	ctx context.Context,
	pn insolar.PulseNumber,
	jetID insolar.JetID,
) (*message.HeavyPayload, error) {
	dr, err := d.dropAccessor.ForPulse(ctx, jetID, pn)
	if err != nil {
		inslogger.FromContext(ctx).Error("synchronize: can't fetch a drop")
		return nil, errors.Wrap(err, "failed to fetch drop")
	}

	bls := d.blobsAccessor.ForPulse(ctx, jetID, pn)
	records := d.recsAccessor.ForPulse(ctx, jetID, pn)

	indexes := d.indexReplicaAccessor.ForPNAndJet(ctx, pn, jetID)

	return &message.HeavyPayload{
		JetID:        jetID,
		PulseNum:     pn,
		IndexBuckets: convertIndexBuckets(ctx, indexes),
		Drop:         drop.MustEncode(&dr),
		Blobs:        convertBlobs(bls),
		Records:      convertRecords(ctx, records),
	}, nil
}

func convertIndexBuckets(ctx context.Context, buckets []object.IndexBucket) [][]byte {
	convertedBucks := make([][]byte, len(buckets))
	for i, buck := range buckets {
		buff, err := buck.Marshal()
		if err != nil {
			inslogger.FromContext(ctx).Errorf("problems with marshaling bucket - %v", err)
			continue
		}
		convertedBucks[i] = buff
	}

	return convertedBucks
}

func convertBlobs(blobs []blob.Blob) [][]byte {
	res := make([][]byte, len(blobs))
	for i, b := range blobs {
		temp := b
		res[i] = blob.MustEncode(&temp)
	}
	return res
}

func convertRecords(ctx context.Context, records []record.Material) [][]byte {
	res := make([][]byte, len(records))
	for i, r := range records {
		data, err := r.Marshal()
		if err != nil {
			inslogger.FromContext(ctx).Error("Can't serialize record", r)
		}
		res[i] = data
	}
	return res
}
