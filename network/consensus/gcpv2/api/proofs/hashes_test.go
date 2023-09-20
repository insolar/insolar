package proofs

import (
	"testing"
)

func TestNewNodeStateHashEvidence(t *testing.T) {
	t.Skipped()
	// sd := cryptkit.NewSignedDigest(cryptkit.Digest{}, cryptkit.Signature{})
	// sh := NewNodeStateHashEvidence(sd)
	// require.Equal(t, sd, sh.(*nodeStateHashEvidence).SignedDigest)
}

func TestGetNodeStateHash(t *testing.T) {
	t.Skipped()
	// fr := longbits.NewFoldableReaderMock(t)
	// sd := cryptkit.NewSignedDigest(cryptkit.NewDigest(fr, cryptkit.DigestMethod("testDigest")), cryptkit.NewSignature(fr, cryptkit.SignatureMethod("testSignature")))
	// sh := NewNodeStateHashEvidence(sd)
	// require.Equal(t, sh.GetNodeStateHash().GetDigestMethod(), sd.GetDigest().AsDigestHolder().GetDigestMethod())
}

func TestGetGlobulaNodeStateSignature(t *testing.T) {
	t.Skipped()
	// fr := longbits.NewFoldableReaderMock(t)
	// sd := cryptkit.NewSignedDigest(cryptkit.NewDigest(fr, cryptkit.DigestMethod("testDigest")), cryptkit.NewSignature(fr, cryptkit.SignatureMethod("testSignature")))
	// sh := NewNodeStateHashEvidence(sd)
	// require.Equal(t, sh.GetGlobulaNodeStateSignature().GetSignatureMethod(), sd.GetSignature().AsSignatureHolder().GetSignatureMethod())
}
