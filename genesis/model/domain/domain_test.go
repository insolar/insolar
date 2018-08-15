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

package domain

import (
	"testing"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/contract"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/stretchr/testify/assert"
)

type mockParent struct {
	ContextStorage storage.Storage
}

func (p *mockParent) GetClassID() string {
	return "mockParent"
}

func (p *mockParent) GetClass() object.Factory {
	return nil
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

type mockProxy struct {
	reference object.Reference
}

func (p *mockProxy) GetClassID() string {
	return "mockProxy"
}

func (p *mockProxy) GetClass() object.Factory {
	return nil
}

func (p *mockProxy) GetReference() object.Reference {
	return p.reference
}

func (p *mockProxy) SetReference(reference object.Reference) {
	p.reference = reference
}

type mockChildProxy struct {
	mockProxy
	ContextStorage storage.Storage
	parent         object.Parent
}

func (c *mockChildProxy) GetClassID() string {
	return "mockChild"
}

func (c *mockChildProxy) GetParent() object.Parent {
	return c.parent
}

var child = &mockChildProxy{}

type mockFactory struct {
}

func (f *mockFactory) Create(parent object.Parent) (object.Proxy, error) {
	return &mockChildProxy{
		parent: parent,
	}, nil
}

func (f *mockFactory) GetClassID() string {
	return "mockFactory"
}

func (f *mockFactory) GetClass() object.Factory {
	return &mockFactory{}
}

func (f *mockFactory) GetReference() object.Reference {
	return nil
}

func (f *mockFactory) GetParent() object.Parent {
	return nil
}

func (f *mockFactory) SetReference(reference object.Reference) {

}

func TestNewBaseDomain(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}

	domain := NewBaseDomain(parent, factory, "NewDomain")

	sc := contract.BaseSmartContract{
		CompositeMap: make(map[string]object.Composite),
		ChildStorage: storage.NewMapStorage(),
		Parent:       parent,
	}

	assert.Equal(t, &BaseDomain{
		BaseSmartContract: sc,
		Name:              "NewDomain",
		class:             factory,
	}, domain)
}

func TestBaseDomain_GetClassID(t *testing.T) {
	parent := &mockParent{}
	domain := NewBaseDomain(parent, &mockFactory{}, "NewDomain")

	classID := domain.GetClassID()

	assert.Equal(t, class.DomainID, classID)
}

func TestBaseDomain_GetClass(t *testing.T) {
	factory := &mockFactory{}
	parent := &mockParent{}
	domain := NewBaseDomain(parent, factory, "NewDomain")

	assert.Equal(t, factory, domain.GetClass())
}

func TestBaseDomain_GetName(t *testing.T) {
	parent := &mockParent{}
	domain := NewBaseDomain(parent, &mockFactory{}, "NewDomain")

	name := domain.GetName()

	assert.Equal(t, "NewDomain", name)
}
