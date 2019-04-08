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

package genesis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/platformpolicy/commoncrypto"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/nodedomain"
	"github.com/insolar/insolar/application/contract/noderecord"
	"github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/application/contract/wallet"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/pkg/errors"
)

const (
	nodeDomain        = "nodedomain"
	nodeRecord        = "noderecord"
	rootDomain        = "rootdomain"
	walletContract    = "wallet"
	memberContract    = "member"
	allowanceContract = "allowance"
	nodeAmount        = 32
)

var contractNames = []string{walletContract, memberContract, allowanceContract, rootDomain, nodeDomain, nodeRecord}

type messageBusLocker interface {
	Lock(ctx context.Context)
	Unlock(ctx context.Context)
}

type nodeInfo struct {
	privateKey platformpolicy.PrivateKey
	publicKey  string
	ref        *insolar.Reference
}

// Genesis is a component for precreation insolar contracts types and RootDomain instance
type Genesis struct {
	rootDomainRef   *insolar.Reference
	nodeDomainRef   *insolar.Reference
	rootMemberRef   *insolar.Reference
	prototypeRefs   map[string]*insolar.Reference
	isGenesis       bool
	config          *Config
	keyOut          string
	ArtifactManager artifacts.Client `inject:""`
	MBLock          messageBusLocker `inject:""`
}

// NewGenesis creates new Genesis
func NewGenesis(isGenesis bool, genesisConfigPath string, genesisKeyOut string) (*Genesis, error) {
	var err error
	genesis := &Genesis{}
	genesis.rootDomainRef = &insolar.Reference{}
	genesis.isGenesis = isGenesis
	if isGenesis {
		genesis.config, err = ParseGenesisConfig(genesisConfigPath)
		genesis.keyOut = genesisKeyOut
	}
	return genesis, err
}

func buildSmartContracts(ctx context.Context, cb *ContractsBuilder, rootDomainID *insolar.ID) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ buildSmartContracts ] building contracts:", contractNames)
	contracts, err := getContractsMap()
	if err != nil {
		return errors.Wrap(err, "[ buildSmartContracts ] couldn't build contracts")
	}

	inslog.Info("[ buildSmartContracts ] Start building contracts ...")
	err = cb.Build(ctx, contracts, rootDomainID)
	if err != nil {
		return errors.Wrap(err, "[ buildSmartContracts ] couldn't build contracts")
	}
	inslog.Info("[ buildSmartContracts ] Stop building contracts ...")

	return nil
}

func (g *Genesis) activateRootDomain(
	ctx context.Context, cb *ContractsBuilder,
	contractID *insolar.ID,
) (artifacts.ObjectDescriptor, error) {
	rd, err := rootdomain.NewRootDomain()
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	instanceData, err := serializeInstance(rd)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	contract := insolar.NewReference(*contractID, *contractID)
	desc, err := g.ArtifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		*g.ArtifactManager.GenesisRef(),
		*cb.Prototypes[rootDomain],
		false,
		instanceData,
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	_, err = g.ArtifactManager.RegisterResult(ctx, *g.ArtifactManager.GenesisRef(), *contract, nil)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	g.rootDomainRef = contract

	return desc, nil
}

func (g *Genesis) activateNodeDomain(
	ctx context.Context, domain *insolar.ID, cb *ContractsBuilder,
) (artifacts.ObjectDescriptor, error) {
	nd, err := nodedomain.NewNodeDomain()
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateNodeDomain ]")
	}

	instanceData, err := serializeInstance(nd)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateNodeDomain ]")
	}

	contractID, err := g.ArtifactManager.RegisterRequest(ctx, *g.rootDomainRef, &message.Parcel{Msg: &message.GenesisRequest{Name: "NodeDomain"}})

	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}
	contract := insolar.NewReference(*domain, *contractID)
	desc, err := g.ArtifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		*g.rootDomainRef,
		*cb.Prototypes[nodeDomain],
		false,
		instanceData,
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}
	_, err = g.ArtifactManager.RegisterResult(ctx, *g.rootDomainRef, *contract, nil)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}

	g.nodeDomainRef = contract

	return desc, nil
}

func (g *Genesis) activateRootMember(
	ctx context.Context, domain *insolar.ID, cb *ContractsBuilder, rootPubKey string,
) error {

	m, err := member.New("RootMember", rootPubKey)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ]")
	}

	instanceData, err := serializeInstance(m)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ]")
	}

	contractID, err := g.ArtifactManager.RegisterRequest(ctx, *g.rootDomainRef, &message.Parcel{Msg: &message.GenesisRequest{Name: "RootMember"}})

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	contract := insolar.NewReference(*domain, *contractID)
	_, err = g.ArtifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		*g.rootDomainRef,
		*cb.Prototypes[memberContract],
		false,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	_, err = g.ArtifactManager.RegisterResult(ctx, *g.rootDomainRef, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	g.rootMemberRef = contract
	return nil
}

// TODO: this is not required since we refer by request id.
func (g *Genesis) updateRootDomain(
	ctx context.Context, domainDesc artifacts.ObjectDescriptor,
) error {
	updateData, err := serializeInstance(&rootdomain.RootDomain{RootMember: *g.rootMemberRef, NodeDomainRef: *g.nodeDomainRef})
	if err != nil {
		return errors.Wrap(err, "[ updateRootDomain ]")
	}
	_, err = g.ArtifactManager.UpdateObject(
		ctx,
		insolar.Reference{},
		insolar.Reference{},
		domainDesc,
		updateData,
	)
	if err != nil {
		return errors.Wrap(err, "[ updateRootDomain ]")
	}

	return nil
}

func (g *Genesis) activateRootMemberWallet(
	ctx context.Context, domain *insolar.ID, cb *ContractsBuilder,
) error {

	w, err := wallet.New(g.config.RootBalance)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	instanceData, err := serializeInstance(w)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	contractID, err := g.ArtifactManager.RegisterRequest(ctx, *g.rootDomainRef, &message.Parcel{Msg: &message.GenesisRequest{Name: "RootWallet"}})

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}
	contract := insolar.NewReference(*domain, *contractID)
	_, err = g.ArtifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		*g.rootMemberRef,
		*cb.Prototypes[walletContract],
		true,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}
	_, err = g.ArtifactManager.RegisterResult(ctx, *g.rootDomainRef, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}

	return nil
}

func (g *Genesis) activateSmartContracts(
	ctx context.Context, cb *ContractsBuilder, rootPubKey string, rootDomainID *insolar.ID,
) ([]genesisNode, error) {

	rootDomainDesc, err := g.activateRootDomain(ctx, cb, rootDomainID)
	errMsg := "[ ActivateSmartContracts ]"
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	nodeDomainDesc, err := g.activateNodeDomain(ctx, rootDomainID, cb)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	err = g.activateRootMember(ctx, rootDomainID, cb, rootPubKey)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	// TODO: this is not required since we refer by request id.
	err = g.updateRootDomain(ctx, rootDomainDesc)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	err = g.activateRootMemberWallet(ctx, rootDomainID, cb)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	indexMap := make(map[string]string)

	discoveryNodes, indexMap, err := g.addDiscoveryIndex(ctx, cb, indexMap)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)

	}
	indexMap, err = g.addIndex(ctx, cb, indexMap)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)

	}

	err = g.updateNodeDomainIndex(ctx, nodeDomainDesc, indexMap)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}

	return discoveryNodes, nil
}

type genesisNode struct {
	node    certificate.BootstrapNode
	privKey platformpolicy.PrivateKey
	ref     *insolar.Reference
	role    string
}

func (g *Genesis) activateDiscoveryNodes(ctx context.Context, cb *ContractsBuilder, nodesInfo []nodeInfo) ([]genesisNode, error) {
	if len(nodesInfo) != len(g.config.DiscoveryNodes) {
		return nil, errors.New("[ activateDiscoveryNodes ] len of nodesInfo param must be equal to len of DiscoveryNodes in genesis config")
	}

	nodes := make([]genesisNode, len(g.config.DiscoveryNodes))

	for i, discoverNode := range g.config.DiscoveryNodes {
		privKey := nodesInfo[i].privateKey
		nodePubKey := nodesInfo[i].publicKey

		nodeState := &noderecord.NodeRecord{
			Record: noderecord.RecordInfo{
				PublicKey: nodePubKey,
				Role:      insolar.GetStaticRoleFromString(discoverNode.Role),
			},
		}
		contract, err := g.activateNodeRecord(ctx, cb, nodeState, "discoverynoderecord_"+strconv.Itoa(i))
		if err != nil {
			return nil, errors.Wrap(err, "[ activateDiscoveryNodes ] Couldn't activateNodeRecord node instance")
		}

		nodes[i] = genesisNode{
			node: certificate.BootstrapNode{
				PublicKey: nodePubKey,
				Host:      discoverNode.Host,
				NodeRef:   contract.String(),
			},
			privKey: privKey,
			ref:     contract,
			role:    discoverNode.Role,
		}
	}
	return nodes, nil
}

func (g *Genesis) activateNodes(ctx context.Context, cb *ContractsBuilder, nodes []nodeInfo) ([]nodeInfo, error) {
	var updatedNodes []nodeInfo

	for i, node := range nodes {
		nodeState := &noderecord.NodeRecord{
			Record: noderecord.RecordInfo{
				PublicKey: node.publicKey,
				Role:      insolar.StaticRoleVirtual,
			},
		}
		contract, err := g.activateNodeRecord(ctx, cb, nodeState, "noderecord_"+strconv.Itoa(i))
		if err != nil {
			return nil, errors.Wrap(err, "[ activateNodes ] Couldn't activateNodeRecord node instance")
		}
		updatedNode := nodeInfo{
			ref:       contract,
			publicKey: node.publicKey,
		}
		updatedNodes = append(updatedNodes, updatedNode)
	}

	return updatedNodes, nil
}

func (g *Genesis) activateNodeRecord(ctx context.Context, cb *ContractsBuilder, record *noderecord.NodeRecord, name string) (*insolar.Reference, error) {
	nodeData, err := serializeInstance(record)
	if err != nil {
		return nil, errors.Wrap(err, "[ activateNodeRecord ] Couldn't serialize node instance")
	}

	nodeID, err := g.ArtifactManager.RegisterRequest(ctx, *g.rootDomainRef, &message.Parcel{Msg: &message.GenesisRequest{Name: name}})
	if err != nil {
		return nil, errors.Wrap(err, "[ activateNodeRecord ] Couldn't register request to artifact manager")
	}
	contract := insolar.NewReference(*g.rootDomainRef.Record(), *nodeID)
	_, err = g.ArtifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		*g.nodeDomainRef,
		*cb.Prototypes[nodeRecord],
		false,
		nodeData,
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ activateNodeRecord ] Could'n activateNodeRecord node object")
	}
	_, err = g.ArtifactManager.RegisterResult(ctx, *g.rootDomainRef, *contract, nil)
	if err != nil {
		return nil, errors.Wrap(err, "[ activateNodeRecord ] Couldn't register result to artifact manager")
	}
	return contract, nil
}

func (g *Genesis) addDiscoveryIndex(ctx context.Context, cb *ContractsBuilder, indexMap map[string]string) ([]genesisNode, map[string]string, error) {
	errMsg := "[ addDiscoveryIndex ]"
	discoveryKeysPath, err := absPath(g.config.DiscoveryKeysDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, errMsg)
	}
	discoveryKeys, err := g.uploadKeys(ctx, discoveryKeysPath, len(g.config.DiscoveryNodes))
	if err != nil {
		return nil, nil, errors.Wrap(err, errMsg)
	}
	discoveryNodes, err := g.activateDiscoveryNodes(ctx, cb, discoveryKeys)
	if err != nil {
		return nil, nil, errors.Wrap(err, errMsg)
	}
	for _, node := range discoveryNodes {
		indexMap[node.node.PublicKey] = node.ref.String()
	}
	return discoveryNodes, indexMap, nil
}

func (g *Genesis) addIndex(ctx context.Context, cb *ContractsBuilder, indexMap map[string]string) (map[string]string, error) {
	errMsg := "[ addIndex ]"
	nodeKeysPath, err := absPath(g.config.NodeKeysDir)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	userKeys, err := g.uploadKeys(ctx, nodeKeysPath, nodeAmount)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	nodes, err := g.activateNodes(ctx, cb, userKeys)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	for _, node := range nodes {
		indexMap[node.publicKey] = node.ref.String()
	}
	return indexMap, nil
}

func (g *Genesis) createKeys(ctx context.Context, path string, amount int) error {
	err := os.RemoveAll(path)
	if err != nil {
		return errors.Wrap(err, "[ createKeys ] couldn't remove old dir")
	}

	for i := 0; i < amount; i++ {
		ks := commoncrypto.NewKeyProcessor()

		privKey, err := ks.GeneratePrivateKey()
		if err != nil {
			return errors.Wrap(err, "[ createKeys ] couldn't generate private key")
		}

		privKeyStr, err := ks.ExportPrivateKeyPEM(privKey)
		if err != nil {
			return errors.Wrap(err, "[ createKeys ] couldn't export private key")
		}

		pubKeyStr, err := ks.ExportPublicKeyPEM(ks.ExtractPublicKey(privKey))
		if err != nil {
			return errors.Wrap(err, "[ createKeys ] couldn't export public key")
		}

		result, err := json.MarshalIndent(map[string]interface{}{
			"private_key": string(privKeyStr),
			"public_key":  string(pubKeyStr),
		}, "", "    ")
		if err != nil {
			return errors.Wrap(err, "[ createKeys ] couldn't marshal keys")
		}

		name := fmt.Sprintf(g.config.KeysNameFormat, i)
		err = WriteFile(path, name, string(result))
		if err != nil {
			return errors.Wrap(err, "[ createKeys ] couldn't write keys to file")
		}
	}

	return nil
}

func (g *Genesis) uploadKeys(ctx context.Context, path string, amount int) ([]nodeInfo, error) {
	var err error
	if !g.config.ReuseKeys {
		err = g.createKeys(ctx, path, amount)
		if err != nil {
			return nil, errors.Wrap(err, "[ uploadKeys ]")
		}
	}

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, errors.Wrap(err, "[ uploadKeys ] can't read dir")
	}
	if len(files) != amount {
		return nil, errors.New(fmt.Sprintf("[ uploadKeys ] amount of nodes != amount of files in directory: %d != %d", len(files), amount))
	}

	var keys []nodeInfo
	for _, f := range files {
		privKey, nodePubKey, err := getKeysFromFile(ctx, filepath.Join(path, f.Name()))
		if err != nil {
			return nil, errors.Wrap(err, "[ uploadKeys ] can't get keys from file")
		}

		key := nodeInfo{
			publicKey:  nodePubKey,
			privateKey: privKey,
		}
		keys = append(keys, key)
	}

	return keys, nil
}

func (g *Genesis) registerGenesisRequest(ctx context.Context, name string) (*insolar.ID, error) {
	return g.ArtifactManager.RegisterRequest(ctx, *g.ArtifactManager.GenesisRef(), &message.Parcel{Msg: &message.GenesisRequest{Name: name}})
}

// Start creates types and RootDomain instance
func (g *Genesis) Start(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ Genesis ] Starting Genesis ...")

	g.MBLock.Unlock(ctx)
	defer g.MBLock.Lock(ctx)

	rootDomainID, err := g.registerGenesisRequest(ctx, rootDomain)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] Couldn't create rootdomain instance")
	}

	cb := NewContractBuilder(g.ArtifactManager)
	g.prototypeRefs = cb.Prototypes
	defer cb.Clean()

	err = buildSmartContracts(ctx, cb, rootDomainID)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] couldn't build contracts")
	}

	_, rootPubKey, err := getKeysFromFile(ctx, g.config.RootKeysFile)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] couldn't get root keys")
	}

	nodes, err := g.activateSmartContracts(ctx, cb, rootPubKey, rootDomainID)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ]")
	}

	err = g.makeCertificates(nodes)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] Couldn't generate discovery certificates")
	}

	err = utils.SendGracefulStopSignal()
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] Couldn't stop genesis graceful")
	}
	return nil
}

func (g *Genesis) makeCertificates(nodes []genesisNode) error {
	certs := make([]certificate.Certificate, len(nodes))
	for i, node := range nodes {
		certs[i].Role = node.role
		certs[i].Reference = node.ref.String()
		certs[i].PublicKey = node.node.PublicKey
		certs[i].RootDomainReference = g.rootDomainRef.String()
		certs[i].MajorityRule = g.config.MajorityRule
		certs[i].MinRoles.Virtual = g.config.MinRoles.Virtual
		certs[i].MinRoles.HeavyMaterial = g.config.MinRoles.HeavyMaterial
		certs[i].MinRoles.LightMaterial = g.config.MinRoles.LightMaterial
		certs[i].BootstrapNodes = make([]certificate.BootstrapNode, len(nodes))
		for j, node := range nodes {
			certs[i].BootstrapNodes[j] = node.node
		}
	}

	var err error
	for i := range nodes {
		for j, node := range nodes {
			certs[i].BootstrapNodes[j].NetworkSign, err = certs[i].SignNetworkPart(node.privKey)
			if err != nil {
				return errors.Wrapf(err, "[ makeCertificates ] Can't SignNetworkPart for %s", node.ref.String())
			}

			certs[i].BootstrapNodes[j].NodeSign, err = certs[i].SignNodePart(node.privKey)
			if err != nil {
				return errors.Wrapf(err, "[ makeCertificates ] Can't SignNodePart for %s", node.ref.String())
			}
		}

		// save cert to disk
		cert, err := json.MarshalIndent(certs[i], "", "  ")
		if err != nil {
			return errors.Wrapf(err, "[ makeCertificates ] Can't MarshalIndent")
		}

		if len(g.config.DiscoveryNodes[i].CertName) == 0 {
			return errors.New("[ makeCertificates ] cert_name must not be empty for node " + strconv.Itoa(i+1))
		}

		err = ioutil.WriteFile(path.Join(g.keyOut, g.config.DiscoveryNodes[i].CertName), cert, 0644)
		if err != nil {
			return errors.Wrap(err, "[ makeCertificates ] WriteFile")
		}
	}
	return nil
}

func (g *Genesis) updateNodeDomainIndex(ctx context.Context, nodeDomainDesc artifacts.ObjectDescriptor, indexMap map[string]string) error {
	updateData, err := serializeInstance(&nodedomain.NodeDomain{NodeIndexPK: indexMap})
	if err != nil {
		return errors.Wrap(err, "[ updateNodeDomainIndex ]  Couldn't serialize NodeDomain")
	}

	_, err = g.ArtifactManager.UpdateObject(
		ctx,
		*g.rootDomainRef,
		*g.nodeDomainRef,
		nodeDomainDesc,
		updateData,
	)
	if err != nil {
		return errors.Wrap(err, "[ updateNodeDomainIndex ]  Couldn't update NodeDomain")
	}

	return nil
}
