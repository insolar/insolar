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
	"testing"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/contract"
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/insolar/insolar/genesis/model/resolver"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type mockCallable struct {
	reference object.Reference
}

func (c *mockCallable) GetReference() object.Reference {
	return c.reference
}

func (c *mockCallable) SetReference(reference object.Reference) {
	c.reference = reference
}

type mockChild struct {
	mockCallable
	ContextStorage storage.Storage
	parent         object.Parent
}

func (c *mockChild) GetClassID() string {
	return "mockChild"
}

func (c *mockChild) GetParent() object.Parent {
	return c.parent
}

var child = &mockChild{}

type mockParent struct {
	mockCallable
	ContextStorage storage.Storage
	parent         object.Parent
}

func (p *mockParent) GetParent() object.Parent {
	return p.parent
}

func (p *mockParent) GetClassID() string {
	return "mockParent"
}

func (p *mockParent) GetChildStorage() storage.Storage {
	return nil
}

func (p *mockParent) AddChild(child object.Child) (string, error) {
	return "", nil
}

func (p *mockParent) GetChild(key string) (object.Child, error) {
	return child, nil
}

func (p *mockParent) GetContext() []string {
	return []string{}
}

func (p *mockParent) GetContextStorage() storage.Storage {
	return p.ContextStorage
}

type mockParentWithError struct {
	mockParent
}

func (p *mockParentWithError) AddChild(child object.Child) (string, error) {
	return "", fmt.Errorf("add child error")
}

var globalParent = &mockParent{}

type mockDomain struct {
	mockCallable
	mockParent
	mockChild
}

func (d *mockDomain) GetClassID() string {
	return "mockDomain"
}

func (d *mockDomain) GetParent() object.Parent {
	return d.mockChild.parent
}

// Create map for global resolving
var globalResolverMap = make(map[string]resolver.Proxy)
var domainString = "123"
var initRefObject, _ = object.NewReference(domainString, "1", object.GlobalScope)

// Create referenceDomain and its proxy
var initRefDomain = newReferenceDomain(nil)
var initRefDomainProxy = newReferenceDomainProxy(globalParent)

// Set one map for Handler and ReferenceDomain
func init() {
	globalResolverMap["123"] = &mockDomain{}
	// Create Handler empty instance
	resolverHandler := resolver.NewHandler(nil)
	// Set map to Handler.GlobalResolver
	resolverHandler.InitGlobalMap(&globalResolverMap)
	// Set map to ReferenceDomain.globalResolverMap
	initRefDomain.InitGlobalMap(&globalResolverMap)
	// Set map to ReferenceDomainProxy.ReferenceDomain.globalResolverMap
	initRefDomainProxy.Instance.(ReferenceDomain).InitGlobalMap(&globalResolverMap)

}

func TestNewReferenceDomain(t *testing.T) {
	parent := &mockParent{}
	refDomain := newReferenceDomain(parent)

	assert.Equal(t, &referenceDomain{
		BaseDomain: *domain.NewBaseDomain(parent, ReferenceDomainName),
	}, refDomain)
}

func TestNewReferenceDomain_WithNoParent(t *testing.T) {
	refDomain := newReferenceDomain(nil)
	expected := &referenceDomain{
		BaseDomain: *domain.NewBaseDomain(nil, ReferenceDomainName),
	}
	expected.Parent = refDomain

	assert.Equal(t, expected, refDomain)
}

func TestReferenceDomain_GetClassID(t *testing.T) {
	refDomain := newReferenceDomain(nil)
	domainID := refDomain.GetClassID()
	assert.Equal(t, class.ReferenceDomainID, domainID)
}

func TestReferenceDomain_SetMap(t *testing.T) {
	refDomain := newReferenceDomain(nil)
	refDomain.globalResolverMap = nil

	newMap := make(map[string]resolver.Proxy)
	refDomain.InitGlobalMap(&newMap)

	assert.Equal(t, &newMap, refDomain.globalResolverMap)
}

func TestReferenceDomain_RegisterReference(t *testing.T) {
	registered, err := initRefDomain.RegisterReference(initRefObject, "mockChild")

	assert.NoError(t, err)

	_, err = uuid.FromString(registered)
	assert.NoError(t, err)
}

func TestReferenceDomain_RegisterReference_GetObject_Error(t *testing.T) {
	refObject, err := object.NewReference("234", "1", object.GlobalScope)
	assert.NoError(t, err)

	registered, err := initRefDomain.RegisterReference(refObject, "classID")

	assert.Equal(t, "", registered)
	assert.EqualError(t, err, "reference with address `#234.#1` not found")
}

func TestReferenceDomain_ResolveReference(t *testing.T) {
	refObject, err := object.NewReference("123", "1", object.GlobalScope)
	assert.NoError(t, err)

	registered, err := initRefDomain.RegisterReference(refObject, "mockChild")
	assert.NoError(t, err)

	resolved, err := initRefDomain.ResolveReference(registered)

	assert.NoError(t, err)
	assert.Equal(t, refObject, resolved)
}

func TestReferenceDomain_ResolveReference_IncorrectRef(t *testing.T) {
	refDomain := newReferenceDomain(nil)
	_, err := refDomain.ResolveReference("1")
	assert.EqualError(t, err, "object with record 1 does not exist")
}

func TestNewReferenceDomainProxy(t *testing.T) {
	parent := &mockParent{}
	refDomainProxy := newReferenceDomainProxy(parent)

	assert.Equal(t, &referenceDomainProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: newReferenceDomain(parent),
		},
	}, refDomainProxy)
}

func TestReferenceDomainProxy_SetMap(t *testing.T) {
	refDomain := newReferenceDomainProxy(nil)
	refDomain.Instance.(*referenceDomain).globalResolverMap = nil

	newMap := make(map[string]resolver.Proxy)
	refDomain.InitGlobalMap(&newMap)

	assert.Equal(t, &newMap, refDomain.Instance.(*referenceDomain).globalResolverMap)
}

func TestReferenceDomainProxy_RegisterReference(t *testing.T) {
	record := "1"
	refObject, err := object.NewReference(domainString, record, object.GlobalScope)
	assert.NoError(t, err)

	registered, err := initRefDomainProxy.RegisterReference(refObject, "mockChild")
	assert.NoError(t, err)

	_, err = uuid.FromString(registered)
	assert.NoError(t, err)
}

func TestReferenceDomainProxy_ResolveReference(t *testing.T) {
	refObject, err := object.NewReference(domainString, "1", object.GlobalScope)
	assert.NoError(t, err)

	registered, err := initRefDomainProxy.RegisterReference(refObject, "mockChild")
	assert.NoError(t, err)

	resolved, err := initRefDomainProxy.ResolveReference(registered)

	assert.NoError(t, err)
	assert.Equal(t, refObject, resolved)
}

func TestReferenceDomainProxy_ResolveReference_IncorrectRef(t *testing.T) {
	parent := &mockParent{}
	refDomainProxy := newReferenceDomainProxy(parent)

	_, err := refDomainProxy.ResolveReference("1")
	assert.EqualError(t, err, "object with record 1 does not exist")
}

func TestNewReferenceDomainFactory(t *testing.T) {
	parent := &mockParent{}
	expected := &referenceDomainFactory{parent: parent}
	factory := NewReferenceDomainFactory(parent)

	assert.Equal(t, expected, factory)
}

func TestReferenceDomainFactory_GetClassID(t *testing.T) {
	parent := &mockParent{}
	factory := NewReferenceDomainFactory(parent)
	id := factory.GetClassID()

	assert.Equal(t, class.ReferenceDomainID, id)
}

func TestReferenceDomainFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory := NewReferenceDomainFactory(parent)
	proxy, err := factory.Create(parent)

	assert.NoError(t, err)
	assert.Equal(t, &referenceDomainProxy{
		BaseSmartContractProxy: contract.BaseSmartContractProxy{
			Instance: newReferenceDomain(parent),
		},
	}, proxy)
}

func TestReferenceDomainFactory_CreateWithError(t *testing.T) {
	parent := &mockParentWithError{}
	factory := NewReferenceDomainFactory(parent)
	proxy, err := factory.Create(parent)

	assert.EqualError(t, err, "add child error")
	assert.Nil(t, proxy)
}
