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
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
)

// WalletDomainName is a name for wallet domain.
const WalletDomainName = "WalletDomain"

// WalletDomain is a contract that allows to add new wallets to system.
type WalletDomain interface {
	// Base domain implementation.
	domain.Domain
	// CreateWallet is used to create new wallet as a child to domain storage and inject composite to member
	CreateWallet(m Member) error
}

type walletDomain struct {
	domain.BaseDomain
	walletFactoryReference object.Reference
}

func newWalletDomain(parent object.Parent, class object.Factory) (*walletDomain, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}

	wf := NewWalletFactory(parent)
	wd := &walletDomain{
		BaseDomain: *domain.NewBaseDomain(parent, class, WalletDomainName),
	}

	// Add walletFactory as a child
	record, err := wd.AddChild(wf)
	if err != nil {
		return nil, err
	}

	wd.walletFactoryReference, err = object.NewReference("", record, object.ChildScope)
	if err != nil {
		return nil, err
	}
	return wd, nil
}

func (*walletDomain) GetClassID() string {
	return class.WalletDomainID
}

func (wd *walletDomain) CreateWallet(m Member) error {
	// Get child by walletFactoryReference
	r := wd.GetResolver()
	// TODO: pass specific classID for factory resolving
	child, err := r.GetObject(wd.walletFactoryReference, nil)
	if err != nil {
		return err
	}

	// Check if it CompositeFactory
	wf, ok := child.(object.CompositeFactory)
	if !ok {
		return fmt.Errorf("child by reference `%s` is not CompositeFactory instance", wd.walletFactoryReference)
	}

	_, err = m.GetOrCreateComposite(wf)
	if err != nil {
		return err
	}

	return nil
}

type walletDomainProxy struct {
	contract.BaseSmartContractProxy
}

// newWalletDomainProxy creates new proxy and associates it with new instance of WalletDomain.
func newWalletDomainProxy(parent object.Parent, class object.Factory) (*walletDomainProxy, error) {
	inst, err := newWalletDomain(parent, class)
	if err != nil {
		return nil, err
	}

	return &walletDomainProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: inst,
		},
	}, nil
}

// CreateWallet is a proxy call for instance method.
func (wdp *walletDomainProxy) CreateWallet(mp *memberProxy) error {
	return wdp.Instance.(WalletDomain).CreateWallet(mp)
}

type walletDomainFactory struct {
	object.BaseFactory
	parent object.Parent
}

func (wdf *walletDomainFactory) Create(parent object.Parent) (object.Proxy, error) {
	proxy, err := newWalletDomainProxy(parent, wdf)
	if err != nil {
		return nil, err
	}

	_, err = parent.AddChild(proxy)

	if err != nil {
		return nil, err
	}
	return proxy, nil
}

func (*walletDomainFactory) GetParent() object.Parent {
	// TODO: return real parent, fix tests
	return nil
}

func (*walletDomainFactory) GetClassID() string {
	return class.WalletDomainID
}

func (wdf *walletDomainFactory) GetClass() object.Proxy {
	return wdf
}

func NewWalletDomainFactory(pt object.Parent) object.Factory {
	return &walletDomainFactory{
		parent: pt,
	}
}
