// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package genesisrefs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/application"
	"github.com/insolar/insolar/applicationbase/genesisrefs"
	"github.com/insolar/insolar/insolar"
)

func TestReferences(t *testing.T) {
	pairs := map[string]struct {
		got    insolar.Reference
		expect string
	}{
		application.GenesisNameRootDomain: {
			got:    ContractRootDomain,
			expect: "insolar:1AAEAAciWtcmPVgAcaIvICkgnSsJmp4Clp650xOHjYks",
		},
		application.GenesisNameRootMember: {
			got:    ContractRootMember,
			expect: "insolar:1AAEAAWeNhA_NwKaH6E36IJ-2PLvXnJRxiTTNWq1giOg",
		},
	}

	for n, p := range pairs {
		t.Run(n, func(t *testing.T) {
			require.Equal(t, p.expect, p.got.String(), "reference is stable")
		})
	}
}

func TestRootDomain(t *testing.T) {
	ref1 := ContractRootDomain
	ref2 := genesisrefs.GenesisRef(application.GenesisNameRootDomain)
	require.Equal(t, ref1.String(), ref2.String(), "reference is the same")
}
