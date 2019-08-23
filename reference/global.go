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

import (
	"io"

	"github.com/insolar/insolar/longbits"
)

/* For LIMITED USE ONLY - can only be used by observer/analytical code */
func NewRecordRef(recID Local) Global {
	if recID.getScope() != 0 {
		panic("illegal value")
	}
	return Global{addressLocal: recID}
}

func NewSelfRef(localID Local) Global {
	if localID.getScope() == baseScopeReserved {
		panic("illegal value")
	}
	return Global{addressLocal: localID, addressBase: localID}
}

func NewGlobal(domainID, localID Local) Global {
	return Global{addressLocal: localID, addressBase: domainID}
}

type Global struct {
	addressLocal Local
	addressBase  Local
}

func (v *Global) GetScope() Scope {
	return Scope(v.addressBase.getScope()<<2 | v.addressLocal.getScope())
}

func (v *Global) WriteTo(w io.Writer) (int64, error) {
	n, err := v.addressLocal.WriteTo(w)
	if err != nil {
		return n, err
	}
	n2, err := v.addressBase.WriteTo(w)
	return n + n2, err
}

func (v *Global) Read(p []byte) (int, error) {
	n, err := v.addressLocal.Read(p)
	if err != nil || n < v.addressLocal.len() {
		return n, err
	}
	n2, err := v.addressBase.Read(p[n:])
	return n + n2, err
}

func (v *Global) AsByteString() longbits.ByteString {
	return longbits.NewByteString(v.AsBytes())
}

func (v *Global) AsBytes() []byte {
	prefix := v.addressLocal.len()
	val := make([]byte, prefix+v.addressBase.len())
	_, _ = v.addressLocal.Read(val)
	_, _ = v.addressBase.Read(val[prefix:])
	return val
}

func (v *Global) IsEmpty() bool {
	return v.addressLocal.IsEmpty() && v.addressBase.IsEmpty()
}

func (v *Global) IsRecordScope() bool {
	return v.addressBase.IsEmpty() && !v.addressLocal.IsEmpty() && v.addressLocal.getScope() == baseScopeLifeline
}

func (v *Global) IsSelfScope() bool {
	return v.addressBase == v.addressLocal
}

func (v *Global) IsLifelineScope() bool {
	return v.addressBase.getScope() == baseScopeLifeline
}

func (v *Global) IsLocalDomainScope() bool {
	return v.addressBase.getScope() == baseScopeLocalDomain
}

func (v *Global) IsGlobalScope() bool {
	return v.addressBase.getScope() == baseScopeGlobal
}

func (v *Global) GetParity() []byte {
	panic("not implemented")
}

func (v *Global) CheckParity(bytes []byte) error {
	panic("not implemented")
}

/* ONLY for parser */
func (v *Global) tryConvertToSelf() bool {
	if !v.addressBase.IsEmpty() {
		return false
	}
	v.addressBase = v.addressLocal
	return true
}

/* ONLY for parser */
func (v *Global) tryApplyBase(base *Global) bool {
	if !v.addressBase.IsEmpty() {
		return false
	}

	if !base.IsSelfScope() {
		switch base.GetScope() {
		case LocalDomainMember, GlobalDomainMember:
			break
		default:
			return false
		}
	}
	v.addressBase = base.addressLocal
	return true
}
