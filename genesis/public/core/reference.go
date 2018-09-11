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

package core

import (
	"fmt"

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/contract"
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
)

// ReferenceDomainName is a name for reference domain.
const ReferenceDomainName = "ReferenceDomain"

// ReferenceDomain is a contract that allows to publish and resolve global references.
type ReferenceDomain interface {
	// Base domain implementation.
	domain.Domain
	// RegisterReference is used to publish new global references.
	RegisterReference(object.Reference, object.Proxy) (string, error)
	// ResolveReference provides reference instance from record.
	ResolveReference(string) (object.Reference, error)
	// InitGlobalMap sets globalResolverMap for references register/resolving.
	InitGlobalMap(globalInstanceMap *map[string]object.Proxy)
}

type referenceDomain struct {
	domain.BaseDomain
	globalResolverMap *map[string]object.Proxy
}

// newReferenceDomain creates new instance of ReferenceDomain.
func newReferenceDomain(parent object.Parent, class object.Factory) *referenceDomain {
	refDomain := &referenceDomain{
		BaseDomain: *domain.NewBaseDomain(parent, class, ReferenceDomainName),
	}
	// Bootstrap case
	if parent == nil {
		refDomain.Parent = refDomain
	}
	return refDomain
}

// GetClassID returns string representation of ReferenceDomain's class.
func (rd *referenceDomain) GetClassID() string {
	return class.ReferenceDomainID
}

// InitGlobalMap sets globalResolverMap for register/resolve references.
func (rd *referenceDomain) InitGlobalMap(globalInstanceMap *map[string]object.Proxy) {
	if rd.globalResolverMap != nil {
		return
	}
	rd.globalResolverMap = globalInstanceMap
}

// RegisterReference sets new reference as a child to domain storage.
func (rd *referenceDomain) RegisterReference(ref object.Reference, class object.Proxy) (string, error) {
	container := object.NewReferenceContainer(ref)
	record, err := rd.ChildStorage.Set(container)
	if err != nil {
		return "", err
	}
	res := rd.GetResolver()
	obj, err := res.GetObject(ref, class)
	if err != nil {
		return "", err
	}
	proxy, ok := obj.(object.Proxy)
	if !ok {
		return "", fmt.Errorf("object with reference `%s` is not `Proxy` instance", ref)
	}
	(*rd.globalResolverMap)[record] = proxy

	return record, nil
}

// ResolveReference returns reference from its record in domain storage.
func (rd *referenceDomain) ResolveReference(record string) (object.Reference, error) {
	reference, err := rd.ChildStorage.Get(record)
	if err != nil {
		return nil, err
	}

	container, ok := reference.(*object.ReferenceContainer)
	if !ok {
		return nil, fmt.Errorf("object with record `%s` is not `ReferenceContainer` instance", record)
	}

	return container.GetStoredReference(), nil
}

type referenceDomainProxy struct {
	contract.BaseSmartContractProxy
}

// newReferenceDomainProxy creates new proxy and associate it with new instance of ReferenceDomain.
func newReferenceDomainProxy(parent object.Parent, class object.Factory) *referenceDomainProxy {
	return &referenceDomainProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: newReferenceDomain(parent, class),
		},
	}
}

// RegisterReference is a proxy call for instance method.

func (rdp *referenceDomainProxy) RegisterReference(address object.Reference, class object.Proxy) (string, error) {
	return rdp.Instance.(ReferenceDomain).RegisterReference(address, class)
}

// ResolveReference is a proxy call for instance method.
func (rdp *referenceDomainProxy) ResolveReference(record string) (object.Reference, error) {
	return rdp.Instance.(ReferenceDomain).ResolveReference(record)
}

// InitGlobalMap is a proxy call for instance method.
func (rdp *referenceDomainProxy) InitGlobalMap(globalInstanceMap *map[string]object.Proxy) {
	rdp.Instance.(ReferenceDomain).InitGlobalMap(globalInstanceMap)
}

type referenceDomainFactory struct {
	object.BaseFactory
	parent object.Parent
}

// NewReferenceDomainFactory creates new factory for ReferenceDomain.
func NewReferenceDomainFactory(parent object.Parent) object.Factory {
	return &referenceDomainFactory{
		parent: parent,
	}
}

// GetParent returns parent
func (rdf *referenceDomainFactory) GetParent() object.Parent {
	// TODO: return real parent, fix tests
	return nil
}

// GetClassID returns string representation of ReferenceDomain's class.
func (rdf *referenceDomainFactory) GetClassID() string {
	return class.ReferenceDomainID
}

func (rdf *referenceDomainFactory) GetClass() object.Proxy {
	return rdf
}

// Create factory is a method for new ReferenceDomain instances.
func (rdf *referenceDomainFactory) Create(parent object.Parent) (object.Proxy, error) {
	proxy := newReferenceDomainProxy(parent, rdf)
	_, err := parent.AddChild(proxy)
	if err != nil {
		return nil, err
	}
	return proxy, nil
}
