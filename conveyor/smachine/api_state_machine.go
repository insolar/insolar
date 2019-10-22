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
	// Returns an initialization function for the given SM.
	GetInitStateFor(StateMachine) InitFunc

	// Initialization code that SM doesn't need to know.
	// Dependencies injected through DependencyInjector and implementing ShadowMigrator will be invoked during migration.
	InjectDependencies(StateMachine, SlotLink, *injector.DependencyInjector)

	// Returns a shadow migration handler for the given stateMachine, that will be invoked on every migration. SM has no control over it.
	// See ShadowMigrator
	GetShadowMigrateFor(StateMachine) ShadowMigrateFunc

	GetStepLogger(context.Context, StateMachine) StateMachineStepLoggerFunc

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
)

type StepLoggerData struct {
	StepNo      StepLink
	CurrentStep SlotStep
	NextStep    SlotStep
	UpdateType  string
	Flags       StepLoggerFlags
}

type StateMachineStepLoggerFunc func(StepLoggerData)

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

func (s *StateMachineDeclTemplate) GetStepLogger(context.Context, StateMachine) StateMachineStepLoggerFunc {
	return nil
}
