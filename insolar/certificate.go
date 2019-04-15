//
// Copyright 2019 Insolar Technologies GmbH
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
//

package insolar

import (
	"crypto"
)

type NodeMeta interface {
	GetNodeRef() *Reference
	GetPublicKey() crypto.PublicKey
}

//go:generate minimock -i github.com/insolar/insolar/insolar.Certificate -o ../testutils -s _mock.go

// Certificate interface provides methods to manage keys
type Certificate interface {
	AuthorizationCertificate

	GetRootDomainReference() *Reference
	GetDiscoveryNodes() []DiscoveryNode
}

//go:generate minimock -i github.com/insolar/insolar/insolar.DiscoveryNode -o ../testutils -s _mock.go
type DiscoveryNode interface {
	NodeMeta

	GetHost() string
}

// AuthorizationCertificate interface provides methods to manage info about node from it certificate
type AuthorizationCertificate interface {
	NodeMeta

	GetRole() StaticRole
	SerializeNodePart() []byte
	GetDiscoverySigns() map[Reference][]byte
}

//go:generate minimock -i github.com/insolar/insolar/insolar.CertificateManager -o ../testutils -s _mock.go

// CertificateManager interface provides methods to manage nodes certificate
type CertificateManager interface {
	GetCertificate() Certificate
	VerifyAuthorizationCertificate(authCert AuthorizationCertificate) (bool, error)
	NewUnsignedCertificate(pKey string, role string, nodeRef string) (Certificate, error)
}
