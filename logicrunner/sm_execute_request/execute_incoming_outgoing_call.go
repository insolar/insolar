package sm_execute_request

import (
	"fmt"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/conveyor"
	"github.com/insolar/insolar/conveyor/smachine"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/log/logcommon"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/requestresult"
	"github.com/insolar/insolar/logicrunner/s_artifact"
	"github.com/insolar/insolar/logicrunner/s_contract_requester"
	"github.com/insolar/insolar/logicrunner/s_contract_runner"
	"github.com/insolar/insolar/logicrunner/s_contract_runner/outgoing"
)

// Embedded StateMachine for processing Outgoing Calls
type ESMOutgoingCallProcess struct {
	// dependencies
	ContractRequester *s_contract_requester.ContractRequesterServiceAdapter
	ArtifactManager   *s_artifact.ArtifactClientServiceAdapter
	pulseSlot         *conveyor.PulseSlot

	// input
	contractTranscript common.Transcript
	object             insolar.Reference

	outgoingEvent *s_contract_runner.ContractExecutionStateUpdate
	code          []byte
	deactivate    bool
	externalError error

	outgoing            record.OutgoingRequest
	outgoingResult      *record.Result
	outgoingReply       insolar.Reply
	outgoingRequestInfo *common.ParsedRequestInfo

	continueExecutionStepCallback func(map[string]interface{}) smachine.StateFunc
}

func (s *ESMOutgoingCallProcess) Prepare(transcript common.Transcript, object insolar.Reference) {
	s.contractTranscript = transcript
	s.object = object
}

func (s *ESMOutgoingCallProcess) ProcessOutgoing(
	ctx smachine.ExecutionContext,
	outgoingExecutionType *s_contract_runner.ContractExecutionStateUpdate,
	continueExecutionStepCallback func(map[string]interface{}) smachine.StateFunc,
) smachine.StateUpdate {
	s.outgoingEvent = outgoingExecutionType
	s.continueExecutionStepCallback = continueExecutionStepCallback

	switch s.outgoingEvent.Outgoing.(type) {
	case outgoing.DeactivateEvent:
		return ctx.Jump(s.stepFinishDeactivate)

	case outgoing.GetCodeEvent:
		return ctx.Jump(s.stepGetCode)

	case outgoing.RouteCallEvent, outgoing.SaveAsChildEvent:
		return ctx.Jump(s.stepOutgoingRegister)

	default:
		panic(fmt.Sprintf("unknown type of event %T", s.outgoingEvent.Outgoing))
	}
}

func (s *ESMOutgoingCallProcess) stepFinishDeactivate(ctx smachine.ExecutionContext) smachine.StateUpdate {
	nextStep := s.continueExecutionStepCallback(map[string]interface{}{
		"deactivate": true,
	})
	return ctx.Jump(nextStep)
}

func (s *ESMOutgoingCallProcess) stepGetCode(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		goCtx = ctx.GetContext()
		event = s.outgoingEvent.Outgoing.(outgoing.GetCodeEvent)
	)

	return s.ArtifactManager.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		desc, err := svc.GetCode(goCtx, event.CodeReference)

		return func(ctx smachine.AsyncResultContext) {
			s.externalError = err
			if err != nil {
				s.code, s.externalError = desc.Code()
			}
		}
	}).DelayedStart().Sleep().ThenJump(s.stepFinishGetCode)
}

func (s *ESMOutgoingCallProcess) stepFinishGetCode(ctx smachine.ExecutionContext) smachine.StateUpdate {
	nextStep := s.continueExecutionStepCallback(map[string]interface{}{
		"error": s.externalError,
		"code":  s.code,
	})
	return ctx.Jump(nextStep)
}

func (s *ESMOutgoingCallProcess) stepOutgoingRegister(ctx smachine.ExecutionContext) smachine.StateUpdate {
	s.outgoing = s.outgoingEvent.Outgoing.(outgoing.RPCOutgoingConstructor).ConstructOutgoing(s.contractTranscript)

	var (
		outgoingRequest = &s.outgoing
		goCtx           = ctx.GetContext()
	)

	return s.ArtifactManager.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		info, err := svc.RegisterOutgoingRequest(goCtx, outgoingRequest)

		return func(ctx smachine.AsyncResultContext) {
			if err != nil {
				s.externalError = err
				return
			}

			s.outgoingRequestInfo, s.externalError = common.NewParsedRequestInfo(outgoingRequest, info)
			if s.externalError != nil {
				return
			}

			if _, ok := s.outgoingRequestInfo.Request.(*record.OutgoingRequest); !ok {
				s.externalError = errors.Errorf("unexpected request type: %T", s.outgoingRequestInfo.Request)
				return
			}

			if s.outgoingRequestInfo.Result != nil {
				s.outgoingResult = s.outgoingRequestInfo.Result
				s.outgoingReply, err = reply.UnmarshalFromMeta(s.outgoingRequestInfo.Result.Payload)
				if err != nil {
					s.externalError = errors.Wrap(err, "failed to unmarshal reply")
					return
				}
			}
		}
	}).DelayedStart().Sleep().ThenJump(s.stepOutgoingExecute)
}

type ContractCallBefore struct {
	*logcommon.LogObjectTemplate `txt:"before external call"`

	Method string
	Object string
}

type ContractCallAfter struct {
	*logcommon.LogObjectTemplate `txt:"after contract call"`

	CallResultType   insolar.Reply `fmte:"%T"`
	RequestReference string
	Method           string
	Error            error
}

func (s *ESMOutgoingCallProcess) stepOutgoingExecute(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.outgoingRequestInfo.Request.IsDetachedCall() || s.outgoingReply != nil || s.externalError != nil {
		ctx.Log().Trace(fmt.Sprintf("IsDetachedCall: %v", s.outgoingRequestInfo.Request.IsDetachedCall()))
		ctx.Log().Trace(fmt.Sprintf("OutgoingReply: %v", s.outgoingReply))
		ctx.Log().Trace(fmt.Sprintf("ExternalError: %v", s.externalError))
		// nextStep := s.continueExecutionStepCallback(map[string]interface{}{
		// 	"result": s.outgoingResult,
		// 	"error":  s.externalError,
		// })
		// return ctx.Jump(nextStep)
		// saga call, deduplicated outgoing reply OR error while registering outgoing request
		return ctx.Jump(s.stepFinishOutgoing)
	}

	var (
		goCtx       = ctx.GetContext()
		incoming    = outgoing.BuildIncomingRequestFromOutgoing(&s.outgoing)
		pulseNumber = s.pulseSlot.PulseData().PulseNumber
		pl          = &payload.CallMethod{Request: incoming, PulseNumber: pulseNumber}
	)

	ctx.SetLogTracing(true)
	logger := ctx.LogAsync()

	return s.ContractRequester.PrepareAsync(ctx, func(svc s_contract_requester.ContractRequesterService) smachine.AsyncResultFunc {
		var (
			objectReferenceString  string
			requestReferenceString string
		)

		if pl.Request.Object != nil {
			objectReferenceString = pl.Request.Object.String()
		}

		logger.Trace(ContractCallBefore{
			Method: pl.Request.Method,
			Object: objectReferenceString,
		})

		callResult, requestReference, err := svc.SendRequest(goCtx, pl)
		if requestReference != nil {
			requestReferenceString = requestReference.String()
		}

		logger.Trace(ContractCallAfter{
			Method:           pl.Request.Method,
			CallResultType:   callResult,
			Error:            err,
			RequestReference: requestReferenceString,
		})

		return func(ctx smachine.AsyncResultContext) {
			s.externalError = err
			s.outgoingReply = callResult
		}
	}).DelayedStart().Sleep().ThenJump(s.stepOutgoingSaveResult)
}

func (s *ESMOutgoingCallProcess) stepOutgoingSaveResult(ctx smachine.ExecutionContext) smachine.StateUpdate {
	var (
		goCtx            = ctx.GetContext()
		requestReference = s.outgoingRequestInfo.RequestReference
		caller           = s.outgoingRequestInfo.Request.(*record.OutgoingRequest).Caller
		result           []byte
	)

	switch v := s.outgoingReply.(type) {
	case *reply.CallMethod: // regular call
		result = v.Result
		s.outgoingResult = &record.Result{
			Object:  *s.object.GetLocal(),
			Request: requestReference,
			Payload: v.Result,
		}

	default:
		s.externalError = errors.Errorf("contractRequester.Call returned unexpected type %T", s.outgoingReply)
		return ctx.Jump(s.stepFinishOutgoing)
		// nextStep := s.continueExecutionStepCallback(map[string]interface{}{
		// 	"result": (*record.Result)(nil),
		// 	"error":
		// })
		// return ctx.Jump(nextStep)
	}

	// Register result of the outgoing method
	requestResult := requestresult.New(result, caller)

	ctx.Log().Trace("Saving request result")
	return s.ArtifactManager.PrepareAsync(ctx, func(svc s_artifact.ArtifactClientService) smachine.AsyncResultFunc {
		err := svc.RegisterResult(goCtx, requestReference, requestResult)

		return func(ctx smachine.AsyncResultContext) {
			if err != nil {
				s.externalError = errors.Wrap(err, "failed to register result")
			}
		}
	}).DelayedStart().Sleep().ThenJump(s.stepFinishOutgoing)
}

func (s *ESMOutgoingCallProcess) stepFinishOutgoing(ctx smachine.ExecutionContext) smachine.StateUpdate {
	nextStep := s.continueExecutionStepCallback(map[string]interface{}{
		"result": s.outgoingResult,
		"error":  s.externalError,
	})
	return ctx.Jump(nextStep)
}

func (s *ESMOutgoingCallProcess) Reset() {
	s.code = nil
	s.externalError = nil
	s.outgoing = record.OutgoingRequest{}
	s.outgoingResult = nil
	s.outgoingReply = nil
	s.outgoingRequestInfo = nil
	s.continueExecutionStepCallback = nil
	s.outgoingEvent = nil
}
