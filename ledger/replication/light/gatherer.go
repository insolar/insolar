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

package light

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/object"
)

type DataGatherer interface {
	ForPulseAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) (*message.HeavyPayload, error)
}

type dataGatherer struct {
	dropAccessor    drop.Accessor
	blobsAccessor   blob.CollectionAccessor
	recsAccessor    object.RecordCollectionAccessor
	indexesAccessor object.IndexCollectionAccessor
}

func NewDataGatherer(
	dropAccessor drop.Accessor,
	blobsAccessor blob.CollectionAccessor,
	recsAccessor object.RecordCollectionAccessor,
	indexesAccessor object.IndexCollectionAccessor,
) *dataGatherer {
	return &dataGatherer{
		dropAccessor:    dropAccessor,
		blobsAccessor:   blobsAccessor,
		recsAccessor:    recsAccessor,
		indexesAccessor: indexesAccessor,
	}
}

func (d *dataGatherer) ForPulseAndJet(ctx context.Context, pn insolar.PulseNumber, jetID insolar.JetID) (*message.HeavyPayload, error) {
	dr, err := d.dropAccessor.ForPulse(ctx, jetID, pn)
	if err != nil {
		inslogger.FromContext(ctx).Error("synchronize: can't fetch a drop")
		return nil, err
	}

	bls := d.blobsAccessor.ForPulse(ctx, jetID, pn)
	records := d.recsAccessor.ForPulse(ctx, jetID, pn)

	indexes := d.indexesAccessor.ForPulseAndJet(ctx, jetID, pn)
	resIdx := map[insolar.ID][]byte{}
	for id, idx := range indexes {
		resIdx[id] = object.EncodeIndex(idx)
	}

	return &message.HeavyPayload{
		JetID:    jetID,
		PulseNum: pn,
		Indexes:  resIdx,
		Drop:     drop.MustEncode(&dr),
		Blobs:    convertBlobs(bls),
		Records:  convertRecords(records),
	}, nil
}

func convertBlobs(blobs []blob.Blob) [][]byte {
	var res [][]byte
	for _, b := range blobs {
		res = append(res, blob.MustEncode(&b))
	}
	return res
}

func convertRecords(records []record.MaterialRecord) (result [][]byte) {
	for _, r := range records {
		result = append(result, object.EncodeMaterial(r))
	}
	return
}
