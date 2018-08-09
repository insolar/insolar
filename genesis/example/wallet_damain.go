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
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/factory"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/insolar/insolar/genesis/model/resolver"
)

// WalletDomainName is a name for wallet domain.
const WalletDomainName = "WalletDomain"

// WalletDomain is a contract that allows to add new wallets to system.
type WalletDomain interface {
	// Base domain implementation.
	domain.Domain
	// CreateWallet is used to create new wallet as a child to domain storage.
	CreateWallet(factory.Factory) (string, error)
	// GetWallet returns wallet from its record in domain storage.
	GetWallet(string) (resolver.Proxy, error)
}

type walletDomain struct {
	domain.BaseDomain
}

func newWalletDomain(parent object.Parent) (*walletDomain, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}

	return &walletDomain{
		BaseDomain: *domain.NewBaseDomain(parent, WalletDomainName),
	}, nil
}

func (wd *walletDomain) GetClassID() string {
	return class.WalletDomainID
}

func (wd *walletDomain) CreateWallet(fc factory.Factory) (string, error) {
	wallet, err := fc.Create(wd)
	if err != nil {
		return "", err
	}
	if wallet == nil {
		return "", fmt.Errorf("factory returns nil")
	}

	record, err := wd.ChildStorage.Set(wallet)
	if err != nil {
		return "", err
	}

	return record, nil
}

func (wd *walletDomain) GetWallet(record string) (resolver.Proxy, error) {
	wallet, err := wd.GetChild(record)
	if err != nil {
		return nil, err
	}

	result, ok := wallet.(resolver.Proxy)
	if !ok {
		return nil, fmt.Errorf("object with record `%s` is not `Proxy` instance", record)
	}

	return result, nil
}

type walletDomainProxy struct {
	contract.BaseSmartContractProxy
}

// newWalletDomainProxy creates new proxy and associates it with new instance of WalletDomain.
func newWalletDomainProxy(parent object.Parent) (*walletDomainProxy, error) {
	inst, err := newWalletDomain(parent)
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
func (wdp *walletDomainProxy) CreateWallet(fc factory.Factory) (string, error) {
	return wdp.Instance.(WalletDomain).CreateWallet(fc)
}

// GetWallet is a proxy call for instance method.
func (wdp *walletDomainProxy) GetWallet(record string) (resolver.Proxy, error) {
	return wdp.Instance.(WalletDomain).GetWallet(record)
}

type walletDomainFactory struct {
	factory.BaseFactory
	parent object.Parent
}

func (wdf *walletDomainFactory) Create(parent object.Parent) (resolver.Proxy, error) {
	proxy, err := newWalletDomainProxy(parent)
	if err != nil {
		return nil, err
	}

	_, err = parent.AddChild(proxy)

	if err != nil {
		return nil, err
	}
	return proxy, nil
}

func (wdf *walletDomainFactory) GetParent() object.Parent {
	// TODO: return real parent, fix tests
	return nil
}

func (wdf *walletDomainFactory) GetClassID() string {
	return class.WalletDomainID
}

func NewWalletDomainFactory(pt object.Parent) factory.Factory {
	return &walletDomainFactory{
		parent: pt,
	}
}
