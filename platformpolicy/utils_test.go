package platformpolicy

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeys_publicKeyNormalize(t *testing.T) {
	var (
		begin   = "-----BEGIN PUBLIC KEY-----\n"
		end     = "-----END PUBLIC KEY-----\n"
		pubKey1 = begin + "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEf+vsMVU75xH8uj5WRcOqYdHXtaHH\nN0na2RVQ1xbhsVybYPae3ujNHeQCPj+RaJyMVhb6Aj/AOsTTOPFswwIDAQ==\n" + end
		pubKey2 = begin + "\n" + "MFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEf+vsMVU75xH8uj5WRcOqYdHXtaHH\nN0na2RVQ1xbhsVybYPae3ujNHeQCPj+RaJyMVhb6Aj/AOsTTOPFswwIDAQ==\n" + end
	)

	s1 := MustNormalizePublicKey([]byte(pubKey1))
	s2 := MustNormalizePublicKey([]byte(pubKey2))
	require.Equal(t, s1, s2, "the same result for the same public key")
}
