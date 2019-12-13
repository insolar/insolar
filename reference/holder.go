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

type Holder interface {
	GetScope() Scope

	// GetBase returns base portion of a full reference
	GetBase() *Local
	// GetLocal returns local portion of a full reference
	GetLocal() *Local

	IsEmpty() bool
}

func IsRecordScope(ref Holder) bool {
	return ref.GetBase().IsEmpty() && !ref.GetLocal().IsEmpty() && ref.GetLocal().getScope() == baseScopeLifeline
}

func IsObjectReference(ref Holder) bool {
	return !ref.GetBase().IsEmpty() && !ref.GetLocal().IsEmpty() && ref.GetLocal().getScope() == baseScopeLifeline
}

func IsSelfScope(ref Holder) bool {
	return ref.GetBase() == ref.GetLocal() || *ref.GetBase() == *ref.GetLocal()
}

func IsLifelineScope(ref Holder) bool {
	return ref.GetBase().getScope() == baseScopeLifeline
}

func IsLocalDomainScope(ref Holder) bool {
	return ref.GetBase().getScope() == baseScopeLocalDomain
}

func IsGlobalScope(ref Holder) bool {
	return ref.GetBase().getScope() == baseScopeGlobal
}

func Equal(ref0, ref1 Holder) bool {
	if p := ref1.GetLocal(); p == nil || !ref0.GetLocal().Equal(*p) {
		return false
	}
	if p := ref1.GetBase(); p == nil || !ref0.GetBase().Equal(*p) {
		return false
	}
	return true
}

func Compare(ref0, ref1 Holder) int {
	if cmp := ref0.GetBase().Compare(*ref1.GetBase()); cmp != 0 {
		return cmp
	}
	return ref0.GetLocal().Compare(*ref1.GetLocal())
}
