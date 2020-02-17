// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package certificate

import (
	"crypto"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// AuthorizationCertificate holds info about node from it certificate
type AuthorizationCertificate struct {
	PublicKey      string                       `json:"public_key"`
	Reference      string                       `json:"reference"`
	Role           string                       `json:"role"`
	DiscoverySigns map[insolar.Reference][]byte `json:"-" codec:"discoverysigns"`

	nodePublicKey crypto.PublicKey
}

// GetPublicKey returns public key reference from node certificate
func (authCert *AuthorizationCertificate) GetPublicKey() crypto.PublicKey {
	return authCert.nodePublicKey
}

// GetNodeRef returns reference from node certificate
func (authCert *AuthorizationCertificate) GetNodeRef() *insolar.Reference {
	ref, err := insolar.NewReferenceFromString(authCert.Reference)
	if err != nil {
		log.Errorf("Invalid node reference in auth cert: %s\n", authCert.Reference)
		return nil
	}
	return ref
}

// GetRole returns role from node certificate
func (authCert *AuthorizationCertificate) GetRole() insolar.StaticRole {
	return insolar.GetStaticRoleFromString(authCert.Role)
}

// GetDiscoverySigns return map of discovery nodes signs
func (authCert *AuthorizationCertificate) GetDiscoverySigns() map[insolar.Reference][]byte {
	return authCert.DiscoverySigns
}

// SerializeNodePart returns some node info decoded in bytes
func (authCert *AuthorizationCertificate) SerializeNodePart() []byte {
	return []byte(authCert.PublicKey + authCert.Reference + authCert.Role)
}

// SignNodePart signs node part in certificate
func (authCert *AuthorizationCertificate) SignNodePart(key crypto.PrivateKey) ([]byte, error) {
	signer := scheme.DataSigner(key, scheme.IntegrityHasher())
	sign, err := signer.Sign(authCert.SerializeNodePart())
	if err != nil {
		return nil, errors.Wrap(err, "[ SignNodePart ] Can't Sign")
	}
	return sign.Bytes(), nil
}

// Deserialize deserializes data to AuthorizationCertificate interface
func Deserialize(data []byte, keyProc insolar.KeyProcessor) (insolar.AuthorizationCertificate, error) {
	cert := &AuthorizationCertificate{}
	err := insolar.Deserialize(data, cert)

	if err != nil {
		return nil, errors.Wrap(err, "[ AuthorizatonCertificate::Deserialize ] failed to deserialize a data")
	}

	key, err := keyProc.ImportPublicKeyPEM([]byte(cert.PublicKey))

	if err != nil {
		return nil, errors.Wrap(err, "[ AuthorizationCertificate::Deserialize ] failed to import a public key")
	}

	cert.nodePublicKey = key

	return cert, nil
}

// Serialize serializes AuthorizationCertificate interface
func Serialize(authCert insolar.AuthorizationCertificate) ([]byte, error) {
	data, err := insolar.Serialize(authCert)
	if err != nil {
		return nil, errors.Wrap(err, "[ AuthorizationCertificate::Serialize ]")
	}
	return data, nil
}
