///
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
///

package conveyor

import "github.com/insolar/insolar/conveyor/smachine"

type pulseSMTemplate struct {
	smachine.StateMachineDeclTemplate
	ps  *PulseSlot
	psa *PulseServiceAdapter
}

type FuturePulseSM struct {
	pulseSMTemplate
}

func (sm *FuturePulseSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *FuturePulseSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	return sm.Init
}

func (sm *FuturePulseSM) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetMigration(sm.PulseCommitted)
	return ctx.Next(sm.StateWorking)
}

func (sm *FuturePulseSM) PulseCommitted(ctx smachine.MigrationContext) smachine.StateUpdate {
	if sm.ps.State() != Present {
		panic("unexpected state")
	}
	sm.ps.slotMachine.Migrate()
	return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &PresentPulseSM{sm.pulseSMTemplate}
	})
}

func (sm *FuturePulseSM) StateWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return sm.ps.processEvents(ctx, false)
}

type PresentPulseSM struct {
	pulseSMTemplate
}

func (sm *PresentPulseSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *PresentPulseSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	return sm.Init
}

func (sm *PresentPulseSM) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetMigration(sm.PulseCommitted)
	sm.psa.svc.subscribe(ctx, sm.PulsePrepare, sm.PulseCancel)
	return ctx.Next(sm.StateWorking)
}

func (sm *PresentPulseSM) PulsePrepare(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Next(sm.StateSuspending)
}

func (sm *PresentPulseSM) PulseCommitted(ctx smachine.MigrationContext) smachine.StateUpdate {
	if sm.ps.State() != Past {
		panic("unexpected state")
	}
	sm.ps.slotMachine.Migrate()
	return ctx.Replace(func(ctx smachine.ConstructionContext) smachine.StateMachine {
		return &PastPulseSM{sm.pulseSMTemplate}
	})
}

func (sm *PresentPulseSM) PulseCancel(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return ctx.Next(sm.StateWorking)
}

func (sm *PresentPulseSM) StateWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return sm.ps.processEventsAndOperations(ctx) // TODO here we have to somehow say that this state can be replaced
}

func (sm *PresentPulseSM) StateSuspending(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return sm.ps.processEvents(ctx, true)
}

type PastPulseSM struct {
	pulseSMTemplate
}

func (sm *PastPulseSM) GetStateMachineDeclaration() smachine.StateMachineDeclaration {
	return sm
}

func (sm *PastPulseSM) GetInitStateFor(machine smachine.StateMachine) smachine.InitFunc {
	if sm != machine {
		panic("illegal value")
	}
	return sm.Init
}

func (sm *PastPulseSM) Init(ctx smachine.InitializationContext) smachine.StateUpdate {
	ctx.SetMigration(sm.PulseCommitted)
	return ctx.Next(sm.StateWorking)
}

func (sm *PastPulseSM) PulseCommitted(ctx smachine.MigrationContext) smachine.StateUpdate {
	// trigger some events?
	return ctx.Same()
}

func (sm *PastPulseSM) StateWorking(ctx smachine.ExecutionContext) smachine.StateUpdate {
	return sm.ps.processEvents(ctx, false)
}

func (sm *PastPulseSM) StateSuspending(ctx smachine.ExecutionContext) smachine.StateUpdate {
	// TODO check event queue through adapter
	return ctx.Stop()
}
