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

package artifactmanager

import (
	"testing"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/ledger/index"
	"github.com/insolar/insolar/ledger/record"
	"github.com/stretchr/testify/assert"
)

func TestLedgerArtifactManager_GetHistory(t *testing.T) {

	ctx, db, am, be, cleaner := getTestData(t)
	defer cleaner()

	msg := message.GenesisRequest{Name: "my test message"}
	reqRef, err := am.RegisterRequest(ctx, &msg)
	assert.NoError(t, err)
	requestRef := core.NewRecordRef(domainID, *reqRef)

	objID, err := db.SetRecord(
		ctx,
		core.PulseNumber(0),
		&record.ObjectActivateRecord{
			SideEffectRecord: record.SideEffectRecord{
				Domain:  domainRef,
				Request: *requestRef,
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
		*requestRef,
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
		*requestRef,
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

	iterator, err := be.GetHistory(ctx, *genRefWithID(objID), nil)

	rec, err := iterator.Next()
	assert.NoError(t, err)
	assert.Equal(t, updateRec2.(record.ObjectState).PrevStateID(), rec.Record())
	rec, err = iterator.Next()
	assert.NoError(t, err)
	assert.Equal(t, updateRec1.(record.ObjectState).PrevStateID(), rec.Record())
}
