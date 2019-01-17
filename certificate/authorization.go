/*
 *    Copyright 2018 Insolar
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

package certificate

import (
	"crypto"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// AuthorizationCertificate holds info about node from it certificate
type AuthorizationCertificate struct {
	PublicKey      string                     `json:"public_key"`
	Reference      string                     `json:"reference"`
	Role           string                     `json:"role"`
	DiscoverySigns map[*core.RecordRef][]byte `json:"-"`

	nodePublicKey crypto.PublicKey
}

// GetPublicKey returns public key reference from node certificate
func (authCert *AuthorizationCertificate) GetPublicKey() crypto.PublicKey {
	return authCert.nodePublicKey
}

// GetNodeRef returns reference from node certificate
func (authCert *AuthorizationCertificate) GetNodeRef() *core.RecordRef {
	ref, err := core.NewRefFromBase58(authCert.Reference)
	if err != nil {
		log.Errorf("Invalid node reference in auth cert: %s\n", authCert.Reference)
		return nil
	}
	return ref
}

// GetRole returns role from node certificate
func (authCert *AuthorizationCertificate) GetRole() core.StaticRole {
	return core.GetStaticRoleFromString(authCert.Role)
}

// GetDiscoverySigns return map of discovery nodes signs
func (authCert *AuthorizationCertificate) GetDiscoverySigns() map[*core.RecordRef][]byte {
	return authCert.DiscoverySigns
}

// SerializeNodePart returns some node info decoded in bytes
func (authCert *AuthorizationCertificate) SerializeNodePart() []byte {
	return []byte(authCert.PublicKey + authCert.Reference + authCert.Role)
}

// SignNodePart signs node part in certificate
func (authCert *AuthorizationCertificate) SignNodePart(key crypto.PrivateKey) ([]byte, error) {
	signer := scheme.Signer(key)
	sign, err := signer.Sign(authCert.SerializeNodePart())
	if err != nil {
		return nil, errors.Wrap(err, "[ SignNodePart ] Can't Sign")
	}
	return sign.Bytes(), nil
}

// Deserialize deserializes data to AuthorizationCertificate interface
func Deserialize(data []byte, keyProc core.KeyProcessor) (core.AuthorizationCertificate, error) {
	cert := &AuthorizationCertificate{}
	err := core.Deserialize(data, cert)

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
func Serialize(authCert core.AuthorizationCertificate) ([]byte, error) {
	data, err := core.Serialize(authCert)
	if err != nil {
		return nil, errors.Wrap(err, "[ AuthorizationCertificate::Serialize ]")
	}
	return data, nil
}
