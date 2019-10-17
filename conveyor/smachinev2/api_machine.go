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

type MachineCallFunc func(MachineCallContext)
type MachineCallContext interface {
	SlotMachine() *SlotMachine

	AddNew(ctx context.Context, parent SlotLink, sm StateMachine) SlotLink
	AddNewByFunc(ctx context.Context, parent SlotLink, cf CreateFunc) (SlotLink, bool)

	BargeInNow(SlotLink, interface{}, BargeInApplyFunc) bool

	GetPublished(key interface{}) interface{}
	GetPublishedLink(key interface{}) SharedDataLink
	//UseShared(SharedDataAccessor) SharedAccessReport

	Migrate()
	AddMigrationCallback(fn MigrationFunc)

	Cleanup()
	Stop()
}
