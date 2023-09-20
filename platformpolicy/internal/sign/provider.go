package sign

import (
	"crypto"

	"github.com/insolar/insolar/insolar"
)

type ecdsaProvider struct {
}

func NewECDSAProvider() AlgorithmProvider {
	return &ecdsaProvider{}
}

func (p *ecdsaProvider) DataSigner(privateKey crypto.PrivateKey, hasher insolar.Hasher) insolar.Signer {
	return &ecdsaDataSignerWrapper{
		ecdsaDigestSignerWrapper: ecdsaDigestSignerWrapper{
			privateKey: MustConvertPrivateKeyToEcdsa(privateKey),
		},
		hasher: hasher,
	}
}
func (p *ecdsaProvider) DigestSigner(privateKey crypto.PrivateKey) insolar.Signer {
	return &ecdsaDigestSignerWrapper{
		privateKey: MustConvertPrivateKeyToEcdsa(privateKey),
	}
}

func (p *ecdsaProvider) DataVerifier(publicKey crypto.PublicKey, hasher insolar.Hasher) insolar.Verifier {
	return &ecdsaDataVerifyWrapper{
		ecdsaDigestVerifyWrapper: ecdsaDigestVerifyWrapper{
			publicKey: MustConvertPublicKeyToEcdsa(publicKey),
		},
		hasher: hasher,
	}
}

func (p *ecdsaProvider) DigestVerifier(publicKey crypto.PublicKey) insolar.Verifier {
	return &ecdsaDigestVerifyWrapper{
		publicKey: MustConvertPublicKeyToEcdsa(publicKey),
	}
}
