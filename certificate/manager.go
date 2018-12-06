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
	"io"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

// CertificateManager is a component for working with current node certificate
type CertificateManager struct {
	CS          core.CryptographyService `inject:""`
	certificate core.Certificate
}

// NewCertificateManager returns new CertificateManager instance
func NewCertificateManager(cert core.Certificate) *CertificateManager {
	return &CertificateManager{certificate: cert}
}

// GetCertificate returns current node certificate
func (m *CertificateManager) GetCertificate() core.Certificate {
	return m.certificate
}

// VerifyAuthorizationCertificate verifies certificate from some node
func (m *CertificateManager) VerifyAuthorizationCertificate(authCert core.AuthorizationCertificate) (bool, error) {
	discoveryNodes := m.certificate.GetDiscoveryNodes()
	if len(discoveryNodes) != len(authCert.GetDiscoverySigns()) {
		return false, nil
	}
	data := authCert.SerializeNodePart()
	for _, node := range discoveryNodes {
		sign := authCert.GetDiscoverySigns()[node.GetNodeRef()]
		ok := m.CS.Verify(node.GetPublicKey(), core.SignatureFromBytes(sign), data)
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// NewUnsignedCertificate returns new certificate
func (m *CertificateManager) NewUnsignedCertificate(pKey string, ref string, role string) (core.Certificate, error) {
	cert := m.certificate.(*Certificate)
	newCert := Certificate{
		MajorityRule: cert.MajorityRule,
		MinRoles:     cert.MinRoles,
		AuthorizationCertificate: AuthorizationCertificate{
			PublicKey: pKey,
			Reference: ref,
			Role:      role,
		},
		PulsarPublicKeys:    cert.PulsarPublicKeys,
		RootDomainReference: cert.RootDomainReference,
		BootstrapNodes:      make([]BootstrapNode, len(cert.BootstrapNodes)),
	}
	for i, node := range cert.BootstrapNodes {
		newCert.BootstrapNodes[i].Host = node.Host
		newCert.BootstrapNodes[i].PublicKey = node.PublicKey
		newCert.BootstrapNodes[i].NetworkSign = node.NetworkSign
	}
	return &newCert, nil
}

// NewManagerReadCertificate constructor creates new CertificateManager component
func NewManagerReadCertificate(publicKey crypto.PublicKey, keyProcessor core.KeyProcessor, certPath string) (*CertificateManager, error) {
	cert, err := ReadCertificate(publicKey, keyProcessor, certPath)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewManagerReadCertificate ] failed to read certificate:")
	}
	certManager := NewCertificateManager(cert)
	return certManager, nil
}

// NewManagerReadCertificateFromReader constructor creates new CertificateManager component
func NewManagerReadCertificateFromReader(publicKey crypto.PublicKey, keyProcessor core.KeyProcessor, reader io.Reader) (*CertificateManager, error) {
	cert, err := ReadCertificateFromReader(publicKey, keyProcessor, reader)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewManagerReadCertificateFromReader ] failed to read certificate data:")
	}
	certManager := NewCertificateManager(cert)
	return certManager, nil
}

// NewManagerCertificateWithKeys generate manager with certificate from given keys
// DEPRECATED, this method generates invalid certificate
func NewManagerCertificateWithKeys(publicKey crypto.PublicKey, keyProcessor core.KeyProcessor) (*CertificateManager, error) {
	cert, err := NewCertificatesWithKeys(publicKey, keyProcessor)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewManagerCertificateWithKeys ] failed to create certificate:")
	}
	certManager := NewCertificateManager(cert)
	return certManager, nil
}
