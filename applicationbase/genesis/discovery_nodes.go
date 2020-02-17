// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package genesis

import (
	"context"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/applicationbase/builtin/contract/nodedomain"
	"github.com/insolar/insolar/applicationbase/builtin/contract/noderecord"
	"github.com/insolar/insolar/applicationbase/genesisrefs"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/artifact"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/platformpolicy"
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
func (nm *DiscoveryNodeManager) StoreDiscoveryNodes(ctx context.Context, discoveryNodes []DiscoveryNodeRegister, parentDomain string) error {
	if len(discoveryNodes) == 0 {
		return nil
	}

	nodeDomainDesc, err := nm.artifactManager.GetObject(ctx, genesisrefs.ContractNodeDomain)
	if err != nil {
		return errors.Wrap(err, "failed to get node domain contract")
	}

	var ndObj nodedomain.NodeDomain
	insolar.MustDeserialize(nodeDomainDesc.Memory(), &ndObj)
	inslogger.FromContext(ctx).Debugf("get index on the node domain contract: %v", ndObj.NodeIndexPublicKey)

	if len(ndObj.NodeIndexPublicKey) != 0 {
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
	return nm.updateDiscoveryData(ctx, nodesInfo, parentDomain)
}

// nodeInfo carries data for node objects required by DiscoveryNodeManager methods.
type nodeInfo struct {
	role insolar.StaticRole
	key  string
}

func (nm *DiscoveryNodeManager) updateDiscoveryData(
	ctx context.Context,
	nodes []nodeInfo,
	parentDomain string,
) error {
	indexMap, err := nm.addDiscoveryNodes(ctx, nodes)
	if err != nil {
		return errors.Wrap(err, "failed to add discovery nodes")
	}

	err = nm.updateNodeDomainIndex(ctx, indexMap, parentDomain)
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

		canonicalKey, err := foundation.ExtractCanonicalPublicKey(n.key)
		if err != nil {
			return nil, errors.Wrapf(err, "extracting canonical pk failed, current value %v", n.key)
		}

		indexMap[canonicalKey] = contract.String()
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
		record.IncomingRequest{
			CallType: record.CTGenesis,
			Method:   node.Record.PublicKey,
		},
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to register request for node record")
	}

	contract := insolar.NewReference(*nodeID)
	err = nm.artifactManager.ActivateObject(
		ctx,
		*insolar.NewEmptyReference(),
		*contract,
		genesisrefs.ContractNodeDomain,
		genesisrefs.GenesisRef(genesisrefs.GenesisNameNodeRecord+"_proto"),
		nodeData,
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed to activate object of node record")
	}

	_, err = nm.artifactManager.RegisterResult(ctx, genesisrefs.ContractNodeDomain, *contract, nil)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't register result for new node object")
	}
	return contract, nil
}

func (nm *DiscoveryNodeManager) updateNodeDomainIndex(
	ctx context.Context,
	indexMap map[string]string,
	parentDomain string,
) error {
	nodeDomainDesc, err := nm.artifactManager.GetObject(ctx, genesisrefs.ContractNodeDomain)
	if err != nil {
		return err
	}

	indexStableMap := foundation.StableMap(indexMap)

	updateData, err := insolar.Serialize(
		&nodedomain.NodeDomain{
			NodeIndexPublicKey: indexStableMap,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to serialize index for node domain")
	}

	err = nm.artifactManager.UpdateObject(
		ctx,
		genesisrefs.GenesisRef(parentDomain),
		genesisrefs.ContractNodeDomain,
		nodeDomainDesc,
		updateData,
	)
	return errors.Wrap(err, "failed to update node domain")
}
