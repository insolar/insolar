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

package resolver

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/stretchr/testify/assert"
)

type mockProxyReference struct {
}

func (c *mockProxyReference) GetClassID() string {
	return class.ReferenceID
}

func (c *mockProxyReference) GetReference() *object.Reference {
	return nil
}

func (c *mockProxyReference) GetParent() object.Parent {
	return nil
}
func (c *mockProxyReference) GetChildStorage() storage.Storage {
	return nil
}

func (c *mockProxyReference) AddChild(child object.Child) (string, error) {
	return "", nil
}

func (c *mockProxyReference) GetChild(key string) (object.Child, error) {
	return child, nil
}

func (c *mockProxyReference) GetContext() []string {
	return []string{}
}

func (c *mockProxyReference) GetContextStorage() storage.Storage {
	return nil
}

type mockParentNotChild struct {
	Reference      *object.Reference
	ContextStorage storage.Storage
}

func (p *mockParentNotChild) GetClassID() string {
	return "mockParentNotChild"
}

func (p *mockParentNotChild) GetReference() *object.Reference {
	return p.Reference
}

func (p *mockParentNotChild) GetChildStorage() storage.Storage {
	return nil
}

func (p *mockParentNotChild) AddChild(child object.Child) (string, error) {
	return "", nil
}

func (p *mockParentNotChild) GetChild(key string) (object.Child, error) {
	return child, nil
}

func (p *mockParentNotChild) GetContext() []string {
	return []string{}
}

func (p *mockParentNotChild) GetContextStorage() storage.Storage {
	return p.ContextStorage
}

func TestNewContextResolver(t *testing.T) {
	mockParent := &mockParent{}
	resolver := newContextResolver(mockParent)

	assert.Equal(t, &contextResolver{
		parent: mockParent,
	}, resolver)
}

func TestContextResolver_GetObject_No_Object(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	mockParent := &mockParent{
		ContextStorage: contextStorage,
	}
	resolver := newContextResolver(mockParent)
	ref, _ := object.NewReference("1", "1", object.ContextScope)

	obj, err := resolver.GetObject(ref, "someClass")

	assert.EqualError(t, err, "object with record 1 does not exist")
	assert.Nil(t, obj)
}

func TestContextResolver_GetObject_Wrong_classID(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	record, _ := contextStorage.Set(child)
	mockParent := &mockParent{
		ContextStorage: contextStorage,
	}
	resolver := newContextResolver(mockParent)
	ref, _ := object.NewReference(record, "1", object.ContextScope)

	obj, err := resolver.GetObject(ref, "someClass")

	assert.EqualError(t, err, "instance class is not `someClass`")
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
	ref, _ := object.NewReference(record, "1", object.ContextScope)

	obj, err := resolver.GetObject(ref, "someClass")

	assert.EqualError(t, err, fmt.Sprintf("object with name #1.#%s does not exist", record))
	assert.Nil(t, obj)
}

func TestContextResolver_GetObject_Not_Reference(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	mockParent := &mockParent{
		ContextStorage: contextStorage,
	}
	resolver := newContextResolver(mockParent)

	obj, err := resolver.GetObject("not reference", "mockChild")

	assert.EqualError(t, err, "reference is not Reference class object")
	assert.Nil(t, obj)
}

func TestContextResolver_GetObject_ClassID_Not_Str(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	record, _ := contextStorage.Set(child)
	mockParent := &mockParent{
		ContextStorage: contextStorage,
	}
	resolver := newContextResolver(mockParent)
	ref, _ := object.NewReference(record, "1", object.ContextScope)

	obj, err := resolver.GetObject(ref, ref)

	assert.EqualError(t, err, "classID is not string")
	assert.Nil(t, obj)
}

func TestContextResolver_GetObject(t *testing.T) {
	contextStorage := storage.NewMapStorage()
	record, _ := contextStorage.Set(child)
	mockParent := &mockParent{
		ContextStorage: contextStorage,
	}
	resolver := newContextResolver(mockParent)
	ref, _ := object.NewReference(record, "1", object.ContextScope)

	obj, err := resolver.GetObject(ref, "mockChild")

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}
