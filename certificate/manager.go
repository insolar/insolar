// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package certificate

import (
	"crypto"

	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

// CertificateManager is a component for working with current node certificate
type CertificateManager struct { // nolint: golint
	certificate insolar.Certificate
}

// NewCertificateManager returns new CertificateManager instance
func NewCertificateManager(cert insolar.Certificate) *CertificateManager {
	return &CertificateManager{certificate: cert}
}

// GetCertificate returns current node certificate
func (m *CertificateManager) GetCertificate() insolar.Certificate {
	return m.certificate
}

// VerifyAuthorizationCertificate verifies certificate from some node
func VerifyAuthorizationCertificate(cs insolar.CryptographyService, discoveryNodes []insolar.DiscoveryNode, authCert insolar.AuthorizationCertificate) (bool, error) {
	if len(discoveryNodes) != len(authCert.GetDiscoverySigns()) {
		return false, nil
	}
	data := authCert.SerializeNodePart()
	for _, node := range discoveryNodes {
		sign := authCert.GetDiscoverySigns()[*node.GetNodeRef()]
		ok := cs.Verify(node.GetPublicKey(), insolar.SignatureFromBytes(sign), data)
		if !ok {
			return false, nil
		}
	}
	return true, nil
}

// NewUnsignedCertificate creates new unsigned certificate by copying
func NewUnsignedCertificate(baseCert insolar.Certificate, pKey string, role string, ref string) (insolar.Certificate, error) {
	cert := baseCert.(*Certificate)
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
		newCert.BootstrapNodes[i].NodeRef = node.NodeRef
		newCert.BootstrapNodes[i].PublicKey = node.PublicKey
		newCert.BootstrapNodes[i].NetworkSign = node.NetworkSign
		newCert.BootstrapNodes[i].NodeRole = node.NodeRole
	}
	return &newCert, nil
}

// NewManagerReadCertificate constructor creates new CertificateManager component
func NewManagerReadCertificate(publicKey crypto.PublicKey, keyProcessor insolar.KeyProcessor, certPath string) (*CertificateManager, error) {
	cert, err := ReadCertificate(publicKey, keyProcessor, certPath)
	if err != nil {
		return nil, errors.Wrap(err, "[ NewManagerReadCertificate ] failed to read certificate:")
	}
	certManager := NewCertificateManager(cert)
	return certManager, nil
}
