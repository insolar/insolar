/*
 *    Copyright 2018 INS Ecosystem
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

package jetcoordinator

import (
	"fmt"

	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
	"github.com/insolar/insolar/ledger/storage"
)

// JetCoordinator is responsible for all jet interactions
type JetCoordinator struct {
	storage storage.LedgerStorer
}

// Pulse creates new jet drop and ends current slot.
// This should be called when receiving a new pulse from pulsar.
func (jc *JetCoordinator) Pulse(new record.PulseNum) (*jetdrop.JetDrop, error) {
	current := jc.storage.GetCurrentPulse()
	if new-current != 1 {
		panic(fmt.Sprintf("Wrong pulse, got %v, but current is %v\n", new, current))
	}
	// TODO: increment stored pulse number and wait for all records from previous pulse to store
	return jc.CreateDrop(current)
}
