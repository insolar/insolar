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

import (
	"context"

	"github.com/insolar/insolar/conveyor/injector"
)

type StateMachine interface {
	// Returns a meta-type / declaration of a SM.
	// Must be non-nil.
	GetStateMachineDeclaration() StateMachineDeclaration
}

type StateMachineDeclaration interface {
	// Initialization code that SM doesn't need to know.
	// Is called once per SM right after GetStateMachineDeclaration().
	// Dependencies injected through DependencyInjector and implementing ShadowMigrator will be invoked during migration.
	InjectDependencies(StateMachine, SlotLink, *injector.DependencyInjector)

	// Provides per-SM logger.
	// When isOutput=true, only the returned logger will be applied and with logger=nil there will be no output.
	// When isOutput=false, the returned logger is applied to StepLoggerData, and then the default logger from SlotMachine is executed.
	// When result is (nil, false) then StateMachine reference will be remembered to be added into every StepLoggerData.
	// Is called once per SM after InjectDependencies().
	GetStepLogger(context.Context, StateMachine) (lf StepLoggerFunc, isOutput bool)

	// Returns an initialization function for the given SM.
	// Is called once per SM after InjectDependencies().
	GetInitStateFor(StateMachine) InitFunc

	// Returns a shadow migration handler for the given stateMachine, that will be invoked on every migration. SM has no control over it.
	// Is called once per SM after InjectDependencies().
	// See ShadowMigrator
	GetShadowMigrateFor(StateMachine) ShadowMigrateFunc

	// WARNING! DO NOT EVER return "true" here without CLEAR understanding of internal mechanics.
	// Returning "true" blindly will LIKELY lead to infinite loops.
	IsConsecutive(cur, next StateFunc) bool
}

// See ShadowMigrator
type ShadowMigrateFunc func(migrationCount, migrationDelta uint32)

// Provides assistance to injected and other objects handle migration events.
type ShadowMigrator interface {
	// Called on migration of a related slot BEFORE every call to a normal migration handler with migrationDelta=1.
	// When there is no migration handler is present or SkipMultipleMigrations() was used, then an additional call is made
	// with migrationDelta > 0 to indicate how many migration steps were skipped.
	ShadowMigrate(migrationCount, migrationDelta uint32)
}

type StepLoggerFlags uint8

const (
	StepLoggerMigrate StepLoggerFlags = 1 << iota
	StepLoggerDetached

	// logger should not log these without necessity
	StepLoggerInternal
)

type StepLoggerData struct {
	StepNo      StepLink
	CurrentStep SlotStep
	NextStep    SlotStep
	UpdateType  string
	Flags       StepLoggerFlags
	Error       error

	// NB! This field can't be provided by SlotMachine and will be nil
	// but it can be filled in by custom wrappers of StepLoggerFunc
	SM StateMachine
}

type StepLoggerFunc func(*StepLoggerData)
type StepLoggerFactoryFunc func(context.Context) StepLoggerFunc

// A template to include into SM to avoid hassle of creation of any methods but GetInitStateFor()
type StateMachineDeclTemplate struct {
}

//var _ StateMachineDeclaration = &StateMachineDeclTemplate{}
//
//func (s *StateMachineDeclTemplate) GetInitStateFor(StateMachine) InitFunc {
//	panic("implement me")
//}

func (s *StateMachineDeclTemplate) IsConsecutive(cur, next StateFunc) bool {
	return false
}

func (s *StateMachineDeclTemplate) GetShadowMigrateFor(StateMachine) ShadowMigrateFunc {
	return nil
}

func (s *StateMachineDeclTemplate) InjectDependencies(StateMachine, SlotLink, *injector.DependencyInjector) {
}

func (s *StateMachineDeclTemplate) GetStepLogger(context.Context, StateMachine) (StepLoggerFunc, bool) {
	return nil, false
}

type TerminationHandlerFunc func(context.Context, TerminationData)

type TerminationData struct {
	Slot   StepLink
	Parent SlotLink
	Result interface{}
	Error  error

	// ===============
	worker FixedSlotWorker
}

// See mergeDefaultValues() and prepareNewSlotWithDefaults()
type CreateDefaultValues struct {
	Context                context.Context
	Parent                 SlotLink
	OverriddenDependencies map[string]interface{}
	TerminationHandler     TerminationHandlerFunc
	StepLogger             StepLoggerFunc
}

func (p *CreateDefaultValues) PutOverride(id string, v interface{}) {
	if id == "" {
		panic("illegal value")
	}
	if p.OverriddenDependencies == nil {
		p.OverriddenDependencies = map[string]interface{}{id: v}
	} else {
		p.OverriddenDependencies[id] = v
	}
}
