package proc_test

import (
	"context"
	"math/rand"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/gojuno/minimock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bus"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/heavy/proc"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/pulse"
	"github.com/insolar/insolar/testutils"
)

func TestSendRequests_Proceed(t *testing.T) {
	mc := minimock.NewController(t)
	ctx := inslogger.TestContext(t)

	var (
		sender  *bus.SenderMock
		records *object.RecordMemory
		indexes *object.IndexAccessorMock
		p       *proc.SendRequests
	)

	pcs := testutils.NewPlatformCryptographyScheme()

	resetComponents := func() {
		sender = bus.NewSenderMock(mc)
		indexes = object.NewIndexAccessorMock(mc)
		records = object.NewRecordMemory()
	}

	newProc := func(msg payload.Meta) *proc.SendRequests {
		p := proc.NewSendRequests(msg)
		p.Dep(sender, records, indexes)
		return p
	}

	resetComponents()
	t.Run("object does not exist", func(t *testing.T) {
		p = newProc(payload.Meta{})

		indexes.ForIDMock.Return(record.Index{}, object.ErrNotFound)

		err := p.Proceed(ctx)
		require.Error(t, err)
	})

	resetComponents()
	t.Run("empty response", func(t *testing.T) {
		msg := payload.GetFilament{
			ObjectID:  gen.ID(),
			StartFrom: gen.ID(),
		}
		buf, err := msg.Marshal()
		require.NoError(t, err)
		receivedMeta := payload.Meta{Payload: buf}
		p = newProc(receivedMeta)

		indexes.ForIDMock.Set(func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) (record.Index, error) {
			require.Equal(t, msg.StartFrom.Pulse(), pn)
			require.Equal(t, msg.ObjectID, id)
			return record.Index{}, nil
		})

		err = p.Proceed(ctx)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		b := newFilamentBuilder(ctx, pcs, records)
		rec1 := b.Append(pulse.MinTimePulse+1, &record.IncomingRequest{Nonce: rand.Uint64()})
		rec2 := b.Append(pulse.MinTimePulse+2, &record.IncomingRequest{Nonce: rand.Uint64()})
		rec3 := b.Append(pulse.MinTimePulse+4, &record.IncomingRequest{Nonce: rand.Uint64()})
		b.Append(pulse.MinTimePulse+5, &record.IncomingRequest{Nonce: rand.Uint64()})

		msg := payload.GetFilament{
			ObjectID:  gen.ID(),
			StartFrom: rec3.MetaID,
		}
		buf, err := msg.Marshal()
		require.NoError(t, err)
		receivedMeta := payload.Meta{Payload: buf}
		p = newProc(receivedMeta)

		indexes.ForIDMock.Set(func(_ context.Context, pn insolar.PulseNumber, id insolar.ID) (record.Index, error) {
			require.Equal(t, msg.StartFrom.Pulse(), pn)
			require.Equal(t, msg.ObjectID, id)
			return record.Index{}, nil
		})

		sender.ReplyMock.Set(func(_ context.Context, origin payload.Meta, rep *message.Message) {
			require.Equal(t, receivedMeta, origin)

			resp, err := payload.Unmarshal(rep.Payload)
			require.NoError(t, err)

			filaments, ok := resp.(*payload.FilamentSegment)
			require.True(t, ok)
			assert.Equal(t, msg.ObjectID, filaments.ObjectID)
			assert.Equal(t, []record.CompositeFilamentRecord{rec3, rec2, rec1}, filaments.Records)
		})

		err = p.Proceed(ctx)
		assert.NoError(t, err)

		mc.Wait(10 * time.Minute)
		mc.Finish()
	})
}

type filamentBuilder struct {
	records   object.AtomicRecordModifier
	currentID insolar.ID
	ctx       context.Context
	pcs       insolar.PlatformCryptographyScheme
}

func newFilamentBuilder(
	ctx context.Context,
	pcs insolar.PlatformCryptographyScheme,
	records object.AtomicRecordModifier,
) *filamentBuilder {
	return &filamentBuilder{
		ctx:     ctx,
		records: records,
		pcs:     pcs,
	}
}

func (b *filamentBuilder) Append(pn insolar.PulseNumber, rec record.Record) record.CompositeFilamentRecord {
	return b.append(pn, rec, true)
}

func (b *filamentBuilder) AppendNoPersist(pn insolar.PulseNumber, rec record.Record) record.CompositeFilamentRecord {
	return b.append(pn, rec, false)
}

func (b *filamentBuilder) append(pn insolar.PulseNumber, rec record.Record, persist bool) record.CompositeFilamentRecord {
	var composite record.CompositeFilamentRecord
	{
		virtual := record.Wrap(rec)
		hash := record.HashVirtual(b.pcs.ReferenceHasher(), virtual)
		id := *insolar.NewID(pn, hash)
		material := record.Material{
			Virtual: virtual,
			ID:      id,
			JetID:   insolar.ZeroJetID,
		}
		if persist {
			err := b.records.SetAtomic(b.ctx, material)
			if err != nil {
				panic(err)
			}
		}
		composite.RecordID = id
		composite.Record = material
	}

	{
		rec := record.PendingFilament{RecordID: composite.RecordID}
		if !b.currentID.IsEmpty() {
			curr := b.currentID
			rec.PreviousRecord = &curr
		}
		virtual := record.Wrap(&rec)
		hash := record.HashVirtual(b.pcs.ReferenceHasher(), virtual)
		id := *insolar.NewID(pn, hash)
		material := record.Material{
			Virtual: virtual,
			ID:      id,
			JetID:   insolar.ZeroJetID,
		}
		if persist {
			err := b.records.SetAtomic(b.ctx, material)
			if err != nil {
				panic(err)
			}
		}
		composite.MetaID = id
		composite.Meta = material
	}

	b.currentID = composite.MetaID

	return composite
}
