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

package pulse

import (
	"testing"

	"github.com/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeStorage_ForPulseNumber(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	pn := gen.PulseNumber()
	pulse := insolar.Pulse{PulseNumber: pn}
	storage := NewStorageMem()
	storage.storage[pn] = &memNode{pulse: pulse}

	t.Run("returns error when no pulse", func(t *testing.T) {
		res, err := storage.ForPulseNumber(ctx, gen.PulseNumber())
		assert.Equal(t, ErrNotFound, err)
		assert.Equal(t, insolar.Pulse{}, res)
	})

	t.Run("returns correct pulse", func(t *testing.T) {
		res, err := storage.ForPulseNumber(ctx, pn)
		assert.NoError(t, err)
		assert.Equal(t, pulse, res)
	})
}

func TestNodeStorage_Latest(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)

	t.Run("returns error when no pulse", func(t *testing.T) {
		storage := NewStorageMem()
		res, err := storage.Latest(ctx)
		assert.Equal(t, ErrNotFound, err)
		assert.Equal(t, insolar.Pulse{}, res)
	})

	t.Run("returns correct pulse", func(t *testing.T) {
		storage := NewStorageMem()
		pulse := insolar.Pulse{PulseNumber: gen.PulseNumber()}
		storage.head = &memNode{pulse: pulse}
		res, err := storage.Latest(ctx)
		assert.NoError(t, err)
		assert.Equal(t, pulse, res)
	})
}

func TestNodeStorage_Append(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	pn := gen.PulseNumber()
	pulse := insolar.Pulse{PulseNumber: pn}

	t.Run("appends to an empty storage", func(t *testing.T) {
		storage := NewStorageMem()

		err := storage.Append(ctx, pulse)
		require.NoError(t, err)
		require.NotNil(t, storage.head)
		require.NotNil(t, storage.storage[pulse.PulseNumber])
		assert.Equal(t, storage.storage[pulse.PulseNumber], storage.head)
		assert.Equal(t, storage.tail, storage.head)
		assert.Equal(t, memNode{pulse: pulse}, *storage.head)
	})

	t.Run("returns error if pulse number is equal or less", func(t *testing.T) {
		storage := NewStorageMem()
		head := &memNode{pulse: pulse}
		storage.storage[pn] = head
		storage.head = head
		storage.tail = head

		{
			err := storage.Append(ctx, insolar.Pulse{PulseNumber: pn})
			assert.Equal(t, ErrBadPulse, err)
		}
		{
			err := storage.Append(ctx, insolar.Pulse{PulseNumber: pn - 1})
			assert.Equal(t, ErrBadPulse, err)
		}
	})

	t.Run("appends to a filled storage", func(t *testing.T) {
		storage := NewStorageMem()
		head := &memNode{pulse: pulse}
		storage.storage[pn] = head
		storage.head = head
		storage.tail = head
		pulse := pulse
		pulse.PulseNumber += 1

		err := storage.Append(ctx, pulse)
		require.NoError(t, err)
		require.NotNil(t, storage.head)
		require.NotNil(t, storage.storage[pulse.PulseNumber])
		assert.Equal(t, storage.storage[pulse.PulseNumber], storage.head)
		assert.NotEqual(t, storage.tail, storage.head)
		assert.Equal(t, memNode{pulse: pulse, prev: head}, *storage.head)
	})
}

func TestMemoryStorage_Shift(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	pn := gen.PulseNumber()
	pulse := insolar.Pulse{PulseNumber: pn}

	t.Run("returns error if empty", func(t *testing.T) {
		storage := NewStorageMem()
		shifted, err := storage.Shift(ctx)
		assert.Error(t, err)
		assert.Equal(t, insolar.Pulse{}, shifted)
	})

	t.Run("shifts if one in storage", func(t *testing.T) {
		storage := NewStorageMem()
		head := &memNode{pulse: pulse}
		storage.storage[pn] = head
		storage.head = head
		storage.tail = head

		shifted, err := storage.Shift(ctx)
		assert.NoError(t, err)
		assert.Equal(t, pulse, shifted)
		assert.Nil(t, storage.head)
		assert.Nil(t, storage.tail)
		assert.Empty(t, storage.storage)
	})

	t.Run("shifts if two in storage", func(t *testing.T) {
		storage := NewStorageMem()
		tailPulse := pulse
		headPulse := pulse
		headPulse.PulseNumber += 1
		head := &memNode{pulse: headPulse}
		tail := &memNode{pulse: tailPulse}
		head.prev = tail
		tail.next = head
		storage.storage[headPulse.PulseNumber] = head
		storage.storage[tailPulse.PulseNumber] = tail
		storage.head = head
		storage.tail = tail

		shifted, err := storage.Shift(ctx)
		assert.NoError(t, err)
		assert.Equal(t, tail.pulse, shifted)
		assert.Equal(t, storage.head, storage.tail)
		assert.Equal(t, head, storage.storage[head.pulse.PulseNumber])
		assert.Equal(t, memNode{pulse: headPulse}, *head)
	})
}

func TestMemoryStorage_ForwardsBackwards(t *testing.T) {
	t.Parallel()

	ctx := inslogger.TestContext(t)
	storage := NewStorageMem()
	tailPulse := insolar.Pulse{PulseNumber: gen.PulseNumber()}
	headPulse := insolar.Pulse{PulseNumber: tailPulse.PulseNumber + 1}
	head := &memNode{pulse: headPulse}
	tail := &memNode{pulse: tailPulse}
	head.prev = tail
	tail.next = head
	storage.storage[headPulse.PulseNumber] = head
	storage.storage[tailPulse.PulseNumber] = tail
	storage.head = head
	storage.tail = tail

	t.Run("forwards returns itself if zero steps", func(t *testing.T) {
		pulse, err := storage.Forwards(ctx, tailPulse.PulseNumber, 0)
		assert.NoError(t, err)
		assert.Equal(t, pulse, tailPulse)
	})
	t.Run("forwards returns next if one step", func(t *testing.T) {
		pulse, err := storage.Forwards(ctx, tailPulse.PulseNumber, 1)
		assert.NoError(t, err)
		assert.Equal(t, pulse, headPulse)
	})
	t.Run("forwards returns error if forward overflow", func(t *testing.T) {
		pulse, err := storage.Forwards(ctx, tailPulse.PulseNumber, 2)
		assert.Equal(t, ErrNotFound, err)
		assert.Equal(t, insolar.Pulse{}, pulse)
	})
	t.Run("forwards returns error if backward overflow", func(t *testing.T) {
		pulse, err := storage.Forwards(ctx, tailPulse.PulseNumber-1, 1)
		assert.Equal(t, ErrNotFound, err)
		assert.Equal(t, insolar.Pulse{}, pulse)
	})

	t.Run("backwards returns itself if zero steps", func(t *testing.T) {
		pulse, err := storage.Backwards(ctx, headPulse.PulseNumber, 0)
		assert.NoError(t, err)
		assert.Equal(t, pulse, headPulse)
	})
	t.Run("backwards returns next if one step", func(t *testing.T) {
		pulse, err := storage.Backwards(ctx, headPulse.PulseNumber, 1)
		assert.NoError(t, err)
		assert.Equal(t, pulse, tailPulse)
	})
	t.Run("backwards returns error if backward overflow", func(t *testing.T) {
		pulse, err := storage.Backwards(ctx, headPulse.PulseNumber, 2)
		assert.Equal(t, ErrNotFound, err)
		assert.Equal(t, insolar.Pulse{}, pulse)
	})
	t.Run("backwards returns error if forward overflow", func(t *testing.T) {
		pulse, err := storage.Backwards(ctx, headPulse.PulseNumber-1, 1)
		assert.Equal(t, ErrNotFound, err)
		assert.Equal(t, insolar.Pulse{}, pulse)
	})
}
