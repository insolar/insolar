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

package logicrunner

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner.RequestsExecutor -o ./ -s _mock.go

type RequestsExecutor interface {
	ExecuteAndSave(ctx context.Context, current *Transcript) (insolar.Reply, error)
	Execute(ctx context.Context, current *Transcript) (*RequestResult, error)
	Save(ctx context.Context, current *Transcript, res *RequestResult) (insolar.Reply, error)
	SendReply(ctx context.Context, current *Transcript, re insolar.Reply, err error)
}

type requestsExecutor struct {
	MessageBus      insolar.MessageBus  `inject:""`
	NodeNetwork     insolar.NodeNetwork `inject:""`
	LogicExecutor   LogicExecutor       `inject:""`
	ArtifactManager artifacts.Client    `inject:""`
}

func NewRequestsExecutor() RequestsExecutor {
	return &requestsExecutor{}
}

func (e *requestsExecutor) ExecuteAndSave(
	ctx context.Context, transcript *Transcript,
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

	return repl, nil
}

func (e *requestsExecutor) Execute(
	ctx context.Context, transcript *Transcript,
) (
	*RequestResult, error,
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
	ctx context.Context, transcript *Transcript, res *RequestResult,
) (
	insolar.Reply, error,
) {
	inslogger.FromContext(ctx).Debug("Saving result")

	am := e.ArtifactManager
	request := transcript.Request
	reqRef := transcript.RequestRef

	switch {
	case res.Activation:
		err := e.ArtifactManager.ActivateObject(
			ctx, reqRef, *request.Base, *request.Prototype,
			request.CallType == record.CTSaveAsDelegate,
			res.NewMemory,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't activate object")
		}
		return &reply.CallConstructor{Object: &reqRef}, err
	case res.Deactivation:
		err := am.DeactivateObject(
			ctx, reqRef, transcript.ObjectDescriptor, res.Result,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't deactivate object")
		}
		return &reply.CallMethod{Result: res.Result}, nil
	case res.NewMemory != nil:
		err := am.UpdateObject(
			ctx, reqRef, transcript.ObjectDescriptor, res.NewMemory, res.Result,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't update object")
		}
		return &reply.CallMethod{Result: res.Result}, nil
	default:
		_, err := am.RegisterResult(ctx, *request.Object, reqRef, res.Result)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't save results")
		}
		return &reply.CallMethod{Result: res.Result}, nil
	}
}

func (e *requestsExecutor) SendReply(
	ctx context.Context, transcript *Transcript, re insolar.Reply, err error,
) {
	if transcript.Request.ReturnMode != record.ReturnResult {
		return
	}

	inslogger.FromContext(ctx).Debug("Returning result")

	target := transcript.Request.Sender

	errstr := ""
	if err != nil {
		errstr = err.Error()
	}
	_, err = e.MessageBus.Send(
		ctx,
		&message.ReturnResults{
			Target:     target,
			RequestRef: transcript.RequestRef,
			Reply:      re,
			Error:      errstr,
		},
		&insolar.MessageSendOptions{
			Receiver: &target,
		},
	)
	if err != nil {
		inslogger.FromContext(ctx).Error("couldn't deliver results: ", err)
	}
}
