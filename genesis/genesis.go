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

package genesis

import (
	"context"
	"crypto"
	"encoding/json"
	"io/ioutil"
	"log"
	"path"
	"strconv"

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/nodedomain"
	"github.com/insolar/insolar/application/contract/noderecord"
	"github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/application/contract/wallet"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/utils"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/goplugin/goplugintestutils"
	"github.com/pkg/errors"
)

const (
	nodeDomain        = "nodedomain"
	nodeRecord        = "noderecord"
	rootDomain        = "rootdomain"
	walletContract    = "wallet"
	memberContract    = "member"
	allowanceContract = "allowance"
)

var contractNames = []string{walletContract, memberContract, allowanceContract, rootDomain, nodeDomain, nodeRecord}

// Genesis is a component for precreation core contracts types and RootDomain instance
type Genesis struct {
	rootDomainRef   *core.RecordRef
	nodeDomainRef   *core.RecordRef
	rootMemberRef   *core.RecordRef
	prototypeRefs   map[string]*core.RecordRef
	isGenesis       bool
	config          *genesisConfig
	keyOut          string
	ArtifactManager core.ArtifactManager `inject:""`
	PulseManager    core.PulseManager    `inject:""`
	JetCoordinator  core.JetCoordinator  `inject:""`
	Network         core.Network         `inject:""`
}

// NewGenesis creates new Genesis
func NewGenesis(isGenesis bool, genesisConfigPath string, genesisKeyOut string) (*Genesis, error) {
	var err error
	genesis := &Genesis{}
	genesis.rootDomainRef = &core.RecordRef{}
	genesis.isGenesis = isGenesis
	if isGenesis {
		genesis.config, err = parseGenesisConfig(genesisConfigPath)
		genesis.keyOut = genesisKeyOut
	}
	return genesis, err
}

func buildSmartContracts(ctx context.Context, cb *goplugintestutils.ContractsBuilder) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ buildSmartContracts ] building contracts:", contractNames)
	contracts, err := getContractsMap()
	if err != nil {
		return errors.Wrap(err, "[ buildSmartContracts ] couldn't build contracts")
	}

	inslog.Info("[ buildSmartContracts ] Start building contracts ...")
	err = cb.Build(contracts)
	if err != nil {
		return errors.Wrap(err, "[ buildSmartContracts ] couldn't build contracts")
	}
	inslog.Info("[ buildSmartContracts ] Stop building contracts ...")

	return nil
}

func (g *Genesis) activateRootDomain(
	ctx context.Context, cb *goplugintestutils.ContractsBuilder,
) (*core.RecordID, core.ObjectDescriptor, error) {
	rd, err := rootdomain.NewRootDomain()
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	instanceData, err := serializeInstance(rd)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	contractID, err := g.ArtifactManager.RegisterRequest(ctx, &message.Parcel{Msg: &message.GenesisRequest{Name: "RootDomain"}})

	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	contract := core.NewRecordRef(*contractID, *contractID)
	desc, err := g.ArtifactManager.ActivateObject(
		ctx,
		core.RecordRef{},
		*contract,
		*g.ArtifactManager.GenesisRef(),
		*cb.Prototypes[rootDomain],
		false,
		instanceData,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	g.rootDomainRef = contract

	return contractID, desc, nil
}

func (g *Genesis) activateNodeDomain(
	ctx context.Context, domain *core.RecordID, cb *goplugintestutils.ContractsBuilder,
) error {
	nd, err := nodedomain.NewNodeDomain()
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ]")
	}

	instanceData, err := serializeInstance(nd)
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ]")
	}

	contractID, err := g.ArtifactManager.RegisterRequest(ctx, &message.Parcel{Msg: &message.GenesisRequest{Name: "NodeDomain"}})

	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}
	contract := core.NewRecordRef(*domain, *contractID)
	_, err = g.ArtifactManager.ActivateObject(
		ctx,
		core.RecordRef{},
		*contract,
		*g.rootDomainRef,
		*cb.Prototypes[nodeDomain],
		false,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}

	g.nodeDomainRef = contract

	return nil
}

func (g *Genesis) activateRootMember(
	ctx context.Context, domain *core.RecordID, cb *goplugintestutils.ContractsBuilder, rootPubKey string,
) error {

	m, err := member.New("RootMember", rootPubKey)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ]")
	}

	instanceData, err := serializeInstance(m)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ]")
	}

	contractID, err := g.ArtifactManager.RegisterRequest(ctx, &message.Parcel{Msg: &message.GenesisRequest{Name: "RootMember"}})

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	contract := core.NewRecordRef(*domain, *contractID)
	_, err = g.ArtifactManager.ActivateObject(
		ctx,
		core.RecordRef{},
		*contract,
		*g.rootDomainRef,
		*cb.Prototypes[memberContract],
		false,
		instanceData,
	)

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	g.rootMemberRef = contract
	return nil
}

// TODO: this is not required since we refer by request id.
func (g *Genesis) updateRootDomain(
	ctx context.Context, cb *goplugintestutils.ContractsBuilder, domainDesc core.ObjectDescriptor,
) error {
	updateData, err := serializeInstance(&rootdomain.RootDomain{RootMember: *g.rootMemberRef, NodeDomainRef: *g.nodeDomainRef})
	if err != nil {
		return errors.Wrap(err, "[ updateRootDomain ]")
	}
	_, err = g.ArtifactManager.UpdateObject(
		ctx,
		core.RecordRef{},
		core.RecordRef{},
		domainDesc,
		updateData,
	)
	if err != nil {
		return errors.Wrap(err, "[ updateRootDomain ]")
	}

	return nil
}

func (g *Genesis) activateRootMemberWallet(
	ctx context.Context, domain *core.RecordID, cb *goplugintestutils.ContractsBuilder,
) error {
	w, err := wallet.New(g.config.RootBalance)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	instanceData, err := serializeInstance(w)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	contractID, err := g.ArtifactManager.RegisterRequest(ctx, &message.Parcel{Msg: &message.GenesisRequest{Name: "RootWallet"}})

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}
	contract := core.NewRecordRef(*domain, *contractID)
	_, err = g.ArtifactManager.ActivateObject(
		ctx,
		core.RecordRef{},
		*contract,
		*g.rootMemberRef,
		*cb.Prototypes[walletContract],
		true,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}

	return nil
}

func (g *Genesis) activateSmartContracts(ctx context.Context, cb *goplugintestutils.ContractsBuilder, rootPubKey string) error {
	domain, domainDesc, err := g.activateRootDomain(ctx, cb)
	errMsg := "[ ActivateSmartContracts ]"
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = g.activateNodeDomain(ctx, domain, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = g.activateRootMember(ctx, domain, cb, rootPubKey)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	// TODO: this is not required since we refer by request id.
	err = g.updateRootDomain(ctx, cb, domainDesc)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = g.activateRootMemberWallet(ctx, domain, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

type genesisNode struct {
	node    certificate.BootstrapNode
	privKey crypto.PrivateKey
	ref     *core.RecordRef
	role    string
}

func (g *Genesis) registerDiscoveryNodes(ctx context.Context, cb *goplugintestutils.ContractsBuilder) ([]genesisNode, error) {

	nodes := make([]genesisNode, len(g.config.DiscoveryNodes))

	for i, discoverNode := range g.config.DiscoveryNodes {
		privKey, nodePubKey, err := getKeysFromFile(ctx, discoverNode.KeysFile)
		if err != nil {
			log.Fatal(err)
		}

		nodeState := &noderecord.NodeRecord{
			Record: noderecord.RecordInfo{
				PublicKey: nodePubKey,
				Role:      core.GetStaticRoleFromString(discoverNode.Role),
			},
		}
		nodeData, err := serializeInstance(nodeState)
		if err != nil {
			return nil, errors.Wrap(err, "[ registerDiscoveryNodes ] Couldn't serialize discovery node instance")
		}

		nodeID, err := g.ArtifactManager.RegisterRequest(ctx, &message.Parcel{Msg: &message.GenesisRequest{Name: "noderecord_" + strconv.Itoa(i)}})
		if err != nil {
			return nil, errors.Wrap(err, "[ registerDiscoveryNodes ] Couldn't register request to artifact manager")
		}
		contract := core.NewRecordRef(*g.rootDomainRef.Record(), *nodeID)
		_, err = g.ArtifactManager.ActivateObject(
			ctx,
			core.RecordRef{},
			*contract,
			*g.nodeDomainRef,
			*cb.Prototypes[nodeRecord],
			false,
			nodeData,
		)
		if err != nil {
			return nil, errors.Wrap(err, "[ registerDiscoveryNodes ] Could'n activate discovery node object")
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

// Start creates types and RootDomain instance
func (g *Genesis) Start(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ Genesis ] Starting Genesis ...")

	_, insgocc, err := goplugintestutils.Build()
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] couldn't build insgocc")
	}

	cb := goplugintestutils.NewContractBuilder(g.ArtifactManager, insgocc)
	g.prototypeRefs = cb.Prototypes
	defer cb.Clean()

	err = buildSmartContracts(ctx, cb)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] couldn't build contracts")
	}

	_, rootPubKey, err := getKeysFromFile(ctx, g.config.RootKeysFile)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] couldn't get root keys")
	}

	err = g.activateSmartContracts(ctx, cb, rootPubKey)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ]")
	}

	nodes, err := g.registerDiscoveryNodes(ctx, cb)
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
				return errors.Wrap(err, "[ makeCertificates ]")
			}
			certs[i].BootstrapNodes[j].NodeSign, err = certs[i].SignNodePart(node.privKey)
			if err != nil {
				return errors.Wrap(err, "[ makeCertificates ]")
			}
		}

		// save cert to disk
		cert, err := json.MarshalIndent(certs[i], "", "  ")
		if err != nil {
			return errors.Wrap(err, "[ makeCertificates ]")
		}
		err = ioutil.WriteFile(path.Join(g.keyOut, "discovery_cert_"+strconv.Itoa(i+1)+".json"), cert, 0644)
		if err != nil {
			return errors.Wrap(err, "[ makeCertificates ]")
		}
	}
	return nil
}
