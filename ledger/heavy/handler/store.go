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

package handler

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/drop"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

func storeIndexes(
	ctx context.Context,
	indexes object.IndexModifier,
	rawIndexes map[insolar.ID][]byte,
) error {
	for id, rwi := range rawIndexes {
		idx, err := object.DecodeIndex(rwi)
		if err != nil {
			inslogger.FromContext(ctx).Error(err)
			continue
		}

		err = indexes.Set(ctx, id, idx)
		if err != nil {
			return errors.Wrapf(err, "heavyserver: index storing failed")
		}
	}

	return nil
}

func storeDrop(
	ctx context.Context,
	drops drop.Modifier,
	rawDrop []byte,
) error {
	d, err := drop.Decode(rawDrop)
	if err != nil {
		inslogger.FromContext(ctx).Error(err)
		return err
	}
	err = drops.Set(ctx, *d)
	if err != nil {
		return errors.Wrapf(err, "heavyserver: drop storing failed")
	}

	return nil
}

func storeBlobs(
	ctx context.Context,
	blobs blob.Modifier,
	pcs insolar.PlatformCryptographyScheme,
	pn insolar.PulseNumber,
	rawBlobs [][]byte,
) {
	inslog := inslogger.FromContext(ctx)

	for _, rwb := range rawBlobs {
		b, err := blob.Decode(rwb)
		if err != nil {
			inslog.Error(err, "heavyserver: deserialize blob failed")
			continue
		}

		blobID := object.CalculateIDForBlob(pcs, pn, b.Value)
		err = blobs.Set(ctx, *blobID, *b)
		if err != nil {
			inslog.Error(err, "heavyserver: blob storing failed")
			continue
		}
	}
}

func storeRecords(
	ctx context.Context,
	records object.RecordModifier,
	pcs insolar.PlatformCryptographyScheme,
	pn insolar.PulseNumber,
	rawRecords [][]byte,
) {
	inslog := inslogger.FromContext(ctx)

	for _, rawRec := range rawRecords {
		rec, err := object.DecodeMaterial(rawRec)
		if err != nil {
			inslog.Error(err, "heavyserver: deserialize record failed")
			continue
		}

		virtRec := rec.Record

		id := object.NewRecordIDFromRecord(pcs, pn, virtRec)
		err = records.Set(ctx, *id, rec)
		if err != nil {
			inslog.Error(err, "heavyserver: store record failed")
			continue
		}
	}
}
