package artifactmanager

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLedgerArtifactManager_GetHistory(t *testing.T) {

	ctx, db, am, cleaner := getTestData(t)
	defer cleaner()

	objID, err := db.SetRecord(
		ctx,
		core.PulseNumber(0),
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain: domainRef,
			},
		},
	)
	db.SetObjectIndex(ctx, objID, &index.ObjectLifeline{
		State:               record.StateActivation,
		LatestState:         objID,
		LatestStateApproved: objID,
	})

	memory := []byte("123")
	prototype := genRandomRef(0)
	obj1, err := am.UpdateObject(
		ctx,
		domainRef,
		requestRef,
		&ObjectDescriptor{
			ctx:       ctx,
			head:      *genRefWithID(objID),
			state:     *objID,
			prototype: prototype,
		},
		memory,
	)
	assert.Nil(t, err)
	updateRec1, err := db.GetRecord(ctx, obj1.StateID())
	assert.Equal(t, record.CalculateIDForBlob(core.GenesisPulse.PulseNumber, memory), updateRec1.(record.ObjectState).GetMemory())
	assert.Nil(t, err)

	memory = []byte("456")
	prototype = genRandomRef(1)
	obj2, err := am.UpdateObject(
		ctx,
		domainRef,
		requestRef,
		&ObjectDescriptor{
			ctx:       ctx,
			head:      *genRefWithID(objID),
			state:     *obj1.StateID(),
			prototype: prototype,
		},
		memory,
	)
	assert.Nil(t, err)

	db.SetObjectIndex(ctx, objID, &index.ObjectLifeline{
		State:               record.StateAmend,
		LatestState:         obj2.StateID(),
		LatestStateApproved: obj2.StateID(),
	})

	updateRec2, err := db.GetRecord(ctx, obj2.StateID())
	assert.Nil(t, err)
	assert.Equal(t, record.CalculateIDForBlob(core.GenesisPulse.PulseNumber, memory), updateRec2.(record.ObjectState).GetMemory())

	iterator, err := am.GetHistory(ctx, *genRefWithID(objID), nil)

	rec, err := iterator.Next()
	assert.NoError(t, err)
	assert.Equal(t, updateRec2.(record.ObjectState).PrevStateID(), rec.Record())
	rec, err = iterator.Next()
	assert.NoError(t, err)
	assert.Equal(t, updateRec1.(record.ObjectState).PrevStateID(), rec.Record())
}
