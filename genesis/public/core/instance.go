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

// InstanceDomainName is a name for instance domain.
const InstanceDomainName = "InstanceDomain"

// InstanceDomain is a contract that...
type InstanceDomain interface {
	// Base domain implementation.
	domain.Domain
	// CreateInstance is used to...
	CreateInstance(*factory.Factory) (string, error)
	// GetInstance provides...
	GetInstance(string) (*factory.Factory, error)
}

type instanceDomain struct {
	domain.BaseDomain
}

// newInstanceDomain creates new instance of InstanceDomain
func newInstanceDomain(parent object.Parent) *instanceDomain {
	instDomain := &instanceDomain{
		BaseDomain: *domain.NewBaseDomain(parent, InstanceDomainName),
	}
	return instDomain
}

// GetClassID return string representation of InstanceDomain's class.
func (id *instanceDomain) GetClassID() string {
	return class.InstanceDomainID
}

// CreateInstance create new instance as a child to domain storage.
func (id *instanceDomain) CreateInstance(fc factory.Factory) (string, error) {
	instance := fc.Create(id)
	record, err := id.ChildStorage.Set(instance)
	if err != nil {
		return "", err
	}

	return record, nil
}

// GetInstance returns instance from its record in domain storage.
func (id *instanceDomain) GetInstance(record string) (object.Proxy, error) {
	instance, err := id.ChildStorage.Get(record)
	if err != nil {
		return nil, err
	}

	result, ok := instance.(object.Proxy)
	if !ok {
		return nil, fmt.Errorf("object with record `%s` is not `Proxy` instance", record)
	}

	return result, nil
}
