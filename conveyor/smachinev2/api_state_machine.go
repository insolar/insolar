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

type StateMachine interface {
	//GetStateMachineName() string
	GetStateMachineDeclaration() StateMachineDeclaration
}

type DependencyRegistry interface {
	FindDependency(id string) interface{}
}

type StateMachineDeclaration interface {
	IsConsecutive(cur, next StateFunc) bool
	GetInitStateFor(StateMachine) InitFunc
	InjectDependencies(StateMachine, SlotLink, DependencyRegistry) bool
}

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

func (s *StateMachineDeclTemplate) GetMigrateFn(StateFunc) MigrateFunc {
	return nil
}

func (s *StateMachineDeclTemplate) InjectDependencies(StateMachine, SlotLink, *SlotMachine, DependencyRegistry) bool {
	return false
}
