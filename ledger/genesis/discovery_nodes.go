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

	"github.com/insolar/insolar/application/contract/nodedomain"
	"github.com/insolar/insolar/application/contract/noderecord"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/genesisrefs"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/rootdomain"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/internal/ledger/artifact"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/pkg/errors"
)

// DiscoveryNodeManager should be created by NewDiscoveryNodeManager.
type DiscoveryNodeManager struct {
	artifactManager artifact.Manager
}

// NewDiscoveryNodeManager creates new DiscoveryNodeManager instance.
func NewDiscoveryNodeManager(
	am artifact.Manager,
) *DiscoveryNodeManager {
	return &DiscoveryNodeManager{
		artifactManager: am,
	}
}

// StoreDiscoveryNodes saves discovery nodes objects and saves discovery nodes index in node domain index.
// If node domain index not empty this method does nothing.
func (nm *DiscoveryNodeManager) StoreDiscoveryNodes(ctx context.Context, discoveryNodes []insolar.DiscoveryNodeRegister) error {
	if len(discoveryNodes) == 0 {
		return nil
	}

	nodeDomainDesc, err := nm.artifactManager.GetObject(ctx, genesisrefs.ContractNodeDomain)
	if err != nil {
		return errors.Wrap(err, "failed to get node domain contract")
	}

	var ndObj nodedomain.NodeDomain
	insolar.MustDeserialize(nodeDomainDesc.Memory(), &ndObj)
	inslogger.FromContext(ctx).Debugf("get index on the node domain contract: %v", ndObj.NodeIndexPK)

	if len(ndObj.NodeIndexPK) != 0 {
		inslogger.FromContext(ctx).Info("discovery nodes already saved in the node domain index.")
		return nil
	}

	nodesInfo := make([]nodeInfo, 0, len(discoveryNodes))
	for _, n := range discoveryNodes {
		nodesInfo = append(nodesInfo, nodeInfo{
			role: insolar.GetStaticRoleFromString(n.Role),
			key:  platformpolicy.MustNormalizePublicKey([]byte(n.PublicKey)),
		})
	}
	return nm.updateDiscoveryData(ctx, nodesInfo)
}

// nodeInfo carries data for node objects required by DiscoveryNodeManager methods.
type nodeInfo struct {
	role insolar.StaticRole
	key  string
}

func (nm *DiscoveryNodeManager) updateDiscoveryData(
	ctx context.Context,
	nodes []nodeInfo,
) error {
	indexMap, err := nm.addDiscoveryNodes(ctx, nodes)
	if err != nil {
		return errors.Wrap(err, "failed to add discovery nodes")
	}

	err = nm.updateNodeDomainIndex(ctx, indexMap)
	return errors.Wrap(err, "failed to update node domain index")
}

// addDiscoveryNodes adds discovery nodeInfo's objects on ledger, returns index to store on nodeInfo domain.
func (nm *DiscoveryNodeManager) addDiscoveryNodes(
	ctx context.Context,
	nodes []nodeInfo,
) (map[string]string, error) {
	indexMap := map[string]string{}
	for _, n := range nodes {
		nodeState := &noderecord.NodeRecord{
			Record: noderecord.RecordInfo{
				PublicKey: n.key,
				Role:      n.role,
			},
		}

		contract, err := nm.activateNodeRecord(ctx, nodeState)
		if err != nil {
			return nil, errors.Wrap(err, "[ activateDiscoveryNodes ] Couldn't activateNodeRecord nodeInfo instance")
		}

		indexMap[n.key] = contract.String()
	}
	return indexMap, nil
}

func (nm *DiscoveryNodeManager) activateNodeRecord(
	ctx context.Context,
	node *noderecord.NodeRecord,
) (*insolar.Reference, error) {
	nodeData, err := insolar.Serialize(node)
	if err != nil {
		return nil, errors.Wrap(err, "failed to serialize node record data")
	}

	nodeID, err := nm.artifactManager.RegisterRequest(
		ctx,
		record.Request{
			CallType: record.CTGenesis,
			Method:   node.Record.PublicKey,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register request for node record")
	}

	contract := insolar.NewReference(*nodeID)
	_, err = nm.artifactManager.ActivateObject(
		ctx,
		insolar.Reference{},
		*contract,
		genesisrefs.ContractNodeDomain,
		rootdomain.GenesisRef(insolar.GenesisNameNodeRecord+"_proto"),
		false,
		nodeData,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to activate object of node record")
	}

	_, err = nm.artifactManager.RegisterResult(ctx, genesisrefs.ContractRootDomain, *contract, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't register result for new node object")
	}
	return contract, nil
}

func (nm *DiscoveryNodeManager) updateNodeDomainIndex(
	ctx context.Context,
	indexMap map[string]string,
) error {
	nodeDomainDesc, err := nm.artifactManager.GetObject(ctx, genesisrefs.ContractNodeDomain)
	if err != nil {
		return err
	}

	updateData, err := insolar.Serialize(
		&nodedomain.NodeDomain{
			NodeIndexPK: indexMap,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to serialize index for node domain")
	}

	_, err = nm.artifactManager.UpdateObject(
		ctx,
		genesisrefs.ContractRootDomain,
		genesisrefs.ContractNodeDomain,
		nodeDomainDesc,
		updateData,
	)
	return errors.Wrap(err, "failed to update node domain")
}
