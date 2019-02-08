/*
 *    Copyright 2019 Insolar
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
	"io/ioutil"
	"os"
	"testing"

	"github.com/insolar/insolar/application/contract/noderecord"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/testutils"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
)

func mockArtifactManager(t *testing.T) *testutils.ArtifactManagerMock {
	amMock := testutils.NewArtifactManagerMock(t)
	amMock.RegisterRequestFunc = func(p context.Context, p1 core.RecordRef, p2 core.Parcel) (r *core.RecordID, r1 error) {
		id := testutils.RandomID()
		return &id, nil
	}
	amMock.ActivateObjectFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) (r core.ObjectDescriptor, r1 error) {
		return testutils.NewObjectDescriptorMock(t), nil
	}
	amMock.RegisterResultFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) (r *core.RecordID, r1 error) {
		id := testutils.RandomID()
		return &id, nil
	}
	return amMock
}

func mockArtifactManagerWithRegisterRequestError(t *testing.T) *testutils.ArtifactManagerMock {
	amMock := testutils.NewArtifactManagerMock(t)
	amMock.RegisterRequestFunc = func(p context.Context, p1 core.RecordRef, p2 core.Parcel) (r *core.RecordID, r1 error) {
		return nil, errors.New("test reasons")
	}
	return amMock
}

func TestCreateKeys(t *testing.T) {
	g := Genesis{}
	ctx := inslogger.TestContext(t)
	path := "gentestdata"
	amount := 5
	defer os.RemoveAll(path)

	err := g.createKeys(ctx, path, amount)
	require.Nil(t, err)

	files, _ := ioutil.ReadDir(path)
	require.Equal(t, amount, len(files))
}

func TestUploadKeys_DontReuse(t *testing.T) {
	g := Genesis{
		config: &Config{
			ReuseKeys: false,
		},
	}
	ctx := inslogger.TestContext(t)
	path := "gentestdata"
	amount := 5
	defer os.RemoveAll(path)

	info, err := g.uploadKeys(ctx, path, amount)
	require.Nil(t, err)

	require.Equal(t, amount, len(info))
}

func TestUploadKeys_Reuse(t *testing.T) {
	g := Genesis{
		config: &Config{
			ReuseKeys: true,
		},
	}
	ctx := inslogger.TestContext(t)
	path := "gentestdata"
	amount := 5
	err := g.createKeys(ctx, path, amount)
	require.Nil(t, err)

	info, err := g.uploadKeys(ctx, path, amount)
	require.Nil(t, err)

	require.Equal(t, amount, len(info))
}

func TestUploadKeys_Reuse_WrongAmount(t *testing.T) {
	g := Genesis{
		config: &Config{
			ReuseKeys: true,
		},
	}
	ctx := inslogger.TestContext(t)
	path := "gentestdata"
	amount := 5
	err := g.createKeys(ctx, path, amount+5)
	defer os.RemoveAll(path)
	require.Nil(t, err)

	_, err = g.uploadKeys(ctx, path, amount)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "[ uploadKeys ] amount of nodes != amount of files in directory")
}

func TestUploadKeys_Reuse_DirNotExist(t *testing.T) {
	g := Genesis{
		config: &Config{
			ReuseKeys: true,
		},
	}
	ctx := inslogger.TestContext(t)
	path := "gentestdata"
	amount := 5

	_, err := g.uploadKeys(ctx, path, amount)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "[ uploadKeys ] dir is not exist")
}

func TestActivateNodeRecord_RegisterRequest_Err(t *testing.T) {
	am := mockArtifactManagerWithRegisterRequestError(t)
	ref := testutils.RandomRef()
	g := Genesis{
		config: &Config{
			ReuseKeys: true,
		},
		ArtifactManager: am,
		rootDomainRef:   &ref,
	}
	cb := NewContractBuilder(g.ArtifactManager)
	ctx := inslogger.TestContext(t)
	publicKey := "fancy_public_key"
	name := "fancy_name"
	record := &noderecord.NodeRecord{
		Record: noderecord.RecordInfo{
			PublicKey: publicKey,
			Role:      core.StaticRoleVirtual,
		},
	}

	_, err := g.activateNodeRecord(ctx, cb, record, name)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "[ activateNodes ] Couldn't register request to artifact manager: test reasons")
}

func TestActivateNodeRecord_Activate_Err(t *testing.T) {
	am := testutils.NewArtifactManagerMock(t)
	am.RegisterRequestFunc = func(p context.Context, p1 core.RecordRef, p2 core.Parcel) (r *core.RecordID, r1 error) {
		id := testutils.RandomID()
		return &id, nil
	}
	am.ActivateObjectFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) (r core.ObjectDescriptor, r1 error) {
		return nil, errors.New("test reasons")
	}

	ref := testutils.RandomRef()
	g := Genesis{
		config: &Config{
			ReuseKeys: true,
		},
		ArtifactManager: am,
		rootDomainRef:   &ref,
		nodeDomainRef:   &ref,
	}
	cb := NewContractBuilder(g.ArtifactManager)
	cb.Prototypes[nodeRecord] = &ref
	ctx := inslogger.TestContext(t)
	publicKey := "fancy_public_key"
	name := "fancy_name"
	record := &noderecord.NodeRecord{
		Record: noderecord.RecordInfo{
			PublicKey: publicKey,
			Role:      core.StaticRoleVirtual,
		},
	}

	_, err := g.activateNodeRecord(ctx, cb, record, name)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "[ activateNodes ] Could'n activateNodeRecord node object: test reasons")
}

func TestActivateNodeRecord_RegisterResult_Err(t *testing.T) {
	am := testutils.NewArtifactManagerMock(t)
	am.RegisterRequestFunc = func(p context.Context, p1 core.RecordRef, p2 core.Parcel) (r *core.RecordID, r1 error) {
		id := testutils.RandomID()
		return &id, nil
	}
	am.ActivateObjectFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 core.RecordRef, p4 core.RecordRef, p5 bool, p6 []byte) (r core.ObjectDescriptor, r1 error) {
		return testutils.NewObjectDescriptorMock(t), nil
	}
	am.RegisterResultFunc = func(p context.Context, p1 core.RecordRef, p2 core.RecordRef, p3 []byte) (r *core.RecordID, r1 error) {
		return nil, errors.New("test reasons")
	}

	ref := testutils.RandomRef()
	g := Genesis{
		config: &Config{
			ReuseKeys: true,
		},
		ArtifactManager: am,
		rootDomainRef:   &ref,
		nodeDomainRef:   &ref,
	}
	cb := NewContractBuilder(g.ArtifactManager)
	cb.Prototypes[nodeRecord] = &ref
	ctx := inslogger.TestContext(t)
	publicKey := "fancy_public_key"
	name := "fancy_name"
	record := &noderecord.NodeRecord{
		Record: noderecord.RecordInfo{
			PublicKey: publicKey,
			Role:      core.StaticRoleVirtual,
		},
	}

	_, err := g.activateNodeRecord(ctx, cb, record, name)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "[ activateNodes ] Couldn't register result to artifact manager: test reasons")
}

func TestActivateNodeRecord(t *testing.T) {
	am := mockArtifactManager(t)
	ref := testutils.RandomRef()
	g := Genesis{
		config: &Config{
			ReuseKeys: true,
		},
		ArtifactManager: am,
		rootDomainRef:   &ref,
		nodeDomainRef:   &ref,
	}
	cb := NewContractBuilder(g.ArtifactManager)
	cb.Prototypes[nodeRecord] = &ref
	ctx := inslogger.TestContext(t)
	publicKey := "fancy_public_key"
	name := "fancy_name"
	record := &noderecord.NodeRecord{
		Record: noderecord.RecordInfo{
			PublicKey: publicKey,
			Role:      core.StaticRoleVirtual,
		},
	}

	contract, err := g.activateNodeRecord(ctx, cb, record, name)
	require.Nil(t, err)
	require.NotNil(t, contract)
}

func TestActivateNodes_Err(t *testing.T) {
	am := mockArtifactManagerWithRegisterRequestError(t)
	ref := testutils.RandomRef()
	g := Genesis{
		config: &Config{
			ReuseKeys: true,
		},
		ArtifactManager: am,
		rootDomainRef:   &ref,
		nodeDomainRef:   &ref,
	}
	cb := NewContractBuilder(g.ArtifactManager)
	cb.Prototypes[nodeRecord] = &ref
	ctx := inslogger.TestContext(t)

	var nodes []nodeInfo
	nodes = append(nodes,
		nodeInfo{
			publicKey: "test_pk_1",
		},
		nodeInfo{
			publicKey: "test_pk_2",
		},
	)

	_, err := g.activateNodes(ctx, cb, nodes)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "[ activateNodes ] Couldn't activateNodeRecord node instance:")
}

func TestActivateNodes(t *testing.T) {
	am := mockArtifactManager(t)
	ref := testutils.RandomRef()
	g := Genesis{
		config: &Config{
			ReuseKeys: true,
		},
		ArtifactManager: am,
		rootDomainRef:   &ref,
		nodeDomainRef:   &ref,
	}
	cb := NewContractBuilder(g.ArtifactManager)
	cb.Prototypes[nodeRecord] = &ref
	ctx := inslogger.TestContext(t)

	var nodes []nodeInfo
	nodes = append(nodes,
		nodeInfo{
			publicKey: "test_pk_1",
		},
		nodeInfo{
			publicKey: "test_pk_2",
		},
	)

	updatedNodes, err := g.activateNodes(ctx, cb, nodes)
	require.Nil(t, err)
	require.Len(t, updatedNodes, len(nodes))
	for i := 0; i < len(nodes); i++ {
		require.Equal(t, nodes[i].publicKey, updatedNodes[i].publicKey)
		require.NotNil(t, updatedNodes[i].ref)
	}
}

func TestActivateDiscoveryNodes_DiffLen(t *testing.T) {
	am := mockArtifactManager(t)
	ref := testutils.RandomRef()
	var discoveryNodes []Node
	discoveryNodes = append(discoveryNodes,
		Node{
			Role: "virtual",
		},
	)
	g := Genesis{
		config: &Config{
			ReuseKeys:      true,
			DiscoveryNodes: discoveryNodes,
		},
		ArtifactManager: am,
		rootDomainRef:   &ref,
		nodeDomainRef:   &ref,
	}
	cb := NewContractBuilder(g.ArtifactManager)
	cb.Prototypes[nodeRecord] = &ref
	ctx := inslogger.TestContext(t)

	var nodes []nodeInfo
	nodes = append(nodes,
		nodeInfo{
			publicKey: "test_pk_1",
		},
		nodeInfo{
			publicKey: "test_pk_2",
		},
	)

	_, err := g.activateDiscoveryNodes(ctx, cb, nodes)
	require.EqualError(t, err, "[ activateDiscoveryNodes ] len of nodesInfo param must be equal to len of DiscoveryNodes in genesis config")
}

func TestActivateDiscoveryNodes_Err(t *testing.T) {
	am := mockArtifactManagerWithRegisterRequestError(t)
	ref := testutils.RandomRef()
	var discoveryNodes []Node
	discoveryNodes = append(discoveryNodes,
		Node{
			Role: "virtual",
		},
		Node{
			Role: "light_material",
		},
	)
	g := Genesis{
		config: &Config{
			ReuseKeys:      true,
			DiscoveryNodes: discoveryNodes,
		},
		ArtifactManager: am,
		rootDomainRef:   &ref,
		nodeDomainRef:   &ref,
	}
	cb := NewContractBuilder(g.ArtifactManager)
	cb.Prototypes[nodeRecord] = &ref
	ctx := inslogger.TestContext(t)

	var nodes []nodeInfo
	nodes = append(nodes,
		nodeInfo{
			publicKey: "test_pk_1",
		},
		nodeInfo{
			publicKey: "test_pk_2",
		},
	)
	require.Len(t, nodes, len(discoveryNodes))

	_, err := g.activateDiscoveryNodes(ctx, cb, nodes)
	require.NotNil(t, err)
	require.Contains(t, err.Error(), "[ activateDiscoveryNodes ] Couldn't activateNodeRecord node instance:")
}

func TestActivateDiscoveryNodes(t *testing.T) {
	am := mockArtifactManager(t)
	ref := testutils.RandomRef()
	var discoveryNodes []Node
	discoveryNodes = append(discoveryNodes,
		Node{
			Role: "virtual",
		},
		Node{
			Role: "light_material",
		},
	)
	g := Genesis{
		config: &Config{
			ReuseKeys:      true,
			DiscoveryNodes: discoveryNodes,
		},
		ArtifactManager: am,
		rootDomainRef:   &ref,
		nodeDomainRef:   &ref,
	}
	cb := NewContractBuilder(g.ArtifactManager)
	cb.Prototypes[nodeRecord] = &ref
	ctx := inslogger.TestContext(t)

	var nodes []nodeInfo
	nodes = append(nodes,
		nodeInfo{
			publicKey: "test_pk_1",
		},
		nodeInfo{
			publicKey: "test_pk_2",
		},
	)
	require.Len(t, nodes, len(discoveryNodes))

	genesisNodes, err := g.activateDiscoveryNodes(ctx, cb, nodes)
	require.Nil(t, err)
	require.Len(t, genesisNodes, len(discoveryNodes))
	for i := 0; i < len(discoveryNodes); i++ {
		require.Equal(t, discoveryNodes[i].Role, genesisNodes[i].role)
		require.Equal(t, nodes[i].publicKey, genesisNodes[i].node.PublicKey)
		require.NotNil(t, genesisNodes[i].ref)
	}
}

func TestAddDiscoveryIndex(t *testing.T) {
	am := mockArtifactManager(t)
	ref := testutils.RandomRef()
	var discoveryNodes []Node
	discoveryNodes = append(discoveryNodes,
		Node{
			Role: "virtual",
		},
		Node{
			Role: "light_material",
		},
	)
	g := Genesis{
		config: &Config{
			ReuseKeys:      false,
			DiscoveryNodes: discoveryNodes,
		},
		ArtifactManager: am,
		rootDomainRef:   &ref,
		nodeDomainRef:   &ref,
	}
	cb := NewContractBuilder(g.ArtifactManager)
	cb.Prototypes[nodeRecord] = &ref
	ctx := inslogger.TestContext(t)

	indexMap := make(map[string]string)

	genesisNodes, resIndexMap, err := g.addDiscoveryIndex(ctx, cb, indexMap)
	require.Nil(t, err)
	require.Len(t, genesisNodes, len(discoveryNodes))
	require.Len(t, resIndexMap, len(discoveryNodes))
}
