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

type StepSyncCatalog struct {
	keys map[string]*SortedSlotDependencies
}

func (p *StepSyncCatalog) Join(slot *Slot, key string, weight int32, receiveFunc BroadcastReceiveFunc) Syncronizer {
	panic("not implemented")
	//head := p.keys[key]
	//if head == nil {
	//	p.keys[key] = slot
	//	return
	//}
	//p.keys[key] = head
}

//var _ Syncronizer = &stepSync{}
//
//type stepSync struct {
//
//}
