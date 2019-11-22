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
	"fmt"

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

	// Provides per-SM logger. Zero implementation must return (nil, false).
	// Is called once per SM after InjectDependencies().
	// When result is (_, false) then StepLoggerFactory will be used.
	// When result is (nil, true) then any logging will be disabled.
	GetStepLogger(context.Context, StateMachine, TracerId, StepLoggerFactoryFunc) (StepLogger, bool)

	// Returns an initialization function for the given SM.
	// Is called once per SM after InjectDependencies().
	GetInitStateFor(StateMachine) InitFunc

	// Returns a shadow migration handler for the given stateMachine, that will be invoked on every migration. SM has no control over it.
	// Is called once per SM after InjectDependencies().
	// See ShadowMigrator
	GetShadowMigrateFor(StateMachine) ShadowMigrateFunc

	// Returns a StepDeclaration for the given step. Return nil when implementation is not available.
	GetStepDeclaration(StateFunc) *StepDeclaration

	// This function is only invoked when GetStepDeclaration() is not available for the current step.
	// WARNING! DO NOT EVER return "true" here without CLEAR understanding of internal mechanics.
	// Returning "true" blindly will LIKELY lead to infinite loops.
	IsConsecutive(cur, next StateFunc) (bool, *StepDeclaration)
}

type stepDeclExt struct {
	SeqId int
	Name  string
}

type StepDeclaration struct {
	SlotStep
	stepDeclExt
}

func (v StepDeclaration) GetStepName() string {
	switch {
	case len(v.Name) > 0:
		if v.SeqId != 0 {
			return fmt.Sprintf("%s[%d]", v.Name, v.SeqId)
		}
		return fmt.Sprintf("%s", v.Name)

	case v.SeqId != 0:
		return fmt.Sprintf("#[%d]", v.SeqId)

	case v.Transition == nil:
		return "<nil>"

	default:
		return fmt.Sprintf("%p", v.Transition)
	}
}

func (v StepDeclaration) IsNameless() bool {
	return v.SeqId == 0 && len(v.Name) == 0
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

// A template to include into SM to avoid hassle of creation of any methods but GetInitStateFor()
type StateMachineDeclTemplate struct {
}

//var _ StateMachineDeclaration = &StateMachineDeclTemplate{}
//
//func (s *StateMachineDeclTemplate) GetInitStateFor(StateMachine) InitFunc {
//	panic("implement me")
//}

func (s *StateMachineDeclTemplate) GetStepDeclaration(StateFunc) *StepDeclaration {
	return nil
}

func (s *StateMachineDeclTemplate) IsConsecutive(StateFunc, StateFunc) (bool, *StepDeclaration) {
	return false, nil
}

func (s *StateMachineDeclTemplate) GetShadowMigrateFor(StateMachine) ShadowMigrateFunc {
	return nil
}

func (s *StateMachineDeclTemplate) InjectDependencies(StateMachine, SlotLink, *injector.DependencyInjector) {
}

func (s *StateMachineDeclTemplate) GetStepLogger(context.Context, StateMachine, TracerId, StepLoggerFactoryFunc) (StepLogger, bool) {
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
	TracerId               TracerId
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
