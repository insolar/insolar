package genesisrefs

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
)

func TestReferences(t *testing.T) {
	pairs := map[string]struct {
		got    insolar.Reference
		expect string
	}{
		GenesisNameNodeDomain: {
			got:    ContractNodeDomain,
			expect: "insolar:1AAEAAUdxJyWoY-IjQatMYpOk51MZx9tEThkqd1dSB1U",
		},
		GenesisNameNodeRecord: {
			got:    ContractNodeRecord,
			expect: "insolar:1AAEAAQy4dc1JKDJGNd5YfU7ow3DFrW_9j7v772siVMQ",
		},
	}

	for n, p := range pairs {
		t.Run(n, func(t *testing.T) {
			require.Equal(t, p.expect, p.got.String(), "reference is stable")
		})
	}
}

func TestGenesisRef(t *testing.T) {
	var (
		pubKey    = "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEf+vsMVU75xH8uj5WRcOqYdHXtaHH\nN0na2RVQ1xbhsVybYPae3ujNHeQCPj+RaJyMVhb6Aj/AOsTTOPFswwIDAQ==\n-----END PUBLIC KEY-----\n"
		pubKeyRef = "insolar:1AAEAAcEp7HwQByGOr6rZwkyiRA3wR2POYCrDIhqBJyY"
	)
	genesisRef := GenesisRef(pubKey)
	require.Equal(t, pubKeyRef, genesisRef.String(), "reference by name always the same")
}
