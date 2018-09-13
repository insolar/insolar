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
	"github.com/stretchr/testify/assert"
)

func TestCalculateNextNodes(t *testing.T) {
	nodeIds := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L"}
	c := SendData{
		NodeIds:           nodeIds,
		Entropy:           core.Entropy{0},
		ReplicationFactor: 2,
	}
	r, _ := CalculateNextNodes(c, "")
	assert.Equal(t, []string{"J", "F"}, r)
	r, _ = CalculateNextNodes(c, "J")
	assert.Equal(t, []string{"H", "D"}, r)
	r, _ = CalculateNextNodes(c, "H")
	assert.Equal(t, []string{"C", "L"}, r)
}

func Test_geometricProgressionSum(t *testing.T) {
	assert.Equal(t, 1022, geometricProgressionSum(2, 2, 9))
	assert.Equal(t, 39, geometricProgressionSum(3, 3, 3))
}

func Test_calcHash(t *testing.T) {
	str := "AAAAAAAAAAAAAAAA"
	c, _ := hex.DecodeString("445215f965178aa7bb7dda56286fe51e45b6e4724dd6a33d4872057c")
	assert.Equal(t, c, calcHash(str, core.Entropy{0}))
}

func Test_getNextCascadeLayerIndexes(t *testing.T) {
	nodeIds := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	startIndex, endIndex := getNextCascadeLayerIndexes(nodeIds, "4", 2)
	assert.Equal(t, 10, startIndex)
	assert.Equal(t, 12, endIndex)
	startIndex, endIndex = getNextCascadeLayerIndexes(nodeIds, "1", 2)
	assert.Equal(t, 4, startIndex)
	assert.Equal(t, 6, endIndex)
	startIndex, endIndex = getNextCascadeLayerIndexes(nodeIds, "2", 3)
	assert.Equal(t, 9, startIndex)
	assert.Equal(t, 12, endIndex)
	startIndex, endIndex = getNextCascadeLayerIndexes(nodeIds, "x", 2)
	assert.Equal(t, len(nodeIds), startIndex)
	assert.Equal(t, len(nodeIds), endIndex)
}
