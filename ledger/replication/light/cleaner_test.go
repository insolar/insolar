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

package light

import (
	"context"
	"testing"

	"github.com/gojuno/minimock"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/jet"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/ledger/recentstorage"
	"github.com/insolar/insolar/ledger/storage/blob"
	"github.com/insolar/insolar/ledger/storage/drop"
	"github.com/insolar/insolar/ledger/storage/node"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/ledger/storage/pulse"
	"github.com/stretchr/testify/require"
)

func TestCleaner_getExcludedIndexes(t *testing.T) {
	ctx := inslogger.TestContext(t)
	pn := insolar.PulseNumber(333)
	ja := jet.NewAccessorMock(t)
	rsp := recentstorage.NewProviderMock(t)

	fJet := gen.JetID()
	sJet := gen.JetID()

	fStorage := recentstorage.NewRecentIndexStorageMock(t)
	fStorage.GetObjectsMock.Return(map[insolar.ID]int{
		*insolar.NewID(334, nil): 999,
		*insolar.NewID(332, nil): 0,
		*insolar.NewID(331, nil): 999,
	})

	sStorage := recentstorage.NewRecentIndexStorageMock(t)
	sStorage.GetObjectsMock.Return(map[insolar.ID]int{
		*insolar.NewID(334, nil): 999,
		*insolar.NewID(32, nil):  0,
		*insolar.NewID(31, nil):  999,
	})

	rsp.GetIndexStorageFunc = func(p context.Context, p1 insolar.ID) (r recentstorage.RecentIndexStorage) {
		if p1 == insolar.ID(fJet) {
			return fStorage
		}
		if p1 == insolar.ID(sJet) {
			return sStorage
		}
		panic("test is totally broken")
	}

	ja.AllMock.Return([]insolar.JetID{fJet, sJet})

	cleaner := cleaner{
		jetAccessor:    ja,
		recentProvider: rsp,
	}

	res := cleaner.getExcludedIndexes(ctx, pn)

	require.Equal(t, 2, len(res))
	_, ok := res[*insolar.NewID(331, nil)]
	require.Equal(t, true, ok)
	_, ok = res[*insolar.NewID(31, nil)]
	require.Equal(t, true, ok)
}

func TestCleaner_Clean(t *testing.T) {
	ctx := inslogger.TestContext(t)

	pn := insolar.PulseNumber(98765)

	ctrl := minimock.NewController(t)

	jm := jet.NewModifierMock(ctrl)
	jm.DeleteMock.Expect(ctx, pn)

	ja := jet.NewAccessorMock(ctrl)
	ja.AllMock.Return([]insolar.JetID{})

	nm := node.NewModifierMock(ctrl)
	nm.DeleteMock.Expect(pn)

	dc := drop.NewCleanerMock(ctrl)
	dc.DeleteMock.Expect(pn)

	bc := blob.NewCleanerMock(ctrl)
	bc.DeleteMock.Expect(ctx, pn)

	rc := object.NewRecordCleanerMock(ctrl)
	rc.RemoveMock.Expect(ctx, pn)

	ic := object.NewIndexCleanerMock(ctrl)
	ic.RemoveUntilMock.Expect(ctx, pn, map[insolar.ID]struct{}{})

	ps := pulse.NewShifterMock(ctrl)
	ps.ShiftMock.Expect(ctx, pn).Return(nil)

	rp := recentstorage.NewProviderMock(ctrl)

	cleaner := NewCleaner(jm, ja, nm, dc, bc, rc, ic, rp, ps)

	cleaner.Clean(ctx, pn)

	ctrl.Finish()
}
