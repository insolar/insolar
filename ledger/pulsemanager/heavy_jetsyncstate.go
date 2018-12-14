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

package pulsemanager

import (
	"sync"

	"github.com/insolar/insolar/core"
)

func newJetSyncStates() *jetSyncStates {
	return &jetSyncStates{
		states: map[core.RecordID][]core.PulseNumber{},
	}
}

type jetSyncStates struct {
	sync.Mutex
	states map[core.RecordID][]core.PulseNumber
}

func (jss *jetSyncStates) setJetPulses(jet core.RecordID, pns []core.PulseNumber) {
	jss.Lock()
	jss.states[jet] = pns
	jss.Unlock()
}

func (jss *jetSyncStates) addPulseToJets(jets []core.RecordID, pn core.PulseNumber) {
	jss.Lock()
	for _, jet := range jets {
		jss.states[jet] = append(jss.states[jet], pn)
	}
	jss.Unlock()
}

func (jss *jetSyncStates) puJetPulse(jet core.RecordID, pn core.PulseNumber) {
	jss.Lock()
	jss.states[jet] = append(jss.states[jet], pn)
	jss.Unlock()
}

func (jss *jetSyncStates) unshiftJetPulse(jet core.RecordID) *core.PulseNumber {
	jss.Lock()
	defer jss.Unlock()
	return jss.unshiftJetPulseNoLock(jet)
}

func (jss *jetSyncStates) unshiftJetPulseNoLock(jet core.RecordID) *core.PulseNumber {
	l := jss.states[jet]
	if len(l) == 0 {
		return nil
	}
	result := l[0]
	newl := l[:len(l)-1]
	copy(newl, l[1:])
	jss.states[jet] = newl
	return &result
}

type jetpulse struct {
	jet core.RecordID
	pn  core.PulseNumber
}

func (jss *jetSyncStates) unshiftJetPulses() (pairs []jetpulse) {
	jss.Lock()
	defer jss.Unlock()
	for jet := range jss.states {
		pnRef := jss.unshiftJetPulseNoLock(jet)
		if pnRef == nil {
			continue
		}
		pairs = append(pairs, jetpulse{jet: jet, pn: *pnRef})
	}
	return
}
