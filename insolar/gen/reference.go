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
)

// ID generates random id.
func ID() (id insolar.ID) {
	fuzz.New().NilChance(0).Fuzz(&id)
	return
}

// JetID generates random jet id.
func JetID() (jetID insolar.JetID) {
	f := fuzz.New().Funcs(func(jet *insolar.JetID, c fuzz.Continue) {
		var id insolar.ID
		c.Fuzz(&id)
		copy(jet[:], id[:])
		// set special pulse number
		copy(jet[:insolar.PulseNumberSize], insolar.PulseNumberJet.Bytes())
		// set depth
		// adds 1 because Intn returns [0,n)
		jet[insolar.PulseNumberSize] = byte(c.Intn(insolar.JetMaximumDepth + 1))

	})
	f.Fuzz(&jetID)
	return
}

// UniqueJetIDs generates several different jet ids
func UniqueJetIDs(jets ...*insolar.JetID) {
	seen := make(map[insolar.JetID]struct{})

	for _, j := range jets {
		for {
			*j = JetID()
			if _, ok := seen[*j]; !ok {
				break
			}
		}
		seen[*j] = struct{}{}
	}
}

// Reference generates random reference.
func Reference() (ref insolar.Reference) {
	fuzz.New().NilChance(0).Fuzz(&ref)
	return
}
