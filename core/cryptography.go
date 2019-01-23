/*
 *    Copyright 2019 Insolar
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

//go:generate minimock -i github.com/insolar/insolar/core.CryptographyService -o ../testutils -s _mock.go
type CryptographyService interface {
	GetPublicKey() (crypto.PublicKey, error)
	Sign([]byte) (*Signature, error)
	Verify(crypto.PublicKey, Signature, []byte) bool
}
