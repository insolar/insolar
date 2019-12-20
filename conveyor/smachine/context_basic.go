///
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
///

package smachine

import (
	"context"
	"fmt"
	"math"
	"unsafe"
)

type contextTemplate struct {
	mode updCtxMode
}

// protects from taking a copy of a context
func (p *contextTemplate) getMarker() ContextMarker {
	return ContextMarker(unsafe.Pointer(p))
	// ContextMarker(unsafe.Pointer(p)) ^ atomicCounter.AddUint32(1) << 16
}

func (p *contextTemplate) ensureAndPrepare(s *Slot, stateUpdate StateUpdate) StateUpdate {
	stateUpdate.ensureMarker(p.getMarker())

	sut := typeOfStateUpdateForPrepare(p.mode, stateUpdate)
	sut.Prepare(s, &stateUpdate)

	return stateUpdate
}

func (p *contextTemplate) setMode(mode updCtxMode) {
	if mode == updCtxInactive {
		panic("illegal value")
	}
	if p.mode != updCtxInactive {
		panic("illegal state")
	}
	p.mode = mode
}

func (p *contextTemplate) ensureAtLeast(mode updCtxMode) {
	if p.mode < mode {
		panic("illegal state")
	}
}

func (p *contextTemplate) ensure(mode0 updCtxMode) {
	if p.mode != mode0 {
		panic("illegal state")
	}
}

func (p *contextTemplate) ensureAny2(mode0, mode1 updCtxMode) {
	if p.mode != mode0 && p.mode != mode1 {
		panic("illegal state")
	}
}

func (p *contextTemplate) ensureAny3(mode0, mode1, mode2 updCtxMode) {
	if p.mode != mode0 && p.mode != mode1 && p.mode != mode2 {
		panic("illegal state")
	}
}

func (p *contextTemplate) ensureValid() {
	if p.mode <= updCtxDiscarded {
		panic("illegal state")
	}
}

func (p *contextTemplate) template(updType stateUpdKind) StateUpdateTemplate {
	return newStateUpdateTemplate(p.mode, p.getMarker(), updType)
}

func (p *contextTemplate) setDiscarded() {
	p.mode = updCtxDiscarded
}

func (p *contextTemplate) discardAndCapture(msg string, recovered interface{}, err *error) {
	p.mode = updCtxDiscarded
	if recovered == nil {
		return
	}
	*err = RecoverSlotPanic(msg, recovered, *err)
}

func (p *contextTemplate) discardAndUpdate(msg string, recovered interface{}, update *StateUpdate) {
	p.mode = updCtxDiscarded
	recoverSlotPanicAsUpdate(update, msg, recovered, nil)
}

/* ========================================================================= */

type slotContext struct {
	contextTemplate
	s *Slot
	w DetachableSlotWorker
}

func (p *slotContext) clone(mode updCtxMode) slotContext {
	p.ensureValid()
	return slotContext{s: p.s, w: p.w, contextTemplate: contextTemplate{mode: mode}}
}

func (p *slotContext) SlotLink() SlotLink {
	p.ensureValid()
	return p.s.NewLink()
}

func (p *slotContext) StepLink() StepLink {
	p.ensure(updCtxExec)
	return p.s.NewStepLink()
}

func (p *slotContext) GetContext() context.Context {
	p.ensureValid()
	return p.s.ctx
}

func (p *slotContext) ParentLink() SlotLink {
	p.ensureValid()
	return p.s.parent
}

func (p *slotContext) SetDefaultErrorHandler(fn ErrorHandlerFunc) {
	p.ensureAtLeast(updCtxInit)
	p.s.defErrorHandler = fn
}

func (p *slotContext) SetDefaultMigration(fn MigrateFunc) {
	p.ensureAtLeast(updCtxInit)
	p.s.defMigrate = fn
}

func (p *slotContext) SetDefaultFlags(flags StepFlags) {
	p.ensureAtLeast(updCtxInit)
	if flags&StepResetAllFlags != 0 {
		p.s.defFlags = flags &^ StepResetAllFlags
	} else {
		p.s.defFlags |= flags
	}
}

func (p *slotContext) SetDefaultTerminationResult(v interface{}) {
	p.ensureAtLeast(updCtxInit)
	p.s.defResult = v
}

func (p *slotContext) GetDefaultTerminationResult() interface{} {
	p.ensureAtLeast(updCtxInit)
	return p.s.defResult
}

func (p *slotContext) SetDynamicBoost(boosted bool) {
	p.ensureAtLeast(updCtxInit)
	if boosted {
		p.s.boost = activeBoost
	} else {
		p.s.boost = inactiveBoost
	}
}

func (p *slotContext) JumpExt(step SlotStep) StateUpdate {
	return p.template(stateUpdNext).newStep(step, nil)
}

func (p *slotContext) Jump(fn StateFunc) StateUpdate {
	return p.template(stateUpdNext).newStep(SlotStep{Transition: fn}, nil)
}

func (p *slotContext) Stop() StateUpdate {
	return p.template(stateUpdStop).newNoArg()
}

func (p *slotContext) Error(err error) StateUpdate {
	return p.template(stateUpdError).newError(err)
}

func (p *slotContext) Errorf(msg string, a ...interface{}) StateUpdate {
	return p.Error(fmt.Errorf(msg, a...))
}

func (p *slotContext) _prepareReplacementData() prepareSlotValue {
	return prepareSlotValue{
		slotReplaceData: p.s.slotReplaceData.takeOutForReplace(),
		isReplacement:   true,
		tracerId:        p.s.getTracerId(),
	}
}

func (p *slotContext) Replace(fn CreateFunc) StateUpdate {
	tmpl := p.template(stateUpdReplace) // ensures state of this context

	def := prepareReplaceData{fn: fn,
		def: p._prepareReplacementData(),
	}
	return tmpl.newVar(def)
}

func (p *slotContext) ReplaceExt(fn CreateFunc, defValues CreateDefaultValues) StateUpdate {
	tmpl := p.template(stateUpdReplace) // ensures state of this context

	def := prepareReplaceData{fn: fn,
		def: p._prepareReplacementData(),
	}
	mergeDefaultValues(&def.def, defValues)

	return tmpl.newVar(def)
}

func (p *slotContext) ReplaceWith(sm StateMachine) StateUpdate {
	tmpl := p.template(stateUpdReplaceWith) // ensures state of this context

	def := prepareReplaceData{sm: sm,
		def: p._prepareReplacementData(),
	}
	return tmpl.newVar(def)
}

func (p *slotContext) Repeat(limit int) StateUpdate {
	ulimit := uint32(0)
	switch {
	case limit <= 0:
	case limit > math.MaxUint32:
		ulimit = math.MaxUint32
	default:
		ulimit = uint32(limit)
	}

	return p.template(stateUpdRepeat).newUint(ulimit)
}

func (p *slotContext) Stay() StateUpdate {
	return p.template(stateUpdNoChange).newNoArg()
}

func (p *slotContext) WakeUp() StateUpdate {
	return p.template(stateUpdWakeup).newNoArg()
}

func (p *slotContext) AffectedStep() SlotStep {
	p.ensureAny3(updCtxMigrate, updCtxBargeIn, updCtxFail)
	r := p.s.step
	r.Flags |= StepResetAllFlags
	return p.s.step
}

func (p *slotContext) NewChild(fn CreateFunc) SlotLink {
	return p._newChild(fn, false, CreateDefaultValues{Context: p.s.ctx, Parent: p.s.NewLink()})
}

func (p *slotContext) NewChildExt(fn CreateFunc, defValues CreateDefaultValues) SlotLink {
	return p._newChild(fn, false, defValues)
}

func (p *slotContext) InitChild(fn CreateFunc) SlotLink {
	return p._newChild(fn, true, CreateDefaultValues{Context: p.s.ctx, Parent: p.s.NewLink()})
}

func (p *slotContext) InitChildExt(fn CreateFunc, defValues CreateDefaultValues) SlotLink {
	return p._newChild(fn, true, defValues)
}

func (p *slotContext) _newChild(fn CreateFunc, runInit bool, defValues CreateDefaultValues) SlotLink {
	p.ensureAny2(updCtxExec, updCtxFail)
	if fn == nil {
		panic("illegal value")
	}
	if len(defValues.TracerId) == 0 {
		defValues.TracerId = p.s.getTracerId()
	}

	m := p.s.machine
	link, ok := m.prepareNewSlotWithDefaults(p.s, fn, nil, defValues)
	if ok {
		m.startNewSlotByDetachable(link.s, runInit, p.w)
	}
	return link
}

func (p *slotContext) Log() Logger {
	p.ensureAtLeast(updCtxInit)
	return p._newLogger()
}

func (p *slotContext) LogAsync() Logger {
	p.ensure(updCtxExec)
	logger, _ := p._newLoggerAsync()
	return logger
}

func (p *slotContext) _newLogger() Logger {
	// TODO make newStepLoggerData() call lazy
	return Logger{p.s.ctx, p}
}

func (p *slotContext) _newLoggerAsync() (Logger, uint32) {
	logger := Logger{p.s.ctx, nil}
	stepLogger, level, _ := p.getStepLogger()
	if stepLogger == nil {
		return logger, 0
	}

	fsl := fixedSlotLogger{logger: stepLogger, level: level, data: p.getStepLoggerData()}
	logger.ctx, fsl.logger = stepLogger.CreateAsyncLogger(logger.ctx, &fsl.data)
	if fsl.logger == nil || logger.ctx == nil {
		panic("illegal state - logger doesnt support async")
	}
	fsl.data.Flags |= StepLoggerDetached

	logger.loggerFn = fsl
	return logger, fsl.data.StepNo.step
}

func (p *slotContext) getStepLogger() (StepLogger, StepLogLevel, uint32) {
	return p.s.stepLogger, p.s.getStepLogLevel(), 0
}

func (p *slotContext) getStepLoggerData() StepLoggerData {
	return p.s.newStepLoggerData(StepLoggerTrace, p.s.NewStepLink())
}

func (p *slotContext) SetLogTracing(b bool) {
	p.ensureAtLeast(updCtxInit)
	p.s.setTracing(b)
}

func (p *slotContext) UpdateDefaultStepLogger(updateFn StepLoggerUpdateFunc) {
	p.ensureAtLeast(updCtxInit)
	p.s.setStepLoggerAfterInit(updateFn)
}

func (p *slotContext) BargeInWithParam(applyFn BargeInApplyFunc) BargeInParamFunc {
	p.ensureAtLeast(updCtxInit)
	return p.s.machine.createBargeIn(p.s.NewStepLink().AnyStep(), applyFn)
}

func (p *slotContext) BargeIn() BargeInBuilder {
	p.ensureAtLeast(updCtxInit)
	return &bargeInBuilder{p.clone(updCtxBargeIn), p, p.s.NewStepLink().AnyStep()}
}

func (p *slotContext) BargeInThisStepOnly() BargeInBuilder {
	p.ensureAtLeast(updCtxExec)
	return &bargeInBuilder{p.clone(updCtxBargeIn), p, p.s.NewStepLink()}
}

func (p *slotContext) Check(link SyncLink) BoolDecision {
	p.ensureAtLeast(updCtxInit)

	if link.controller == nil {
		panic("illegal value")
	}

	dep := p.s.dependency
	if dep != nil {
		if d, ok := link.controller.CheckDependency(dep).AsValid(); ok {
			return d
		}
	}
	return link.controller.CheckState()
}

func (p *slotContext) AcquireForThisStep(link SyncLink) BoolDecision {
	return p.acquire(link, false, syncForOneStep)
}

func (p *slotContext) Acquire(link SyncLink) BoolDecision {
	return p.acquire(link, false, 0)
}

func (p *slotContext) AcquireAndRelease(link SyncLink) BoolDecision {
	return p.acquire(link, true, 0)
}

func (p *slotContext) AcquireForThisStepAndRelease(link SyncLink) BoolDecision {
	return p.acquire(link, true, 0)
}

func (p *slotContext) acquire(link SyncLink, autoRelease bool, flags SlotDependencyFlags) (d BoolDecision) {
	p.ensureAtLeast(updCtxInit)

	switch {
	case p.s.isPriority():
		flags |= syncPriorityHigh
	case p.s.isBoosted():
		flags |= syncPriorityBoosted
	}

	dep := p.s.dependency
	if dep == nil {
		d, p.s.dependency = link.controller.CreateDependency(p.s.NewLink(), flags)
		return d
	}
	if d, ok := link.controller.UseDependency(dep, flags).AsValid(); ok {
		return d
	}

	if !autoRelease {
		panic("SM has already acquired another sync or the same one but with an incompatible mode")
	}

	slotLink := p.s.NewLink()
	p.s.dependency = nil
	d, p.s.dependency = link.controller.CreateDependency(slotLink, flags)

	postponed, released := dep.ReleaseAll()
	released = PostponedList(postponed).PostponedActivate(released)

	p.s.machine.activateDependantByDetachable(released, slotLink, p.w)

	return d
}

func (p *slotContext) ReleaseLast() bool {
	p.ensureAtLeast(updCtxInit)
	return p.release(nil)
}

func (p *slotContext) Release(link SyncLink) bool {
	p.ensureAtLeast(updCtxInit)

	if link.IsZero() {
		panic("illegal value")
	}
	return p.release(link.controller)
}

func (p *slotContext) ReleaseAll() bool {
	p.ensureAtLeast(updCtxInit)

	dep := p.s.dependency
	if dep == nil {
		return false
	}

	p.s.dependency = nil
	postponed, released := dep.ReleaseAll()
	released = PostponedList(postponed).PostponedActivate(released)

	p.s.machine.activateDependantByDetachable(released, p.s.NewLink(), p.w)
	return true
}

func (p *slotContext) release(controller DependencyController) bool {
	dep := p.s.dependency
	if dep == nil {
		return false
	}
	if controller != nil && !controller.CheckDependency(dep).IsValid() {
		// mismatched sync object
		return false
	}

	released := p.s._releaseDependency()
	p.s.machine.activateDependantByDetachable(released, p.s.NewLink(), p.w)
	return true
}

func (p *slotContext) ApplyAdjustment(adj SyncAdjustment) bool {
	p.ensureAtLeast(updCtxInit)

	if adj.controller == nil {
		panic("illegal value")
	}

	released, activate := adj.controller.AdjustLimit(adj.adjustment, adj.isAbsolute)
	if activate {
		return p.s.machine.activateDependantByDetachable(released, p.s.NewLink(), p.w)
	}

	// actually, we MUST NOT stop a slot from outside
	return len(released) > 0
}
