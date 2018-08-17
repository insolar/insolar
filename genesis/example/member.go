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

type Member interface {
	object.ComposingContainer
	contract.SmartContract
	GetUsername() string
	GetPublicKey() string
	// AddChild(child object.Child) (string, error)
	// GetChild(key string) (object.Child, error)
	// GetChildStorage() storage.Storage
	// GetContext() []string
	// GetContextStorage() storage.Storage
}

type member struct {
	contract.BaseSmartContract
	username  string
	publicKey string
}

// newMember creates new instance of member.
func newMember(parent object.Parent) (Member, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}
	return &member{
		BaseSmartContract: *contract.NewBaseSmartContract(parent),
	}, nil
}

// GetClassID returns string representation of member's class.
func (m *member) GetClassID() string {
	return class.MemberID
}

// GetUsername returns member's username.
func (m *member) GetUsername() string {
	return m.username
}

// GetPublicKey returns member's public key.
func (m *member) GetPublicKey() string {
	return m.publicKey
}

type memberProxy struct {
	contract.BaseSmartContractProxy
}

// newMemberProxy creates new proxy and associates it with new instance of Member.
func newMemberProxy(parent object.Parent) (*memberProxy, error) {
	instance, err := newMember(parent)
	if err != nil {
		return nil, err
	}
	return &memberProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: instance,
		},
	}, nil
}

// GetUsername is a proxy call for instance method.
func (mp *memberProxy) GetUsername() string {
	return mp.Instance.(Member).GetUsername()
}

// GetPublicKey is a proxy call for instance method.
func (mp *memberProxy) GetPublicKey() string {
	return mp.Instance.(Member).GetPublicKey()
}

// CreateComposite is a proxy call for instance method.
func (mp *memberProxy) CreateComposite(compositeFactory object.CompositeFactory) (object.Composite, error) {
	return mp.Instance.(Member).CreateComposite(compositeFactory)
}

// GetComposite is a proxy call for instance method.
func (mp *memberProxy) GetComposite(interfaceKey string, classID string) (object.Composite, error) {
	return mp.Instance.(Member).GetComposite(interfaceKey, classID)
}

// GetOrCreateComposite is a proxy call for instance method.
func (mp *memberProxy) GetOrCreateComposite(compositeFactory object.CompositeFactory) (object.Composite, error) {
	return mp.Instance.(Member).GetOrCreateComposite(compositeFactory)
}

type memberFactory struct {
	object.BaseProxy
	parent object.Parent
}

// NewMemberFactory creates new factory for Member.
func NewMemberFactory(parent object.Parent) object.Factory {
	return &memberFactory{
		parent: parent,
	}
}

// GetClassID returns string representation of Member's class.
func (mf *memberFactory) GetClassID() string {
	return class.MemberID
}

// GetParent returns parent.
func (mf *memberFactory) GetParent() object.Parent {
	return mf.parent
}

// Create is a factory method for new Member instances.
func (mf *memberFactory) Create(parent object.Parent) (object.Proxy, error) {
	proxy, err := newMemberProxy(parent)
	if err != nil {
		return nil, err
	}

	_, err = parent.AddChild(proxy)
	if err != nil {
		return nil, err
	}
	return proxy, nil
}
