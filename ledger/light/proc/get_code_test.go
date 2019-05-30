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

package proc_test

import (
	"context"
	"testing"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/reply"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/stretchr/testify/require"
)

func TestGetCode_Proceed(t *testing.T) {
	mc := minimock.NewController(t)
	a := require.New(t)
	ctx := inslogger.TestContext(t)

	replyChan := make(chan *message.Message, 1)
	requestMessage := message.NewMessage(watermill.NewUUID(), nil)
	blobValue := blob.Blob{Value: []byte{1, 2, 3}}
	blobID := gen.ID()
	codeRec := codeRecord(blobID)
	codeRef := gen.Reference()
	getCode := proc.NewGetCode(codeRef, requestMessage)
	records := object.NewRecordAccessorMock(mc)
	records.ForIDFunc = func(c context.Context, id insolar.ID) (record.Material, error) {
		a.Equal(*codeRef.Record(), id)
		return codeRec, nil
	}
	blobs := blob.NewAccessorMock(mc)
	blobs.ForIDFunc = func(c context.Context, id insolar.ID) (blob.Blob, error) {
		a.Equal(blobID, id)
		return blobValue, nil
	}
	sender := bus.NewSenderMock(mc)
	sender.ReplyFunc = func(p context.Context, p1 *message.Message, p2 *message.Message) {
		replyChan <- p2
	}
	getCode.Dep.RecordAccessor = records
	getCode.Dep.BlobAccessor = blobs
	getCode.Dep.Sender = sender

	err := getCode.Proceed(ctx)
	a.NoError(err)

	unwrappedCodeRec := record.Unwrap(codeRec.Virtual)

	expectedRep := bus.ReplyAsMessage(ctx, &reply.Code{
		Code:        blobValue.Value,
		MachineType: unwrappedCodeRec.(*record.Code).MachineType,
	})
	rep := <-replyChan
	a.Equal(expectedRep.Metadata, rep.Metadata)
	a.Equal(expectedRep.Payload, rep.Payload)
}

func codeRecord(codeID insolar.ID) record.Material {
	return record.Material{
		Virtual: &record.Virtual{
			Union: &record.Virtual_Code{
				Code: &record.Code{
					Code:        codeID,
					MachineType: insolar.MachineTypeBuiltin,
				},
			},
		},
	}
}
