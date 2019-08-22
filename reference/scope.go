//
// Copyright 2019 Insolar Technologies GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package reference

type Scope uint8

const ( // super-scopes
	baseScopeLifeline = iota
	baseScopeLocalDomain
	baseScopeReserved
	baseScopeGlobal
)

const ( // super-scopes
	LifelineSuperScope    Scope = 0x04 * baseScopeLifeline
	LocalDomainSuperScope Scope = 0x04 * baseScopeLocalDomain
	GlobalSuperScope      Scope = 0x04 * baseScopeGlobal
)

const SuperScopeMask = 0x0C
const SubScopeMask = 0x03

const (
	LifelineRecordOrSelf Scope = LifelineSuperScope + iota
	LifelinePrivateChild
	LifelinePublicChild
	LifelineDelegate
)

const (
	LocalDomainMember Scope = LocalDomainSuperScope + iota
	LocalDomainPrivatePolicy
	LocalDomainPublicPolicy
	_
)

const (
	RemoteDomainMember Scope = GlobalSuperScope + iota
	_
	GlobalDomainPublicPolicy
	GlobalDomainMember
)

func (v Scope) IsLocal() bool {
	return v&SuperScopeMask <= LocalDomainSuperScope
}

func (v Scope) IsOfLifeline() bool {
	return v&SuperScopeMask == LifelineSuperScope
}

func (v Scope) IsOfLocalDomain() bool {
	return v&SuperScopeMask == LocalDomainSuperScope
}

func (v Scope) IsGlobal() bool {
	return v&SuperScopeMask == GlobalSuperScope
}
