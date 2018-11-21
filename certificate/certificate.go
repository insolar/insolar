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
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
)

// BootstrapNode holds info about bootstrap nodes
type BootstrapNode struct {
	PublicKey   string `json:"public_key"`
	Host        string `json:"host"`
	NetworkSign []byte `json:"network_sign"`
	NodeSign    []byte `json:"node_sign"`
}

// Certificate holds info about certificate
type Certificate struct {
	MajorityRule int `json:"majority_rule"`
	MinRoles     struct {
		Virtual       uint `json:"virtual"`
		HeavyMaterial uint `json:"heavy_material"`
		LightMaterial uint `json:"light_material"`
	} `json:"min_roles"`
	PublicKey           string          `json:"public_key"`
	Reference           string          `json:"reference"`
	PulsarPublicKeys    []string        `json:"pulsar_public_keys"`
	Role                string          `json:"role"`
	BootstrapNodes      []BootstrapNode `json:"bootstrap_nodes"`
	RootDomainReference string          `json:"root_domain_ref"`
}

func (cert *Certificate) serializeNetworkPart() []byte {
	out := strconv.Itoa(cert.MajorityRule) + strconv.Itoa(int(cert.MinRoles.Virtual)) +
		strconv.Itoa(int(cert.MinRoles.HeavyMaterial)) + strconv.Itoa(int(cert.MinRoles.LightMaterial)) +
		cert.RootDomainReference

	sort.Strings(cert.PulsarPublicKeys)
	out += strings.Join(cert.PulsarPublicKeys, "")
	sort.Slice(cert.BootstrapNodes, func(i, j int) bool {
		return strings.Compare(cert.BootstrapNodes[i].PublicKey, cert.BootstrapNodes[j].PublicKey) == -1
	})

	for _, node := range cert.BootstrapNodes {
		out += node.PublicKey + node.Host
	}

	return []byte(out)
}

// SignNetworkPart signs network part in certificate
func (cert *Certificate) SignNetworkPart(key crypto.PrivateKey) ([]byte, error) {
	signer := platformpolicy.NewPlatformCryptographyScheme().Signer(key)
	sign, err := signer.Sign(cert.serializeNetworkPart())
	if err != nil {
		return nil, err
	}
	return sign.Bytes(), nil
}

func (cert *Certificate) serializeNodePart() []byte {
	return []byte(cert.PublicKey + cert.Reference + cert.Role)
}

// SignNodePart signs node part in certificate
func (cert *Certificate) SignNodePart(key crypto.PrivateKey) ([]byte, error) {
	signer := platformpolicy.NewPlatformCryptographyScheme().Signer(key)
	sign, err := signer.Sign(cert.serializeNodePart())
	if err != nil {
		return nil, err
	}
	return sign.Bytes(), nil
}

// GetBootstrapNodes return bootstrap nodes array
func (cert *Certificate) GetBootstrapNodes() []BootstrapNode {
	return cert.BootstrapNodes
}

// ReadCertificate constructor creates new Certificate component
func ReadCertificate(publicKey crypto.PublicKey, keyProcessor core.KeyProcessor, certPath string) (*Certificate, error) {
	data, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return nil, errors.New("[ ReadCertificate ] failed to read certificate from: " + certPath)
	}
	cert := Certificate{}
	err = json.Unmarshal(data, &cert)
	if err != nil {
		return nil, errors.Wrap(err, "[ ReadCertificate ] failed to parse certificate json")
	}

	pub, err := keyProcessor.ExportPublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ReadCertificate ] failed to retrieve public key from node private key")
	}

	if cert.PublicKey != string(pub) {
		return nil, errors.New("[ ReadCertificate ] Different public keys. Cert path: " + certPath + ".")
	}

	return &cert, nil
}

func (cert *Certificate) reset() {
	cert.PublicKey = ""
	cert.BootstrapNodes = []BootstrapNode{}
	cert.Reference = ""
	cert.MajorityRule = 0
}

func (cert *Certificate) GetRole() core.NodeRole {
	return core.GetRoleFromString(cert.Role)
}

// GetRootDomainReference returns RootDomain reference as string
func (cert *Certificate) GetRootDomainReference() string {
	return cert.RootDomainReference
}

// SetRootDomainReference sets RootDomain reference for certificate
func (cert *Certificate) SetRootDomainReference(ref *core.RecordRef) {
	cert.RootDomainReference = ref.String()
}

// NewCertificatesWithKeys generate certificate from given keys
func NewCertificatesWithKeys(publicKey crypto.PublicKey, keyProcessor core.KeyProcessor) (*Certificate, error) {
	cert := Certificate{}
	cert.reset()

	cert.Reference = testutils.RandomRef().String()

	keyBytes, err := keyProcessor.ExportPublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ReadCertificate ] failed to retrieve public key from node private key")
	}

	cert.PublicKey = string(keyBytes)
	return &cert, nil
}

func (cert *Certificate) Dump() (string, error) {
	result, err := json.MarshalIndent(cert, "", "    ")
	if err != nil {
		return "", errors.Wrap(err, "[ Certificate::Dump ]")
	}

	return string(result), nil
}
