// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/logicexecutor"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.RequestsExecutor -o ./ -s _mock.go -g

type RequestsExecutor interface {
	ExecuteAndSave(ctx context.Context, current *common.Transcript) (artifacts.RequestResult, error)
	Execute(ctx context.Context, current *common.Transcript) (artifacts.RequestResult, error)
	Save(ctx context.Context, current *common.Transcript, res artifacts.RequestResult) error
	SendReply(ctx context.Context, reqRef insolar.Reference, req record.IncomingRequest, re insolar.Reply, err error)
}

type requestsExecutor struct {
	Sender          bus.Sender                  `inject:""`
	LogicExecutor   logicexecutor.LogicExecutor `inject:""`
	ArtifactManager artifacts.Client            `inject:""`
	PulseAccessor   pulse.Accessor              `inject:""`
}

func NewRequestsExecutor() RequestsExecutor {
	return &requestsExecutor{}
}

func (e *requestsExecutor) ExecuteAndSave(
	ctx context.Context, transcript *common.Transcript,
) (
	artifacts.RequestResult, error,
) {
	ctx, span := instracer.StartSpan(ctx, "RequestsExecutor.ExecuteAndSave")
	defer span.Finish()

	result, err := e.Execute(ctx, transcript)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't execute request")
	}

	err = e.Save(ctx, transcript, result)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't save request result")
	}

	return result, nil
}

func (e *requestsExecutor) Execute(
	ctx context.Context, transcript *common.Transcript,
) (
	artifacts.RequestResult, error,
) {
	ctx, span := instracer.StartSpan(ctx, "LogicRunner.executeLogic")
	defer span.Finish()

	inslogger.FromContext(ctx).Debug("executing request")

	result, err := e.LogicExecutor.Execute(ctx, transcript)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't execute method")
	}

	return result, nil
}

func (e *requestsExecutor) Save(
	ctx context.Context, transcript *common.Transcript, res artifacts.RequestResult,
) error {
	inslogger.FromContext(ctx).Debug("registering IncomingRequest result")

	err := e.ArtifactManager.RegisterResult(ctx, transcript.RequestRef, res)
	if err != nil {
		return errors.Wrapf(err, "couldn't save result with %s side effect", res.Type().String())
	}

	inslogger.FromContext(ctx).Debug("registered IncomingRequest result")

	return nil
}

func (e *requestsExecutor) SendReply(
	ctx context.Context,
	reqRef insolar.Reference, req record.IncomingRequest,
	re insolar.Reply, err error,
) {
	logger := inslogger.FromContext(ctx)
	if rm := req.ReturnMode; rm != record.ReturnResult {
		logger.Debug(
			"Not sending result, return mode: ", rm.String(),
		)
		return
	}

	if reqRef.IsEmpty() {
		logger.Error(
			"Not sending result, empty request reference, request: ", req.String(),
		)
		return
	}

	logger.Debug("returning result")

	var errStr string
	if err != nil {
		errStr = err.Error()
	}

	var replyBytes []byte
	if re != nil {
		replyBytes = reply.ToBytes(re)
	}

	if len(errStr) == 0 && len(replyBytes) == 0 {
		logger.Error("both reply and error are empty, this is wrong, not sending")
		return
	}

	if req.APINode.IsEmpty() {
		go e.sendToCaller(ctx, reqRef, req, replyBytes, errStr)
		return
	}

	go e.sendToAPINode(ctx, reqRef, req, replyBytes, errStr)
}

func (e *requestsExecutor) sendToCaller(
	ctx context.Context,
	reqRef insolar.Reference, req record.IncomingRequest,
	re []byte, errStr string,
) {
	sender := bus.NewWaitOKWithRetrySender(e.Sender, e.PulseAccessor, 1)

	msg, err := payload.NewResultMessage(&payload.ReturnResults{
		Target:     req.Caller,
		RequestRef: reqRef,
		Reason:     req.Reason,
		Reply:      re,
		Error:      errStr,
	})
	if err != nil {
		inslogger.FromContext(ctx).Error("couldn't serialize message: ", err)
		return
	}

	inslogger.FromContext(ctx).Debug("sending to caller")

	sender.SendRole(ctx, msg, insolar.DynamicRoleVirtualExecutor, req.Caller)
}

func (e *requestsExecutor) sendToAPINode(
	ctx context.Context,
	reqRef insolar.Reference, req record.IncomingRequest,
	re []byte, errStr string,
) {
	sender := bus.NewWaitOKWithRetrySender(e.Sender, e.PulseAccessor, 1)

	msg, err := payload.NewResultMessage(&payload.ReturnResults{
		RequestRef: reqRef,
		Reply:      re,
		Error:      errStr,
	})
	if err != nil {
		inslogger.FromContext(ctx).Error("couldn't serialize message: ", err)
		return
	}

	inslogger.FromContext(ctx).Debug("sending to APINode")

	sender.SendTarget(ctx, msg, req.APINode)
}
