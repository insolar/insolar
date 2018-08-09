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
	"fmt"

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/contract"
	"github.com/insolar/insolar/genesis/model/object"
)

type Wallet interface {
	object.Composite
	contract.SmartContract
}

type wallet struct {
	contract.BaseSmartContract
}

func (w *wallet) GetClassID() string {
	return class.WalletID
}

func (w *wallet) GetInterfaceKey() string {
	return w.GetClassID()
}

func newWallet(parent object.Parent) (Wallet, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}

	return &wallet{
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
	}, nil
}
