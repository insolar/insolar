/*
 *    Copyright 2018 Insolar
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

type Allowance interface {
	object.Composite
	contract.SmartContract
	GetAmount() int
	GetSender() object.Reference
	GetReceiver() object.Reference
	MarkCompleted()
	IsCompleted() bool
}

type allowance struct {
	contract.BaseSmartContract
	sender    object.Reference
	receiver  object.Reference
	amount    int
	completed bool
}

func newAllowance(parent object.Parent, class object.CompositeFactory) (Allowance, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}

	//TODO: add posibility to init allowance fields
	return &allowance{
		BaseSmartContract: *contract.NewBaseSmartContract(parent, class.(object.Proxy)),
		completed:         false,
	}, nil
}

func NewAllowanceWithParams(parent object.Parent, class object.CompositeFactory, sender object.Reference, receiver object.Reference, amount int) (Allowance, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}

	return &allowance{
		BaseSmartContract: *contract.NewBaseSmartContract(parent, class.(object.Proxy)),
		sender:            sender,
		receiver:          receiver,
		amount:            amount,
		completed:         false,
	}, nil
}

func (a *allowance) GetInterfaceKey() string {
	return class.AllowanceID
}

func (a *allowance) GetAmount() int {
	return a.amount
}

func (a *allowance) IsCompleted() bool {
	return a.completed
}

func (a *allowance) MarkCompleted() {
	a.completed = true
}

func (a *allowance) GetSender() object.Reference {
	return a.sender
}

func (a *allowance) GetReceiver() object.Reference {
	return a.receiver
}

type allowanceFactory struct {
	object.BaseProxy
	parent object.Parent
}

func NewAllowanceFactory(parent object.Parent) (object.CompositeFactory, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}

	return &allowanceFactory{
		parent: parent,
	}, nil
}

func (*allowanceFactory) GetInterfaceKey() string {
	return class.AllowanceID
}

type AllowanceCompositeCollection struct {
	contract.BaseCompositeCollection
	parent object.Parent
	class  object.Proxy
}

func (acc *AllowanceCompositeCollection) GetClass() object.Proxy {
	return acc.class
}

func (acc *AllowanceCompositeCollection) GetParent() object.Parent {
	return acc.parent
}

func (*AllowanceCompositeCollection) GetInterfaceKey() string {
	return class.AllowanceID
}

func (acc *AllowanceCompositeCollection) GetClassID() string {
	return class.AllowanceID
}

func newAllowanceCollectionProxy(parent object.Parent, class object.CompositeFactory) (*contract.BaseCompositeCollectionProxy, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}
	alCollection := AllowanceCompositeCollection{
		parent: parent,
		class:  class,
	}

	cProxy := &contract.BaseCompositeCollectionProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: &alCollection,
		},
	}

	return cProxy, nil
}

func (af *allowanceFactory) Create(parent object.Parent) (object.Composite, error) {
	proxy, err := newAllowanceCollectionProxy(parent, af)
	if err != nil {
		return nil, err
	}

	_, err = parent.AddChild(proxy)
	if err != nil {
		return nil, err
	}

	return proxy, nil
}

func (*allowanceFactory) GetClassID() string {
	return class.AllowanceID
}

func (af *allowanceFactory) GetClass() object.Proxy {
	return af
}

func (af *allowanceFactory) GetParent() object.Parent {
	return af.parent
}

type allowanceProxy struct {
	contract.BaseSmartContractProxy
}

func newAllowanceProxy(parent object.Parent, class object.CompositeFactory) (*allowanceProxy, error) {
	inst, err := newAllowance(parent, class)
	if err != nil {
		return nil, err
	}

	return &allowanceProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: inst,
		},
	}, nil
}

func (ap *allowanceProxy) GetAmount() int {
	return ap.Instance.(Allowance).GetAmount()
}

func (ap *allowanceProxy) GetSender() object.Reference {
	return ap.Instance.(Allowance).GetSender()
}

func (ap *allowanceProxy) GetReceiver() object.Reference {
	return ap.Instance.(Allowance).GetReceiver()
}

func (ap *allowanceProxy) MarkCompleted() {
	ap.Instance.(Allowance).MarkCompleted()
}

func (ap *allowanceProxy) IsCompleted() bool {
	return ap.Instance.(Allowance).IsCompleted()
}

func (ap *allowanceProxy) GetInterfaceKey() string {
	return ap.Instance.(Allowance).GetInterfaceKey()
}
