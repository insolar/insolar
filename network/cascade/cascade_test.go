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

package cascade

import (
	"encoding/hex"
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/require"
)

const (
	domainStr = ".4F7BsTMVPKFshM1MwLf6y23cid6fL3xMpazVoF9krzUw"
	id1Str    = "4K2V1kpVycZ6qSFsNdz2FtpNxnJs17eBNzf9rdCMcKoe"
	id2Str    = "4NwnA4HWZurKyXWNowJwYmb9CwX4gBKzwQKov1ExMf8M"
	id3Str    = "4Ss5JMkXAD9Z7cktFEdrqeMuT6jGMF1pVozTyPHZ6zT4"
	id4Str    = "4WnNSfDXkWSnFi1PgXxn8X8fhFwU2Jhe4Df82mL9rKmm"
	id5Str    = "4ahfaxgYLok1PoFu7qHhRPuRwR9fhNPTcdKn69Nkbf6U"
	id6Str    = "4ecxjG9Yw73EXtWQZ8cciGgCBaMsNS5HB2zS9XRMLzRB"
	id7Str    = "4iYFsZcZXQLTfykuzRwY19SxRja53Vm6jSf6CuTx6Kjt"
	id8Str    = "4nTZ1s5a7hdgp51RRjGTJ2DiftnGiZSvHrKkGHWYqf4b"
	id9Str    = "4rNrAAYahzvuxAFvs2bNatzUv3zUPd8jrFzQKfZ9azPJ"
	id10Str   = "4K1b7kbvUPB935DdMuLqpfmG23zMhxKcHQ9gbdmydPVZ"
	id11Str   = "4K2UQtex1jnjN2Vx8yCMcsmf1HNuMJ4NeA7TgNeVs7kk"
	id12Str   = "4K3Mi2hyZ6QKgynGv33sR5n3zWmSzdo8zv5Em7X26r1w"
)

func TestCalculateNextNodes(t *testing.T) {
	//	t.Skip()
	nodeIds := make([]core.RecordRef, 0)

	ref, err := core.NewRefFromBase58(id1Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)
	ref, err = core.NewRefFromBase58(id2Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)
	ref, err = core.NewRefFromBase58(id3Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)
	ref, err = core.NewRefFromBase58(id4Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)
	ref, err = core.NewRefFromBase58(id5Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)
	ref, err = core.NewRefFromBase58(id6Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)
	ref, err = core.NewRefFromBase58(id7Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)
	ref, err = core.NewRefFromBase58(id8Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)
	ref, err = core.NewRefFromBase58(id9Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)
	ref, err = core.NewRefFromBase58(id10Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)
	ref, err = core.NewRefFromBase58(id11Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)
	ref, err = core.NewRefFromBase58(id12Str + domainStr)
	require.NoError(t, err)
	nodeIds = append(nodeIds, *ref)

	c := core.Cascade{
		NodeIds:           nodeIds,
		Entropy:           core.Entropy{0},
		ReplicationFactor: 2,
	}
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	r, _ := CalculateNextNodes(pcs, c, nil)
	require.Equal(t, []core.RecordRef{nodeIds[3], nodeIds[8]}, r)
	r, _ = CalculateNextNodes(pcs, c, &nodeIds[3])
	require.Equal(t, []core.RecordRef{nodeIds[1], nodeIds[6]}, r)
	r, _ = CalculateNextNodes(pcs, c, &nodeIds[1])
	require.Equal(t, []core.RecordRef{nodeIds[4], nodeIds[7]}, r)
}

func Test_geometricProgressionSum(t *testing.T) {
	require.Equal(t, 1022, geometricProgressionSum(2, 2, 9))
	require.Equal(t, 39, geometricProgressionSum(3, 3, 3))
}

func Test_calcHash(t *testing.T) {
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	ref, err := core.NewRefFromBase58("4SxZ6BSx6qBP41nqQgtsFW5EF3JLDxYscZeVQnviPUGZ.4P3FnY89dNJKiR3qTzMvSbgryVwnxMspus6JrwrAYiVG")
	require.NoError(t, err)
	c, _ := hex.DecodeString("b87473cbc572fd5afc8f165c6b554939068e870490bac3077daf6f9f9ec83df6dbe2d50c061d7034617c24f2d841bc9800121f9f43f90ba91bd8b646df0642ad")
	require.Equal(t, c, calcHash(pcs, *ref, core.Entropy{0}))
}

func Test_getNextCascadeLayerIndexes(t *testing.T) {
	// nodeIds := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	nodeIds := make([]core.RecordRef, 0, 12)
	for i := 0; i < 11; i++ {
		nodeIds = append(nodeIds, testutils.RandomRef())
	}
	startIndex, endIndex := getNextCascadeLayerIndexes(nodeIds, nodeIds[4], 2)
	require.Equal(t, 10, startIndex)
	require.Equal(t, 12, endIndex)
	startIndex, endIndex = getNextCascadeLayerIndexes(nodeIds, nodeIds[1], 2)
	require.Equal(t, 4, startIndex)
	require.Equal(t, 6, endIndex)
	startIndex, endIndex = getNextCascadeLayerIndexes(nodeIds, nodeIds[2], 3)
	require.Equal(t, 9, startIndex)
	require.Equal(t, 12, endIndex)
	startIndex, endIndex = getNextCascadeLayerIndexes(nodeIds, testutils.RandomRef(), 2)
	require.Equal(t, len(nodeIds), startIndex)
	require.Equal(t, len(nodeIds), endIndex)
}
