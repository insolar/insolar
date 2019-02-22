package cmp_test

import (
	"testing"

	"github.com/google/gofuzz"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/gen"
	"github.com/insolar/insolar/ins"
	"github.com/insolar/insolar/ledger/storage/nodes"
	"github.com/stretchr/testify/assert"
)

func TestNodes(t *testing.T) {
	storage := nodes.NewStorage()

	var (
		virtuals  []ins.Node
		materials []ins.Node
		all       []ins.Node
	)
	{
		f := fuzz.New().Funcs(func(e *ins.Node, c fuzz.Continue) {
			e.ID = gen.Reference()
			e.Role = core.StaticRoleVirtual
		})
		f.NumElements(5, 10).NilChance(0).Fuzz(&virtuals)
	}
	{
		f := fuzz.New().Funcs(func(e *ins.Node, c fuzz.Continue) {
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
