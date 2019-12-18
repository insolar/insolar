package sm_execute_request

import (
	"github.com/insolar/insolar/conveyor/injector"
	"github.com/insolar/insolar/conveyor/smachine"
)

type SMPreExecuteImmutable struct {
	smachine.StateMachineDeclTemplate

	*ExecuteIncomingCommon
}

/* -------- Declaration ------------- */

func (s *SMPreExecuteImmutable) GetInitStateFor(smachine.StateMachine) smachine.InitFunc {
	return s.Init
}

func (s *SMPreExecuteImmutable) InjectDependencies(_ smachine.StateMachine, _ smachine.SlotLink, injector *injector.DependencyInjector) {
}

/* -------- Instance ------------- */

func (s *SMPreExecuteImmutable) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return s
}

func (s *SMPreExecuteImmutable) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	return ctx.Jump(s.stepTakeLock)
}

func (s *SMPreExecuteImmutable) stepTakeLock(ctx smachine.ExecutionContext) smachine.StateUpdate {
	if s.RequestInfo.Result != nil {
		return ctx.Jump(s.stepSpawnExecution)
	}

	if !ctx.AcquireAndRelease(s.objectInfo.ImmutableExecute) {
		return ctx.Sleep().ThenRepeat()
	}

	return ctx.Jump(s.stepSpawnExecution)
}

func (s *SMPreExecuteImmutable) stepSpawnExecution(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &SMExecute{
			ExecuteIncomingCommon: s.ExecuteIncomingCommon,
		}
	})
}
