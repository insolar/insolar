/*
 *    Copyright 2019 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package conveyor

import (
	"github.com/insolar/insolar/core"
)

type StateID uint32

type AdapterTask interface {
}

type Event interface {
	DefaultTarget() *core.RecordRef
}

type Helper interface {
	SendTask(task AdapterTask)
}

type SlotItem interface {
	Event() Event
	SetPayload(interface{})
}

type Handler interface {
}

func RegisterActive(id StateID, handler Handler) {
}

func RegisterInactive(id StateID, handler Handler) {
}

func RegisterReply(id StateID, handler Handler) {
}
