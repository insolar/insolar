package insolar

import "crypto"

type KeyStore interface {
	GetPrivateKey(string) (crypto.PrivateKey, error)
}
