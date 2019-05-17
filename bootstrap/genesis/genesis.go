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
	walletcontract "github.com/insolar/insolar/application/contract/wallet"
	"github.com/insolar/insolar/bootstrap"
	"github.com/insolar/insolar/bootstrap/rootdomain"
	"github.com/insolar/insolar/certificate"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/artifact"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

var contractNames = []string{
	insolar.GenesisNameRootDomain,
	insolar.GenesisNameNodeDomain,
	insolar.GenesisNameNodeRecord,
	insolar.GenesisNameRootMember,
	insolar.GenesisNameWallet,
	insolar.GenesisNameAllowance,
}

type nodeInfo struct {
	privateKey crypto.PrivateKey
	publicKey  string
}

func (ni nodeInfo) reference() insolar.Reference {
	return rootdomain.GenesisRef(ni.publicKey)
}

// Generator is a component for generating RootDomain instance and genesis contracts.
type Generator struct {
	config          *Config
	artifactManager artifact.Manager

	keyOut string
}

// NewGenerator creates new Generator.
func NewGenerator(
	config *Config,
	am artifact.Manager,
	genesisKeyOut string,
) *Generator {
	return &Generator{
		config:          config,
		artifactManager: am,

		keyOut: genesisKeyOut,
	}
}

// Run generates genesis data.
func (g *Generator) Run(ctx context.Context) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ Genesis ] Starting  ...")

	inslog.Info("[ Genesis ] newContractBuilder ...")
	cb := newContractBuilder(g.artifactManager)
	defer cb.clean()

	inslog.Info("[ Genesis ] buildSmartContracts ...")
	prototypes, err := cb.buildPrototypes(ctx, contractNames)
	if err != nil {
		panic(errors.Wrap(err, "[ Genesis ] couldn't build contracts"))
	}

	inslog.Info("[ Genesis ] ReadKeysFile ...")
	pair, err := secrets.ReadKeysFile(g.config.RootKeysFile)
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

	err = g.activateSmartContracts(ctx, platformpolicy.MustPublicKeyToString(pair.Public), prototypes)
	if err != nil {
		panic(errors.Wrap(err, "[ Genesis ] could't activate smart contracts"))
	}

	err = g.saveDiscoveryNodes(ctx, *cb.prototypes[insolar.GenesisNameNodeRecord], discoveryNodes)
	if err != nil {
		panic(errors.Wrap(err, "[ Genesis ] could't save discovery nodes data"))
	}

	inslog.Info("[ Genesis ] makeCertificates ...")
	err = g.makeCertificates(ctx, discoveryNodes)
	if err != nil {
		return errors.Wrap(err, "[ Genesis ] Couldn't generate discovery certificates")
	}

	inslog.Info("[ Genesis ] Finished.")
	return nil
}

func (g *Generator) saveDiscoveryNodes(ctx context.Context, nodeRecordProto insolar.Reference, discoveryNodes []nodeInfo) error {
	err := g.activateDiscoveryNodes(ctx, nodeRecordProto, discoveryNodes)
	if err != nil {
		return errors.Wrap(err, "failed on adding discovery index")
	}

	err = g.updateNodeDomainIndex(ctx, discoveryNodes)
	return errors.Wrap(err, "failed update node domain index")
}

func (g *Generator) activateRootDomain(
	ctx context.Context,
	rootDomainProto insolar.Reference,
) error {
	inslog := inslogger.FromContext(ctx)
	inslog.Info("[ activateRootDomain ] createKeysInDir ...")

	data, err := insolar.Serialize(&rootdomaincontract.RootDomain{
		RootMember:    bootstrap.ContractRootMember,
		NodeDomainRef: bootstrap.ContractNodeDomain,
	})
	if err != nil {
		return errors.Wrap(err, "[ updateRootDomain ]")
	}

	regID, err := g.artifactManager.RegisterRequest(
		ctx,
		insolar.GenesisRecord.Ref(),
		&message.Parcel{
			Msg: &message.GenesisRequest{
				Name: insolar.GenesisNameRootDomain,
			},
		},
	)
	if err != nil {
		panic(errors.Wrap(err, "[ Genesis ] Couldn't create rootdomain instance"))
	}
	inslog.Debugf("[ ActivateRootDomain ] register id=%v", regID.String())

	rootDomainRef := rootdomain.RootDomain.Ref()
	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		rootDomainRef,
		insolar.GenesisRecord.Ref(),
		rootDomainProto,
		false,
		data,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}

	_, err = g.artifactManager.RegisterResult(ctx, insolar.GenesisRecord.Ref(), rootDomainRef, nil)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootDomain ] Couldn't create rootdomain instance")
	}

	return nil
}

func (g *Generator) activateNodeDomain(
	ctx context.Context, nodeDomainProto insolar.Reference,
) error {
	nd, _ := nodedomain.NewNodeDomain()

	instanceData, err := insolar.Serialize(nd)
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] node domain serialization")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		bootstrap.ContractRootDomain,
		&message.Parcel{
			Msg: &message.GenesisRequest{Name: insolar.GenesisNameNodeDomain},
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}
	contract := insolar.NewReference(rootdomain.RootDomain.ID(), *contractID)

	inslogger.FromContext(ctx).Debugf("[activateNodeDomain] Ref: %v", contract)

	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		bootstrap.ContractRootDomain,
		nodeDomainProto,
		false,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}
	_, err = g.artifactManager.RegisterResult(ctx, bootstrap.ContractRootDomain, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ ActivateNodeDomain ] couldn't create nodedomain instance")
	}

	inslogger.FromContext(ctx).Debugf("%v contract ref=%v", bootstrap.ContractNodeDomain, contract)

	return nil
}

func (g *Generator) activateRootMember(
	ctx context.Context,
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
		bootstrap.ContractRootDomain,
		&message.Parcel{
			Msg: &message.GenesisRequest{Name: insolar.GenesisNameRootMember},
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root rootMember instance")
	}
	contract := insolar.NewReference(rootdomain.RootDomain.ID(), *contractID)
	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		bootstrap.ContractRootDomain,
		memberContractProto,
		false,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root rootMember instance")
	}
	_, err = g.artifactManager.RegisterResult(ctx, bootstrap.ContractRootDomain, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootMember ] couldn't create root rootMember instance")
	}
	return nil
}

func (g *Generator) activateRootMemberWallet(
	ctx context.Context, walletContractProto insolar.Reference,
) error {

	w, err := walletcontract.New(g.config.RootBalance)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	instanceData, err := insolar.Serialize(w)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ]")
	}

	contractID, err := g.artifactManager.RegisterRequest(
		ctx,
		bootstrap.ContractRootDomain,
		&message.Parcel{
			Msg: &message.GenesisRequest{Name: "RootWallet"},
		},
	)

	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}
	contract := insolar.NewReference(rootdomain.RootDomain.ID(), *contractID)
	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		bootstrap.ContractRootMember,
		walletContractProto,
		true,
		instanceData,
	)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}
	_, err = g.artifactManager.RegisterResult(ctx, bootstrap.ContractRootDomain, *contract, nil)
	if err != nil {
		return errors.Wrap(err, "[ ActivateRootWallet ] couldn't create root wallet")
	}

	return nil
}

func (g *Generator) activateSmartContracts(
	ctx context.Context,
	rootPubKey string,
	prototypes prototypes,
) error {
	var err error

	err = g.activateRootDomain(ctx, *prototypes[insolar.GenesisNameRootDomain])
	// errMsg := "[ ActivateSmartContracts ]"
	if err != nil {
		return errors.Wrap(err, "failed to store root domain contract")
	}

	err = g.activateNodeDomain(ctx, *prototypes[insolar.GenesisNameNodeDomain])
	if err != nil {
		return errors.Wrap(err, "failed to store node domain contract")
	}

	err = g.activateRootMember(ctx, rootPubKey, *prototypes[insolar.GenesisNameRootMember])
	if err != nil {
		return errors.Wrap(err, "failed to store root GenesisNameRootMember contract")
	}

	err = g.activateRootMemberWallet(ctx, *prototypes[insolar.GenesisNameWallet])
	if err != nil {
		return errors.Wrap(err, "failed to store root rootMemberWallet contract")
	}

	return nil
}

// activateDiscoveryNodes activates discoverynoderecord_{N} objects.
//
// It returns list of genesisNode structures (for node domain save and certificates generation at the end of genesis).
func (g *Generator) activateDiscoveryNodes(
	ctx context.Context,
	nodeRecordProto insolar.Reference,
	discoveryNodes []nodeInfo,
) error {
	if len(discoveryNodes) != len(g.config.DiscoveryNodes) {
		return errors.Errorf(
			"[ activateDiscoveryNodes ] len of discoveryNodes (%v) must be equal to len of DiscoveryNodes (%v) in genesis config",
			len(discoveryNodes), len(g.config.DiscoveryNodes))
	}

	for i, discoverNode := range g.config.DiscoveryNodes {
		nodePubKey := discoveryNodes[i].publicKey

		nodeState := &noderecord.NodeRecord{
			Record: noderecord.RecordInfo{
				PublicKey: nodePubKey,
				Role:      insolar.GetStaticRoleFromString(discoverNode.Role),
			},
		}

		_, err := g.activateNodeRecord(ctx, nodeState, discoveryNodes[i], nodeRecordProto)
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
		bootstrap.ContractRootDomain,
		&message.Parcel{
			Msg: &message.GenesisRequest{Name: node.publicKey},
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ activateNodeRecord ] Couldn't register request")
	}
	contract := insolar.NewReference(*bootstrap.ContractRootDomain.Record(), *nodeID)
	_, err = g.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		bootstrap.ContractNodeDomain,
		nodeRecordProto,
		false,
		nodeData,
	)
	if err != nil {
		return nil, errors.Wrap(err, "[ activateNodeRecord ] Could'n activateNodeRecord node object")
	}
	_, err = g.artifactManager.RegisterResult(ctx, bootstrap.ContractRootDomain, *contract, nil)
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

			RootDomainReference: bootstrap.ContractRootDomain.String(),
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
	nodeDomainDesc, err := g.artifactManager.GetObject(ctx, bootstrap.ContractNodeDomain)
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
		bootstrap.ContractRootDomain,
		bootstrap.ContractNodeDomain,
		nodeDomainDesc,
		updateData,
	)
	if err != nil {
		return errors.Wrap(err, "[ updateNodeDomainIndex ]  Couldn't update NodeDomain")
	}

	return nil
}
