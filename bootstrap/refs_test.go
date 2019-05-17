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
			expect: "1tJEEHUZYWYXXEa1KCTKKR8G6MrteZqCFPdcejHinR.1tJEEHUZYWYXXEa1KCTKKR8G6MrteZqCFPdcejHinR",
		},
		insolar.GenesisNameNodeDomain: {
			got:    ContractNodeDomain,
			expect: "1tJCNCWXKnkKxE8DuNnquDXN4uQXw9UxWopL41ME8Y.1tJEEHUZYWYXXEa1KCTKKR8G6MrteZqCFPdcejHinR",
		},
		insolar.GenesisNameNodeRecord: {
			got:    ContractNodeRecord,
			expect: "1tJBvdpcKPQpcnpxWTYyDyqYrizXJrHsTGDLvwK39P.1tJEEHUZYWYXXEa1KCTKKR8G6MrteZqCFPdcejHinR",
		},
		insolar.GenesisNameRootMember: {
			got:    ContractRootMember,
			expect: "1tJCdDBJ2dsjB9QdJddvYNsb2vamoagc3B2CNVrgr3.1tJEEHUZYWYXXEa1KCTKKR8G6MrteZqCFPdcejHinR",
		},
		insolar.GenesisNameWallet: {
			got:    ContractWallet,
			expect: "1tJCTxgPKZuh7K8nSy8Z4L7UUtLrS4ZwQrbFgBELoa.1tJEEHUZYWYXXEa1KCTKKR8G6MrteZqCFPdcejHinR",
		},
		insolar.GenesisNameAllowance: {
			got:    ContractAllowance,
			expect: "1tJDCCNyX2kMvt8t2x4bPkhzYrZ3aWBi9bKzNL3SXe.1tJEEHUZYWYXXEa1KCTKKR8G6MrteZqCFPdcejHinR",
		},
	}

	for n, p := range pairs {
		t.Run(n, func(t *testing.T) {
			require.Equal(t, p.expect, p.got.String(), "reference is stable")
		})
	}
}
