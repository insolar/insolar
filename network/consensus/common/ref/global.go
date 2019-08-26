///
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
///

package ref

import (
	"github.com/insolar/insolar/longbits"
	"io"
)

/* For LIMITED USE ONLY - can only be used by observer/analytical code */
func NewRecordRef(recID Local) Global {
	if recID.getScope() != 0 {
		panic("illegal value")
	}
	return Global{addressLocal: recID}
}

func NewSelfRef(localID Local) Global {
	if localID.getScope() != baseScopeReserved {
		panic("illegal value")
	}
	return Global{addressLocal: localID, addressBase: localID}
}

func New(domainID, localID Local) Global {
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
