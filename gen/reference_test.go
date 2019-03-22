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
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
)

func TestGen_JetID(t *testing.T) {
	for i := 0; i < 10000; i++ {
		jetID := JetID()
		recID := (*insolar.RecordID)(&jetID)
		require.Equalf(t,
			insolar.PulseNumberJet, recID.Pulse(),
			"pulse number should be insolar.PulseNumberJet. jet: %v", recID.DebugString())
		require.GreaterOrEqualf(t,
			uint8(insolar.JetMaximumDepth), jetID.Depth(),
			"jet depth %v should be less than maximum value %v. jet: %v",
			jetID.Depth(), insolar.JetMaximumDepth, jetID.DebugString(),
		)
	}
}
