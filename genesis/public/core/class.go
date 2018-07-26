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
	"github.com/insolar/insolar/genesis/model/domain"
	"github.com/insolar/insolar/genesis/model/factory"
	"github.com/insolar/insolar/genesis/model/object"
)

// ClassDomainName is a name for class domain.
const ClassDomainName = "ClassDomain"

// ClassDomain is a contract that allows to publish new classes (e.g. new contract types).
type ClassDomain interface {
	// Base domain implementation.
	domain.Domain
	// RegisterClass is used to publish new .
	RegisterClass(factory.Factory) (string, error)
	// GetClass provides factory instance from record.
	GetClass(string) (factory.Factory, error)
}

type classDomain struct {
	domain.BaseDomain
}

// newClassDomain creates new instance of ClassDomain
func newClassDomain(parent object.Parent) *classDomain {
	classDomain := &classDomain{
		BaseDomain: *domain.NewBaseDomain(parent, ClassDomainName),
	}
	return classDomain
}
