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
	"github.com/insolar/insolar/reference"
)

// ID generates random id.
func ID() insolar.ID {
	var id insolar.ID

	f := fuzz.New().NilChance(0).Funcs(func(id *insolar.ID, c fuzz.Continue) {
		var hash [reference.LocalBinaryHashSize]byte
		c.Fuzz(&hash)

		pn := PulseNumber()

		*id = *insolar.NewID(pn, hash[:])
	})
	f.Fuzz(&id)

	return id
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
func IDWithPulse(pn insolar.PulseNumber) insolar.ID {
	hash := make([]byte, reference.LocalBinaryHashSize)

	fuzz.New().
		NilChance(0).
		NumElements(insolar.RecordHashSize, insolar.RecordHashSize).
		Fuzz(&hash)
	return *insolar.NewID(pn, hash)
}

// JetID generates random jet id.
func JetID() insolar.JetID {
	var jetID insolar.JetID
	f := fuzz.New().Funcs(func(jet *insolar.JetID, c fuzz.Continue) {
		prefix := make([]byte, insolar.JetPrefixSize)
		c.Fuzz(&prefix)
		depth := c.Intn(insolar.JetMaximumDepth + 1)

		*jet = *insolar.NewJetID(uint8(depth), prefix)
	})
	f.Fuzz(&jetID)

	return jetID
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

// Reference generates random reference.
func RecordReference() insolar.Reference {
	return *insolar.NewRecordReference(ID())
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
