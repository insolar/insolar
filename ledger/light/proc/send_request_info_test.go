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

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/proc"
	"github.com/insolar/insolar/ledger/object"
)

func TestSendRequestInfo_Proceed(t *testing.T) {
	ctx := inslogger.TestContext(t)
	mc := minimock.NewController(t)

	var (
		filament *executor.FilamentCalculatorMock
		sender   *bus.SenderMock
		locker   *object.IndexLockerMock
	)

	setup := func() {
		filament = executor.NewFilamentCalculatorMock(mc)
		sender = bus.NewSenderMock(mc)
		locker = object.NewIndexLockerMock(mc)
	}

	t.Run("basic fail", func(t *testing.T) {
		setup()
		defer mc.Finish()

		p := proc.NewSendRequestInfo(payload.Meta{}, gen.ID(), insolar.ID{}, insolar.FirstPulseNumber)
		p.Dep(filament, sender, locker)
		err := p.Proceed(ctx)
		assert.Error(t, err, "expected error 'requestID is empty'")

		p = proc.NewSendRequestInfo(payload.Meta{}, insolar.ID{}, gen.ID(), insolar.FirstPulseNumber)
		p.Dep(filament, sender, locker)
		err = p.Proceed(ctx)
		assert.Error(t, err, "expected error 'objectID is empty'")

		p = proc.NewSendRequestInfo(payload.Meta{}, gen.ID(), gen.ID(), insolar.FirstPulseNumber-10)
		p.Dep(filament, sender, locker)
		err = p.Proceed(ctx)
		assert.Error(t, err, "expected error 'pulse is wrong'")

	})

	t.Run("basic ok", func(t *testing.T) {
		setup()
		defer mc.Finish()

		reqID := gen.ID()
		resID := gen.ID()
		objID := reqID
		msg := payload.Meta{}

		request := record.CompositeFilamentRecord{
			Record:   record.Material{ID: reqID, ObjectID: objID},
			RecordID: reqID,
		}
		result := record.CompositeFilamentRecord{
			Record:   record.Material{ID: resID, ObjectID: objID},
			RecordID: resID,
		}
		reqBuf, err := request.Record.Marshal()
		resBuf, err := result.Record.Marshal()

		replyMsg, _ := payload.NewMessage(&payload.RequestInfo{
			ObjectID:  objID,
			RequestID: reqID,
			Request:   reqBuf,
			Result:    resBuf,
		})

		filament.RequestInfoMock.Return(&request, &result, nil)

		sender.ReplyMock.Inspect(func(ctx context.Context, origin payload.Meta, reply *message.Message) {
			assert.Equal(t, reply.Payload, replyMsg.Payload)
			assert.Equal(t, reply.Metadata, replyMsg.Metadata)
		}).Return()

		locker.UnlockMock.Return()
		locker.LockMock.Return()

		p := proc.NewSendRequestInfo(msg, objID, reqID, insolar.FirstPulseNumber)
		p.Dep(filament, sender, locker)
		err = p.Proceed(ctx)

		assert.NoError(t, err)
	})
}
