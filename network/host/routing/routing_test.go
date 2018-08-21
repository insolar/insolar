/*
 *    Copyright 2018 INS Ecosystem
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

package routing

import (
	"math/big"
	"sort"
	"testing"

	"github.com/insolar/insolar/network/host/id"
	"github.com/insolar/insolar/network/host/node"

	"github.com/stretchr/testify/assert"
)

func TestDistanceMetric(t *testing.T) {
	n := NewRouteNode(&node.Node{})
	n.ID = getIDWithValues(0)
	assert.Equal(t, 20, len(n.ID.GetHash()))

	value := getDistance(n.ID.GetHash(), getIDWithValues(0).GetHash())
	assert.Equal(t, 0, value.Cmp(new(big.Int).SetInt64(int64(0))))

	v := getIDWithValues(0)
	v.GetHash()[19] = byte(1)
	value = getDistance(n.ID.GetHash(), v.GetHash())
	assert.Equal(t, big.NewInt(1), value)

	v = getIDWithValues(0)
	v.GetHash()[18] = byte(1)
	value = getDistance(n.ID.GetHash(), v.GetHash())
	assert.Equal(t, big.NewInt(256), value)

	v = getIDWithValues(255)
	value = getDistance(n.ID.GetHash(), v.GetHash())

	// (2^160)-1 = max possible distance
	maxDistance := new(big.Int).Exp(big.NewInt(2), big.NewInt(160), nil)
	maxDistance.Sub(maxDistance, big.NewInt(1))

	assert.Equal(t, maxDistance, value)
}

func TestHasBit(t *testing.T) {
	for i := uint8(0); i < 8; i++ {
		assert.Equal(t, true, hasBit(byte(255), i))
	}

	assert.Equal(t, true, hasBit(byte(1), 7))

	for i := uint8(0); i < 8; i++ {
		assert.Equal(t, false, hasBit(byte(0), i))
	}
}

func TestRouteSet(t *testing.T) {
	nl := NewRouteSet()
	comparator := getIDWithValues(0)
	n1 := &node.Node{ID: getZerodIDWithNthByte(19, 1)}
	n2 := &node.Node{ID: getZerodIDWithNthByte(18, 1)}
	n3 := &node.Node{ID: getZerodIDWithNthByte(17, 1)}
	n4 := &node.Node{ID: getZerodIDWithNthByte(16, 1)}

	nl.nodes = []*node.Node{n3, n2, n4, n1}
	nl.comparator = comparator.GetHash()

	sort.Sort(nl)

	assert.Equal(t, n1, nl.nodes[0])
	assert.Equal(t, n2, nl.nodes[1])
	assert.Equal(t, n3, nl.nodes[2])
	assert.Equal(t, n4, nl.nodes[3])
}

func getZerodIDWithNthByte(n int, v byte) id.ID {
	id := getIDWithValues(0)
	id.GetHash()[n] = v
	return id
}

func getIDWithValues(b byte) id.ID {
	id1, _ := id.NewID(nil)
	id1.SetHash([]byte{b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b, b})
	return id1
}
