package adapters

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/rand"
	"io"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/network/consensus/common/cryptkit"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/require"
)

var (
	processor  = platformpolicy.NewKeyProcessor()
	scheme     = platformpolicy.NewPlatformCryptographyScheme()
	pk, _      = processor.GeneratePrivateKey()
	privateKey = pk.(*ecdsa.PrivateKey)
	publicKey  = &privateKey.PublicKey
)

func TestNewSha3512Digester(t *testing.T) {
	digester := NewSha3512Digester(scheme)

	require.Implements(t, (*cryptkit.DataDigester)(nil), digester)

	require.Equal(t, digester.scheme, scheme)
}

func TestSha3512Digester_GetDigestOf(t *testing.T) {
	digester := NewSha3512Digester(scheme)

	b := make([]byte, 120)
	_, _ = rand.Read(b)
	reader := bytes.NewReader(b)

	digest := digester.GetDigestOf(reader)
	require.Equal(t, digest.FixedByteSize(), scheme.IntegrityHashSize())

	expected := scheme.IntegrityHasher().Hash(b)

	require.Equal(t, expected, digest.AsBytes())
}

func TestSha3512Digester_GetDigestMethod(t *testing.T) {
	digester := NewSha3512Digester(scheme)

	require.Equal(t, digester.GetDigestMethod(), SHA3512Digest)
}

func TestNewECDSAPublicKeyStore(t *testing.T) {
	ks := NewECDSAPublicKeyStore(publicKey)

	require.Implements(t, (*cryptkit.PublicKeyStore)(nil), ks)

	require.Equal(t, ks.publicKey, publicKey)
}

func TestECDSAPublicKeyStore_PublicKeyStore(t *testing.T) {
	ks := NewECDSAPublicKeyStore(publicKey)

	ks.PublicKeyStore()
}

func TestNewECDSASecretKeyStore(t *testing.T) {
	ks := NewECDSASecretKeyStore(privateKey)

	require.Implements(t, (*cryptkit.SecretKeyStore)(nil), ks)

	require.Equal(t, ks.privateKey, privateKey)
}

func TestECDSASecretKeyStore_PrivateKeyStore(t *testing.T) {
	ks := NewECDSASecretKeyStore(privateKey)

	ks.PrivateKeyStore()
}

func TestECDSASecretKeyStore_AsPublicKeyStore(t *testing.T) {
	ks := NewECDSASecretKeyStore(privateKey)

	expected := NewECDSAPublicKeyStore(publicKey)

	require.Equal(t, expected, ks.AsPublicKeyStore())
}

func TestNewECDSADigestSigner(t *testing.T) {
	ds := NewECDSADigestSigner(privateKey, scheme)

	require.Implements(t, (*cryptkit.DigestSigner)(nil), ds)

	require.Equal(t, ds.privateKey, privateKey)
	require.Equal(t, ds.scheme, scheme)
}

func TestECDSADigestSigner_SignDigest(t *testing.T) {
	ds := NewECDSADigestSigner(privateKey, scheme)
	digester := NewSha3512Digester(scheme)

	verifier := scheme.DigestVerifier(publicKey)

	b := make([]byte, 120)
	_, _ = rand.Read(b)
	reader := bytes.NewReader(b)

	digest := digester.GetDigestOf(reader)
	digestBytes := digest.AsBytes()

	signature := ds.SignDigest(digest)
	require.Equal(t, scheme.SignatureSize(), signature.FixedByteSize())
	require.Equal(t, signature.GetSignatureMethod(), SHA3512Digest.SignedBy(SECP256r1Sign))

	signatureBytes := signature.AsBytes()

	require.True(t, verifier.Verify(insolar.SignatureFromBytes(signatureBytes), digestBytes))
}

func TestECDSADigestSigner_GetSignMethod(t *testing.T) {
	ds := NewECDSADigestSigner(privateKey, scheme)

	require.Equal(t, ds.GetSignMethod(), SECP256r1Sign)
}

func TestNewECDSASignatureVerifier(t *testing.T) {
	digester := NewSha3512Digester(scheme)
	dv := NewECDSASignatureVerifier(digester, scheme, publicKey)

	require.Implements(t, (*cryptkit.SignatureVerifier)(nil), dv)

	require.Equal(t, dv.digester, digester)
	require.Equal(t, dv.scheme, scheme)
	require.Equal(t, dv.publicKey, publicKey)
}

func TestECDSASignatureVerifier_IsDigestMethodSupported(t *testing.T) {
	digester := NewSha3512Digester(scheme)
	dv := NewECDSASignatureVerifier(digester, scheme, publicKey)

	require.True(t, dv.IsDigestMethodSupported(SHA3512Digest))
	require.False(t, dv.IsDigestMethodSupported("SOME DIGEST METHOD"))
}

func TestECDSASignatureVerifier_IsSignMethodSupported(t *testing.T) {
	digester := NewSha3512Digester(scheme)
	dv := NewECDSASignatureVerifier(digester, scheme, publicKey)

	require.True(t, dv.IsSignMethodSupported(SECP256r1Sign))
	require.False(t, dv.IsSignMethodSupported("SOME SIGN METHOD"))
}

func TestECDSASignatureVerifier_IsSignOfSignatureMethodSupported(t *testing.T) {
	digester := NewSha3512Digester(scheme)
	dv := NewECDSASignatureVerifier(digester, scheme, publicKey)

	require.True(t, dv.IsSignOfSignatureMethodSupported(SHA3512Digest.SignedBy(SECP256r1Sign)))
	require.False(t, dv.IsSignOfSignatureMethodSupported("SOME SIGNATURE METHOD"))
	require.False(t, dv.IsSignOfSignatureMethodSupported(SHA3512Digest.SignedBy("SOME SIGN METHOD")))
	require.True(t, dv.IsSignOfSignatureMethodSupported(cryptkit.DigestMethod("SOME DIGEST METHOD").SignedBy(SECP256r1Sign)))
}

func TestECDSASignatureVerifier_IsValidDigestSignature(t *testing.T) {
	digester := NewSha3512Digester(scheme)
	dv := NewECDSASignatureVerifier(digester, scheme, publicKey)

	signer := scheme.DigestSigner(privateKey)

	b := make([]byte, 120)
	_, _ = rand.Read(b)
	reader := bytes.NewReader(b)

	digest := digester.GetDigestOf(reader)
	digestBytes := digest.AsBytes()

	signature, _ := signer.Sign(digestBytes)

	sig := cryptkit.NewSignature(longbits.NewBits512FromBytes(signature.Bytes()), SHA3512Digest.SignedBy(SECP256r1Sign))

	require.True(t, dv.IsValidDigestSignature(digest.AsDigestHolder(), sig.AsSignatureHolder()))
}

func TestECDSASignatureVerifier_IsValidDigestSignature_InvalidMethod(t *testing.T) {
	digester := NewSha3512Digester(scheme)
	dv := NewECDSASignatureVerifier(digester, scheme, publicKey)

	signer := scheme.DigestSigner(privateKey)

	b := make([]byte, 120)
	_, _ = rand.Read(b)
	reader := bytes.NewReader(b)

	digest := digester.GetDigestOf(reader)
	digestBytes := digest.AsBytes()

	signature, _ := signer.Sign(digestBytes)
	bits := longbits.NewBits512FromBytes(signature.Bytes())

	sig1 := cryptkit.NewSignature(bits, SHA3512Digest.SignedBy(SECP256r1Sign))
	require.True(t, dv.IsValidDigestSignature(digest.AsDigestHolder(), sig1.AsSignatureHolder()))

	sig2 := cryptkit.NewSignature(bits, "SOME DIGEST METHOD")
	require.False(t, dv.IsValidDigestSignature(digest.AsDigestHolder(), sig2.AsSignatureHolder()))

	sig3 := cryptkit.NewSignature(bits, SHA3512Digest.SignedBy("SOME SIGN METHOD"))
	require.False(t, dv.IsValidDigestSignature(digest.AsDigestHolder(), sig3.AsSignatureHolder()))

	sig4 := cryptkit.NewSignature(bits, cryptkit.DigestMethod("SOME DIGEST METHOD").SignedBy(SECP256r1Sign))
	require.False(t, dv.IsValidDigestSignature(digest.AsDigestHolder(), sig4.AsSignatureHolder()))
}

func TestECDSASignatureVerifier_IsValidDataSignature(t *testing.T) {
	digester := NewSha3512Digester(scheme)
	dv := NewECDSASignatureVerifier(digester, scheme, publicKey)

	signer := scheme.DigestSigner(privateKey)

	b := make([]byte, 120)
	_, _ = rand.Read(b)
	reader := bytes.NewReader(b)

	digest := digester.GetDigestOf(reader)
	digestBytes := digest.AsBytes()

	signature, _ := signer.Sign(digestBytes)

	sig := cryptkit.NewSignature(longbits.NewBits512FromBytes(signature.Bytes()), SHA3512Digest.SignedBy(SECP256r1Sign))

	_, _ = reader.Seek(0, io.SeekStart)
	require.True(t, dv.IsValidDataSignature(reader, sig.AsSignatureHolder()))
}

func TestECDSASignatureVerifier_IsValidDataSignature_InvalidMethod(t *testing.T) {
	digester := NewSha3512Digester(scheme)
	dv := NewECDSASignatureVerifier(digester, scheme, publicKey)

	signer := scheme.DigestSigner(privateKey)

	b := make([]byte, 120)
	_, _ = rand.Read(b)
	reader := bytes.NewReader(b)

	digest := digester.GetDigestOf(reader)
	digestBytes := digest.AsBytes()

	signature, _ := signer.Sign(digestBytes)

	bits := longbits.NewBits512FromBytes(signature.Bytes())

	_, _ = reader.Seek(0, io.SeekStart)
	sig1 := cryptkit.NewSignature(bits, SHA3512Digest.SignedBy(SECP256r1Sign))
	require.True(t, dv.IsValidDataSignature(reader, sig1.AsSignatureHolder()))

	_, _ = reader.Seek(0, io.SeekStart)
	sig2 := cryptkit.NewSignature(bits, "SOME DIGEST METHOD")
	require.False(t, dv.IsValidDataSignature(reader, sig2.AsSignatureHolder()))

	_, _ = reader.Seek(0, io.SeekStart)
	sig3 := cryptkit.NewSignature(bits, SHA3512Digest.SignedBy("SOME SIGN METHOD"))
	require.False(t, dv.IsValidDataSignature(reader, sig3.AsSignatureHolder()))

	_, _ = reader.Seek(0, io.SeekStart)
	sig4 := cryptkit.NewSignature(bits, cryptkit.DigestMethod("SOME DIGEST METHOD").SignedBy(SECP256r1Sign))
	require.False(t, dv.IsValidDataSignature(reader, sig4.AsSignatureHolder()))
}
