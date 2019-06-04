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

package bootstrap

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/stretchr/testify/require"
)

func TestReferences(t *testing.T) {
	pairs := map[string]struct {
		got    insolar.Reference
		expect string
	}{
		insolar.GenesisNameRootDomain: {
			got:    ContractRootDomain,
			expect: "1tJD1WuE5mNrtKEdqCQUMTGSKJzdUjp21HNitWcMFz.11111111111111111111111111111111",
		},
		insolar.GenesisNameNodeDomain: {
			got:    ContractNodeDomain,
			expect: "1tJE1AHr77o4xpXug1kRSvYN2hckA38tdehYtsLeLj.11111111111111111111111111111111",
		},
		insolar.GenesisNameNodeRecord: {
			got:    ContractNodeRecord,
			expect: "1tJBqEWVDz6CzF5jGVB3K48FXLjsqRmAoTgJCJBk39.11111111111111111111111111111111",
		},
		insolar.GenesisNameRootMember: {
			got:    ContractRootMember,
			expect: "1tJBsrjLEU53BnKgComgJrWYX3hp7UwMYKmeEZwm6G.11111111111111111111111111111111",
		},
		insolar.GenesisNameRootWallet: {
			got:    ContractWallet,
			expect: "1tJDJepz5dRok7CD6orVs2baCVhHF4KgW8CK5ufifa.11111111111111111111111111111111",
		},
		insolar.GenesisNameAllowance: {
			got:    ContractAllowance,
			expect: "1tJCJAFCEor2Fm3r9gSUiCQiugzPPUvnk8uk8zz9aT.11111111111111111111111111111111",
		},
	}

	for n, p := range pairs {
		t.Run(n, func(t *testing.T) {
			require.Equal(t, p.expect, p.got.String(), "reference is stable")
		})
	}
}
