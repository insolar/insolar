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

package contract

import (
	"fmt"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/insolar/insolar/genesis/model/resolver"
)

// SmartContract marks that object is smart contract.
// TODO: Composite work interface
type SmartContract interface {
	object.Child
	GetResolver() resolver.Resolver
}

// BaseSmartContract is a base implementation of ComposingContainer, Callable and TypedObject interfaces.
type BaseSmartContract struct {
	Reference      *object.Reference
	CompositeMap   map[string]object.Composite
	ChildStorage   storage.Storage
	ContextStorage storage.Storage
	Parent         object.Parent
	resolver       resolver.Resolver
}

// NewBaseSmartContract creates new BaseSmartContract instance with empty CompositeMap, ChildStorage and specific parent.
func NewBaseSmartContract(parent object.Parent) *BaseSmartContract {
	// TODO: NewCompositeHolder
	sc := BaseSmartContract{
		CompositeMap: make(map[string]object.Composite),
		ChildStorage: storage.NewMapStorage(),
		Parent:       parent,
	}
	sc.GetResolver()
	return &sc
}

// GetResolver return instance or create it if necessary.
func (sc *BaseSmartContract) GetResolver() resolver.Resolver {
	if sc.resolver == nil {
		sc.resolver = resolver.NewHandler(sc)
	}
	return sc.resolver
}

// GetClassID return string representation of object's class.
func (sc *BaseSmartContract) GetClassID() string {
	return class.SmartContractID
}

// GetReference return reference to BaseSmartContract instance.
func (sc *BaseSmartContract) GetReference() *object.Reference {
	// TODO should return actual reference
	return sc.Reference
}

// CreateComposite allows to create composites inside smart contract.
func (sc *BaseSmartContract) CreateComposite(compositeFactory object.CompositeFactory) (object.Composite, error) {
	composite := compositeFactory.Create()
	interfaceKey := composite.GetInterfaceKey()
	_, isExist := sc.CompositeMap[interfaceKey]
	if isExist {
		return nil, fmt.Errorf("delegate with name %s already exist", interfaceKey)
	}
	sc.CompositeMap[interfaceKey] = composite
	return composite, nil
}

// GetComposite return composite by its key (if its exist inside smart contract).
func (sc *BaseSmartContract) GetComposite(key string) (object.Composite, error) {
	composite, isExist := sc.CompositeMap[key]
	if !isExist {
		return nil, fmt.Errorf("delegate with name %s does not exist", key)
	}
	return composite, nil
}

// GetOrCreateComposite return composite by its key if its exist inside smart contract and create new one otherwise.
func (sc *BaseSmartContract) GetOrCreateComposite(interfaceKey string, compositeFactory object.CompositeFactory) (object.Composite, error) {
	composite, err := sc.GetComposite(interfaceKey)
	if err != nil {
		composite, err = sc.CreateComposite(compositeFactory)
		if err != nil {
			return nil, err
		}
		return composite, nil
	}
	return composite, nil
}

// GetChildStorage return storage with children of smart contract.
func (sc *BaseSmartContract) GetChildStorage() storage.Storage {
	return sc.ChildStorage
}

// AddChild add new child to smart contract's ChildStorage.
func (sc *BaseSmartContract) AddChild(child object.Child) (string, error) {
	key, err := sc.ChildStorage.Set(child)
	if err != nil {
		return "", err
	}
	return key, nil
}

// GetChild get child from smart contract's ChildStorage.
func (sc *BaseSmartContract) GetChild(key string) (object.Child, error) {
	child, err := sc.ChildStorage.Get(key)
	if err != nil {
		return nil, err
	}
	return child.(object.Child), nil
}

// GetContextStorage return storage with objects, which smart contract's children will have access to.
func (sc *BaseSmartContract) GetContextStorage() storage.Storage {
	return sc.ContextStorage
}

// GetContext return list of keys in ContextStorage.
func (sc *BaseSmartContract) GetContext() []string {
	return sc.GetContextStorage().GetKeys()
}

// GetParent return parent of smart contract.
func (sc *BaseSmartContract) GetParent() object.Parent {
	return sc.Parent
}
