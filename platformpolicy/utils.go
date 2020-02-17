// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
