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

package rootdomain

import (
	"encoding/hex"
	"testing"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/platformpolicy"
	"github.com/stretchr/testify/require"
)

var (
	idHex  = "0001000149cb5c5f47e921162ecea421c82ca9af969484df62c991271f0b4868"
	refHex = idHex + "0000000000000000000000000000000000000000000000000000000000000000"
)

func TestID(t *testing.T) {
	rootRecord := &Record{
		PCS: initPCS(),
	}
	require.Equal(t, idHex, hex.EncodeToString(rootRecord.ID().Bytes()), "root domain ID should always be the same")
}

func TestReference(t *testing.T) {
	rootRecord := &Record{
		PCS: initPCS(),
	}
	require.Equal(t, refHex, hex.EncodeToString(rootRecord.Ref().Bytes()), "root domain Ref should always be the same")

}

func TestGenesisRef(t *testing.T) {
	var (
		pubKey    = "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEf+vsMVU75xH8uj5WRcOqYdHXtaHH\nN0na2RVQ1xbhsVybYPae3ujNHeQCPj+RaJyMVhb6Aj/AOsTTOPFswwIDAQ==\n-----END PUBLIC KEY-----\n"
		pubKeyRef = "1tJCoBmGrM8TbHkX1bCUAxnAYAUdScRWZ53a8EFLVG.11111111111111111111111111111111"
	)
	genesisRef := GenesisRef(pubKey)
	require.Equal(t, pubKeyRef, genesisRef.String(), "reference by name always the same")
}

func initPCS() insolar.PlatformCryptographyScheme {
	return platformpolicy.NewPlatformCryptographyScheme()
}
