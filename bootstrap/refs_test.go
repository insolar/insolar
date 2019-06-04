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
		insolar.GenesisNameMember: {
			got:    ContractRootMember,
			expect: "1tJCspdy6s9Ve9vASgdmgnRzdZk3xskhqBcfrLNCGr.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
		},
		insolar.GenesisNameWallet: {
			got:    ContractRootWallet,
			expect: "1tJD4tvUdCMzKFxXUu6Q8AXB2Ubdy4PVqJ4G7TZEKt.1tJDJLGWcX3TCXZMzZodTYWZyJGVdsajgGqyq8Vidw",
		},
	}

	for n, p := range pairs {
		t.Run(n, func(t *testing.T) {
			require.Equal(t, p.expect, p.got.String(), "reference is stable")
		})
	}
}
