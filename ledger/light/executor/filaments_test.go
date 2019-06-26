package executor_test

import (
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/executor"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFilamentModifier_SetRequest(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)

	var (
		pcs     insolar.PlatformCryptographyScheme
		indexes object.IndexStorage
		records object.RecordStorage
		manager *executor.FilamentModifierDefault
	)
	resetComponents := func() {
		pcs = testutils.NewPlatformCryptographyScheme()
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		manager = executor.NewFilamentModifier(indexes, records, pcs, nil)
	}

	objRef := gen.Reference()
	validRequest := record.Request{Object: &objRef}

	resetComponents()
	t.Run("object id is empty", func(t *testing.T) {
		err := manager.SetRequest(ctx, insolar.ID{}, gen.JetID(), validRequest)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("jet is not valid", func(t *testing.T) {
		err := manager.SetRequest(ctx, gen.ID(), insolar.JetID{}, validRequest)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("index does not exist", func(t *testing.T) {
		err := manager.SetRequest(ctx, gen.ID(), gen.JetID(), validRequest)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("request from the past", func(t *testing.T) {
		reqID := gen.ID()
		reqID.SetPulse(insolar.FirstPulseNumber + 1)
		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 2)

		err := indexes.SetIndex(ctx, reqID.Pulse(), object.FilamentIndex{
			Lifeline: object.Lifeline{
				PendingPointer: &latestPendingID,
			},
		})
		require.NoError(t, err)

		err = manager.SetRequest(ctx, reqID, gen.JetID(), validRequest)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		reqID := gen.ID()
		reqID.SetPulse(insolar.FirstPulseNumber + 2)
		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 1)
		jetID := gen.JetID()

		err := indexes.SetIndex(ctx, reqID.Pulse(), object.FilamentIndex{
			ObjID: *validRequest.Object.Record(),
			Lifeline: object.Lifeline{
				PendingPointer: &latestPendingID,
			},
		})
		require.NoError(t, err)

		err = manager.SetRequest(ctx, reqID, jetID, validRequest)
		assert.NoError(t, err)

		idx, err := indexes.ForID(ctx, reqID.Pulse(), *validRequest.Object.Record())
		require.NoError(t, err)

		expectedFilamentRecord := record.PendingFilament{
			RecordID:       reqID,
			PreviousRecord: &latestPendingID,
		}
		virtual := record.Wrap(expectedFilamentRecord)
		hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
		expectedFilamentRecordID := *insolar.NewID(reqID.Pulse(), hash)

		assert.Equal(t, expectedFilamentRecordID, *idx.Lifeline.PendingPointer)
		assert.Equal(t, reqID.Pulse(), *idx.Lifeline.EarliestOpenRequest)

		rec, err := records.ForID(ctx, expectedFilamentRecordID)
		require.NoError(t, err)
		virtual = record.Wrap(expectedFilamentRecord)
		assert.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

		rec, err = records.ForID(ctx, reqID)
		require.NoError(t, err)
		virtual = record.Wrap(validRequest)
		assert.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

		mc.Finish()
	})
}

func TestFilamentModifier_SetResult(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	defer mc.Finish()

	ctx := inslogger.TestContext(t)

	var (
		pcs     insolar.PlatformCryptographyScheme
		indexes object.IndexStorage
		records object.RecordStorage
		manager *executor.FilamentModifierDefault
	)
	resetComponents := func() {
		pcs = testutils.NewPlatformCryptographyScheme()
		indexes = object.NewIndexStorageMemory()
		records = object.NewRecordMemory()
		manager = executor.NewFilamentModifier(indexes, records, pcs, nil)
	}

	validResult := record.Result{Object: gen.ID()}

	resetComponents()
	t.Run("object id is empty", func(t *testing.T) {
		err := manager.SetResult(ctx, insolar.ID{}, gen.JetID(), validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("jet is not valid", func(t *testing.T) {
		err := manager.SetResult(ctx, gen.ID(), insolar.JetID{}, validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("index does not exist", func(t *testing.T) {
		err := manager.SetResult(ctx, gen.ID(), gen.JetID(), validResult)
		assert.Error(t, err)

		mc.Finish()
	})

	resetComponents()
	t.Run("happy basic", func(t *testing.T) {
		t.Skip()
		reqID := gen.ID()
		reqID.SetPulse(insolar.FirstPulseNumber + 2)
		latestPendingID := gen.ID()
		latestPendingID.SetPulse(insolar.FirstPulseNumber + 1)
		jetID := gen.JetID()

		err := indexes.SetIndex(ctx, reqID.Pulse(), object.FilamentIndex{
			ObjID: validResult.Object,
			Lifeline: object.Lifeline{
				PendingPointer: &latestPendingID,
			},
		})
		require.NoError(t, err)

		err = manager.SetResult(ctx, reqID, jetID, validResult)
		assert.NoError(t, err)

		idx, err := indexes.ForID(ctx, reqID.Pulse(), validResult.Object)
		require.NoError(t, err)

		expectedFilamentRecord := record.PendingFilament{
			RecordID:       reqID,
			PreviousRecord: &latestPendingID,
		}
		virtual := record.Wrap(expectedFilamentRecord)
		hash := record.HashVirtual(pcs.ReferenceHasher(), virtual)
		expectedFilamentRecordID := *insolar.NewID(reqID.Pulse(), hash)

		assert.Equal(t, expectedFilamentRecordID, *idx.Lifeline.PendingPointer)
		assert.Equal(t, reqID.Pulse(), *idx.Lifeline.EarliestOpenRequest)

		rec, err := records.ForID(ctx, expectedFilamentRecordID)
		require.NoError(t, err)
		virtual = record.Wrap(expectedFilamentRecord)
		assert.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

		rec, err = records.ForID(ctx, reqID)
		require.NoError(t, err)
		virtual = record.Wrap(validResult)
		assert.Equal(t, record.Material{Virtual: &virtual, JetID: jetID}, rec)

		mc.Finish()
	})
}
