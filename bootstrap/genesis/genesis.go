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
	"io/ioutil"
	"path"
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
	publicKey  crypto.PublicKey
}

func (ni nodeInfo) publicKeyString() string {
	ks := platformpolicy.NewKeyProcessor()

	pubKeyStr, err := ks.ExportPublicKeyPEM(ni.publicKey)
	if err != nil {
		panic(err)
	}
	return string(pubKeyStr)
}

func (ni nodeInfo) reference() insolar.Reference {
	return refByName(ni.publicKeyString())
}

// Generator is a component for generating RootDomain instance and genesis contracts.
type Generator struct {
	config          *Config
	artifactManager artifact.Manager

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
		config:          config,
		artifactManager: am,

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

	inslog.Info("[ Genesis ] createKeysInDir ...")
	discoveryNodes, err := createKeysInDir(
		ctx,
		g.config.DiscoveryKeysDir,
		g.config.KeysNameFormat,
		len(g.config.DiscoveryNodes),
		g.config.ReuseKeys,
	)
	if err != nil {
		return errors.Wrapf(err, "[ Genesis ] create keys step failed")
	}

	err = g.activateSmartContracts(ctx, cb, rootPubKey, discoveryNodes, &rootDomainID, prototypes)
	if err != nil {
		panic(errors.Wrap(err, "[ Genesis ] could't activate smart contracts"))
	}

	err = g.activateDiscoveryNodes(ctx, *cb.prototypes[nodeRecord], discoveryNodes)
	if err != nil {
		return errors.Wrap(err, "failed on adding discovery index")
	}

	err = g.updateNodeDomainIndex(ctx, discoveryNodes)
	if err != nil {
		return errors.Wrap(err, "failed update nodedomai ")
	}

	inslog.Info("[ Genesis ] makeCertificates ...")
	err = g.makeCertificates(ctx, discoveryNodes)
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

	inslogger.FromContext(ctx).Debugf("[activateNodeDomain] Ref: %v", contract)

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
	inslogger.FromContext(ctx).Debugf("%v contract ref=%v", "NodeDomain", contract)

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
	discoveryNodes []nodeInfo,
	rootDomainID *insolar.ID,
	prototypes prototypes,
) error {
	// TODO: merge root domain activation with update (no need two-phase here)
	rootDomainDesc, err := g.activateRootDomain(ctx, cb, rootDomainID)
	errMsg := "[ ActivateSmartContracts ]"
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	_, err = g.activateNodeDomain(ctx, rootDomainID, *prototypes[nodeDomain])
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = g.activateRootMember(ctx, rootDomainID, rootPubKey, *cb.prototypes[memberContract])
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	// TODO: this is not required since we refer by request id.
	err = g.updateRootDomain(ctx, rootDomainDesc)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = g.activateRootMemberWallet(ctx, rootDomainID, *cb.prototypes[walletContract])
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

// activateDiscoveryNodes activates discoverynoderecord_{N} objects.
//
// It returns list of genesisNode structures (for node domain save and certificates generation at the end of genesis).
func (g *Generator) activateDiscoveryNodes(
	ctx context.Context,
	nodeRecordProto insolar.Reference,
	nodesInfo []nodeInfo,
) error {
	if len(nodesInfo) != len(g.config.DiscoveryNodes) {
		return errors.New("[ activateDiscoveryNodes ] len of nodesInfo param must be equal to len of DiscoveryNodes in genesis config")
	}

	for i, discoverNode := range g.config.DiscoveryNodes {
		nodePubKey := nodesInfo[i].publicKeyString()

		nodeState := &noderecord.NodeRecord{
			Record: noderecord.RecordInfo{
				PublicKey: nodePubKey,
				Role:      insolar.GetStaticRoleFromString(discoverNode.Role),
			},
		}

		_, err := g.activateNodeRecord(ctx, nodeState, nodesInfo[i], nodeRecordProto)
		if err != nil {
			return errors.Wrap(err, "[ activateDiscoveryNodes ] Couldn't activateNodeRecord node instance")
		}
	}
	return nil
}

func (g *Generator) activateNodeRecord(
	ctx context.Context,
	record *noderecord.NodeRecord,
	node nodeInfo,
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
			Msg: &message.GenesisRequest{Name: node.publicKeyString()},
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

func (g *Generator) makeCertificates(ctx context.Context, discoveryNodes []nodeInfo) error {
	var certs []certificate.Certificate
	for i, node := range g.config.DiscoveryNodes {
		pubKey := discoveryNodes[i].publicKeyString()
		ref := discoveryNodes[i].reference()

		c := certificate.Certificate{
			AuthorizationCertificate: certificate.AuthorizationCertificate{
				PublicKey: pubKey,
				Role:      node.Role,
				Reference: ref.String(),
			},
			MajorityRule: g.config.MajorityRule,

			RootDomainReference: g.rootDomainContract.String(),
		}
		c.MinRoles.Virtual = g.config.MinRoles.Virtual
		c.MinRoles.HeavyMaterial = g.config.MinRoles.HeavyMaterial
		c.MinRoles.LightMaterial = g.config.MinRoles.LightMaterial
		c.BootstrapNodes = []certificate.BootstrapNode{}

		for j, n2 := range g.config.DiscoveryNodes {
			pk := discoveryNodes[j].publicKeyString()
			ref := discoveryNodes[j].reference()
			c.BootstrapNodes = append(c.BootstrapNodes, certificate.BootstrapNode{
				PublicKey: pk,
				Host:      n2.Host,
				NodeRef:   ref.String(),
			})
		}

		certs = append(certs, c)
	}

	var err error
	for i, node := range g.config.DiscoveryNodes {
		for j := range g.config.DiscoveryNodes {
			dn := discoveryNodes[j]

			certs[i].BootstrapNodes[j].NetworkSign, err = certs[i].SignNetworkPart(dn.privateKey)
			if err != nil {
				return errors.Wrapf(err, "[ makeCertificates ] Can't SignNetworkPart for %s",
					dn.reference())
			}

			certs[i].BootstrapNodes[j].NodeSign, err = certs[i].SignNodePart(dn.privateKey)
			if err != nil {
				return errors.Wrapf(err, "[ makeCertificates ] Can't SignNodePart for %s",
					dn.reference())
			}
		}

		// save cert to disk
		cert, err := json.MarshalIndent(certs[i], "", "  ")
		if err != nil {
			return errors.Wrapf(err, "[ makeCertificates ] Can't MarshalIndent")
		}

		if len(node.CertName) == 0 {
			return errors.New("[ makeCertificates ] cert_name must not be empty for node number " + strconv.Itoa(i+1))
		}

		certFile := path.Join(g.keyOut, node.CertName)
		err = ioutil.WriteFile(certFile, cert, 0644)
		if err != nil {
			return errors.Wrapf(err, "[ makeCertificates ] filed create ceritificate: %v", certFile)
		}
	}
	return nil
}

// updateNodeDomainIndex saves in node domain contract's object discovery nodes map.
func (g *Generator) updateNodeDomainIndex(ctx context.Context, discoveryNodes []nodeInfo) error {
	nodeDomainDesc, err := g.artifactManager.GetObject(ctx, *g.nodeDomainContract)
	if err != nil {
		return errors.Wrap(err, "failed to get domain contract")
	}

	indexMap := map[string]string{}
	for _, node := range discoveryNodes {
		indexMap[node.publicKeyString()] = node.reference().String()
	}

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
