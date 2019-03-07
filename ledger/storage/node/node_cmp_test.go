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

package node_test

import (
	"testing"

	"github.com/google/gofuzz"
	"github.com/insolar/insolar"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/stretchr/testify/assert"
)

func TestNodes(t *testing.T) {
	storage := node.NewStorage()

	var (
		virtuals  []insolar.Node
		materials []insolar.Node
		all       []insolar.Node
	)
	{
		f := fuzz.New().Funcs(func(e *insolar.Node, c fuzz.Continue) {
			e.ID = gen.Reference()
			e.Role = core.StaticRoleVirtual
		})
		f.NumElements(5, 10).NilChance(0).Fuzz(&virtuals)
	}
	{
		f := fuzz.New().Funcs(func(e *insolar.Node, c fuzz.Continue) {
			e.ID = gen.Reference()
			e.Role = core.StaticRoleLightMaterial
		})
		f.NumElements(5, 10).NilChance(0).Fuzz(&materials)
	}
	all = append(virtuals, materials...)
	pulse := gen.PulseNumber()

	t.Run("saves nodes", func(t *testing.T) {
		err := storage.Set(pulse, all)
		assert.NoError(t, err)
	})
	t.Run("returns all nodes", func(t *testing.T) {
		result, err := storage.All(pulse)
		assert.NoError(t, err)
		assert.Equal(t, all, result)
	})
	t.Run("returns in role nodes", func(t *testing.T) {
		result, err := storage.InRole(pulse, core.StaticRoleVirtual)
		assert.NoError(t, err)
		assert.Equal(t, virtuals, result)
	})
	t.Run("deletes nodes", func(t *testing.T) {
		storage.Delete(pulse)
		result, err := storage.All(pulse)
		assert.Equal(t, core.ErrNoNodes, err)
		assert.Nil(t, result)
	})
}
