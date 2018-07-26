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
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type mockParent struct {
	Reference      *object.Reference
	ContextStorage storage.Storage
	parent         object.Parent
}

func (p *mockParent) GetParent() object.Parent {
	return p.parent
}

func (p *mockParent) GetClassID() string {
	return "mockParent"
}

func (p *mockParent) GetReference() *object.Reference {
	return p.Reference
}

func (p *mockParent) GetChildStorage() storage.Storage {
	return nil
}

func (p *mockParent) AddChild(child object.Child) (string, error) {
	return "", nil
}

func (p *mockParent) GetChild(key string) (object.Child, error) {
	return nil, nil
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

func TestReferenceDomain_RegisterReference(t *testing.T) {
	refObject, err := object.NewReference("1", "1", object.GlobalScope)
	assert.NoError(t, err)

	refDomain := newReferenceDomain(nil)
	registered, err := refDomain.RegisterReference(refObject)

	assert.NoError(t, err)

	_, err = uuid.FromString(registered)
	assert.NoError(t, err)
}

func TestReferenceDomain_ResolveReference(t *testing.T) {
	refObject, err := object.NewReference("1", "1", object.GlobalScope)
	assert.NoError(t, err)

	refDomain := newReferenceDomain(nil)
	registered, err := refDomain.RegisterReference(refObject)
	assert.NoError(t, err)

	resolved, err := refDomain.ResolveReference(registered)

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
		instance: newReferenceDomain(parent),
	}, refDomainProxy)
}

func TestReferenceDomainProxy_RegisterReference(t *testing.T) {
	domain := "1"
	record := "1"
	refObject, err := object.NewReference(record, domain, object.GlobalScope)
	assert.NoError(t, err)

	parent := &mockParent{}
	refDomainProxy := newReferenceDomainProxy(parent)

	registered, err := refDomainProxy.RegisterReference(refObject)
	assert.NoError(t, err)

	_, err = uuid.FromString(registered)
	assert.NoError(t, err)
}

func TestReferenceDomainProxy_ResolveReference(t *testing.T) {
	refObject, err := object.NewReference("1", "1", object.GlobalScope)
	assert.NoError(t, err)

	parent := &mockParent{}
	refDomainProxy := newReferenceDomainProxy(parent)

	registered, err := refDomainProxy.RegisterReference(refObject)
	assert.NoError(t, err)

	resolved, err := refDomainProxy.ResolveReference(registered)

	assert.NoError(t, err)
	assert.Equal(t, refObject, resolved)
}

func TestReferenceDomainProxy_ResolveReference_IncorrectRef(t *testing.T) {
	parent := &mockParent{}
	refDomainProxy := newReferenceDomainProxy(parent)

	_, err := refDomainProxy.ResolveReference("1")
	assert.EqualError(t, err, "object with record 1 does not exist")
}

func TestReferenceDomainProxy_GetReference(t *testing.T) {
	parent := &mockParent{}
	refDomainProxy := newReferenceDomainProxy(parent)

	reference := refDomainProxy.GetReference()
	// TODO should return actual reference
	assert.Nil(t, reference)
}

func TestReferenceDomainProxy_GetParent(t *testing.T) {
	parent := &mockParent{}
	refDomainProxy := newReferenceDomainProxy(parent)

	returnedParent := refDomainProxy.GetParent()
	assert.Equal(t, parent, returnedParent)
}

func TestReferenceDomainProxy_GetClassID(t *testing.T) {
	parent := &mockParent{}
	refDomainProxy := newReferenceDomainProxy(parent)

	id := refDomainProxy.GetClassID()
	assert.Equal(t, class.ReferenceDomainID, id)
}

func TestNewReferenceDomainFactory(t *testing.T) {
	expected := &referenceDomainFactory{}
	factory := NewReferenceDomainFactory()

	assert.Equal(t, expected, factory)
}

func TestReferenceDomainFactory_GetClassID(t *testing.T) {
	factory := NewReferenceDomainFactory()
	id := factory.GetClassID()

	assert.Equal(t, class.ReferenceDomainID, id)
}

func TestReferenceDomainFactory_GetReference(t *testing.T) {
	factory := NewReferenceDomainFactory()
	reference := factory.GetReference()

	assert.Nil(t, reference)
}

func TestReferenceDomainFactory_Create(t *testing.T) {
	parent := &mockParent{}
	factory := NewReferenceDomainFactory()
	refDomainProxy := factory.Create(parent)

	assert.Equal(t, &referenceDomainProxy{
		instance: newReferenceDomain(parent),
	}, refDomainProxy)
}

func TestReferenceDomainFactory_CreateWithError(t *testing.T) {
	parent := &mockParentWithError{}
	factory := NewReferenceDomainFactory()
	refDomainProxy := factory.Create(parent)

	assert.Nil(t, refDomainProxy)
}
