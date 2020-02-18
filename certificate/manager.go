// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
		BootstrapNodes: make([]BootstrapNode, len(cert.BootstrapNodes)),
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
