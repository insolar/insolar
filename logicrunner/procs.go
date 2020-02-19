// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logicrunner

import (
	"context"

	"github.com/pkg/errors"
	"go.opencensus.io/stats"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/instrumentation/instracer"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/metrics"
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
	defer span.Finish()

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
	defer span.Finish()

	reqInfo, err := r.ArtifactManager.RegisterIncomingRequest(ctx, &r.request)
	if err != nil {
		return err
	}

	r.setResult(reqInfo)
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
	defer span.Finish()

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

func ProcessLogicalError(ctx context.Context, err error) bool {
	e, ok := errors.Cause(err).(*payload.CodedError)
	if ok {
		switch e.Code {
		case payload.CodeNotFound:
			return true
		case payload.CodeLoopDetected:
			stats.Record(ctx, metrics.CallMethodLoopDetected.M(1))
			return true
		}
	}
	return false
}
