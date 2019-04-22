/*
 *    Copyright 2019 Insolar Technologies
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

package object

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/stretchr/testify/require"
)

func TestNewPulseIndex(t *testing.T) {
	pi := NewPulseIndex()

	require.NotNil(t, pi)
}

func TestPulseIndex_Add(t *testing.T) {
	pi := NewPulseIndex()

	fID := gen.ID()
	sID := gen.ID()
	tID := gen.ID()
	fthID := gen.ID()

	pn := gen.PulseNumber()

	pi.Add(fID, pn)
	pi.Add(sID, pn)
	pi.Add(tID, pn)
	pi.Add(fthID, pn)

	res := pi.ForPN(pn)

	require.Equal(t, 4, len(res))
	_, ok := res[fID]
	require.Equal(t, true, ok)
	_, ok = res[sID]
	require.Equal(t, true, ok)
	_, ok = res[tID]
	require.Equal(t, true, ok)
	_, ok = res[fthID]
	require.Equal(t, true, ok)
}

func TestPulseIndex_Add_ChangePn(t *testing.T) {
	pi := pulseIndex{
		lastUsagePn: map[insolar.ID]insolar.PulseNumber{},
		idsByPulse:  map[insolar.PulseNumber]map[insolar.ID]struct{}{},
	}

	id := gen.ID()
	fPn := gen.PulseNumber()
	sPn := gen.PulseNumber()

	pi.Add(id, fPn)
	pi.Add(id, sPn)

	_, inFpn := pi.idsByPulse[fPn][id]
	require.Equal(t, false, inFpn)
	_, inSpn := pi.idsByPulse[sPn][id]
	require.Equal(t, true, inSpn)

}

func TestPulseIndex_DeleteForPulse(t *testing.T) {
	pi := pulseIndex{
		lastUsagePn: map[insolar.ID]insolar.PulseNumber{},
		idsByPulse:  map[insolar.PulseNumber]map[insolar.ID]struct{}{},
	}

	fID := gen.ID()
	sID := gen.ID()
	tID := gen.ID()
	fthID := gen.ID()

	fPn := gen.PulseNumber()
	sPn := gen.PulseNumber()

	pi.Add(fID, fPn)
	pi.Add(sID, sPn)
	pi.Add(tID, fPn)
	pi.Add(fthID, sPn)

	pi.DeleteForPulse(fPn)

	require.Equal(t, 1, len(pi.idsByPulse))
	idsStor := pi.idsByPulse[sPn]
	require.Equal(t, 2, len(idsStor))
	_, ok := idsStor[sID]
	require.Equal(t, true, ok)
	_, ok = idsStor[fthID]
	require.Equal(t, true, ok)

	require.Equal(t, 2, len(pi.lastUsagePn))
	pn := pi.lastUsagePn[sID]
	require.Equal(t, sPn, pn)
	pn = pi.lastUsagePn[fthID]
	require.Equal(t, sPn, pn)
}

func TestPulseIndex_ForPN(t *testing.T) {
	pi := pulseIndex{
		lastUsagePn: map[insolar.ID]insolar.PulseNumber{},
		idsByPulse:  map[insolar.PulseNumber]map[insolar.ID]struct{}{},
	}

	fID := gen.ID()
	sID := gen.ID()
	tID := gen.ID()
	fthID := gen.ID()

	fPn := gen.PulseNumber()
	sPn := gen.PulseNumber()

	pi.Add(fID, fPn)
	pi.Add(sID, sPn)
	pi.Add(tID, fPn)
	pi.Add(fthID, sPn)

	res := pi.ForPN(fPn)

	require.Equal(t, 2, len(res))
	_, ok := res[fID]
	require.Equal(t, true, ok)
	_, ok = res[tID]
	require.Equal(t, true, ok)
}

func TestPulseIndex_LastUsage(t *testing.T) {
	pi := pulseIndex{
		lastUsagePn: map[insolar.ID]insolar.PulseNumber{},
		idsByPulse:  map[insolar.PulseNumber]map[insolar.ID]struct{}{},
	}

	id := gen.ID()
	fPn := gen.PulseNumber()
	sPn := gen.PulseNumber()

	pi.Add(id, fPn)
	pi.Add(id, sPn)

	res, ok := pi.LastUsage(id)

	require.Equal(t, true, ok)
	require.Equal(t, res, sPn)
}

func TestPulseIndex_LastUsage_NoIds(t *testing.T) {
	pi := NewPulseIndex()

	id := gen.ID()
	_, ok := pi.LastUsage(id)

	require.Equal(t, false, ok)
}

func TestPulseIndex_ForPN_NoData(t *testing.T) {
	pi := NewPulseIndex()

	res := pi.ForPN(gen.PulseNumber())

	require.Equal(t, 0, len(res))
}
