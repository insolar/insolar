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
	"crypto"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strconv"

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/nodedomain"
	"github.com/insolar/insolar/application/contract/noderecord"
	rootdomaincontract "github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/application/contract/wallet"
	"github.com/insolar/insolar/bootstrap/rootdomain"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/artifact"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

const (
	nodeDomain        = "nodedomain"
	nodeRecord        = "noderecord"
	rootDomain        = rootdomain.Name
	walletContract    = "wallet"
	memberContract    = "member"
	allowanceContract = "allowance"
)

var contractNames = []string{walletContract, memberContract, allowanceContract, rootDomain, nodeDomain, nodeRecord}

type nodeInfo struct {
	privateKey crypto.PrivateKey
	publicKey  string
}

// Generator is a component for generating RootDomain instance and genesis contracts.
type Generator struct {
	artifactManager artifact.Manager
	config          *Config

	rootRecord         *rootdomain.Record
	rootDomainContract *insolar.Reference
	nodeDomainContract *insolar.Reference
	rootMemberContract *insolar.Reference

	keyOut string
}

// NewGenerator creates new Generator.
func NewGenerator(
	config *Config,
	am artifact.Manager,
	rootRecord *rootdomain.Record,
	genesisKeyOut string,
) *Generator {
	return &Generator{
		artifactManager: am,
		config:          config,

		rootRecord: rootRecord,

		keyOut: genesisKeyOut,
	}
}

// Run generates genesis data.
func (g *Generator) Run(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ Genesis ] Starting  ...")
	defer inslog.Info("[ Genesis ] Finished.")

	rootDomainID := g.rootRecord.ID()

	inslog.Info("[ Genesis ] newContractBuilder ...")
	cb := newContractBuilder(g.artifactManager)
	defer cb.clean()

	// TODO: don't build prototypes, just get they references from builtins
	inslog.Info("[ Genesis ] buildSmartContracts ...")
	prototypes, err := cb.buildPrototypes(ctx, &rootDomainID)
	if err != nil {
		panic(errors.Wrap(err, "[ Genesis ] couldn't build contracts"))
	}

	inslog.Info("[ Genesis ] getKeysFromFile ...")
	_, rootPubKey, err := getKeysFromFile(ctx, g.config.RootKeysFile)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] couldn't get root keys")
	}

	inslog.Info("[ Genesis ] activateSmartContracts ...")
	nodes, err := g.activateSmartContracts(ctx, cb, rootPubKey, &rootDomainID, prototypes)
	if err != nil {
		panic(errors.Wrap(err, "[ Genesis ] could't activate smart contracts"))
	}

	inslog.Info("[ Genesis ] makeCertificates ...")
	err = g.makeCertificates(ctx, nodes)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] Couldn't generate discovery certificates")
	}

	return nil
}

func (g *Generator) activateRootDomain(
	ctx context.Context,
	cb *contractsBuilder,
	id *insolar.ID,
) (artifact.ObjectDescriptor, error) {
	rdContract, err := rootdomaincontract.NewRootDomain()
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	instanceData, err := insolar.Serialize(rdContract)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	_, err = g.artifactManager.RegisterRequest(
		ctx,
		insolar.GenesisRecord.Ref(),
		&message.Parcel{
			Msg: &message.GenesisRequest{
				Name: rootDomain,
			},
		},
	)
	if err != nil {
		panic(errors.Wrap(err, "[ Genesis ] Couldn't create rootdomain instance"))
	}

	rootDomainRef := g.rootRecord.Ref()
	desc, err := g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		rootDomainRef,
		insolar.GenesisRecord.Ref(),
		*cb.prototypes[rootDomain],
		false,
		instanceData,
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	_, err = g.artifactManager.RegisterResult(ctx, insolar.GenesisRecord.Ref(), rootDomainRef, nil)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	g.rootDomainContract = &rootDomainRef

	return desc, nil
}

func (g *Generator) activateNodeDomain(
	ctx context.Context, domain *insolar.ID, nodeDomainProto insolar.Reference,
) (artifact.ObjectDescriptor, error) {
	nd, _ := nodedomain.NewNodeDomain()

	instanceData, err := insolar.Serialize(nd)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateNodeDomain ] node domain serialization")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		*g.rootDomainContract,
		&message.Parcel{
			Msg: &message.GenesisRequest{Name: "NodeDomain"},
		},
	)

	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}
	contract := insolar.NewReference(*domain, *contractID)
	desc, err := g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		*g.rootDomainContract,
		nodeDomainProto,
		false,
		instanceData,
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}
	_, err = g.artifactManager.RegisterResult(ctx, *g.rootDomainContract, *contract, nil)
	if err != nil {
		return nil, errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}

	g.nodeDomainContract = contract

	return desc, nil
}

func (g *Generator) activateRootMember(
	ctx context.Context,
	domain *insolar.ID,
	rootPubKey string,
	memberContractProto insolar.Reference,
) error {

	m, err := member.New("RootMember", rootPubKey)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ]")
	}

	instanceData, err := insolar.Serialize(m)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ]")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		*g.rootDomainContract,
		&message.Parcel{
			Msg: &message.GenesisRequest{Name: "RootMember"},
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	contract := insolar.NewReference(*domain, *contractID)
	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		*g.rootDomainContract,
		memberContractProto,
		false,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	_, err = g.artifactManager.RegisterResult(ctx, *g.rootDomainContract, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	g.rootMemberContract = contract
	return nil
}

// TODO: this is not required since we refer by request id.
func (g *Generator) updateRootDomain(
	ctx context.Context, domainDesc artifact.ObjectDescriptor,
) error {
	updateData, err := insolar.Serialize(&rootdomaincontract.RootDomain{
		RootMember:    *g.rootMemberContract,
		NodeDomainRef: *g.nodeDomainContract,
	})
	if err != nil {
		return errors.Wrap(err, "[ updateRootDomain ]")
	}
	_, err = g.artifactManager.UpdateObject(
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

func (g *Generator) activateRootMemberWallet(
	ctx context.Context, domain *insolar.ID, walletContractProto insolar.Reference,
) error {

	w, err := wallet.New(g.config.RootBalance)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	instanceData, err := insolar.Serialize(w)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		*g.rootDomainContract,
		&message.Parcel{
			Msg: &message.GenesisRequest{Name: "RootWallet"},
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}
	contract := insolar.NewReference(*domain, *contractID)
	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		*g.rootMemberContract,
		walletContractProto,
		true,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}
	_, err = g.artifactManager.RegisterResult(ctx, *g.rootDomainContract, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}

	return nil
}

func (g *Generator) activateSmartContracts(
	ctx context.Context,
	cb *contractsBuilder,
	rootPubKey string,
	rootDomainID *insolar.ID,
	prototypes prototypes,
) ([]genesisNode, error) {

	rootDomainDesc, err := g.activateRootDomain(ctx, cb, rootDomainID)
	errMsg := "[ ActivateSmartContracts ]"
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	nodeDomainDesc, err := g.activateNodeDomain(ctx, rootDomainID, *prototypes[nodeDomain])
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	err = g.activateRootMember(ctx, rootDomainID, rootPubKey, *cb.prototypes[memberContract])
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	// TODO: this is not required since we refer by request id.
	err = g.updateRootDomain(ctx, rootDomainDesc)
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	err = g.activateRootMemberWallet(ctx, rootDomainID, *cb.prototypes[walletContract])
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	indexMap := make(map[string]string)

	discoveryNodes, indexMap, err := g.addDiscoveryIndex(ctx, indexMap, *cb.prototypes[nodeRecord])
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
	privKey crypto.PrivateKey
	ref     *insolar.Reference
	role    string
}

func (g *Generator) activateDiscoveryNodes(
	ctx context.Context,
	nodeRecordProto insolar.Reference,
	nodesInfo []nodeInfo,
) ([]genesisNode, error) {
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
		contract, err := g.activateNodeRecord(ctx, nodeState, "discoverynoderecord_"+strconv.Itoa(i), nodeRecordProto)
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

func (g *Generator) activateNodeRecord(
	ctx context.Context,
	record *noderecord.NodeRecord,
	name string,
	nodeRecordProto insolar.Reference,
) (*insolar.Reference, error) {
	nodeData, err := insolar.Serialize(record)
	if err != nil {
		return nil, errors.Wrap(err, "[ activateNodeRecord ] Couldn't serialize node instance")
	}

	nodeID, err := g.artifactManager.RegisterRequest(
		ctx,
		*g.rootDomainContract,
		&message.Parcel{
			Msg: &message.GenesisRequest{Name: name},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ activateNodeRecord ] Couldn't register request")
	}
	contract := insolar.NewReference(*g.rootDomainContract.Record(), *nodeID)
	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		*g.nodeDomainContract,
		nodeRecordProto,
		false,
		nodeData,
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ activateNodeRecord ] Could'n activateNodeRecord node object")
	}
	_, err = g.artifactManager.RegisterResult(ctx, *g.rootDomainContract, *contract, nil)
	if err != nil {
		return nil, errors.Wrap(err, "[ activateNodeRecord ] Couldn't register result")
	}
	return contract, nil
}

func (g *Generator) addDiscoveryIndex(
	ctx context.Context,
	indexMap map[string]string,
	nodeRecordProto insolar.Reference,
) ([]genesisNode, map[string]string, error) {
	errMsg := "[ addDiscoveryIndex ]"
	discoveryKeysPath, err := filepath.Abs(g.config.DiscoveryKeysDir)
	if err != nil {
		return nil, nil, errors.Wrap(err, errMsg)
	}

	discoveryKeys, err := g.uploadKeys(ctx, discoveryKeysPath, len(g.config.DiscoveryNodes))
	if err != nil {
		return nil, nil, errors.Wrap(err, errMsg)
	}

	discoveryNodes, err := g.activateDiscoveryNodes(ctx, nodeRecordProto, discoveryKeys)
	if err != nil {
		return nil, nil, errors.Wrap(err, errMsg)
	}

	for _, node := range discoveryNodes {
		indexMap[node.node.PublicKey] = node.ref.String()
	}
	return discoveryNodes, indexMap, nil
}

func (g *Generator) createKeys(ctx context.Context, dir string, amount int) error {
	fmt.Println("createKeys, skip RemoveAll of", dir)
	// err := os.RemoveAll(dir)
	// if err != nil {
	// 	return errors.Wrap(err, "[ createKeys ] couldn't remove old dir")
	// }

	for i := 0; i < amount; i++ {
		ks := platformpolicy.NewKeyProcessor()

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
		err = makeFileWithDir(dir, name, result)
		if err != nil {
			return errors.Wrap(err, "[ createKeys ] couldn't write keys to file")
		}
	}

	return nil
}

func (g *Generator) uploadKeys(ctx context.Context, dir string, amount int) ([]nodeInfo, error) {
	var err error
	if !g.config.ReuseKeys {
		err = g.createKeys(ctx, dir, amount)
		if err != nil {
			return nil, errors.Wrap(err, "[ uploadKeys ]")
		}
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrap(err, "[ uploadKeys ] can't read dir")
	}
	if len(files) != amount {
		return nil, errors.New(fmt.Sprintf("[ uploadKeys ] amount of nodes != amount of files in directory: %d != %d", len(files), amount))
	}

	var keys []nodeInfo
	for _, f := range files {
		privKey, nodePubKey, err := getKeysFromFile(ctx, filepath.Join(dir, f.Name()))
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

func (g *Generator) makeCertificates(ctx context.Context, nodes []genesisNode) error {
	certs := make([]certificate.Certificate, len(nodes))
	for i, node := range nodes {
		certs[i].Role = node.role
		certs[i].Reference = node.ref.String()
		certs[i].PublicKey = node.node.PublicKey
		certs[i].RootDomainReference = g.rootDomainContract.String()
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

		certFile := path.Join(g.keyOut, g.config.DiscoveryNodes[i].CertName)
		err = ioutil.WriteFile(certFile, cert, 0644)
		if err != nil {
			return errors.Wrap(err, "[ makeCertificates ] makeFileWithDir")
		}
	}
	return nil
}

func (g *Generator) updateNodeDomainIndex(ctx context.Context, nodeDomainDesc artifact.ObjectDescriptor, indexMap map[string]string) error {
	updateData, err := insolar.Serialize(
		&nodedomain.NodeDomain{
			NodeIndexPK: indexMap,
		},
	)
	if err != nil {
		return errors.Wrap(err, "[ updateNodeDomainIndex ]  Couldn't serialize NodeDomain")
	}

	_, err = g.artifactManager.UpdateObject(
		ctx,
		*g.rootDomainContract,
		*g.nodeDomainContract,
		nodeDomainDesc,
		updateData,
	)
	if err != nil {
		return errors.Wrap(err, "[ updateNodeDomainIndex ]  Couldn't update NodeDomain")
	}

	return nil
}

func getKeysFromFile(ctx context.Context, file string) (crypto.PrivateKey, string, error) {
	absPath, err := filepath.Abs(file)
	if err != nil {
		return nil, "", errors.Wrap(err, "[ getKeyFromFile ] couldn't get abs path")
	}
	data, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, "", errors.Wrap(err, "[ getKeyFromFile ] couldn't read keys file "+absPath)
	}
	var keys map[string]string
	err = json.Unmarshal(data, &keys)
	if err != nil {
		return nil, "", errors.Wrapf(err, "[ getKeyFromFile ] couldn't unmarshal data from %s", absPath)
	}
	if keys["private_key"] == "" {
		return nil, "", errors.New("[ getKeyFromFile ] empty private key")
	}
	if keys["public_key"] == "" {
		return nil, "", errors.New("[ getKeyFromFile ] empty public key")
	}
	kp := platformpolicy.NewKeyProcessor()
	key, err := kp.ImportPrivateKeyPEM([]byte(keys["private_key"]))
	if err != nil {
		return nil, "", errors.Wrapf(err, "[ getKeyFromFile ] couldn't import private key")
	}
	return key, keys["public_key"], nil
}
