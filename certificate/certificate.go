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
	"io"
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
	NodeRef     string `json:"node_ref"`

	// preprocessed fields
	nodePublicKey crypto.PublicKey
}

// GetNodeRef returns reference of bootstrap node
func (bn *BootstrapNode) GetNodeRef() *core.RecordRef {
	ref := core.NewRefFromBase58(bn.NodeRef)
	return &ref
}

// GetPublicKey returns public key reference of bootstrap node
func (bn *BootstrapNode) GetPublicKey() crypto.PublicKey {
	return bn.nodePublicKey
}

// GetHost returns host of bootstrap node
func (bn *BootstrapNode) GetHost() string {
	return bn.Host
}

// NodeSign returns signed information about some node
func (bn *BootstrapNode) GetNodeSign() []byte {
	return bn.NodeSign
}

var scheme = platformpolicy.NewPlatformCryptographyScheme()

// Certificate holds info about certificate
type Certificate struct {
	AuthorizationCertificate
	MajorityRule int `json:"majority_rule"`
	MinRoles     struct {
		Virtual       uint `json:"virtual"`
		HeavyMaterial uint `json:"heavy_material"`
		LightMaterial uint `json:"light_material"`
	} `json:"min_roles"`
	PulsarPublicKeys    []string        `json:"pulsar_public_keys"`
	RootDomainReference string          `json:"root_domain_ref"`
	BootstrapNodes      []BootstrapNode `json:"bootstrap_nodes"`

	// preprocessed fields
	pulsarPublicKey []crypto.PublicKey
}

func newCertificate(publicKey crypto.PublicKey, keyProcessor core.KeyProcessor, data []byte) (*Certificate, error) {
	cert := Certificate{}
	err := json.Unmarshal(data, &cert)
	if err != nil {
		return nil, errors.Wrap(err, "[ newCertificate ] failed to parse certificate json")
	}

	pub, err := keyProcessor.ExportPublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ newCertificate ] failed to retrieve public key from node private key")
	}

	if cert.PublicKey != string(pub) {
		return nil, errors.New("[ newCertificate ] Different public keys")
	}

	err = cert.fillExtraFields(keyProcessor)
	if err != nil {
		return nil, errors.Wrap(err, "[ newCertificate ] Incorrect fields")
	}

	cert.DiscoverySigns = make(map[*core.RecordRef][]byte)
	for _, node := range cert.BootstrapNodes {
		cert.DiscoverySigns[node.GetNodeRef()] = node.NodeSign
	}

	return &cert, nil
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
	signer := scheme.Signer(key)
	sign, err := signer.Sign(cert.serializeNetworkPart())
	if err != nil {
		return nil, errors.Wrap(err, "[ SignNetworkPart ] Can't Sign")
	}
	return sign.Bytes(), nil
}

func (cert *Certificate) fillExtraFields(keyProcessor core.KeyProcessor) error {
	importedNodePubKey, err := keyProcessor.ImportPublicKey([]byte(cert.PublicKey))
	if err != nil {
		return errors.Wrapf(err, "[ fillExtraFields ] Bad PublicKey: %s", cert.PublicKey)
	}
	cert.nodePublicKey = importedNodePubKey

	for _, pulsarKey := range cert.PulsarPublicKeys {
		importedPulsarPubKey, err := keyProcessor.ImportPublicKey([]byte(pulsarKey))
		if err != nil {
			return errors.Wrapf(err, "[ fillExtraFields ] Bad pulsarKey: %s", pulsarKey)
		}
		cert.pulsarPublicKey = append(cert.pulsarPublicKey, importedPulsarPubKey)
	}

	for i := 0; i < len(cert.BootstrapNodes); i++ {
		currentNode := &cert.BootstrapNodes[i]
		importedBNodePubKey, err := keyProcessor.ImportPublicKey([]byte(currentNode.PublicKey))
		if err != nil {
			return errors.Wrapf(err, "[ fillExtraFields ] Bad Bootstrap PublicKey: %s", currentNode.PublicKey)
		}
		currentNode.nodePublicKey = importedBNodePubKey
	}

	return nil
}

// GetRootDomainReference returns RootDomain reference
func (cert *Certificate) GetRootDomainReference() *core.RecordRef {
	ref := core.NewRefFromBase58(cert.RootDomainReference)
	return &ref
}

// GetDiscoveryNodes return bootstrap nodes array
func (cert *Certificate) GetDiscoveryNodes() []core.DiscoveryNode {
	result := make([]core.DiscoveryNode, 0)
	for i := 0; i < len(cert.BootstrapNodes); i++ {
		// we get node by pointer, so ranged for loop does not suite
		result = append(result, &cert.BootstrapNodes[i])
	}
	return result
}

// Dump returns all info about certificate in json format
func (cert *Certificate) Dump() (string, error) {
	result, err := json.MarshalIndent(cert, "", "    ")
	if err != nil {
		return "", errors.Wrap(err, "[ Certificate::Dump ]")
	}

	return string(result), nil
}

// ReadCertificate constructor creates new Certificate component
func ReadCertificate(publicKey crypto.PublicKey, keyProcessor core.KeyProcessor, certPath string) (*Certificate, error) {
	data, err := ioutil.ReadFile(filepath.Clean(certPath))
	if err != nil {
		return nil, errors.Wrapf(err, "[ ReadCertificate ] failed to read certificate from: %s", certPath)
	}
	cert, err := newCertificate(publicKey, keyProcessor, data)
	if err != nil {
		return nil, errors.Wrap(err, "[ ReadCertificate ]")
	}
	return cert, nil
}

// ReadCertificateFromReader constructor creates new Certificate component
func ReadCertificateFromReader(publicKey crypto.PublicKey, keyProcessor core.KeyProcessor, reader io.Reader) (*Certificate, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, errors.Wrapf(err, "[ ReadCertificateFromReader ] failed to read certificate data")
	}
	cert, err := newCertificate(publicKey, keyProcessor, data)
	if err != nil {
		return nil, errors.Wrap(err, "[ ReadCertificateFromReader ]")
	}
	return cert, nil
}

// NewCertificatesWithKeys generate certificate from given keys
// DEPRECATED, this method generates invalid certificate
func NewCertificatesWithKeys(publicKey crypto.PublicKey, keyProcessor core.KeyProcessor) (*Certificate, error) {
	cert := Certificate{}

	cert.Reference = testutils.RandomRef().String()

	keyBytes, err := keyProcessor.ExportPublicKey(publicKey)
	if err != nil {
		return nil, errors.Wrap(err, "[ ReadCertificate ] failed to retrieve public key from node private key")
	}

	cert.PublicKey = string(keyBytes)
	cert.nodePublicKey = publicKey
	return &cert, nil
}
