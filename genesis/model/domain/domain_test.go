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
	Reference      *object.Reference
	ContextStorage storage.Storage
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

func TestNewBaseDomain(t *testing.T) {
	parent := &mockParent{}

	domain := NewBaseDomain(parent, "NewDomain")

	sc := contract.BaseSmartContract{
		CompositeMap: make(map[string]object.Composite),
		ChildStorage: storage.NewMapStorage(),
		Parent:       parent,
	}
	sc.GetResolver()
	assert.Equal(t, &BaseDomain{
		BaseSmartContract: sc,
		Name:              "NewDomain",
	}, domain)
}

func TestBaseDomain_GetClassID(t *testing.T) {
	parent := &mockParent{}
	domain := NewBaseDomain(parent, "NewDomain")

	classID := domain.GetClassID()

	assert.Equal(t, class.DomainID, classID)
}

func TestBaseDomain_GetName(t *testing.T) {
	parent := &mockParent{}
	domain := NewBaseDomain(parent, "NewDomain")

	name := domain.GetName()

	assert.Equal(t, "NewDomain", name)
}
