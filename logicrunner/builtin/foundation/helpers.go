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

// GetPulseNumber returns current pulse from context.
func GetPulseNumber() (insolar.PulseNumber, error) {
	req := GetLogicalContext().Request
	if req == nil {
		return insolar.PulseNumber(0), errors.New("request from LogicCallContext is nil, get pulse is failed")
	}
	return req.Record().Pulse(), nil
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

// Get reference on  MigrationAdmin contract.
func GetMigrationAdmin() insolar.Reference {
	return genesisrefs.ContractMigrationAdmin
}

// Get reference on  RootDomain contract.
func GetRootDomain() insolar.Reference {
	return genesisrefs.ContractRootDomain
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
