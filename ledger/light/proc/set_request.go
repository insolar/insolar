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

package proc

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/light/hot"
	"github.com/insolar/insolar/ledger/object"
	"github.com/pkg/errors"
)

type SetRequest struct {
	message   payload.Meta
	request   record.Request
	requestID insolar.ID
	jetID     insolar.JetID

	dep struct {
		writer      hot.WriteAccessor
		filament    executor.FilamentModifier
		sender      bus.Sender
		locker      object.IndexLocker
		index       object.IndexStorage
		coordinator jet.Coordinator
	}
}

func NewSetRequest(
	msg payload.Meta,
	rec record.Request,
	recID insolar.ID,
	jetID insolar.JetID,
) *SetRequest {
	return &SetRequest{
		message:   msg,
		request:   rec,
		requestID: recID,
		jetID:     jetID,
	}
}

func (p *SetRequest) Dep(
	w hot.WriteAccessor,
	f executor.FilamentModifier,
	s bus.Sender,
	l object.IndexLocker,
	i object.IndexStorage,
	c jet.Coordinator,
) {
	p.dep.writer = w
	p.dep.filament = f
	p.dep.sender = s
	p.dep.locker = l
	p.dep.index = i
	p.dep.coordinator = c
}

func (p *SetRequest) Proceed(ctx context.Context) error {
	logger := inslogger.FromContext(ctx).WithField("request_id", p.requestID.DebugString())
	logger.Debug("trying to save request")

	done, err := p.dep.writer.Begin(ctx, flow.Pulse(ctx))
	if err != nil {
		if err == hot.ErrWriteClosed {
			return flow.ErrCancelled
		}
		return err
	}
	defer done()

	var objectID insolar.ID
	if p.request.IsCreationRequest() {
		objectID = p.requestID
	} else {
		objectID = *p.request.AffinityRef().Record()
		idx, err := p.dep.index.ForID(ctx, flow.Pulse(ctx), objectID)
		if err != nil {
			return errors.Wrap(err, "failed to check an object state")
		}
		if idx.Lifeline.StateID == record.StateDeactivation {
			msg, err := payload.NewMessage(&payload.Error{Text: "object is deactivated", Code: payload.CodeDeactivated})
			if err != nil {
				return errors.Wrap(err, "failed to create reply")
			}

			p.dep.sender.Reply(ctx, p.message, msg)
			return nil
		}
	}

	virtNode, err := p.dep.coordinator.VirtualExecutorForObject(ctx, objectID, flow.Pulse(ctx))
	if err != nil {
		return err
	}
	if p.message.Sender != *virtNode {
		logger.Errorf("sender isn't the executor. sender - %v, executor - %v", p.message.Sender, *virtNode)
		//return ErrExecutorMismatch
	}

	p.dep.locker.Lock(objectID)
	defer p.dep.locker.Unlock(objectID)

	req, res, err := p.dep.filament.SetRequest(ctx, p.requestID, p.jetID, p.request)
	if err != nil {
		return errors.Wrap(err, "failed to store record")
	}

	var (
		reqBuf []byte
		resBuf []byte
	)

	requestID := p.requestID
	if req != nil {
		reqBuf, err = req.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal stored record")
		}
		requestID = req.RecordID
	}

	if res != nil {
		resBuf, err = res.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal stored record")
		}
	}

	msg, err := payload.NewMessage(&payload.RequestInfo{
		ObjectID:  objectID,
		RequestID: requestID,
		Request:   reqBuf,
		Result:    resBuf,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	logger.WithFields(map[string]interface{}{
		"duplicate":   req != nil,
		"has_result":  res != nil,
		"is_creation": p.request.IsCreationRequest(),
	}).Debug("request saved")
	p.dep.sender.Reply(ctx, p.message, msg)
	return nil
}
