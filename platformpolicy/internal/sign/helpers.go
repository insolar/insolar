package sign

import (
	"crypto"
	"crypto/ecdsa"
)

func MustConvertPublicKeyToEcdsa(publicKey crypto.PublicKey) *ecdsa.PublicKey {
	ecdsaPublicKey, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		panic("Failed to convert public key to ecdsa public key")
	}
	return ecdsaPublicKey
}

func MustConvertPrivateKeyToEcdsa(privateKey crypto.PrivateKey) *ecdsa.PrivateKey {
	ecdsaPrivateKey, ok := privateKey.(*ecdsa.PrivateKey)
	if !ok {
		panic("Failed to convert private key to ecdsa private key")
	}
	return ecdsaPrivateKey
}
