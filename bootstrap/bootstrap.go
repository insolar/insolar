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

package bootstrap

import (
	"encoding/json"

	"github.com/insolar/insolar/application/contract/member"
	"github.com/insolar/insolar/application/contract/nodedomain"
	"github.com/insolar/insolar/application/contract/rootdomain"
	"github.com/insolar/insolar/application/contract/wallet"
	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/inscontext"
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

// Bootstrapper is a component for precreation core contracts types and RootDomain instance
type Bootstrapper struct {
	rootDomainRef *core.RecordRef
	nodeDomainRef *core.RecordRef
	rootMemberRef *core.RecordRef
	rootKeysFile  string
	rootPubKey    string
	rootBalance   uint
	classRefs     map[string]*core.RecordRef
}

// Info returns json with references for info api endpoint
func (b *Bootstrapper) Info() ([]byte, error) {
	classes := map[string]string{}
	for class, ref := range b.classRefs {
		classes[class] = ref.String()
	}
	return json.MarshalIndent(map[string]interface{}{
		"root_domain": b.rootDomainRef.String(),
		"root_member": b.rootMemberRef.String(),
		"classes":     classes,
	}, "", "   ")
}

// GetRootDomainRef returns reference to RootDomain instance
func (b *Bootstrapper) GetRootDomainRef() *core.RecordRef {
	return b.rootDomainRef
}

// NewBootstrapper creates new Bootstrapper
func NewBootstrapper(cfg configuration.Bootstrap) (*Bootstrapper, error) {
	bootstrapper := &Bootstrapper{}
	bootstrapper.rootKeysFile = cfg.RootKeys
	bootstrapper.rootBalance = cfg.RootBalance
	bootstrapper.rootDomainRef = &core.RecordRef{}
	return bootstrapper, nil
}

func buildSmartContracts(cb *goplugintestutils.ContractsBuilder) error {
	log.Info("[ buildSmartContracts ] building contracts:", contractNames)
	contracts, err := getContractsMap()
	if err != nil {
		return errors.Wrap(err, "[ buildSmartContracts ] couldn't build contracts")
	}

	log.Info("[ buildSmartContracts ] Start building contracts ...")
	err = cb.Build(contracts)
	if err != nil {
		return errors.Wrap(err, "[ buildSmartContracts ] couldn't build contracts")
	}
	log.Info("[ buildSmartContracts ] Stop building contracts ...")

	return nil
}

func (b *Bootstrapper) activateRootDomain(
	am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder,
) (*core.RecordID, core.ObjectDescriptor, error) {
	rd, err := rootdomain.NewRootDomain()
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	instanceData, err := serializeInstance(rd)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	ctx := inscontext.TODO()
	contractID, err := am.RegisterRequest(ctx, &message.BootstrapRequest{Name: "RootDomain"})
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	contract := core.ComposeRecordRef(*contractID, *contractID)
	desc, err := am.ActivateObject(
		ctx,
		core.RecordRef{}, contract,
		*cb.Classes[rootDomain],
		*am.GenesisRef(),
		false,
		instanceData,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	b.rootDomainRef = &contract

	return contractID, desc, nil
}

func (b *Bootstrapper) activateNodeDomain(
	domain *core.RecordID, am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder,
) error {
	nd, err := nodedomain.NewNodeDomain()
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ]")
	}

	instanceData, err := serializeInstance(nd)
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ]")
	}

	ctx := inscontext.TODO()
	contractID, err := am.RegisterRequest(ctx, &message.BootstrapRequest{Name: "NodeDomain"})
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}
	contract := core.ComposeRecordRef(*domain, *contractID)
	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{}, contract,
		*cb.Classes[nodeDomain],
		*b.rootDomainRef,
		false,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}

	b.nodeDomainRef = &contract

	return nil
}

func (b *Bootstrapper) activateRootMember(
	domain *core.RecordID, am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder,
) error {
	m, err := member.New("RootMember", b.rootPubKey)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ]")
	}

	instanceData, err := serializeInstance(m)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ]")
	}

	ctx := inscontext.TODO()
	contractID, err := am.RegisterRequest(ctx, &message.BootstrapRequest{Name: "RootMember"})
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	contract := core.ComposeRecordRef(*domain, *contractID)
	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{}, contract,
		*cb.Classes[memberContract],
		*b.rootDomainRef,
		false,
		instanceData,
	)

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	b.rootMemberRef = &contract
	return nil
}

// TODO: this is not required since we refer by request id.
func (b *Bootstrapper) updateRootDomain(
	am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder, domainDesc core.ObjectDescriptor,
) error {
	ctx := inscontext.TODO()
	updateData, err := serializeInstance(&rootdomain.RootDomain{RootMember: *b.rootMemberRef, NodeDomainRef: *b.nodeDomainRef})
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

func (b *Bootstrapper) activateRootMemberWallet(
	domain *core.RecordID, am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder,
) error {
	w, err := wallet.New(b.rootBalance)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	instanceData, err := serializeInstance(w)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	ctx := inscontext.TODO()
	contractID, err := am.RegisterRequest(ctx, &message.BootstrapRequest{Name: "RootWallet"})
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}
	contract := core.ComposeRecordRef(*domain, *contractID)
	_, err = am.ActivateObject(
		ctx,
		core.RecordRef{}, contract,
		*cb.Classes[walletContract],
		*b.rootMemberRef,
		true,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}

	return nil
}

func (b *Bootstrapper) activateSmartContracts(am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder) error {
	domain, domainDesc, err := b.activateRootDomain(am, cb)
	errMsg := "[ ActivateSmartContracts ]"
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = b.activateNodeDomain(domain, am, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = b.activateRootMember(domain, am, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	// TODO: this is not required since we refer by request id.
	err = b.updateRootDomain(am, cb, domainDesc)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = b.activateRootMemberWallet(domain, am, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

// Start creates types and RootDomain instance
func (b *Bootstrapper) Start(ctx core.Context, c core.Components) error {
	log.Info("[ Bootstrapper ] Starting Bootstrap ...")

	rootDomainRef, err := getRootDomainRef(c)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't get ref of rootDomain")
	}
	if rootDomainRef != nil {
		b.rootDomainRef = rootDomainRef
		log.Info("[ Bootstrapper ] RootDomain was found in ledger. Don't do bootstrap")
		return nil
	}

	b.rootPubKey, err = getRootMemberPubKey(b.rootKeysFile)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't get root member keys")
	}

	isLightExecutor, err := isLightExecutor(c)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't check if node is light executor")
	}
	if !isLightExecutor {
		log.Info("[ Bootstrapper ] Node is not light executor. Don't do bootstrap")
		return nil
	}

	_, insgocc, err := goplugintestutils.Build()
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't build insgocc")
	}

	am := c.Ledger.GetArtifactManager()
	cb := goplugintestutils.NewContractBuilder(am, insgocc)
	b.classRefs = cb.Classes
	defer cb.Clean()

	err = buildSmartContracts(cb)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ] couldn't build contracts")
	}

	err = b.activateSmartContracts(am, cb)
	if err != nil {
		return errors.Wrap(err, "[ Bootstrapper ]")
	}

	return nil
}

// Stop implements core.Component method
func (b *Bootstrapper) Stop(ctx core.Context) error {
	return nil
}
