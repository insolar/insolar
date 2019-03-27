/*
 *    Copyright 2019 Insolar Technologies
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

package handler

import (
	"github.com/insolar/insolar/conveyor/interfaces/fsm"
	"github.com/insolar/insolar/conveyor/interfaces/slot"
)

// Types below describes different types of raw handlers
type TransitHandler func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error)
type MigrationHandler func(element slot.SlotElementHelper) (interface{}, fsm.ElementState, error)
type AdapterResponseHandler func(element slot.SlotElementHelper, response interface{}) (interface{}, fsm.ElementState, error)
type NestedHandler func(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState)
type TransitionErrorHandler func(element slot.SlotElementHelper, err error) (interface{}, fsm.ElementState)
type ResponseErrorHandler func(element slot.SlotElementHelper, response interface{}, err error) (interface{}, fsm.ElementState)
