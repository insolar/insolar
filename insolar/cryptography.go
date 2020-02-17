// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insolar

import (
	"crypto"
)

type Signature struct {
	raw []byte
}

func SignatureFromBytes(raw []byte) Signature {
	return Signature{raw: raw}
}

func (s *Signature) Bytes() []byte {
	return s.raw
}

//go:generate minimock -i github.com/insolar/insolar/insolar.CryptographyService -o ../testutils -s _mock.go -g
type CryptographyService interface {
	Signer
	GetPublicKey() (crypto.PublicKey, error)
	Verify(crypto.PublicKey, Signature, []byte) bool
}
