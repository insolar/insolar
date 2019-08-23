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

package gen

import (
	fuzz "github.com/google/gofuzz"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/bits"
)

// ID generates random id.
func ID() (id insolar.ID) {
	fuzz.New().NilChance(0).Fuzz(&id)
	return
}

// UniqueIDs generates multiple random unique IDs.
func UniqueIDs(a int) []insolar.ID {
	ids := make([]insolar.ID, a)
	seen := make(map[insolar.ID]struct{})

	for i := 0; i < a; i++ {
		for {
			ids[i] = ID()
			if _, ok := seen[ids[i]]; !ok {
				break
			}
		}
		seen[ids[i]] = struct{}{}
	}
	return ids
}

// IDWithPulse generates random id with provided pulse.
func IDWithPulse(pn insolar.PulseNumber) (id insolar.ID) {
	copy(id[:insolar.PulseNumberSize], pn.Bytes())
	fill := id[insolar.PulseNumberSize:]
	fuzz.New().
		NilChance(0).
		NumElements(insolar.RecordHashSize, insolar.RecordHashSize).
		Fuzz(&fill)
	copy(id[insolar.PulseNumberSize:], fill)
	return
}

// JetID generates random jet id.
func JetID() (jetID insolar.JetID) {
	f := fuzz.New().Funcs(func(jet *insolar.JetID, c fuzz.Continue) {
		id := ID()
		copy(jet[:], id[:])
		// set special pulse number
		copy(jet[:insolar.PulseNumberSize], insolar.PulseNumberJet.Bytes())
		// set depth
		// adds 1 because Intn returns [0,n)
		depth := byte(c.Intn(insolar.JetMaximumDepth + 1))
		jet[insolar.PulseNumberSize] = depth

		resetJet := bits.ResetBits(jet[:], depth+insolar.PulseNumberSize*8)
		copy(jet[:], resetJet)
	})
	f.Fuzz(&jetID)
	return
}

// UniqueJetIDs generates several different jet ids
func UniqueJetIDs(a int) []insolar.JetID {
	ids := make([]insolar.JetID, a)
	seen := make(map[insolar.JetID]struct{})

	for i := 0; i < a; i++ {
		for {
			ids[i] = JetID()
			if _, ok := seen[ids[i]]; !ok {
				break
			}
		}
		seen[ids[i]] = struct{}{}
	}
	return ids
}

// Reference generates random reference.
func Reference() insolar.Reference {
	return *insolar.NewReference(ID())
}

// UniqueReferences generates multiple random unique References.
func UniqueReferences(a int) []insolar.Reference {
	refs := make([]insolar.Reference, a)
	seen := make(map[insolar.Reference]struct{})

	for i := 0; i < a; i++ {
		for {
			refs[i] = Reference()
			if _, ok := seen[refs[i]]; !ok {
				break
			}
		}
		seen[refs[i]] = struct{}{}
	}
	return refs
}
