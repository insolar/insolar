package sm_execute_request

import (
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
)

type SMPreExecuteMutable struct {
	smachine.StateMachineDeclTemplate

	*ExecuteIncomingCommon
}

/* -------- Declaration ------------- */

func (s *SMPreExecuteMutable) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *SMPreExecuteMutable) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
}

/* -------- Instance ------------- */

func (s *SMPreExecuteMutable) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *SMPreExecuteMutable) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepWaitObjectIsReady)
}

func (s *SMPreExecuteMutable) stepWaitObjectIsReady(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.RequestInfo.Result != nil {
		return ctx.Jump(s.stepSpawnExecution)
	}

	if !ctx.AcquireForThisStep(s.objectInfo.ReadyToWork) {
		return ctx.Sleep().ThenRepeat()
	}

	return ctx.Jump(s.stepTakeLock)
}

func (s *SMPreExecuteMutable) stepTakeLock(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if !ctx.AcquireAndRelease(s.objectInfo.MutableExecute) {
		return ctx.Sleep().ThenRepeat()
	}

	return ctx.Jump(s.stepSpawnExecution)
}

func (s *SMPreExecuteMutable) stepCheckOrdering(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// passing right now
	return ctx.Jump(s.stepSpawnExecution)
}

func (s *SMPreExecuteMutable) stepSpawnExecution(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &SMExecute{
			ExecuteIncomingCommon: s.ExecuteIncomingCommon,
		}
	})
}
