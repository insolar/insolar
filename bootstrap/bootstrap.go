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
	"io/ioutil"
	"path/filepath"
	"runtime"

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
	"github.com/ugorji/go/codec"
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
	classes := map[string][]byte{}
	for class, ref := range b.classRefs {
		classes[class] = ref[:]
	}
	return json.MarshalIndent(map[string]interface{}{
		"root_domain": b.rootDomainRef[:],
		"root_member": b.rootMemberRef[:],
		"classes":     classes,
	}, "", "   ")
}

// GetRootDomainRef returns reference to RootDomain instance
func (b *Bootstrapper) GetRootDomainRef() *core.RecordRef {
	return b.rootDomainRef
}

// NewBootstrapper creates new Bootstrapper
func NewBootstrapper(cfg configuration.Configuration) (*Bootstrapper, error) {
	bootstrapper := &Bootstrapper{}
	bootstrapper.rootKeysFile = cfg.Bootstrap.RootKeys
	bootstrapper.rootBalance = cfg.Bootstrap.RootBalance
	bootstrapper.rootDomainRef = &core.RecordRef{}
	return bootstrapper, nil
}

var pathToContracts = "application/contract/"

func getAbsolutePath(relativePath string) (string, error) {
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.Wrap(nil, "[ getFullPath ] couldn't find info about current file")
	}
	rootDir := filepath.Dir(filepath.Dir(currentFile))
	return filepath.Join(rootDir, relativePath), nil
}

func getContractPath(name string) (string, error) {
	contractDir, err := getAbsolutePath(pathToContracts)
	if err != nil {
		return "", errors.Wrap(nil, "[ getContractPath ] couldn't get absolute path to contracts")
	}
	contractFile := name + ".go"
	return filepath.Join(contractDir, name, contractFile), nil
}

func getContractsMap() (map[string]string, error) {
	contracts := make(map[string]string)
	for _, name := range contractNames {
		contractPath, err := getContractPath(name)
		if err != nil {
			return nil, errors.Wrap(err, "[ contractsMap ] couldn't get path to contracts: ")
		}
		code, err := ioutil.ReadFile(filepath.Clean(contractPath))
		if err != nil {
			return nil, errors.Wrap(err, "[ contractsMap ] couldn't read contract: ")
		}
		contracts[name] = string(code)
	}
	return contracts, nil
}

func isLightExecutor(c core.Components) (bool, error) {
	am := c.Ledger.GetArtifactManager()
	jc := c.Ledger.GetJetCoordinator()
	pm := c.Ledger.GetPulseManager()
	currentPulse, err := pm.Current()
	if err != nil {
		return false, errors.Wrap(err, "[ isLightExecutor ] couldn't get current pulse")
	}

	network := c.Network
	nodeID := network.GetNodeID()

	isLightExecutor, err := jc.IsAuthorized(core.RoleLightExecutor, *am.GenesisRef(), currentPulse.PulseNumber, nodeID)
	if err != nil {
		return false, errors.Wrap(err, "[ isLightExecutor ] couldn't authorized node")
	}
	if !isLightExecutor {
		log.Info("[ isLightExecutor ] Is not light executor. Don't build contracts")
		return false, nil
	}
	return true, nil
}

func getRootDomainRef(c core.Components) (*core.RecordRef, error) {
	am := c.Ledger.GetArtifactManager()
	rootObj, err := am.GetObject(*am.GenesisRef(), nil)
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

func serializeInstance(contractInstance interface{}) ([]byte, error) {
	var instanceData []byte

	ch := new(codec.CborHandle)
	err := codec.NewEncoderBytes(&instanceData, ch).Encode(
		contractInstance,
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ serializeInstance ] Problem with CBORing")
	}

	return instanceData, nil
}

func (b *Bootstrapper) activateRootDomain(am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder) error {
	instanceData, err := serializeInstance(rootdomain.NewRootDomain())
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootDomain ]")
	}

	ctx := inscontext.TODO()
	contract, err := am.RegisterRequest(ctx, &message.BootstrapRequest{Name: "RootDomain"})
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	_, err = am.ActivateObject(
		core.RecordRef{}, *contract,
		*cb.Classes[rootDomain],
		*am.GenesisRef(),
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}
	b.rootDomainRef = contract

	return nil
}

func (b *Bootstrapper) activateNodeDomain(am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder) error {
	instanceData, err := serializeInstance(nodedomain.NewNodeDomain())
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ]")
	}

	ctx := inscontext.TODO()
	contract, err := am.RegisterRequest(ctx, &message.BootstrapRequest{Name: "NodeDomain"})
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}
	_, err = am.ActivateObject(
		core.RecordRef{}, *contract,
		*cb.Classes[nodeDomain],
		*b.rootDomainRef,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}

	b.nodeDomainRef = contract

	return nil
}

func (b *Bootstrapper) activateRootMember(am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder) error {

	instanceData, err := serializeInstance(member.New("RootMember", b.rootPubKey))
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ]")
	}

	ctx := inscontext.TODO()
	contract, err := am.RegisterRequest(ctx, &message.BootstrapRequest{Name: "RootMember"})
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	_, err = am.ActivateObject(
		core.RecordRef{}, *contract,
		*cb.Classes[memberContract],
		*b.rootDomainRef,
		instanceData,
	)

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root member instance")
	}
	b.rootMemberRef = contract
	return nil
}

func (b *Bootstrapper) updateRootDomain(am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder) error {
	updateData, err := serializeInstance(&rootdomain.RootDomain{RootMember: *b.rootMemberRef, NodeDomainRef: *b.nodeDomainRef})
	if err != nil {
		return errors.Wrap(err, "[ updateRootDomain ]")
	}
	_, err = am.UpdateObject(
		core.RecordRef{}, core.RecordRef{},
		*b.rootDomainRef, updateData,
	)
	if err != nil {
		return errors.Wrap(err, "[ updateRootDomain ]")
	}

	return nil
}

func (b *Bootstrapper) activateRootMemberWallet(am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder) error {
	instanceData, err := serializeInstance(wallet.New(b.rootBalance))
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	ctx := inscontext.TODO()
	contract, err := am.RegisterRequest(ctx, &message.BootstrapRequest{Name: "RootMember"})
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}
	_, err = am.ActivateObjectDelegate(
		core.RecordRef{}, *contract,
		*cb.Classes[walletContract],
		*b.rootMemberRef,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}

	return nil
}

func (b *Bootstrapper) activateSmartContracts(am core.ArtifactManager, cb *goplugintestutils.ContractsBuilder) error {
	err := b.activateRootDomain(am, cb)
	errMsg := "[ ActivateSmartContracts ]"
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = b.activateNodeDomain(am, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = b.activateRootMember(am, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = b.updateRootDomain(am, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}
	err = b.activateRootMemberWallet(am, cb)
	if err != nil {
		return errors.Wrap(err, errMsg)
	}

	return nil
}

func getRootMemberPubKey(file string) (string, error) {
	fileWithPath, err := getAbsolutePath(file)
	if err != nil {
		return "", errors.Wrap(err, "[ getRootMemberPubKey ] couldn't find absolute path for root keys")
	}
	data, err := ioutil.ReadFile(filepath.Clean(fileWithPath))
	if err != nil {
		return "", errors.New("couldn't read rootkeys file")
	}
	var keys map[string]string
	err = json.Unmarshal(data, &keys)
	if err != nil {
		return "", err
	}
	if keys["public_key"] == "" {
		return "", errors.New("empty root public key")
	}
	return keys["public_key"], nil
}

// Start creates types and RootDomain instance
func (b *Bootstrapper) Start(c core.Components) error {
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
func (b *Bootstrapper) Stop() error {
	return nil
}
