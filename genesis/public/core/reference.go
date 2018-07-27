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

package core

import (
	"fmt"

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/factory"
	"github.com/insolar/insolar/genesis/model/object"
)

// ReferenceDomainName is a name for reference domain.
const ReferenceDomainName = "ReferenceDomain"

// ReferenceDomain is a contract that allows to publish and resolve global references.
type ReferenceDomain interface {
	// Base domain implementation.
	domain.Domain
	// RegisterReference is used to publish new global references
	RegisterReference(*object.Reference) (string, error)
	// ResolveReference provides reference instance from record
	ResolveReference(string) (*object.Reference, error)
}

type referenceDomain struct {
	domain.BaseDomain
}

// newReferenceDomain creates new instance of ReferenceDomain
func newReferenceDomain(parent object.Parent) *referenceDomain {
	refDomain := &referenceDomain{
		BaseDomain: *domain.NewBaseDomain(parent, ReferenceDomainName),
	}
	// Bootstrap case
	if parent == nil {
		refDomain.Parent = refDomain
	}
	return refDomain
}

// GetClassID return string representation of ReferenceDomain's class.
func (rd *referenceDomain) GetClassID() string {
	return class.ReferenceDomainID
}

// RegisterReference set new reference as a child to domain storage.
func (rd *referenceDomain) RegisterReference(ref *object.Reference) (string, error) {
	record, err := rd.ChildStorage.Set(ref)
	if err != nil {
		return "", err
	}

	return record, nil
}

// ResolveReference returns reference from its record in domain storage.
func (rd *referenceDomain) ResolveReference(record string) (*object.Reference, error) {
	reference, err := rd.ChildStorage.Get(record)
	if err != nil {
		return nil, err
	}

	result, ok := reference.(*object.Reference)
	if !ok {
		return nil, fmt.Errorf("object with record `%s` is not `Reference` instance", record)
	}

	return result, nil
}

type referenceDomainProxy struct {
	instance *referenceDomain
}

// newReferenceDomainProxy creates new proxy and associate it with new instance of ReferenceDomain.
func newReferenceDomainProxy(parent object.Parent) *referenceDomainProxy {
	return &referenceDomainProxy{
		instance: newReferenceDomain(parent),
	}
}

// RegisterReference proxy call for instance method.
func (rdp *referenceDomainProxy) RegisterReference(address *object.Reference) (string, error) {
	return rdp.instance.RegisterReference(address)
}

// ResolveReference proxy call for instance method.
func (rdp *referenceDomainProxy) ResolveReference(record string) (*object.Reference, error) {
	return rdp.instance.ResolveReference(record)
}

// GetReference proxy call for instance method.
func (rdp *referenceDomainProxy) GetReference() *object.Reference {
	return rdp.instance.GetReference()
}

// GetParent proxy call for instance method.
func (rdp *referenceDomainProxy) GetParent() object.Parent {
	return rdp.instance.GetParent()
}

// GetClassID proxy call for instance method.
func (rdp *referenceDomainProxy) GetClassID() string {
	return class.ReferenceDomainID
}

type referenceDomainFactory struct {
	parent object.Parent
}

// NewReferenceDomainFactory creates new factory for ReferenceDomain.
func NewReferenceDomainFactory(parent object.Parent) factory.Factory {
	return &referenceDomainFactory{
		parent: parent,
	}
}

func (adf *referenceDomainFactory) GetParent() object.Parent {
	// TODO: return real parent, fix tests
	return nil
}

// GetClassID return string representation of ReferenceDomain's class.
func (adf *referenceDomainFactory) GetClassID() string {
	return class.ReferenceDomainID
}

// GetReference returns nil for not published factory
func (adf *referenceDomainFactory) GetReference() *object.Reference {
	return nil
}

// Create factory method for new ReferenceDomain instances.
func (adf *referenceDomainFactory) Create(parent object.Parent) object.Proxy {
	proxy := newReferenceDomainProxy(parent)
	_, err := parent.AddChild(proxy)
	if err != nil {
		return nil
	}
	return proxy
}
