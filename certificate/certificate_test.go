// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package certificate

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/insolar/insolar/cryptography"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/require"
)

const TestCert = "testdata/cert.json"
const TestBadCert = "testdata/bad_cert.json"
const TestInvalidFileCert = "testdata/bad_cert11111.json"

const TestKeys = "testdata/keys.json"
const TestDifferentKeys = "testdata/different_keys.json"

func TestNewCertificate_NoCert(t *testing.T) {
	_, err := ReadCertificate(nil, nil, TestInvalidFileCert)
	require.EqualError(t, err, "[ ReadCertificate ] failed to read certificate from: "+
		"testdata/bad_cert11111.json: open testdata/bad_cert11111.json: no such file or directory")
}

func TestNewCertificate_BadCert(t *testing.T) {
	_, err := ReadCertificate(nil, nil, TestBadCert)
	require.Contains(t, err.Error(), "failed to parse certificate json")
}

func checkKeys(cert *Certificate, cs insolar.CryptographyService, t *testing.T) {
	kp := platformpolicy.NewKeyProcessor()

	pubKey, err := cs.GetPublicKey()
	require.NoError(t, err)

	pubKeyString, err := kp.ExportPublicKeyPEM(pubKey)
	require.NoError(t, err)

	require.Equal(t, string(pubKeyString), cert.PublicKey)
}

func TestReadCertificate(t *testing.T) {
	cs, _ := cryptography.NewStorageBoundCryptographyService(TestKeys)
	kp := platformpolicy.NewKeyProcessor()
	pk, _ := cs.GetPublicKey()

	cert, err := ReadCertificate(pk, kp, TestCert)
	require.NoError(t, err)
	require.NotEmpty(t, cert.PublicKey)
	require.Equal(t, "virtual", cert.Role)
	require.Equal(t, "insolar:1AMuaHS6ENi_58HSwvLmeNg-26BHLcCuz36c0SzEJT1I",
		cert.Reference)
	require.Equal(t, 7, cert.MajorityRule)
	require.Equal(t, "insolar:1AAEAAflht_i9JBIq2_moi5L61cLEVxGNpIBJlzA54zc", cert.RootDomainReference)

	testPubKey1 := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFhw9v2vl9OkMedBoKz8GTndZyx5S\n/KHFc3OKOoEhUPZwuNo1q3bXTaeJ1WBcs4MjGBBGuC5w1i3WcNfJHzyyLw==\n-----END PUBLIC KEY-----\n"
	key1, err := kp.ImportPublicKeyPEM([]byte(testPubKey1))
	require.NoError(t, err)
	testPubKey2 := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEtlCnGJbptDc+c/pr9nAJE3SlWSCP\n3SIQ9Q9iOFKzZFf71c6fCrLyquAl+lCD/S3ch1v/y42d1peGAWiYujmOuw==\n-----END PUBLIC KEY-----\n"
	key2, err := kp.ImportPublicKeyPEM([]byte(testPubKey2))
	require.NoError(t, err)
	testPubKey3 := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEMs0nrI/to+AfRs+AuXr+ri/fHNrQ\nY6A83jNYlQfIUHRFxmv9Oowi6aDGkmsSzRDGVKrqabjgBuVYYTznqm/s9g==\n-----END PUBLIC KEY-----\n"
	key3, err := kp.ImportPublicKeyPEM([]byte(testPubKey3))
	require.NoError(t, err)

	bootstrapNodes := []BootstrapNode{
		BootstrapNode{
			PublicKey:     testPubKey1,
			Host:          "localhost:22001",
			nodePublicKey: key1,
			// NetworkSign:   base58.Decode("ICDTj1ev/wEBewXBWeNhsmByZ9enmfOAM+ltB7pA6s02sUdGr8n1w2STkp5YsADDj2SudxC18enFEDx8Nh3J0zeR"),
			NodeRef: "insolar:1AAEAAbP0lFkTHKV3QBOUfu3ORaa7fIcSEbgwUMEB09M",
		},
		BootstrapNode{
			PublicKey:     testPubKey2,
			Host:          "127.0.0.1:23832",
			nodePublicKey: key2,
			// NetworkSign:   base58.Decode("ICA0L77lNk2LlAIlWM727b621vPXMbKzFHHDFVtrXKRBilnMiuGq44QpidwB8Ps3YIhZ2XzElajYo2eQ88M9hKDh"),
			NodeRef: "insolar:1AAEAASl0Gzn5SYZH9oTf1aMjakbF3SLcyMuvTHV1I4A",
		},
		BootstrapNode{
			PublicKey:     testPubKey3,
			Host:          "127.0.0.1:33833",
			nodePublicKey: key3,
			// NetworkSign:   base58.Decode("ICAgWvCWLtEJSmITTfa0bKiCL1NwsKnAl8Yt6WNYbXGJVMHCTmbSmpTYqDhvRnUNEwq+J1q3E+nKiO3ZbxZBGjLB"),
			NodeRef: "insolar:1AAEAAaYz0ZzIvcNHXrAsY3_iucph-gU_q461IaIQelw",
		},
	}

	require.Equal(t, bootstrapNodes, cert.BootstrapNodes)

	checkKeys(cert, cs, t)
}

func TestReadCertificate_BadBootstrapPublicKey(t *testing.T) {
	cs, _ := cryptography.NewStorageBoundCryptographyService(TestKeys)
	kp := platformpolicy.NewKeyProcessor()
	pk, _ := cs.GetPublicKey()

	_, err := ReadCertificate(pk, kp, "testdata/cert_bad_bootstrap_key.json")
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "Incorrect fields: [ fillExtraFields ] Bad Bootstrap PublicKey")
}

func TestReadPrivateKey_BadJson(t *testing.T) {
	keyProcessor := platformpolicy.NewKeyProcessor()
	_, err := ReadCertificate(nil, keyProcessor, TestBadCert)
	require.Contains(t, err.Error(), "failed to parse certificate json")
}

func TestReadPrivateKey_BadKeyPair(t *testing.T) {
	cs, _ := cryptography.NewStorageBoundCryptographyService(TestDifferentKeys)
	kp := platformpolicy.NewKeyProcessor()
	pk, _ := cs.GetPublicKey()

	_, err := ReadCertificate(pk, kp, TestCert)
	require.Contains(t, err.Error(), "Different public keys")
}

func TestReadCertificateFromReader(t *testing.T) {
	kp := platformpolicy.NewKeyProcessor()
	privateKey, _ := kp.GeneratePrivateKey()
	nodePublicKey := kp.ExtractPublicKey(privateKey)
	publicKey, _ := kp.ExportPublicKeyPEM(nodePublicKey)

	type сertInfo map[string]interface{}
	info := сertInfo{
		"majority_rule": 7,
		"min_roles": map[string]interface{}{
			"virtual":        1,
			"heavy_material": 2,
			"light_material": 3,
		},
		"public_key":      string(publicKey[:]),
		"reference":       "insolar:1AAABAAGmM9GcyL3DR16wLGN_4rnKYfoFP6uOtSGiEHo",
		"role":            "virtual",
		"root_domain_ref": "insolar:1AAABAAH5Ybf4vSQSKtv5qIuS-tXCxFcRjaSASZcwOeM",
		"bootstrap_nodes": []map[string]interface{}{
			{
				"public_key": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFhw9v2vl9OkMedBoKz8GTndZyx5S\n/KHFc3OKOoEhUPZwuNo1q3bXTaeJ1WBcs4MjGBBGuC5w1i3WcNfJHzyyLw==\n-----END PUBLIC KEY-----\n",
				"host":       "localhost:22001",
				// "node_sign":  "ICBcAGEnW9gxvLYUgjTFwsUXh5sISXLKka3gZSVXQumpqQxKMTO4+Q7bF9iQc20+rlBdgMcnh2JYtFkzSOgEdBfH",
				"node_ref": "insolar:1AAABAAGz9JRZExyld0ATlH7tzkWmu3yHEhG4MFDBAdM",
			},
			{
				"public_key": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEtlCnGJbptDc+c/pr9nAJE3SlWSCP\n3SIQ9Q9iOFKzZFf71c6fCrLyquAl+lCD/S3ch1v/y42d1peGAWiYujmOuw==\n-----END PUBLIC KEY-----\n",
				"host":       "localhost:22002",
				// "node_sign":  "ICAMco6BKxQgpIVVrvijNF+IewMavWKS8YnnZRwsqOxpgjtUe+BAR58efoGnWNVK4IBetuJ0tw0ZJqKHOp2NqdAy",
				"node_ref": "insolar:1AAABAAEpdBs5-UmGR_aE39WjI2pGxd0i3MjLr0x1dSM",
			},
			{
				"public_key": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEMs0nrI/to+AfRs+AuXr+ri/fHNrQ\nY6A83jNYlQfIUHRFxmv9Oowi6aDGkmsSzRDGVKrqabjgBuVYYTznqm/s9g==\n-----END PUBLIC KEY-----\n",
				"host":       "localhost:22003",
				// "node_sign":  "ICDX3f1kb1sdJDQyXiFOWR9X4+jDoCFxNN5WfoEM99ZdJMnDMaFkVk62NTkik+nadlOcNEAe/gBNG1ezuxpZVv8q",
				"node_ref": "insolar:1AAABAAGmM9GcyL3DR16wLGN_4rnKYfoFP6uOtSGiEHo",
			},
		},
	}
	certJson, err := json.Marshal(info)
	require.NoError(t, err)

	r := bytes.NewReader(certJson)
	cert, err := ReadCertificateFromReader(nodePublicKey, kp, r)

	require.NoError(t, err)
	require.NotEmpty(t, cert.PublicKey)
	require.Equal(t, "virtual", cert.Role)
	require.Equal(t, "insolar:1AAABAAGmM9GcyL3DR16wLGN_4rnKYfoFP6uOtSGiEHo",
		cert.Reference)
	require.Equal(t, 7, cert.MajorityRule)
	require.Equal(t, uint(1), cert.MinRoles.Virtual)
	require.Equal(t, uint(2), cert.MinRoles.HeavyMaterial)
	require.Equal(t, uint(3), cert.MinRoles.LightMaterial)
	require.Equal(t, "insolar:1AAABAAH5Ybf4vSQSKtv5qIuS-tXCxFcRjaSASZcwOeM", cert.RootDomainReference)
	require.Equal(t, nodePublicKey, cert.nodePublicKey)

	testPubKey1 := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEFhw9v2vl9OkMedBoKz8GTndZyx5S\n/KHFc3OKOoEhUPZwuNo1q3bXTaeJ1WBcs4MjGBBGuC5w1i3WcNfJHzyyLw==\n-----END PUBLIC KEY-----\n"
	key1, err := kp.ImportPublicKeyPEM([]byte(testPubKey1))
	require.NoError(t, err)
	testPubKey2 := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEtlCnGJbptDc+c/pr9nAJE3SlWSCP\n3SIQ9Q9iOFKzZFf71c6fCrLyquAl+lCD/S3ch1v/y42d1peGAWiYujmOuw==\n-----END PUBLIC KEY-----\n"
	key2, err := kp.ImportPublicKeyPEM([]byte(testPubKey2))
	require.NoError(t, err)
	testPubKey3 := "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEMs0nrI/to+AfRs+AuXr+ri/fHNrQ\nY6A83jNYlQfIUHRFxmv9Oowi6aDGkmsSzRDGVKrqabjgBuVYYTznqm/s9g==\n-----END PUBLIC KEY-----\n"
	key3, err := kp.ImportPublicKeyPEM([]byte(testPubKey3))
	require.NoError(t, err)

	bootstrapNodes := []BootstrapNode{
		BootstrapNode{
			PublicKey:     testPubKey1,
			Host:          "localhost:22001",
			nodePublicKey: key1,
			NodeRef:       "insolar:1AAABAAGz9JRZExyld0ATlH7tzkWmu3yHEhG4MFDBAdM",
		},
		BootstrapNode{
			PublicKey:     testPubKey2,
			Host:          "localhost:22002",
			nodePublicKey: key2,
			NodeRef:       "insolar:1AAABAAEpdBs5-UmGR_aE39WjI2pGxd0i3MjLr0x1dSM",
		},
		BootstrapNode{
			PublicKey:     testPubKey3,
			Host:          "localhost:22003",
			nodePublicKey: key3,
			NodeRef:       "insolar:1AAABAAGmM9GcyL3DR16wLGN_4rnKYfoFP6uOtSGiEHo",
		},
	}

	require.Equal(t, bootstrapNodes, cert.BootstrapNodes)
}

func TestSerializeDeserialize(t *testing.T) {
	cert := &AuthorizationCertificate{
		PublicKey: "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEG1XfrtnhPKqO2zSywoi2G8nQG6y8\nyIU7a3NeGzc06ygEaXzWK+DdyeBpeRhop4eUKJdfKFm1mHvZdvEiQwzx4A==\n-----END PUBLIC KEY-----\n",
		Reference: "test_reference",
		Role:      "test_role",
	}

	keyProc := platformpolicy.NewKeyProcessor()
	key, err := keyProc.ImportPublicKeyPEM([]byte(cert.PublicKey))
	require.NoError(t, err)

	cert.nodePublicKey = key

	result, err := Serialize(cert)
	require.NoError(t, err)

	deserializedCert, err := Deserialize(result, keyProc)
	require.NoError(t, err)
	require.Equal(t, cert, deserializedCert)
}
