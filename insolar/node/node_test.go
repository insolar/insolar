// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
