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
	"container/list"
	"context"
	"crypto"
	"encoding/json"
	"io/ioutil"
	"math/big"
	"path"
	"strconv"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/nodedomain"
	"github.com/insolar/insolar/application/contract/noderecord"
	rootdomaincontract "github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/application/contract/wallet"
	"github.com/insolar/insolar/bootstrap/rootdomain"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/artifact"
)

const (
	nodeDomain      = "nodedomain"
	nodeRecord      = "noderecord"
	rootDomain      = rootdomain.Name
	walletContract  = "wallet"
	memberContract  = "member"
	depositContract = "deposit"
)

var contractNames = []string{walletContract, memberContract, depositContract, rootDomain, nodeDomain, nodeRecord}

type nodeInfo struct {
	privateKey crypto.PrivateKey
	publicKey  string
}

func (ni nodeInfo) reference() insolar.Reference {
	return refByName(ni.publicKey)
}

// Generator is a component for generating RootDomain instance and genesis contracts.
type Generator struct {
	config          *Config
	artifactManager artifact.Manager

	rootRecord            *rootdomain.Record
	rootDomainContract    *insolar.Reference
	nodeDomainContract    *insolar.Reference
	rootMemberContract    *insolar.Reference
	oracleMemberContracts map[string]insolar.Reference
	oracleConfirms        map[string]bool
	mdAdminMemberContract *insolar.Reference
	mdWalletContract      *insolar.Reference

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
	_, rootPubKey, err := getKeysFromFile(g.config.RootKeysFile)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] couldn't get root keys")
	}

	inslog.Info("[ Genesis ] getKeysFromFil for mdAdmine ...")
	_, mdAdminPubKey, err := getKeysFromFile(g.config.MDAdminKeysFile)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] couldn't get root keys")
	}

	inslog.Info("[ Genesis ] getKeysFromFile for oracles ...")
	oracleMap := map[string]string{}
	for _, o := range g.config.OracleKeysFiles {
		_, oraclePubKey, err := getKeysFromFile(o.KeysFile)
		if err != nil {
			return errors.Wrap(err, "[ Genesis ] couldn't get oracle keys for oracle: "+o.Name)
		}
		oracleMap[o.Name] = oraclePubKey
	}

	inslog.Info("[ Genesis ] createKeysInDir ...")
	discoveryNodes, err := createKeysInDir(
		ctx,
		g.config.DiscoveryKeysDir,
		g.config.KeysNameFormat,
		g.config.DiscoveryNodes,
		g.config.ReuseKeys,
	)
	if err != nil {
		return errors.Wrapf(err, "[ Genesis ] create keys step failed")
	}

	err = g.activateSmartContracts(ctx, cb, rootPubKey, mdAdminPubKey, oracleMap, &rootDomainID, prototypes)
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
		record.Request{
			CallType: record.CTGenesis,
			Method:   rootDomain,
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
		record.Request{
			CallType: record.CTGenesis,
			Method:   "NodeDomain",
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
		record.Request{
			CallType: record.CTGenesis,
			Method:   "RootMember",
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

func (g *Generator) activateMDAdminMember(
	ctx context.Context,
	domain *insolar.ID,
	mdAdminPubKey string,
	memberContractProto insolar.Reference,
) error {

	m, err := member.New("MDAdminMember", mdAdminPubKey)
	if err != nil {
		return errors.Wrap(err, "[ ActivateMDAdminMember ]")
	}

	instanceData, err := insolar.Serialize(m)
	if err != nil {
		return errors.Wrap(err, "[ ActivateMDAdminMember ]")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   "MDAdminMember",
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ ActivateMDAdminMember ] couldn't create mdAdmin member instance")
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
		return errors.Wrap(err, "[ ActivateMDAdminMember ] couldn't create mdAdmin member instance")
	}
	_, err = g.artifactManager.RegisterResult(ctx, *g.rootDomainContract, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ ActivateMDAdminMember ] couldn't create mdAdmin member instance")
	}
	g.mdAdminMemberContract = contract
	return nil
}

func (g *Generator) activateOracleMembers(
	ctx context.Context,
	domain *insolar.ID,
	oraclePubKeys map[string]string,
	memberContractProto insolar.Reference,
) error {

	g.oracleMemberContracts = map[string]insolar.Reference{}
	g.oracleConfirms = map[string]bool{}
	for name, _ := range oraclePubKeys {
		g.oracleConfirms[name] = false
	}

	for name, key := range oraclePubKeys {
		m, err := member.NewOracleMember(name, key)
		if err != nil {
			return errors.Wrap(err, "[ ActivateOracleMember ]")
		}

		instanceData, err := insolar.Serialize(m)
		if err != nil {
			return errors.Wrap(err, "[ ActivateOracleMember ]")
		}

		contractID, err := g.artifactManager.RegisterRequest(
			ctx,
			record.Request{
				CallType: record.CTGenesis,
				Method:   "OracleMember",
			},
		)

		if err != nil {
			return errors.Wrap(err, "[ ActivateOracleMember ] couldn't create oracle member instance")
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
			return errors.Wrap(err, "[ ActivateOracleMember ] couldn't create oracle member instance")
		}
		_, err = g.artifactManager.RegisterResult(ctx, *g.rootDomainContract, *contract, nil)
		if err != nil {
			return errors.Wrap(err, "[ ActivateOracleMember ] couldn't create oracle member instance")
		}

		g.oracleMemberContracts[name] = *contract
	}

	return nil
}

func (g *Generator) updateRootDomain(
	ctx context.Context, domainDesc artifact.ObjectDescriptor,
) error {
	updateData, err := insolar.Serialize(&rootdomaincontract.RootDomain{
		RootMember:        *g.rootMemberContract,
		OracleMembers:     g.oracleMemberContracts,
		MDAdminMember:     *g.mdAdminMemberContract,
		MDWallet:          *g.mdWalletContract,
		BurnAddressMap:    map[string]insolar.Reference{},
		PublicKeyMap:      map[string]insolar.Reference{},
		FreeBurnAddresses: *list.New(),
		NodeDomain:        *g.nodeDomainContract,
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

func (g *Generator) activateRootWallet(
	ctx context.Context, domain *insolar.ID, walletContractProto insolar.Reference,
) error {

	b := new(big.Int)
	b, ok := b.SetString(g.config.RootBalance, 10)
	if !ok {
		return errors.Errorf("[ ActivateRootWallet ] Failed to parse RootBalance")
	}

	w, err := wallet.New(*b)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	instanceData, err := insolar.Serialize(w)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   "RootWallet",
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

func (g *Generator) activateMDWallet(
	ctx context.Context, domain *insolar.ID, walletContractProto insolar.Reference,
) error {

	b := new(big.Int)
	b, ok := b.SetString(g.config.MDBalance, 10)
	if !ok {
		return errors.Errorf("[ activateMDWallet ] Failed to parse MDBalance")
	}

	w, err := wallet.New(*b)
	if err != nil {
		return errors.Wrap(err, "[ activateMDWallet ]")
	}

	instanceData, err := insolar.Serialize(w)
	if err != nil {
		return errors.Wrap(err, "[ activateMDWallet ]")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   "MDWallet",
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ activateMDWallet ] couldn't create root wallet")
	}
	contract := insolar.NewReference(*domain, *contractID)
	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		*g.mdAdminMemberContract,
		walletContractProto,
		true,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ activateMDWallet ] couldn't create root wallet")
	}
	_, err = g.artifactManager.RegisterResult(ctx, *g.rootDomainContract, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ activateMDWallet ] couldn't create root wallet")
	}

	g.mdWalletContract = contract
	return nil
}

func (g *Generator) activateSmartContracts(
	ctx context.Context,
	cb *contractsBuilder,
	rootPubKey string,
	mdAdminPubKey string,
	oracleMap map[string]string,
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
	err = g.activateRootWallet(ctx, rootDomainID, *cb.prototypes[walletContract])
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = g.activateMDAdminMember(ctx, rootDomainID, mdAdminPubKey, *cb.prototypes[memberContract])
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = g.activateMDWallet(ctx, rootDomainID, *cb.prototypes[walletContract])
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = g.activateOracleMembers(ctx, rootDomainID, oracleMap, *cb.prototypes[memberContract])
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	// TODO: this is not required since we refer by request id.
	err = g.updateRootDomain(ctx, rootDomainDesc)
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
		nodePubKey := nodesInfo[i].publicKey

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
	nRecord *noderecord.NodeRecord,
	node nodeInfo,
	nodeRecordProto insolar.Reference,
) (*insolar.Reference, error) {
	nodeData, err := insolar.Serialize(nRecord)
	if err != nil {
		return nil, errors.Wrap(err, "[ activateNodeRecord ] Couldn't serialize node instance")
	}

	nodeID, err := g.artifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   node.publicKey,
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
	certs := make([]certificate.Certificate, 0, len(g.config.DiscoveryNodes))
	for i, node := range g.config.DiscoveryNodes {
		pubKey := discoveryNodes[i].publicKey
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
			pk := discoveryNodes[j].publicKey
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

		inslogger.FromContext(ctx).Debugf("[makeCertificates] write cert file to %v", certFile)
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
		indexMap[node.publicKey] = node.reference().String()
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
