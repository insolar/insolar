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

package node

import (
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/stretchr/testify/assert"
)

func TestNodeStorage_All(t *testing.T) {
	t.Parallel()

	var all []insolar.Node
	f := fuzz.New().Funcs(func(e *insolar.Node, c fuzz.Continue) {
		e.ID = gen.Reference()
	})
	f.NumElements(5, 10).NilChance(0).Fuzz(&all)
	pulse := gen.PulseNumber()

	t.Run("returns correct nodes", func(t *testing.T) {
		nodeStorage := NewStorage()
		nodeStorage.nodes[pulse] = all
		result, err := nodeStorage.All(pulse)
		assert.NoError(t, err)
		assert.Equal(t, all, result)
	})

	t.Run("returns nil when empty nodes", func(t *testing.T) {
		nodeStorage := NewStorage()
		nodeStorage.nodes[pulse] = nil
		result, err := nodeStorage.All(pulse)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("returns error when no nodes", func(t *testing.T) {
		nodeStorage := NewStorage()
		result, err := nodeStorage.All(pulse)
		assert.Equal(t, ErrNoNodes, err)
		assert.Nil(t, result)
	})
}

func TestNodeStorage_InRole(t *testing.T) {
	t.Parallel()

	var (
		virtuals  []insolar.Node
		materials []insolar.Node
		all       []insolar.Node
	)
	{
		f := fuzz.New().Funcs(func(e *insolar.Node, c fuzz.Continue) {
			e.ID = gen.Reference()
			e.Role = insolar.StaticRoleVirtual
		})
		f.NumElements(5, 10).NilChance(0).Fuzz(&virtuals)
	}
	{
		f := fuzz.New().Funcs(func(e *insolar.Node, c fuzz.Continue) {
			e.ID = gen.Reference()
			e.Role = insolar.StaticRoleLightMaterial
		})
		f.NumElements(5, 10).NilChance(0).Fuzz(&materials)
	}
	all = append(virtuals, materials...)
	pulse := gen.PulseNumber()

	t.Run("returns correct nodes", func(t *testing.T) {
		nodeStorage := NewStorage()
		nodeStorage.nodes[pulse] = all
		{
			result, err := nodeStorage.InRole(pulse, insolar.StaticRoleVirtual)
			assert.NoError(t, err)
			assert.Equal(t, virtuals, result)
		}
		{
			result, err := nodeStorage.InRole(pulse, insolar.StaticRoleLightMaterial)
			assert.NoError(t, err)
			assert.Equal(t, materials, result)
		}
	})

	t.Run("returns nil when empty nodes", func(t *testing.T) {
		nodeStorage := NewStorage()
		nodeStorage.nodes[pulse] = nil
		result, err := nodeStorage.InRole(pulse, insolar.StaticRoleVirtual)
		assert.NoError(t, err)
		assert.Nil(t, result)
	})

	t.Run("returns error when no nodes", func(t *testing.T) {
		nodeStorage := NewStorage()
		result, err := nodeStorage.InRole(pulse, insolar.StaticRoleVirtual)
		assert.Equal(t, ErrNoNodes, err)
		assert.Nil(t, result)
	})
}

func TestStorage_Set(t *testing.T) {
	t.Parallel()

	var nodes []insolar.Node
	f := fuzz.New().Funcs(func(e *insolar.Node, c fuzz.Continue) {
		e.ID = gen.Reference()
	})
	f.NumElements(5, 10).NilChance(0).Fuzz(&nodes)
	pulse := gen.PulseNumber()

	t.Run("saves correct nodes", func(t *testing.T) {
		nodeStorage := NewStorage()
		err := nodeStorage.Set(pulse, nodes)
		assert.NoError(t, err)
		assert.Equal(t, nodes, nodeStorage.nodes[pulse])
	})

	t.Run("saves nil if empty nodes", func(t *testing.T) {
		nodeStorage := NewStorage()
		err := nodeStorage.Set(pulse, []insolar.Node{})
		assert.NoError(t, err)
		assert.Nil(t, nodeStorage.nodes[pulse])
	})

	t.Run("returns error when saving with the same pulse", func(t *testing.T) {
		nodeStorage := NewStorage()
		_ = nodeStorage.Set(pulse, nodes)
		err := nodeStorage.Set(pulse, nodes)
		assert.Equal(t, ErrOverride, err)
		assert.Equal(t, nodes, nodeStorage.nodes[pulse])
	})
}

func TestNewStorage_Delete(t *testing.T) {
	t.Parallel()

	var nodes []insolar.Node
	f := fuzz.New().Funcs(func(e *insolar.Node, c fuzz.Continue) {
		e.ID = gen.Reference()
	})
	f.NumElements(5, 10).NilChance(0).Fuzz(&nodes)
	pulse := gen.PulseNumber()
	nodeStorage := NewStorage()
	nodeStorage.nodes[pulse] = nodes

	t.Run("removes nodes for pulse", func(t *testing.T) {
		{
			result, err := nodeStorage.All(pulse)
			assert.NoError(t, err)
			assert.Equal(t, nodes, result)
		}
		{
			nodeStorage.DeleteForPN(pulse)
			result, err := nodeStorage.All(pulse)
			assert.Equal(t, ErrNoNodes, err)
			assert.Nil(t, result)
		}
	})
}
