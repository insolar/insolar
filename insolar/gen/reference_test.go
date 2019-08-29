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

package gen_test

import (
	"testing"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/pulse"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
)

func TestGen_JetID(t *testing.T) {
	for i := 0; i < 100; i++ {
		jetID := gen.JetID()
		recID := (*insolar.ID)(&jetID)
		require.Equalf(t,
			pulse.Jet, recID.Pulse(),
			"pulse number should be insolar.PulseNumberJet. jet: %v", recID.DebugString())
		require.GreaterOrEqualf(t,
			uint8(insolar.JetMaximumDepth), jetID.Depth(),
			"jet depth %v should be less than maximum value %v. jet: %v",
			jetID.Depth(), insolar.JetMaximumDepth, jetID.DebugString(),
		)
	}
}

func TestGen_IDWithPulse(t *testing.T) {
	// Empty slice for comparison.
	emptySlice := make([]byte, insolar.RecordHashSize)

	for i := 0; i < 100; i++ {
		pulse := gen.PulseNumber()

		idWithPulse := gen.IDWithPulse(pulse)

		require.Equal(t,
			idWithPulse.Pulse().Bytes(),
			pulse.Bytes(), "pulse bytes should be equal pulse bytes from generated ID")

		pulseFromID := idWithPulse.Pulse()
		require.Equal(t,
			pulse, pulseFromID,
			"pulse should be equal pulse from generated ID")

		idHash := idWithPulse.Hash()
		require.NotEqual(t,
			emptySlice, idHash,
			"ID.Hash() should not be empty")
	}
}
