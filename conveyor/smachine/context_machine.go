//
//    Copyright 2019 Insolar Technologies
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.
//

package smachine

import "context"

var _ MachineCallContext = &machineCallContext{}

type machineCallContext struct {
	contextTemplate
	m *SlotMachine
	w FixedSlotWorker
}

func (p *machineCallContext) Migrate(beforeFn func()) {
	p.ensureValid()
	p.m.migrateWithBefore(p.w, beforeFn)
}

func (p *machineCallContext) Cleanup() {
	p.ensureValid()
	p.m.Cleanup(p.w)
}

func (p *machineCallContext) Stop() {
	p.ensureValid()
	p.m.Stop()
}

func (p *machineCallContext) AddNew(ctx context.Context, sm StateMachine, defValues CreateDefaultValues) SlotLink {
	p.ensureValid()
	if sm == nil {
		panic("illegal value")
	}

	switch {
	case ctx != nil:
		defValues.Context = ctx
	case defValues.Context == nil:
		panic("illegal value")
	}

	link, ok := p.m.prepareNewSlotWithDefaults(nil, nil, sm, defValues)
	if ok {
		p.m.startNewSlot(link.s, p.w)
	}
	return link
}

func (p *machineCallContext) AddNewByFunc(ctx context.Context, cf CreateFunc, defValues CreateDefaultValues) (SlotLink, bool) {
	p.ensureValid()

	switch {
	case ctx != nil:
		defValues.Context = ctx
	case defValues.Context == nil:
		panic("illegal value")
	}

	link, ok := p.m.prepareNewSlotWithDefaults(nil, cf, nil, defValues)
	if ok {
		p.m.startNewSlot(link.s, p.w)
	}
	return link, ok
}

func (p *machineCallContext) SlotMachine() *SlotMachine {
	p.ensureValid()
	return p.m
}

func (p *machineCallContext) GetMachineId() string {
	p.ensureValid()
	return p.m.GetMachineId()
}

func (p *machineCallContext) BargeInNow(link SlotLink, param interface{}, fn BargeInApplyFunc) bool {
	p.ensureValid()
	return p.m.bargeInNow(link, param, fn, p.w)
}

func (p *machineCallContext) ApplyAdjustment(adj SyncAdjustment) bool {
	p.ensureValid()

	if adj.controller == nil {
		panic("illegal value")
	}

	released, activate := adj.controller.AdjustLimit(adj.adjustment, adj.isAbsolute)
	if activate {
		p.m.activateDependants(released, SlotLink{}, p.w)
	}

	// actually, we MUST NOT stop a slot from outside
	return len(released) > 0
}

func (p *machineCallContext) Check(link SyncLink) BoolDecision {
	p.ensureValid()

	if link.controller == nil {
		panic("illegal value")
	}

	return link.controller.CheckState()
}

func (p *machineCallContext) executeCall(fn MachineCallFunc) (err error) {
	p.setMode(updCtxMachineCall)
	defer func() {
		p.discardAndCapture("machine call", recover(), &err)
	}()

	fn(p)
	return nil
}
