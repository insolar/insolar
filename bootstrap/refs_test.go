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

	"github.com/insolar/insolar/bootstrap/rootdomain"
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
			expect: "1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
		},
		insolar.GenesisNameNodeDomain: {
			got:    ContractNodeDomain,
			expect: "1tJDPXaGq1h3cHM8zSV4HDqZ5LtCz6EeV5bno3cCda.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
		},
		insolar.GenesisNameNodeRecord: {
			got:    ContractNodeRecord,
			expect: "1tJDsdZGUf7xqYxZcsYKJLKkeQMURphBi2xcR7G3dt.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
		},
		insolar.GenesisNameRootMember: {
			got:    ContractRootMember,
			expect: "1tJDL5m9pKyq2mbanYfgwQ5rSQdrpsXbzc1Dk7a53d.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
		},
		insolar.GenesisNameRootWallet: {
			got:    ContractWallet,
			expect: "1tJCvhfv3caM2VuUJmd3pYc467nqNQhPzh8owVHvwY.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
		},
		insolar.GenesisNameAllowance: {
			got:    ContractAllowance,
			expect: "1tJCxMpe8nTqQq38ByCkdg77LtHhfkcTF1teWWtYwi.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
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
