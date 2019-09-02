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
	"path/filepath"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/genesisrefs"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/insolar/secrets"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/artifact"
	"github.com/insolar/insolar/logicrunner/builtin/contract/nodedomain"
	"github.com/insolar/insolar/logicrunner/builtin/contract/noderecord"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/pulse"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestData_EmptyDomainData(t *testing.T) {
	ctx := inslogger.TestContext(t)

	am := artifact.NewManagerMock(t)
	// should no any calls on empty discovery nodes list
	defer am.MinimockFinish()

	dnm := NewDiscoveryNodeManager(am)
	err := dnm.StoreDiscoveryNodes(ctx, nil)
	require.NoError(t, err, "StoreDiscoveryNodes failed")
}

func TestData_WriteNodeDomainData(t *testing.T) {
	ctx := inslogger.TestContext(t)

	am := initArtifactManager(t)
	dCerts := NewDiscoveryNodeManager(am)
	nodes, err := publicKeysFromDir("testdata/keys", "testdata/keys_meta.json")
	require.NoError(t, err, "can't load keys")

	var networkNodes []insolar.DiscoveryNodeRegister
	for _, n := range nodes {
		networkNodes = append(networkNodes, insolar.DiscoveryNodeRegister{
			Role:      n.role.String(),
			PublicKey: platformpolicy.MustPublicKeyToString(n.key),
		})
	}
	err = dCerts.StoreDiscoveryNodes(ctx, networkNodes)
	require.NoError(t, err, "StoreDiscoveryNodes failed")

	objDesc, err := am.GetObject(ctx, genesisrefs.ContractNodeDomain)
	if err != nil {
		panic(err)
	}
	var ndMemory nodedomain.NodeDomain
	insolar.Deserialize(objDesc.Memory(), &ndMemory)

	expectIndexMap := make(foundation.StableMap)

	for _, n := range nodes {
		pKey := platformpolicy.MustPublicKeyToString(n.key)
		ref := genesisrefs.GenesisRef(pKey)
		expectIndexMap[pKey] = ref.String()

		nodeObjDesc, err := am.GetObject(ctx, ref)
		require.NoError(t, err, "nodeInfo object not found for public key: "+pKey)
		var nodeRec noderecord.NodeRecord
		insolar.Deserialize(nodeObjDesc.Memory(), &nodeRec)
		assert.Equal(t, pKey, nodeRec.Record.PublicKey, "public key is the same")
		assert.Equal(t, n.role, nodeRec.Record.Role, "role is the same")
	}

	assert.Equal(t, expectIndexMap, ndMemory.NodeIndexPublicKey, "NodeDomain memory contains expected map")
}

func initArtifactManager(t *testing.T) artifact.Manager {
	amMock := artifact.NewManagerMock(t)
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	indexMap := make(foundation.StableMap)
	activatedMemory := map[insolar.Reference][]byte{}

	amMock.GetObjectMock.Set(func(_ context.Context, ref insolar.Reference) (artifact.ObjectDescriptor, error) {
		descMock := artifact.NewObjectDescriptorMock(t)
		if ref == genesisrefs.ContractNodeDomain {
			descMock.MemoryMock.Set(func() []byte {
				return insolar.MustSerialize(&nodedomain.NodeDomain{
					NodeIndexPublicKey: indexMap,
				})
			})
		} else {
			descMock.MemoryMock.Set(func() []byte {
				b := activatedMemory[ref]
				return b
			})
		}
		return descMock, nil
	})

	amMock.RegisterRequestMock.Set(func(
		_ context.Context,
		req record.IncomingRequest,
	) (*insolar.ID, error) {
		virtRec := record.Wrap(&req)
		hash := record.HashVirtual(pcs.ReferenceHasher(), virtRec)
		return insolar.NewID(pulse.MinTimePulse, hash), nil
	})

	amMock.UpdateObjectMock.Set(func(
		_ context.Context,
		domain insolar.Reference,
		request insolar.Reference,
		obj artifact.ObjectDescriptor,
		memory []byte,
	) error {
		if domain != genesisrefs.ContractRootDomain {
			return errors.Errorf("domain should be the contract root domain ref")
		}
		if request != genesisrefs.ContractNodeDomain {
			return errors.Errorf("request should be the contract node domain ref")
		}
		var rec nodedomain.NodeDomain
		insolar.MustDeserialize(memory, &rec)
		indexMap = rec.NodeIndexPublicKey
		return nil
	})
	amMock.ActivateObjectMock.Set(func(
		_ context.Context,
		domain, obj, parent, prototype insolar.Reference,
		memory []byte,
	) error {
		activatedMemory[obj] = memory
		return nil
	})

	amMock.RegisterResultMock.Return(nil, nil)
	return amMock
}

type nodeInfoRaw struct {
	role insolar.StaticRole
	key  crypto.PublicKey
}

func publicKeysFromDir(dir string, keysMetaFile string) ([]nodeInfoRaw, error) {
	b, err := ioutil.ReadFile(keysMetaFile)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read file with keys meta info %v", keysMetaFile)
	}
	type info struct {
		Role string
		File string
	}
	var metaInfo []info
	if err := json.Unmarshal(b, &metaInfo); err != nil {
		return nil, errors.Wrapf(err, "can't decode json from file %v", keysMetaFile)
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, errors.Wrapf(err, "can't read dir %v", dir)
	}

	fileSet := map[string]string{}
	// fmt.Println("files:")
	for _, f := range files {
		// fmt.Println(" ", f.Name())
		fileSet[f.Name()] = filepath.Join(dir, f.Name())
	}
	// fmt.Printf("metaInfo: %#v\n", metaInfo)

	nodes := make([]nodeInfoRaw, 0, len(metaInfo))
	for _, meta := range metaInfo {
		f, ok := fileSet[meta.File]
		if !ok {
			return nil, errors.Errorf("not found file %v in directory %v", meta.File, dir)
		}
		pair, err := secrets.ReadKeysFile(f)
		if err != nil {
			return nil, errors.Wrapf(err, "can't read dir %v", dir)
		}
		nodes = append(nodes, nodeInfoRaw{
			role: insolar.GetStaticRoleFromString(meta.Role),
			key:  pair.Public,
		})
	}
	return nodes, nil
}
