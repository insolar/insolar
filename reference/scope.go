// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
