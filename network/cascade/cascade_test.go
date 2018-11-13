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
	"github.com/stretchr/testify/assert"
)

func TestCalculateNextNodes(t *testing.T) {
	//	t.Skip()
	nodeIds := []core.RecordRef{
		core.NewRefFromBase58("4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj"),
		core.NewRefFromBase58("53jNWvey7Nzyh4ZaLdJDf3SRgoD4GpWuwHgrgvVVGLbDkk3A7cwStSmBU2X7s4fm6cZtemEyJbce9dM9SwNxbsxf"),
		core.NewRefFromBase58("9uE5MEWQB2yfKY8kTgTNovWii88D4anmf7GAiovgcxx6Uc6EBmZ212mpyMa1L22u9TcUUd94i8LvULdsdBoG8ed"),
		core.NewRefFromBase58("4qXdYkfL9U4tL3qRPthdbdajtafR4KArcXjpyQSEgEMtpuin3t8aZYmMzKGRnXHBauytaPQ6bfwZyKZzRPpR6gyX"),
		core.NewRefFromBase58("5q5rnvayXyKszoWofxp4YyK7FnLDwhsqAXKxj6H7B5sdEsNn4HKNFoByph4Aj8rGptdWL54ucwMQrySMJgKavxX1"),
		core.NewRefFromBase58("5tsFDwNLMW4GRHxSbBjjxvKpR99G4CSBLRqZAcpqdSk5SaeVcDL3hCiyjjidCRJ7Lu4VZoANWQJN2AgPvSRgCghn"),
		core.NewRefFromBase58("48UWM6w7YKYCHoP7GHhogLvbravvJ6bs4FGETqXfgdhF9aPxiuwDWwHipeiuNBQvx7zyCN9wFxbuRrDYRoAiw5Fj"),
		core.NewRefFromBase58("5owQeqWyHcobFaJqS2BZU2o2ZRQ33GojXkQK6f8vNLgvNx6xeWRwenJMc53eEsS7MCxrpXvAhtpTaNMPr3rjMHA"),
		core.NewRefFromBase58("xF12WfbkcWrjrPXvauSYpEGhkZT2Zha53xpYh5KQdmGHMywJNNgnemfDN2JfPV45aNQobkdma4dsx1N7Xf5wCJ9"),
		core.NewRefFromBase58("4VgDz9o23wmYXN9mEiLnnsGqCEEARGByx1oys2MXtC6M94K85ZpB9sEJwiGDER61gHkBxkwfJqtg9mAFR7PQcssq"),
		core.NewRefFromBase58("48g7C8QnH2CGMa62sNaL1gVVyygkto8EbMRHv168psCBuFR2FXkpTfwk4ZwpY8awFFXKSnWspYWWQ7sMMk5W7s3T"),
		core.NewRefFromBase58("Lvssptdwq7tatd567LUfx2AgsrWZfo4u9q6FJgJ9BgZK8cVooZv2A7F7rrs1FS5VpnTmXhr6XihXuKWVZ8i5YX9"),
	}
	c := core.Cascade{
		NodeIds:           nodeIds,
		Entropy:           core.Entropy{0},
		ReplicationFactor: 2,
	}
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	r, _ := CalculateNextNodes(pcs, c, nil)
	assert.Equal(t, []core.RecordRef{nodeIds[8], nodeIds[5]}, r)
	r, _ = CalculateNextNodes(pcs, c, &nodeIds[8])
	assert.Equal(t, []core.RecordRef{nodeIds[2], nodeIds[0]}, r)
	r, _ = CalculateNextNodes(pcs, c, &nodeIds[2])
	assert.Equal(t, []core.RecordRef{nodeIds[11], nodeIds[10]}, r)
}

func Test_geometricProgressionSum(t *testing.T) {
	assert.Equal(t, 1022, geometricProgressionSum(2, 2, 9))
	assert.Equal(t, 39, geometricProgressionSum(3, 3, 3))
}

func Test_calcHash(t *testing.T) {
	pcs := platformpolicy.NewPlatformCryptographyScheme()

	ref := core.NewRefFromBase58("4gU79K6woTZDvn4YUFHauNKfcHW69X42uyk8ZvRevCiMv3PLS24eM1vcA9mhKPv8b2jWj9J5RgGN9CB7PUzCtBsj")
	c, _ := hex.DecodeString("39e1040cc17bd51bb59803edece23a82a7d8fe01394c337c74c003852f1683e5cba9a396556b6e737e15b54950efb46228cbd1a745c85016b9b3b8fbbe0d94d3")
	assert.Equal(t, c, calcHash(pcs, ref, core.Entropy{0}))
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
