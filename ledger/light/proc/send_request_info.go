// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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

	var (
		reqBuf []byte
		resBuf []byte
	)
	foundReqInfo, err := p.dep.filament.RequestInfo(ctx, p.objectID, p.requestID, p.pulse)
	if err != nil {
		return errors.Wrap(err, "failed to get request info")
	}

	reqBuf, err = foundReqInfo.Request.Record.Marshal()
	if err != nil {
		return errors.Wrap(err, "failed to marshal request record")
	}

	if foundReqInfo.Result != nil {
		resBuf, err = foundReqInfo.Result.Record.Marshal()
		if err != nil {
			return errors.Wrap(err, "failed to marshal result record")
		}
	}

	msg, err := payload.NewMessage(&payload.RequestInfo{
		ObjectID:      p.objectID,
		RequestID:     p.requestID,
		Request:       reqBuf,
		Result:        resBuf,
		OldestMutable: foundReqInfo.OldestMutable,
	})
	if err != nil {
		return errors.Wrap(err, "failed to create reply")
	}

	p.dep.sender.Reply(ctx, p.message, msg)

	logger.WithFields(map[string]interface{}{
		"request":    foundReqInfo.Request != nil,
		"has_result": foundReqInfo.Result != nil,
	}).Debug("send request info finished")
	return nil
}
