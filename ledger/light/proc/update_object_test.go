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

package proc

import (
	"context"
	"crypto/rand"
	"testing"

	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/light/hot"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/insolar/reply"
	"github.com/insolar/insolar/ledger/blob"
	"github.com/insolar/insolar/ledger/light/recentstorage"
	"github.com/insolar/insolar/ledger/object"
	"github.com/insolar/insolar/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	domainID = *genRandomID(0)
)

func genRandomID(pulse insolar.PulseNumber) *insolar.ID {
	buff := [insolar.RecordIDSize - insolar.PulseNumberSize]byte{}
	_, err := rand.Read(buff[:])
	if err != nil {
		panic(err)
	}
	return insolar.NewID(pulse, buff[:])
}

func genRefWithID(id *insolar.ID) *insolar.Reference {
	return insolar.NewReference(*id)
}

func genRandomRef(pulse insolar.PulseNumber) *insolar.Reference {
	return genRefWithID(genRandomID(pulse))
}

func TestMessageHandler_HandleUpdateObject_FetchesIndexFromHeavy(t *testing.T) {
	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	pendingMock := recentstorage.NewPendingStorageMock(t)

	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetPendingStorageMock.Return(pendingMock)

	mb := testutils.NewMessageBusMock(t)
	mb.MustRegisterMock.Return()
	jc := jet.NewCoordinatorMock(t)

	scheme := testutils.NewPlatformCryptographyScheme()
	indexMemoryStor := object.NewInMemoryIndex()

	idLockMock := object.NewIDLockerMock(t)
	idLockMock.LockMock.Return()
	idLockMock.UnlockMock.Return()

	writeManagerMock := hot.NewWriteAccessorMock(t)
	writeManagerMock.BeginFunc = func(context.Context, insolar.PulseNumber) (func(), error) {
		return func() {}, nil
	}

	objIndex := object.Lifeline{LatestState: genRandomID(0), StateID: record.StateActivation}
	amendRecord := record.Amend{
		PrevState: *objIndex.LatestState,
	}
	virtAmend := record.Wrap(amendRecord)
	data, err := virtAmend.Marshal()
	require.NoError(t, err)

	msg := message.UpdateObject{
		Record: data,
		Object: *genRandomRef(0),
	}

	mb.SendFunc = func(c context.Context, gm insolar.Message, o *insolar.MessageSendOptions) (r insolar.Reply, r1 error) {
		if m, ok := gm.(*message.GetObjectIndex); ok {
			assert.Equal(t, msg.Object, m.Object)
			buf := object.EncodeIndex(objIndex)
			return &reply.ObjectIndex{Index: buf}, nil
		}

		panic("unexpected call")
	}

	ctx := inslogger.TestContext(t)
	heavyRef := genRandomRef(0)
	jc.HeavyMock.Return(heavyRef, nil)

	recordStorage := object.NewRecordMemory()

	updateObject := UpdateObject{
		JetID:       insolar.JetID(jetID),
		Message:     &msg,
		PulseNumber: insolar.FirstPulseNumber,
	}
	updateObject.Dep.Bus = mb
	updateObject.Dep.BlobModifier = blob.NewStorageMemory()
	updateObject.Dep.IDLocker = idLockMock
	updateObject.Dep.Coordinator = jc
	updateObject.Dep.LifelineIndex = indexMemoryStor
	updateObject.Dep.PCS = scheme
	updateObject.Dep.RecordModifier = recordStorage
	updateObject.Dep.LifelineStateModifier = indexMemoryStor
	updateObject.Dep.WriteAccessor = writeManagerMock

	rep := updateObject.handle(ctx)
	require.NoError(t, rep.Err)
	objRep, ok := rep.Reply.(*reply.Object)
	require.True(t, ok)

	idx, err := indexMemoryStor.ForID(ctx, insolar.FirstPulseNumber, *msg.Object.Record())
	require.NoError(t, err)
	assert.Equal(t, objRep.State, *idx.LatestState)
}

func TestMessageHandler_HandleUpdateObject_UpdateIndexState(t *testing.T) {
	jetID := insolar.ID(*insolar.NewJetID(0, nil))

	pendingMock := recentstorage.NewPendingStorageMock(t)

	pendingMock.GetRequestsForObjectMock.Return(nil)
	pendingMock.AddPendingRequestMock.Return()
	pendingMock.RemovePendingRequestMock.Return()

	provideMock := recentstorage.NewProviderMock(t)
	provideMock.GetPendingStorageMock.Return(pendingMock)

	writeManagerMock := hot.NewWriteAccessorMock(t)
	writeManagerMock.BeginFunc = func(context.Context, insolar.PulseNumber) (func(), error) {
		return func() {}, nil
	}

	scheme := testutils.NewPlatformCryptographyScheme()
	indexMemoryStor := object.NewInMemoryIndex()
	recordStorage := object.NewRecordMemory()

	idLockMock := object.NewIDLockerMock(t)
	idLockMock.LockMock.Return()
	idLockMock.UnlockMock.Return()

	objIndex := object.Lifeline{
		LatestState:  genRandomID(0),
		StateID:      record.StateActivation,
		LatestUpdate: 0,
		JetID:        insolar.JetID(jetID),
	}
	amendRecord := record.Amend{
		PrevState: *objIndex.LatestState,
	}
	virtAmend := record.Wrap(amendRecord)
	data, err := virtAmend.Marshal()
	require.NoError(t, err)

	msg := message.UpdateObject{
		Record: data,
		Object: *genRandomRef(0),
	}
	ctx := context.Background()
	err = indexMemoryStor.Set(ctx, insolar.FirstPulseNumber, *msg.Object.Record(), objIndex)
	require.NoError(t, err)

	// Act
	updateObject := UpdateObject{
		JetID:       insolar.JetID(jetID),
		Message:     &msg,
		PulseNumber: insolar.FirstPulseNumber,
	}
	updateObject.Dep.BlobModifier = blob.NewStorageMemory()
	updateObject.Dep.IDLocker = idLockMock
	updateObject.Dep.LifelineIndex = indexMemoryStor
	updateObject.Dep.PCS = scheme
	updateObject.Dep.RecordModifier = recordStorage
	updateObject.Dep.LifelineStateModifier = indexMemoryStor
	updateObject.Dep.WriteAccessor = writeManagerMock

	rep := updateObject.handle(ctx)
	require.NoError(t, rep.Err)
	_, ok := rep.Reply.(*reply.Object)
	require.True(t, ok)

	// Arrange
	idx, err := indexMemoryStor.ForID(ctx, insolar.FirstPulseNumber, *msg.Object.Record())
	require.NoError(t, err)
	require.Equal(t, insolar.FirstPulseNumber, int(idx.LatestUpdate))
}
