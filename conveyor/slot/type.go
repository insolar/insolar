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

package slot

import (
	"github.com/insolar/insolar/conveyor/adapter/adapterid"
)

// AdapterResponse gives access to response of adapter
//go:generate minimock -i github.com/insolar/insolar/conveyor.AdapterResponse -o ./ -s _mock.go
type AdapterResponse interface {
	// GetAdapterID returns adapter id
	GetAdapterID() adapterid.ID
	// GetElementID returns element id
	GetElementID() uint32
	// GetHandlerID returns handler id
	GetHandlerID() uint32
	// GetRespPayload returns payload
	GetRespPayload() interface{}
}

// PulseState is the states of pulse inside slot
type PulseState int

//go:generate stringer -type=PulseState
const (
	Unallocated = PulseState(iota)
	Future
	Present
	Past
	Antique
)
