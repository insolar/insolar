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
	rootDomainRef *core.RecordRef
	nodeDomainRef *core.RecordRef
	rootMemberRef *core.RecordRef
	prototypeRefs map[string]*core.RecordRef
	isGenesis     bool
	config        *genesisConfig
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
	genesis.config, err = parseGenesisConfig(genesisConfigPath)
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
	ctx context.Context, am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder,
) (*core.RecordID, core.ObjectDescriptor, error) {
	rd, err := rootdomain.NewRootDomain()
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	instanceData, err := serializeInstance(rd)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	contractID, err := am.RegisterRequest(ctx, &message.GenesisRequest{Name: "RootDomain"})
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	contract := core.NewRecordRef(*contractID, *contractID)
	desc, err := am.ActivateObject(
		ctx,
		core.RecordRef{},
		*contract,
		*am.GenesisRef(),
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
	ctx context.Context, domain *core.RecordID, am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder,
) error {
	nd, err := nodedomain.NewNodeDomain()
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ]")
	}

	instanceData, err := serializeInstance(nd)
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ]")
	}

	contractID, err := am.RegisterRequest(ctx, &message.GenesisRequest{Name: "NodeDomain"})
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}
	contract := core.NewRecordRef(*domain, *contractID)
	_, err = am.ActivateObject(
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
	ctx context.Context, domain *core.RecordID, am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder, rootPubKey string,
) error {

	m, err := member.New("RootMember", rootPubKey)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ]")
	}

	instanceData, err := serializeInstance(m)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ]")
	}

	contractID, err := am.RegisterRequest(ctx, &message.GenesisRequest{Name: "RootMember"})
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	contract := core.NewRecordRef(*domain, *contractID)
	_, err = am.ActivateObject(
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
	ctx context.Context, am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder, domainDesc core.ObjectDescriptor,
) error {
	updateData, err := serializeInstance(&rootdomain.RootDomain{RootMember: *g.rootMemberRef, NodeDomainRef: *g.nodeDomainRef})
	if err != nil {
		return errors.Wrap(err, "[ updateRootDomain ]")
	}
	_, err = am.UpdateObject(
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
	ctx context.Context, domain *core.RecordID, am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder,
) error {
	w, err := wallet.New(g.config.RootBalance)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	instanceData, err := serializeInstance(w)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	contractID, err := am.RegisterRequest(ctx, &message.GenesisRequest{Name: "RootWallet"})
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}
	contract := core.NewRecordRef(*domain, *contractID)
	_, err = am.ActivateObject(
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

func (g *Genesis) activateSmartContracts(ctx context.Context, am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder, rootPubKey string) error {
	domain, domainDesc, err := g.activateRootDomain(ctx, am, cb)
	errMsg := "[ ActivateSmartContracts ]"
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = g.activateNodeDomain(ctx, domain, am, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = g.activateRootMember(ctx, domain, am, cb, rootPubKey)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	// TODO: this is not required since we refer by request id.
	err = g.updateRootDomain(ctx, am, cb, domainDesc)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = g.activateRootMemberWallet(ctx, domain, am, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

// Start creates types and RootDomain instance
func (g *Genesis) Start(ctx context.Context, c core.Components) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ Bootstrapper ] Starting Bootstrap ...")

	rootDomainRef, err := getRootDomainRef(ctx, c)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't get ref of rootDomain")
	}
	if rootDomainRef != nil {
		g.rootDomainRef = rootDomainRef

		rootMemberRef, err := getRootMemberRef(ctx, c, *g.rootDomainRef)
		if err != nil {
			return errors.Wrap(err, "[ Bootstrapper ] couldn't get ref of rootMember")
		}

		g.rootMemberRef = rootMemberRef
		inslog.Info("[ Bootstrapper ] RootDomain was found in ledger. Don't do bootstrap")
		return nil
	}

	_, rootPubKey, err := getKeysFromFile(ctx, g.config.RootKeysFile)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't get root keys")
	}

	isLightExecutor, err := isLightExecutor(ctx, c)
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

	am := c.Ledger.GetArtifactManager()
	cb := goplugintestutils.NewContractBuilder(am, insgocc)
	g.prototypeRefs = cb.Prototypes
	defer cb.Clean()

	err = buildSmartContracts(ctx, cb)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't build contracts")
	}

	err = g.activateSmartContracts(ctx, am, cb, rootPubKey)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ]")
	}

	if g.isGenesis {
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

			nodeID, err := am.RegisterRequest(ctx, &message.GenesisRequest{Name: "noderecord_" + strconv.Itoa(i)})
			if err != nil {
				return errors.Wrap(err, "")
			}
			contract := core.NewRecordRef(*g.rootDomainRef.Record(), *nodeID)
			_, err = am.ActivateObject(
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

// Stop implements core.Component method
func (g *Genesis) Stop(ctx context.Context) error {
	return nil
}
