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

package platformpolicy

import (
	"crypto"
)

// MustNormalizePublicKey parses public key in PEM format, returns normalized (stable) public key value.
// Panics on error.
func MustNormalizePublicKey(b []byte) string {
	ks := NewKeyProcessor()
	pubKey, err := ks.ImportPublicKeyPEM(b)
	if err != nil {
		panic(err)
	}
	return MustPublicKeyToString(pubKey)
}

// MustPublicKeyToBytes returns byte representation of public key.
// Panics on error.
func MustPublicKeyToBytes(key crypto.PublicKey) []byte {
	ks := NewKeyProcessor()
	b, err := ks.ExportPublicKeyPEM(key)
	if err != nil {
		panic(err)
	}
	return b
}

// MustPublicKeyToString returns string representation of public key.
// Panics on error.
func MustPublicKeyToString(key crypto.PublicKey) string {
	return string(MustPublicKeyToBytes(key))
}
