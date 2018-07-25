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
	"testing"

	"github.com/insolar/insolar/genesis/model/class"
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/object"
	"github.com/insolar/insolar/genesis/model/resolver"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

type mockProxy struct {
	parent object.Parent
}

func (p *mockProxy) GetClassID() string {
	return "mockProxy"
}

func (p *mockProxy) GetReference() *object.Reference {
	return nil
}

func (p *mockProxy) GetParent() object.Parent {
	return p.parent
}

type mockFactory struct{}

func (f *mockFactory) Create(parent object.Parent) resolver.Proxy {
	return &mockProxy{
		parent: parent,
	}
}

func (f *mockFactory) GetClassID() string {
	return "mockFactory"
}

func (f *mockFactory) GetReference() *object.Reference {
	return nil
}

func TestNewInstanceDomain(t *testing.T) {
	parent := &mockParent{}
	instDomain := newInstanceDomain(parent)

	assert.Equal(t, &instanceDomain{
		BaseDomain: *domain.NewBaseDomain(parent, InstanceDomainName),
	}, instDomain)
}

func TestInstanceDomain_GetClassID(t *testing.T) {
	parent := &mockParent{}
	instDomain := newInstanceDomain(parent)
	domainID := instDomain.GetClassID()
	assert.Equal(t, class.InstanceDomainID, domainID)
}

func TestCreateInstance(t *testing.T) {
	parent := &mockParent{}
	instDomain := newInstanceDomain(parent)
	factory := &mockFactory{}

	registered, err := instDomain.CreateInstance(factory)
	assert.NoError(t, err)

	_, err = uuid.FromString(registered)
	assert.NoError(t, err)
}

func TestGetInstance(t *testing.T) {
	parent := &mockParent{}
	instDomain := newInstanceDomain(parent)
	factory := &mockFactory{}

	registered, err := instDomain.CreateInstance(factory)
	assert.NoError(t, err)

	resolved, err := instDomain.GetInstance(registered)
	assert.NoError(t, err)

	assert.Equal(t, &mockProxy{
		parent: instDomain,
	}, resolved)
}
