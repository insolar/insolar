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
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/factory"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/insolar/insolar/genesis/model/resolver"
)

// MemberDomainName is a name for member domain.
const MemberDomainName = "MemberDomain"

// MemberDomain is a contract that allows to add new members to system.
type MemberDomain interface {
	// Base domain implementation.
	domain.Domain
	// CreateMember is used to create new member as a child to domain storage.
	CreateMember(factory.Factory) (string, error)
	// GetMember returns member from its record in domain storage.
	GetMember(string) (resolver.Proxy, error)
}

type memberDomain struct {
	domain.BaseDomain
}

// newMemberDomain creates new instance of MemberDomain.
func newMemberDomain(parent object.Parent) (*memberDomain, error) {
	if parent == nil {
		return nil, fmt.Errorf("parent must not be nil")
	}

	return &memberDomain{
		BaseDomain: *domain.NewBaseDomain(parent, MemberDomainName),
	}, nil
}

// GetClassID returns string representation of MemberDomain's class.
func (md *memberDomain) GetClassID() string {
	return class.MemberDomainID
}

// CreateMember creates new member as a child to domain storage.
func (md *memberDomain) CreateMember(fc factory.Factory) (string, error) {
	member, err := fc.Create(md)
	if err != nil {
		return "", err
	}
	if member == nil {
		return "", fmt.Errorf("factory returns nil")
	}

	record, err := md.ChildStorage.Set(member)
	if err != nil {
		return "", err
	}

	return record, nil
}

// GetMember returns member from its record in domain storage.
func (md *memberDomain) GetMember(record string) (resolver.Proxy, error) {
	member, err := md.ChildStorage.Get(record)
	if err != nil {
		return nil, err
	}

	result, ok := member.(resolver.Proxy)
	if !ok {
		return nil, fmt.Errorf("object with record `%s` is not `Proxy` instance", record)
	}

	return result, nil
}

type memberDomainProxy struct {
	resolver.BaseProxy
}

// newMemberDomainProxy creates new proxy and associates it with new instance of MemberDomain.
func newMemberDomainProxy(parent object.Parent) (*memberDomainProxy, error) {
	instance, err := newMemberDomain(parent)
	if err != nil {
		return nil, err
	}
	return &memberDomainProxy{
		BaseProxy: resolver.BaseProxy{
			Instance: instance,
		},
	}, nil
}

// CreateMember is a proxy call for instance method.
func (mdp *memberDomainProxy) CreateMember(fc factory.Factory) (string, error) {
	return mdp.Instance.(MemberDomain).CreateMember(fc)
}

// GetMember is a proxy call for instance method.
func (mdp *memberDomainProxy) GetMember(record string) (resolver.Proxy, error) {
	return mdp.Instance.(MemberDomain).GetMember(record)
}

type memberDomainFactory struct {
	parent object.Parent
}

// NewMemberDomainFactory creates new factory for MemberDomain.
func NewMemberDomainFactory(parent object.Parent) factory.Factory {
	return &memberDomainFactory{
		parent: parent,
	}
}

// GetClassID returns string representation of MemberDomain's class.
func (mdf *memberDomainFactory) GetClassID() string {
	return class.MemberDomainID
}

// GetReference returns nil for not published factory.
func (mdf *memberDomainFactory) GetReference() object.Reference {
	return nil
}

// GetParent returns parent
func (mdf *memberDomainFactory) GetParent() object.Parent {
	// TODO: return real parent, fix tests
	return nil
}

// Create is a factory method for new MemberDomain instances.
func (mdf *memberDomainFactory) Create(parent object.Parent) (resolver.Proxy, error) {
	proxy, err := newMemberDomainProxy(parent)
	if err != nil {
		return nil, err
	}

	_, err = parent.AddChild(proxy)
	if err != nil {
		return nil, err
	}
	return proxy, nil
}
