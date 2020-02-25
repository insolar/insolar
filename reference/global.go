// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package reference

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/longbits"
)

const (
	LocalBinaryHashSize          = 28
	LocalBinaryPulseAndScopeSize = 4
	LocalBinarySize              = LocalBinaryPulseAndScopeSize + LocalBinaryHashSize
	GlobalBinarySize             = 2 * LocalBinarySize
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

// IsEmpty - check for void
func (v Global) IsEmpty() bool {
	return v.addressLocal.IsEmpty() && v.addressBase.IsEmpty()
}

func (v *Global) IsRecordScope() bool {
	return v.addressBase.IsEmpty() && !v.addressLocal.IsEmpty() && v.addressLocal.getScope() == baseScopeLifeline
}

func (v *Global) IsObjectReference() bool {
	return !v.addressBase.IsEmpty() && !v.addressLocal.IsEmpty() && v.addressLocal.getScope() == baseScopeLifeline
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

// GetBase returns base address from Global.
func (v Global) GetBase() *Local {
	return &v.addressBase
}

// GetLocal returns local address from Global.
func (v Global) GetLocal() *Local {
	return &v.addressLocal
}

// Encode encodes Global to string with chosen encoder.
func (v Global) Encode(enc Encoder) string {
	repr, err := enc.Encode(&v)
	if err != nil {
		return NilRef
	}
	return repr
}

// String outputs base64 Reference representation.
func (v Global) String() string {
	return v.Encode(DefaultEncoder())
}

// Bytes returns byte slice of Reference
func (v Global) Bytes() []byte {
	return v.AsBytes()
}

// Equal checks if reference points to the same record.
func (v Global) Equal(other Global) bool {
	return v.addressLocal.Equal(other.addressLocal) && v.addressBase.Equal(other.addressBase)
}

// Compare compares two record references
func (v Global) Compare(other Global) int {
	if cmp := v.addressBase.Compare(other.addressBase); cmp != 0 {
		return cmp
	}
	return v.addressLocal.Compare(other.addressLocal)
}

// MarshalJSON serializes reference into JSONFormat.
func (v *Global) MarshalJSON() ([]byte, error) {
	if v == nil {
		return json.Marshal(nil)
	}
	return json.Marshal(v.String())
}

func (v Global) Marshal() ([]byte, error) {
	return v.AsBytes(), nil
}

func (v Global) MarshalBinary() ([]byte, error) {
	return v.Marshal()
}

func (v Global) MarshalTo(data []byte) (int, error) {
	if len(data) < GlobalBinarySize {
		return 0, errors.New("not enough bytes to marshal reference.Global")
	}
	return copy(data, v.AsBytes()), nil
}

func (v *Global) UnmarshalJSON(data []byte) error {
	var repr interface{}

	err := json.Unmarshal(data, &repr)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal reference.Global")
	}

	switch realRepr := repr.(type) {
	case string:
		*v, err = DefaultDecoder().Decode(realRepr)
		if err != nil {
			return errors.Wrap(err, "failed to unmarshal reference.Global")
		}
	case nil:
	default:
		return errors.Wrapf(err, "unexpected type %T when unmarshal reference.Global", repr)
	}

	return nil
}

func (v *Global) Unmarshal(data []byte) error {
	if len(data) < GlobalBinarySize {
		return errors.New("not enough bytes to unmarshal reference.Global")
	}

	{
		localWriter := v.addressLocal.asWriter()
		for i := 0; i < LocalBinarySize; i++ {
			_ = localWriter.WriteByte(data[i])
		}
	}

	{
		baseWriter := v.addressBase.asWriter()
		for i := 0; i < LocalBinarySize; i++ {
			_ = baseWriter.WriteByte(data[LocalBinarySize+i])
		}
	}

	return nil
}

func (v *Global) UnmarshalBinary(data []byte) error {
	return v.Unmarshal(data)
}

func (v Global) Size() int {
	return GlobalBinarySize
}
