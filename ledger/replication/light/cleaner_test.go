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

// func TestCleaner_Clean(t *testing.T) {
// 	ctx := inslogger.TestContext(t)
//
// 	pn := insolar.PulseNumber(98765)
//
// 	ctrl := minimock.NewController(t)
//
// 	jm := jet.NewModifierMock(ctrl)
// 	jm.DeleteMock.Expect(ctx, pn)
//
// 	ja := jet.NewAccessorMock(ctrl)
// 	ja.AllMock.Return([]insolar.JetID{})
//
// 	nm := node.NewModifierMock(ctrl)
// 	nm.DeleteMock.Expect(pn)
//
// 	dc := drop.NewCleanerMock(ctrl)
// 	dc.DeleteMock.Expect(pn)
//
// 	bc := blob.NewCleanerMock(ctrl)
// 	bc.DeleteMock.Expect(ctx, pn)
//
// 	rc := object.NewRecordCleanerMock(ctrl)
// 	rc.RemoveMock.Expect(ctx, pn)
//
//
// 	ps := pulse.NewShifterMock(ctrl)
// 	ps.ShiftMock.Expect(ctx, pn).Return(nil)
//
// 	rp := recentstorage.NewProviderMock(ctrl)
//
// 	cleaner := NewCleaner(jm, ja, nm, dc, bc, rc, ic, rp, ps)
//
// 	cleaner.NotifyAboutPulse(ctx, pn)
//
// 	ctrl.Finish()
// }
