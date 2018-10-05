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

package object

import (
	"fmt"
	"testing"

	"github.com/insolar/insolar/genesis/mock/storage"
	"github.com/stretchr/testify/assert"
)

type mockProxy struct {
	reference Reference
}

func (p *mockProxy) GetClassID() string {
	return "mockProxy"
}

func (p *mockProxy) GetClass() Proxy {
	return nil
}

func (p *mockProxy) GetReference() Reference {
	return p.reference
}

func (p *mockProxy) SetReference(reference Reference) {
	p.reference = reference
}

type mockFactory struct {
}

func (f *mockFactory) Create(parent Parent) (Proxy, error) {
	return &mockChildProxy{
		parent: parent,
	}, nil
}

func (f *mockFactory) GetClassID() string {
	return "mockFactory"
}

func (f *mockFactory) GetClass() Proxy {
	return f
}

func (f *mockFactory) GetReference() Reference {
	return nil
}

func (f *mockFactory) GetParent() Parent {
	return nil
}

func (f *mockFactory) SetReference(reference Reference) {

}

var factory = &mockFactory{}

type mockChildProxy struct {
	mockProxy
	ContextStorage storage.Storage
	parent         Parent
	class          Proxy
}

func (c *mockChildProxy) GetClassID() string {
	return "mockChild"
}

func (c *mockChildProxy) GetClass() Proxy {
	return c.class
}

func (c *mockChildProxy) GetParent() Parent {
	return c.parent
}

var child = &mockChildProxy{
	class: factory,
}

type mockParentProxy struct {
	mockProxy
	ContextStorage storage.Storage
	parent         Parent
}

func (p *mockParentProxy) GetParent() Parent {
	return p.parent
}

func (p *mockParentProxy) GetClassID() string {
	return "mockParent"
}

func (p *mockParentProxy) GetChildStorage() storage.Storage {
	return nil
}

func (p *mockParentProxy) AddChild(child Child) (string, error) {
	return "", nil
}

func (p *mockParentProxy) GetChild(key string) (Child, error) {
	return child, nil
}

func (p *mockParentProxy) GetContext() []string {
	return []string{}
}

func (p *mockParentProxy) GetContextStorage() storage.Storage {
	return p.ContextStorage
}

type mockParentWithError struct {
	mockParentProxy
}

func (p *mockParentWithError) GetChild(key string) (Child, error) {
	return nil, fmt.Errorf("object with record %s does not exist", key)
}

func TestNewChildResolver(t *testing.T) {
	mockParent := &mockParentProxy{}
	mapStorage := newChildResolver(mockParent)

	assert.Equal(t, &childResolver{
		parent: mockParent,
	}, mapStorage)
}

func TestChildResolver_GetObject_Not_Reference(t *testing.T) {
	mockParent := &mockParentWithError{}
	resolver := newChildResolver(mockParent)

	obj, err := resolver.GetObject("not reference", "mockParent")

	assert.EqualError(t, err, "reference is not Reference class object")
	assert.Nil(t, obj)
}

func TestChildResolver_GetObject_No_Object(t *testing.T) {
	mockParent := &mockParentWithError{}
	resolver := newChildResolver(mockParent)
	ref, _ := NewReference("1", "1", ChildScope)

	obj, err := resolver.GetObject(ref, factory)

	assert.EqualError(t, err, "object with record 1 does not exist")
	assert.Nil(t, obj)
}

func TestChildResolver_GetObject_Wrong_Class(t *testing.T) {
	mockParent := &mockParentProxy{}
	resolver := newChildResolver(mockParent)
	ref, _ := NewReference("1", "1", ChildScope)

	obj, err := resolver.GetObject(ref, ref)

	assert.EqualError(t, err, "instance class is not equal received")
	assert.Nil(t, obj)
}

func TestChildResolver_GetObject(t *testing.T) {
	mockParent := &mockParentProxy{}
	resolver := newChildResolver(mockParent)
	ref, _ := NewReference("1", "1", ChildScope)

	obj, err := resolver.GetObject(ref, factory)

	assert.NoError(t, err)
	assert.Equal(t, child, obj)
}
