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

func (p *machineCallContext) Migrate() {
	p.ensureValid()
	p.m.migrate(p.w)
}

func (p *machineCallContext) Cleanup() {
	p.ensureValid()
	p.m.Cleanup(p.w)
}

func (p *machineCallContext) Stop() {
	p.ensureValid()
	p.m.Stop()
}

func (p *machineCallContext) AddNew(ctx context.Context, parent SlotLink, sm StateMachine) SlotLink {
	p.ensureValid()
	link, ok := p.m._addNew(ctx, parent, sm)
	if ok {
		p.m.startNewSlot(link.s, p.w)
	}
	return link
}

func (p *machineCallContext) AddNewByFunc(ctx context.Context, parent SlotLink, cf CreateFunc) (SlotLink, bool) {
	p.ensureValid()
	link, ok := p.m._addNewWithFunc(ctx, parent, cf)
	if ok {
		p.m.startNewSlot(link.s, p.w)
	}
	return link, ok
}

func (p *machineCallContext) SlotMachine() *SlotMachine {
	p.ensureValid()
	return p.m
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
		p.m.activateDependants(released, p.w)
	}

	// actually, we MUST NOT stop a slot from outside
	return len(released) > 0
}

func (p *machineCallContext) Check(link SyncLink) Decision {
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
