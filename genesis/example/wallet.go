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
	"github.com/insolar/insolar/genesis/model/factory"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/insolar/insolar/genesis/model/resolver"
)

type Wallet interface {
	object.Composite
	contract.SmartContract

	GetBalance() int
}

type wallet struct {
	contract.BaseSmartContract

	balance int
}

func (w *wallet) GetBalance() int {
	return w.balance
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
		balance:           0,
	}, nil
}

type walletProxy struct {
	resolver.BaseProxy
}

func newWalletProxy(parent object.Parent) (*walletProxy, error) {
	instance, err := newWallet(parent)
	if err != nil {
		return nil, err
	}

	return &walletProxy{
		BaseProxy: resolver.BaseProxy{
			Instance: instance,
		},
	}, nil
}

func (wp *walletProxy) GetBalance() int {
	return wp.Instance.(Wallet).GetBalance()
}

type walletFactory struct {
	object.BaseCallable
	parent object.Parent
}

func NewWalletFactory(parent object.Parent) factory.Factory {
	return &walletFactory{
		parent: parent,
	}
}

func (wf *walletFactory) GetClassID() string {
	return class.WalletID
}

func (*walletFactory) Create(parent object.Parent) (resolver.Proxy, error) {
	proxy, err := newWalletProxy(parent)
	if err != nil {
		return nil, err
	}

	_, err = parent.AddChild(proxy)
	if err != nil {
		return nil, err
	}

	return proxy, nil
}

func (wf *walletFactory) GetParent() object.Parent {
	return wf.parent
}
