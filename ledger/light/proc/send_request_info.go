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

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"
)

type SendRequestInfo struct {
	message   payload.Meta
	objectID  insolar.ID
	requestID insolar.ID
	pulse     insolar.PulseNumber

	dep struct {
		filament executor.FilamentCalculator
		sender   bus.Sender
		locker   object.IndexLocker
	}
}

func NewSendRequestInfo(
	msg payload.Meta,
	objectID insolar.ID,
	requestID insolar.ID,
	pulse insolar.PulseNumber,
) *SendRequestInfo {
	return &SendRequestInfo{
		message:   msg,
		objectID:  objectID,
		requestID: requestID,
		pulse:     pulse,
	}
}

func (p *SendRequestInfo) Dep(
	filament executor.FilamentCalculator,
	sender bus.Sender,
	locker object.IndexLocker,
) {
	p.dep.filament = filament
	p.dep.sender = sender
	p.dep.locker = locker
}

func (p *SendRequestInfo) Proceed(ctx context.Context) error {
	if p.requestID.IsEmpty() {
		return errors.New("requestID is empty")
	}
	if p.objectID.IsEmpty() {
		return errors.New("objectID is empty")
	}
	if p.pulse < pulse.MinTimePulse {
		return errors.New("pulse is wrong")
	}

	logger := inslogger.FromContext(ctx).WithFields(map[string]interface{}{
		"request_id":     p.requestID.DebugString(),
		"object_id":      p.objectID.DebugString(),
		"pulse_received": p.pulse,
	})
	logger.Debug("send request info started")

	// Prevent concurrent object modifications.
	p.dep.locker.Lock(p.objectID)
	defer p.dep.locker.Unlock(p.objectID)

	var (
		reqBuf []byte
		resBuf []byte
	)
	foundRequest, foundResult, err := p.dep.filament.RequestInfo(ctx, p.objectID, p.requestID, p.pulse)
	if err != nil {
		return errors.Wrap(err, "failed to get request info")
	}

	if foundRequest != nil {
		reqBuf, err = foundRequest.Record.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal request record")
		}
	}
	if foundResult != nil {
		resBuf, err = foundResult.Record.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal result record")
		}
	}

	msg, err := payload.NewMessage(&payload.RequestInfo{
		ObjectID:  p.objectID,
		RequestID: p.requestID,
		Request:   reqBuf,
		Result:    resBuf,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	p.dep.sender.Reply(ctx, p.message, msg)

	logger.WithFields(map[string]interface{}{
		"request":    foundRequest != nil,
		"has_result": foundResult != nil,
	}).Debug("send request info finished")
	return nil
}
