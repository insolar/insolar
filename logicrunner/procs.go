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
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/requestresult"
)

// ------------- CheckOurRole

type CheckOurRole struct {
	target      insolar.Reference
	role        insolar.DynamicRole
	pulseNumber insolar.PulseNumber

	jetCoordinator jet.Coordinator
}

var ErrCantExecute = errors.New("can't executeAndReply this object")

func (ch *CheckOurRole) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, "CheckOurRole")
	defer span.End()

	// TODO do map of supported objects for pulse, go to jetCoordinator only if map is empty for ref
	isAuthorized, err := ch.jetCoordinator.IsMeAuthorizedNow(ctx, ch.role, *ch.target.GetLocal())
	if err != nil {
		return errors.Wrap(err, "authorization failed with error")
	}
	if !isAuthorized {
		return ErrCantExecute
	}
	return nil
}

// ------------- RegisterIncomingRequest

type RegisterIncomingRequest struct {
	request record.IncomingRequest

	result chan *payload.RequestInfo

	ArtifactManager artifacts.Client
}

func NewRegisterIncomingRequest(request record.IncomingRequest, dep *Dependencies) *RegisterIncomingRequest {
	return &RegisterIncomingRequest{
		request:         request,
		ArtifactManager: dep.ArtifactManager,
		result:          make(chan *payload.RequestInfo, 1),
	}
}

func (r *RegisterIncomingRequest) setResult(result *payload.RequestInfo) { // nolint
	r.result <- result
}

// getResult is blocking
func (r *RegisterIncomingRequest) getResult() *payload.RequestInfo { // nolint
	return <-r.result
}

func (r *RegisterIncomingRequest) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, "RegisterIncomingRequest.Proceed")
	defer span.End()

	inslogger.FromContext(ctx).Debug("registering incoming request")

	reqInfo, err := r.ArtifactManager.RegisterIncomingRequest(ctx, &r.request)
	if err != nil {
		return err
	}

	r.setResult(reqInfo)
	return nil
}

type AddFreshRequest struct {
	broker     ExecutionBrokerI
	requestRef insolar.Reference
	request    record.IncomingRequest
}

func (c *AddFreshRequest) Proceed(ctx context.Context) error {
	tr := common.NewTranscriptCloneContext(ctx, c.requestRef, c.request)
	c.broker.AddFreshRequest(ctx, tr)
	return nil
}

type RecordErrorResult struct {
	artifactManager artifacts.Client

	err        error
	objectRef  insolar.Reference
	requestRef insolar.Reference

	result []byte
}

func (r *RecordErrorResult) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, "RecordErrorResult.Proceed")
	defer span.End()

	inslogger.FromContext(ctx).Debug("recording error result")

	resultWithErr, err := foundation.MarshalMethodErrorResult(r.err)
	if err != nil {
		return errors.Wrap(err, "couldn't marshal result")
	}

	result := requestresult.New(resultWithErr, r.objectRef)

	err = r.artifactManager.RegisterResult(ctx, r.requestRef, result)
	if err != nil {
		return errors.Wrap(err, "couldn't register result")
	}

	r.result = resultWithErr

	return nil
}
