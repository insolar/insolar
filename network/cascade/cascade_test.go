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
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
)

func TestCalculateNextNodes(t *testing.T) {
	nodeIds := []core.RecordRef{
		core.String2Ref("4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"),
		core.String2Ref("53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"),
		core.String2Ref("9uE5MEWQB2yfKY8kTgTNovWii88D4anmf7GAiovgcxx6Uc6EBmZ212mpyMa1L22u9TcUUd94i8LvULdsdBoG8ed"),
		core.String2Ref("4qXdYkfL9U4tL3qRPthdbdajtafR4KArcXjpyQSEgEMtpuin3t8aZYmMzKGRnXHBauytaPQ6bfwZyKZzRPpR6gyX"),
		core.String2Ref("5q5rnvayXyKszoWofxp4YyK7FnLDwhsqAXKxj6H7B5sdEsNn4HKNFoByph4Aj8rGptdWL54ucwMQrySMJgKavxX1"),
		core.String2Ref("5tsFDwNLMW4GRHxSbBjjxvKpR99G4CSBLRqZAcpqdSk5SaeVcDL3hCiyjjidCRJ7Lu4VZoANWQJN2AgPvSRgCghn"),
		core.String2Ref("48UWM6w7YKYCHoP7GHhogLvbravvJ6bs4FGETqXfgdhF9aPxiuwDWwHipeiuNBQvx7zyCN9wFxbuRrDYRoAiw5Fj"),
		core.String2Ref("5owQeqWyHcobFaJqS2BZU2o2ZRQ33GojXkQK6f8vNLgvNx6xeWRwenJMc53eEsS7MCxrpXvAhtpTaNMPr3rjMHA"),
		core.String2Ref("xF12WfbkcWrjrPXvauSYpEGhkZT2Zha53xpYh5KQdmGHMywJNNgnemfDN2JfPV45aNQobkdma4dsx1N7Xf5wCJ9"),
		core.String2Ref("4VgDz9o23wmYXN9mEiLnnsGqCEEARGByx1oys2MXtC6M94K85ZpB9sEJwiGDER61gHkBxkwfJqtg9mAFR7PQcssq"),
		core.String2Ref("48g7C8QnH2CGMa62sNaL1gVVyygkto8EbMRHv168psCBuFR2FXkpTfwk4ZwpY8awFFXKSnWspYWWQ7sMMk5W7s3T"),
		core.String2Ref("Lvssptdwq7tatd567LUfx2AgsrWZfo4u9q6FJgJ9BgZK8cVooZv2A7F7rrs1FS5VpnTmXhr6XihXuKWVZ8i5YX9"),
	}
	c := core.Cascade{
		NodeIds:           nodeIds,
		Entropy:           core.Entropy{0},
		ReplicationFactor: 2,
	}
	r, _ := CalculateNextNodes(c, nil)
	assert.Equal(t, []core.RecordRef{nodeIds[5], nodeIds[8]}, r)
	r, _ = CalculateNextNodes(c, &nodeIds[5])
	assert.Equal(t, []core.RecordRef{nodeIds[10], nodeIds[6]}, r)
	r, _ = CalculateNextNodes(c, &nodeIds[10])
	assert.Equal(t, []core.RecordRef{nodeIds[0], nodeIds[4]}, r)
}

func Test_geometricProgressionSum(t *testing.T) {
	assert.Equal(t, 1022, geometricProgressionSum(2, 2, 9))
	assert.Equal(t, 39, geometricProgressionSum(3, 3, 3))
}

func Test_calcHash(t *testing.T) {
	ref := core.String2Ref("4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj")
	c, _ := hex.DecodeString("65e64988fde08c6dc587f30bbe4a6881e94b7e07ec7c152cfc1aa764")
	assert.Equal(t, c, calcHash(ref, core.Entropy{0}))
}

func Test_getNextCascadeLayerIndexes(t *testing.T) {
	// nodeIds := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}
	nodeIds := make([]core.RecordRef, 0, 12)
	for i := 0; i < 11; i++ {
		nodeIds = append(nodeIds, testutils.RandomRef())
	}
	startIndex, endIndex := getNextCascadeLayerIndexes(nodeIds, nodeIds[4], 2)
	assert.Equal(t, 10, startIndex)
	assert.Equal(t, 12, endIndex)
	startIndex, endIndex = getNextCascadeLayerIndexes(nodeIds, nodeIds[1], 2)
	assert.Equal(t, 4, startIndex)
	assert.Equal(t, 6, endIndex)
	startIndex, endIndex = getNextCascadeLayerIndexes(nodeIds, nodeIds[2], 3)
	assert.Equal(t, 9, startIndex)
	assert.Equal(t, 12, endIndex)
	startIndex, endIndex = getNextCascadeLayerIndexes(nodeIds, testutils.RandomRef(), 2)
	assert.Equal(t, len(nodeIds), startIndex)
	assert.Equal(t, len(nodeIds), endIndex)
}
