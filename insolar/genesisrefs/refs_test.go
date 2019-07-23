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

package genesisrefs

import (
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/rootdomain"
	"github.com/stretchr/testify/require"
)

func TestReferences(t *testing.T) {
	pairs := map[string]struct {
		got    insolar.Reference
		expect string
	}{
		insolar.GenesisNameRootDomain: {
			got:    ContractRootDomain,
			expect: "1tJCSAhchCSKZhV8cLzWQ6d1NQFWeDDnkG9Ubcbj5R.11111111111111111111111111111111",
		},
		insolar.GenesisNameNodeDomain: {
			got:    ContractNodeDomain,
			expect: "1tJECUCoRtbpMeaHrMkzGVpjbUkAXtCjcW57Noubik.11111111111111111111111111111111",
		},
		insolar.GenesisNameNodeRecord: {
			got:    ContractNodeRecord,
			expect: "1tJE3KYJVRHaSbBx4jeXCrpcUgo1kawvNLPH9Ke3VZ.11111111111111111111111111111111",
		},
		insolar.GenesisNameRootMember: {
			got:    ContractRootMember,
			expect: "1tJE6PFMLHeNmQ5gR8dRt6yAhb8NeUv7sNbsRv49uE.11111111111111111111111111111111",
		},
		insolar.GenesisNameRootWallet: {
			got:    ContractRootWallet,
			expect: "1tJDWFghNhav2zu73YCFpK1b2KTTuqEBY1QXQXzHHc.11111111111111111111111111111111",
		},
		insolar.GenesisNameDeposit: {
			got:    ContractDeposit,
			expect: "1tJDxESz1Vn6aDqzupAofBFQJa6yCrr3jTUFoVrGLm.11111111111111111111111111111111",
		},
		insolar.GenesisNameTariff: {
			got:    ContractStandardTariff,
			expect: "1tJBiXErALbFCemdJeub4VstKjbL5is6Vrq6p5SPUW.11111111111111111111111111111111",
		},
		insolar.GenesisNameCostCenter: {
			got:    ContractCostCenter,
			expect: "1tJDFtT3ALAMWcGRnzMcvfJPPNHQH7B64zFj3xjpnQ.11111111111111111111111111111111",
		},
	}

	for n, p := range pairs {
		t.Run(n, func(t *testing.T) {
			require.Equal(t, p.expect, p.got.String(), "reference is stable")
		})
	}
}

func TestRootDomain(t *testing.T) {
	ref1 := rootdomain.RootDomain.Ref()
	ref2 := rootdomain.GenesisRef(insolar.GenesisNameRootDomain)
	require.Equal(t, ref1.String(), ref2.String(), "reference is the same")
}
