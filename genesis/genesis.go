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
	"encoding/json"
	"strconv"

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/nodedomain"
	"github.com/insolar/insolar/application/contract/noderecord"
	"github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/application/contract/wallet"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
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
	ArtifactManager core.ArtifactManager `inject:""`
	PulseManager    core.PulseManager    `inject:""`
	JetCoordinator  core.JetCoordinator  `inject:""`
	Network         core.Network         `inject:""`
}

// Info returns json with references for info api endpoint
func (g *Genesis) Info() ([]byte, error) {
	prototypes := map[string]string{}
	for prototype, ref := range g.prototypeRefs {
		prototypes[prototype] = ref.String()
	}
	return json.MarshalIndent(map[string]interface{}{
		"root_domain": g.rootDomainRef.String(),
		"root_member": g.rootMemberRef.String(),
		"prototypes":  prototypes,
	}, "", "   ")
}

// GetRootDomainRef returns reference to RootDomain instance
func (g *Genesis) GetRootDomainRef() *core.RecordRef {
	return g.rootDomainRef
}

// NewGenesis creates new Genesis
func NewGenesis(isGenesis bool, genesisConfigPath string) /*nodesInfo []map[string]string)*/ (*Genesis, error) {
	var err error
	genesis := &Genesis{}
	genesis.rootDomainRef = &core.RecordRef{}
	genesis.isGenesis = isGenesis
	if isGenesis {
		genesis.config, err = parseGenesisConfig(genesisConfigPath)
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
	//ctx context.Context, domain *core.RecordID, cb *goplugintestutils.ContractsBuilder,
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
	//func (g *Genesis) activateSmartContracts(ctx context.Context, cb *goplugintestutils.ContractsBuilder) error {
	//	domain, domainDesc, err := g.activateRootDomain(ctx, cb)
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

// Start creates types and RootDomain instance
func (g *Genesis) Start(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ Bootstrapper ] Starting Bootstrap ...")

	rootDomainRef, err := g.getRootDomainRef(ctx)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't get ref of rootDomain")
	}
	if rootDomainRef != nil {
		g.rootDomainRef = rootDomainRef

		rootMemberRef, err := g.getRootMemberRef(ctx, *g.rootDomainRef)
		if err != nil {
			return errors.Wrap(err, "[ Bootstrapper ] couldn't get ref of rootMember")
		}

		g.rootMemberRef = rootMemberRef
		inslog.Info("[ Bootstrapper ] RootDomain was found in ledger. Don't do bootstrap")
		return nil
	}

	//g.rootPubKey, err = getRootMemberPubKey(ctx, g.rootKeysFile)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't get root member keys")
	}

	isLightExecutor, err := g.isLightExecutor(ctx)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't check if node is light executor")
	}
	if !isLightExecutor {
		inslog.Info("[ Bootstrapper ] Node is not light executor. Don't do bootstrap")
		return nil
	}

	_, insgocc, err := goplugintestutils.Build()
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't build insgocc")
	}

	cb := goplugintestutils.NewContractBuilder(g.ArtifactManager, insgocc)
	g.prototypeRefs = cb.Prototypes
	defer cb.Clean()

	err = buildSmartContracts(ctx, cb)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't build contracts")
	}

	if g.isGenesis {
		_, rootPubKey, err := getKeysFromFile(ctx, g.config.RootKeysFile)
		if err != nil {
			return errors.Wrap(err, "[ Bootstrapper ] couldn't get root keys")
		}

		err = g.activateSmartContracts(ctx, cb, rootPubKey)
		if err != nil {
			return errors.Wrap(err, "[ Bootstrapper ]")
		}

		for i, discoverNode := range g.config.DiscoveryNodes {
			//_, nodePubKey, err := getKeysFromFile(ctx, discoverNode.KeysFile)
			nodePubKey := ""
			if err != nil {
				log.Fatal(err)
			}

			nodeState := &noderecord.NodeRecord{
				Record: noderecord.RecordInfo{
					PublicKey: nodePubKey,
					Role:      core.GetRoleFromString(discoverNode.Role),
				},
			}
			nodeData, err := serializeInstance(nodeState)
			if err != nil {
				return errors.Wrap(err, "")
			}

			nodeID, err := g.ArtifactManager.RegisterRequest(ctx, &message.Parcel{Msg: &message.GenesisRequest{Name: "noderecord_" + strconv.Itoa(i)}})
			if err != nil {
				return errors.Wrap(err, "")
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
				return errors.Wrap(err, "")
			}

		}
	}

	return nil
}

func (g *Genesis) isLightExecutor(ctx context.Context) (bool, error) {
	currentPulse, err := g.PulseManager.Current(ctx)
	if err != nil {
		return false, errors.Wrap(err, "[ isLightExecutor ] couldn't get current pulse")
	}

	nodeID := g.Network.GetNodeID()

	isLightExecutor, err := g.JetCoordinator.IsAuthorized(
		ctx,
		core.RoleLightExecutor,
		g.ArtifactManager.GenesisRef(),
		currentPulse.PulseNumber,
		nodeID,
	)
	if err != nil {
		return false, errors.Wrap(err, "[ isLightExecutor ] couldn't authorized node")
	}
	if !isLightExecutor {
		inslogger.FromContext(ctx).Info("[ isLightExecutor ] Is not light executor. Don't build contracts")
		return false, nil
	}
	return true, nil
}

func (g *Genesis) getRootDomainRef(ctx context.Context) (*core.RecordRef, error) {
	genesisRef := g.ArtifactManager.GenesisRef()
	if genesisRef == nil {
		return nil, errors.New("[ getRootDomainRef ] Genesis ref is nil")
	}
	rootObj, err := g.ArtifactManager.GetObject(ctx, *genesisRef, nil, true)
	if err != nil {
		return nil, errors.Wrap(err, "[ getRootDomainRef ] couldn't get children of GenesisRef object")
	}
	rootRefChildren, err := rootObj.Children(nil)
	if err != nil {
		return nil, err
	}
	if rootRefChildren.HasNext() {
		rootDomainRef, err := rootRefChildren.Next()
		if err != nil {
			return nil, errors.Wrap(err, "[ getRootDomainRef ] couldn't get next child of GenesisRef object")
		}
		return rootDomainRef, nil
	}
	return nil, nil
}

func (g *Genesis) getRootMemberRef(ctx context.Context, rootDomainRef core.RecordRef) (*core.RecordRef, error) {
	rootDomainObj, err := g.ArtifactManager.GetObject(ctx, rootDomainRef, nil, false)
	if err != nil {
		return nil, errors.Wrap(err, "[ getRootMemberRef ] couldn't get children of RootDomain object")
	}
	rootDomainRefChildren, err := rootDomainObj.Children(nil)
	if err != nil {
		return nil, err
	}
	if rootDomainRefChildren.HasNext() {
		rootMemberRef, err := rootDomainRefChildren.Next()
		if err != nil {
			return nil, errors.Wrap(err, "[ getRootMemberRef ] couldn't get next child of RootDomain object")
		}
		return rootMemberRef, nil
	}
	return nil, nil
}
