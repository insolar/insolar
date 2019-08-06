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

package procs

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
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/executionbroker"
	"github.com/insolar/insolar/logicrunner/transcript"
)

// ------------- CheckOurRole

type CheckOurRole struct {
	Msg         insolar.Message
	Role        insolar.DynamicRole
	PulseNumber insolar.PulseNumber

	JetCoordinator jet.Coordinator
}

var ErrCantExecute = errors.New("can't executeAndReply this object")

func (ch *CheckOurRole) Proceed(ctx context.Context) error {
	ctx, span := instracer.StartSpan(ctx, "CheckOurRole")
	defer span.End()

	// TODO do map of supported objects for pulse, go to jetCoordinator only if map is empty for ref
	target := ch.Msg.DefaultTarget()
	isAuthorized, err := ch.JetCoordinator.IsMeAuthorizedNow(ctx, ch.Role, *target.Record())
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

func NewRegisterIncomingRequest(request record.IncomingRequest, artifactClient artifacts.Client) *RegisterIncomingRequest {
	return &RegisterIncomingRequest{
		request:         request,
		ArtifactManager: artifactClient,
		result:          make(chan *payload.RequestInfo, 1),
	}
}

func (r *RegisterIncomingRequest) SetResult(result *payload.RequestInfo) { // nolint
	r.result <- result
}

// Result is blocking
func (r *RegisterIncomingRequest) Result() *payload.RequestInfo { // nolint
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

	r.SetResult(reqInfo)
	return nil
}

type AddFreshRequest struct {
	Broker     executionbroker.BrokerI
	RequestRef insolar.Reference
	Request    record.IncomingRequest
}

func (c *AddFreshRequest) Proceed(ctx context.Context) error {
	requestCtx := common.FreshContextFromContext(ctx)
	tr := transcript.NewTranscript(requestCtx, c.RequestRef, c.Request)
	c.Broker.AddFreshRequest(ctx, tr)
	return nil
}
