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

package object

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/stretchr/testify/assert"
)

type mockProxyReference struct {
}

func (c *mockProxyReference) GetClassID() string {
	return class.ReferenceID
}

func (c *mockProxyReference) GetClass() Proxy {
	return &ReferenceContainer{}
}

func (c *mockProxyReference) GetReference() Reference {
	return nil
}

func (c *mockProxyReference) SetReference(ref Reference) {
}

func (c *mockProxyReference) GetParent() Parent {
	return nil
}
func (c *mockProxyReference) GetChildStorage() storage.Storage {
	return nil
}

func (c *mockProxyReference) AddChild(child Child) (string, error) {
	return "", nil
}

func (c *mockProxyReference) GetChild(key string) (Child, error) {
	return child, nil
}

func (c *mockProxyReference) GetContext() []string {
	return []string{}
}

func (c *mockProxyReference) GetContextStorage() storage.Storage {
	return nil
}

type mockParentNotChild struct {
	ContextStorage storage.Storage
}

func (p *mockParentNotChild) GetClassID() string {
	return "mockParentNotChild"
}

func (p *mockParentNotChild) GetClass() Proxy {
	return nil
}

func (p *mockParentNotChild) GetChildStorage() storage.Storage {
	return nil
}

func (p *mockParentNotChild) AddChild(child Child) (string, error) {
	return "", nil
}

func (p *mockParentNotChild) GetChild(key string) (Child, error) {
	return child, nil
}

func (p *mockParentNotChild) GetContext() []string {
	return []string{}
}

func (p *mockParentNotChild) GetContextStorage() storage.Storage {
	return p.ContextStorage
}

func TestNewContextResolver(t *testing.T) {
	mockParent := &mockParentProxy{}
	resolver := newContextResolver(mockParent)

	assert.Equal(t, &contextResolver{
		parent: mockParent,
	}, resolver)
}

func TestContextResolver_GetObject_No_Object(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	mockParent := &mockParentProxy{
		ContextStorage: contextStorage,
	}
	resolver := newContextResolver(mockParent)
	ref, _ := NewReference("123", "143", ContextScope)

	obj, err := resolver.GetObject(ref, factory)

	assert.EqualError(t, err, "object with record 143 does not exist")
	assert.Nil(t, obj)
}

func TestContextResolver_GetObject_Not_Child(t *testing.T) {
	proxyRef := &mockProxyReference{}
	parentContextStorage := storage.NewMapStorage()
	record, _ := parentContextStorage.Set(proxyRef)
	parent := &mockParentNotChild{
		ContextStorage: parentContextStorage,
	}

	resolver := newContextResolver(parent)
	ref, _ := NewReference("1", record, ContextScope)

	obj, err := resolver.GetObject(ref, factory)

	assert.EqualError(t, err, fmt.Sprintf("object with name #1.#%s does not exist", record))
	assert.Nil(t, obj)
}

func TestContextResolver_GetObject_Not_Reference(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	mockParent := &mockParentProxy{
		ContextStorage: contextStorage,
	}
	resolver := newContextResolver(mockParent)

	obj, err := resolver.GetObject("not reference", "mockChild")

	assert.EqualError(t, err, "reference is not Reference class object")
	assert.Nil(t, obj)
}

func TestContextResolver_GetObject_Wrong_Class(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	record, _ := contextStorage.Set(child)
	mockParent := &mockParentProxy{
		ContextStorage: contextStorage,
	}
	resolver := newContextResolver(mockParent)
	ref, _ := NewReference("1", record, ContextScope)

	obj, err := resolver.GetObject(ref, ref)

	assert.EqualError(t, err, "instance class is not equal received")
	assert.Nil(t, obj)
}

func TestContextResolver_GetObject(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	record, _ := contextStorage.Set(child)
	mockParent := &mockParentProxy{
		ContextStorage: contextStorage,
	}
	resolver := newContextResolver(mockParent)
	ref, _ := NewReference("1", record, ContextScope)

	obj, err := resolver.GetObject(ref, factory)

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}
