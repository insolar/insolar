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

package example

import (
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/contract"
)

type Allowance interface {
	GetAmount() int
	GetSender() string
}

type allowance struct {
	sender string
	amount int
}

func (a *allowance) GetInterfaceKey() string {
	return class.AllowanceID
}

func (a *allowance) GetAmount() int {
	return a.amount
}

func (a *allowance) GetSender() string {
	return a.sender
}

func newAllowance(sender string, amount int) *allowance {
	return &allowance{
		sender: sender,
		amount: amount,
	}
}

type allowanceCollection struct {
	contract.BaseCompositeCollection
}

func newAllowanceCollection() {

}
