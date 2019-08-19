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

package configuration

import (
	"fmt"
)

// APIRunner holds configuration for api
type APIRunner struct {
	Address  string
	RPC      string
	IsPublic bool
}

// NewAPIRunner creates new api config
func NewAPIRunner(isInternal bool) APIRunner {
	if isInternal {
		return APIRunner{
			Address:  "localhost:19100",
			RPC:      "/admin-api/rpc",
			IsPublic: false,
		}
	}
	return APIRunner{
		Address:  "localhost:19101",
		RPC:      "/api/rpc",
		IsPublic: true,
	}
}

func (ar *APIRunner) String() string {
	res := fmt.Sprintln("Addr ->", ar.Address, ", RPC ->", ar.RPC, ", IsPublic ->", ar.IsPublic)
	return res
}
