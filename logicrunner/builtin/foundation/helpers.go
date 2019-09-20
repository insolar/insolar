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

package foundation

import (
	"errors"
	"strings"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/genesisrefs"
)

const EthAddressLength = 20

// GetPulseNumber returns current pulse from context.
func GetPulseNumber() (insolar.PulseNumber, error) {
	req := GetLogicalContext().Request
	if req == nil {
		return insolar.PulseNumber(0), errors.New("request from LogicCallContext is nil, get pulse is failed")
	}
	return req.GetLocal().Pulse(), nil
}

// GetRequestReference - Returns request reference from context.
func GetRequestReference() insolar.Reference {
	ctx := GetLogicalContext()
	if ctx.Request == nil {
		panic("context has no request set")
	}
	return *ctx.Request
}

// GetObject create proxy by address
// unimplemented
func GetObject(ref insolar.Reference) ProxyInterface {
	panic("not implemented")
}

// Get reference CostCenter contract.
func GetCostCenter() insolar.Reference {
	return genesisrefs.ContractCostCenter
}

// Get reference MigrationAdminMember contract.
func GetMigrationAdminMember() insolar.Reference {
	return genesisrefs.ContractMigrationAdminMember
}

// Get reference RootMember contract.
func GetRootMember() insolar.Reference {
	return genesisrefs.ContractRootMember
}

// Get reference on MigrationAdmin contract.
func GetMigrationAdmin() insolar.Reference {
	return genesisrefs.ContractMigrationAdmin
}

// Get reference on RootDomain contract.
func GetRootDomain() insolar.Reference {
	return genesisrefs.ContractRootDomain
}

// Get reference on  migrationdaemon contract by  migration member.
func GetMigrationDaemon(migrationMember insolar.Reference) (insolar.Reference, error) {
	return genesisrefs.ContractMigrationMap[migrationMember], nil
}

// Check member is migration daemon member or not
func IsMigrationDaemonMember(member insolar.Reference) bool {
	for _, mDaemonMember := range genesisrefs.ContractMigrationDaemonMembers {
		if mDaemonMember.Equal(member) {
			return true
		}
	}
	return false
}

// TrimPublicKey trim public key
func TrimPublicKey(publicKey string) string {
	return TrimAddress(between(publicKey, "KEY-----", "-----END"))
}

// TrimPublicKey trim address
func TrimAddress(address string) string {
	return strings.ToLower(strings.Join(strings.Split(strings.TrimSpace(address), "\n"), ""))
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	pos := strings.Index(value, a)
	if pos == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirst := pos + len(a)
	if posFirst >= posLast {
		return ""
	}
	return value[posFirst:posLast]
}

// IsEthereumAddress Ethereum address format verifier
func IsEthereumAddress(s string) bool {
	if hasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*EthAddressLength && isHex(s)
}

func hasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && (str[1] == 'x' || str[1] == 'X')
}

func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

func isHex(str string) bool {
	if len(str)%2 != 0 {
		return false
	}
	for _, c := range []byte(str) {
		if !isHexCharacter(c) {
			return false
		}
	}
	return true
}
