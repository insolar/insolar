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
	"github.com/insolar/insolar/ledger/jetdrop"
	"github.com/insolar/insolar/ledger/record"
)

// CreateDrop creates jet drop for provided pulse number.
func (jc *JetCoordinator) CreateDrop(pulse record.PulseNum) (*jetdrop.JetDrop, error) {
	prevDrop, err := jc.storage.GetDrop(pulse - 1)
	if err != nil {
		return nil, err
	}
	newDrop, err := jc.storage.SetDrop(pulse, prevDrop)
	if err != nil {
		return nil, err
	}
	return newDrop, nil
}
