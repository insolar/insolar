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

package requestsexecutor

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/logicexecutor"
	"github.com/insolar/insolar/logicrunner/transcript"
	"github.com/insolar/insolar/messagebus"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner/requestsexecutor.RequestsExecutor -o ./ -s _mock.go -g

type RequestsExecutor interface {
	ExecuteAndSave(ctx context.Context, current *transcript.Transcript) (insolar.Reply, error)
	Execute(ctx context.Context, current *transcript.Transcript) (artifacts.RequestResult, error)
	Save(ctx context.Context, current *transcript.Transcript, res artifacts.RequestResult) (insolar.Reply, error)
	SendReply(ctx context.Context, current *transcript.Transcript, re insolar.Reply, err error)
}

type requestsExecutor struct {
	MessageBus      insolar.MessageBus          `inject:""`
	LogicExecutor   logicexecutor.LogicExecutor `inject:""`
	ArtifactManager artifacts.Client            `inject:""`
	PulseAccessor   pulse.Accessor              `inject:""`
}

func NewRequestsExecutor() RequestsExecutor {
	return &requestsExecutor{}
}

func (e *requestsExecutor) ExecuteAndSave(
	ctx context.Context, transcript *transcript.Transcript,
) (
	insolar.Reply, error,
) {
	ctx, span := instracer.StartSpan(ctx, "RequestsExecutor.ExecuteAndSave")
	defer span.End()

	result, err := e.Execute(ctx, transcript)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't execute request")
	}

	repl, err := e.Save(ctx, transcript, result)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't save request result")
	}

	inslogger.FromContext(ctx).Debug("saved result")

	return repl, nil
}

func (e *requestsExecutor) Execute(
	ctx context.Context, transcript *transcript.Transcript,
) (
	artifacts.RequestResult, error,
) {
	ctx, span := instracer.StartSpan(ctx, "LogicRunner.executeLogic")
	defer span.End()

	inslogger.FromContext(ctx).Debug("Executing request")

	if transcript.Request.CallType == record.CTMethod {
		objDesc, err := e.ArtifactManager.GetObject(ctx, *transcript.Request.Object)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get object")
		}
		transcript.ObjectDescriptor = objDesc
	}

	result, err := e.LogicExecutor.Execute(ctx, transcript)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't execute method")
	}

	return result, nil
}

func (e *requestsExecutor) Save(
	ctx context.Context, transcript *transcript.Transcript, res artifacts.RequestResult,
) (
	insolar.Reply, error,
) {
	inslogger.FromContext(ctx).Debug("Saving result")

	err := e.ArtifactManager.RegisterResult(ctx, transcript.RequestRef, res)
	if err != nil {
		return nil, errors.Wrapf(err, "couldn't save result with %s side effect", res.Type().String())
	}

	objRef := res.ObjectReference()
	return &reply.CallMethod{Result: res.Result(), Object: &objRef}, nil
}

func (e *requestsExecutor) SendReply(
	ctx context.Context, transcript *transcript.Transcript, re insolar.Reply, err error,
) {
	if rm := transcript.Request.ReturnMode; rm != record.ReturnResult {
		inslogger.FromContext(ctx).Debug(
			"Not sending result, return mode: ", rm.String(),
		)
		return
	}

	inslogger.FromContext(ctx).Debug("Returning result")

	errstr := ""
	if err != nil {
		errstr = err.Error()
	}
	sender := messagebus.BuildSender(
		e.MessageBus.Send,
		messagebus.RetryIncorrectPulse(e.PulseAccessor),
		messagebus.RetryFlowCancelled(e.PulseAccessor),
	)

	if transcript.Request.APINode.IsEmpty() {
		_, err = sender(
			ctx,
			&message.ReturnResults{
				Target:     transcript.Request.Caller,
				RequestRef: transcript.RequestRef,
				Reason:     transcript.Request.Reason,
				Reply:      re,
				Error:      errstr,
			},
			&insolar.MessageSendOptions{},
		)
	} else {
		_, err = sender(
			ctx,
			&message.ReturnResults{
				RequestRef: transcript.RequestRef,
				Reply:      re,
				Error:      errstr,
			},
			&insolar.MessageSendOptions{
				Receiver: &transcript.Request.APINode,
			},
		)
	}
	if err != nil {
		inslogger.FromContext(ctx).Error("couldn't deliver results: ", err)
	}
}
