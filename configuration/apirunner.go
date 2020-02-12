// Copyright 2020 Insolar Network Ltd.
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

package configuration

import (
	"fmt"
)

// APIRunner holds configuration for api
type APIRunner struct {
	Address string
	RPC     string
	// IsAdmin indicates status of api (internal or external)
	IsAdmin     bool
	SwaggerPath string
}

// NewAPIRunner creates new api config
func NewAPIRunner(admin bool) APIRunner {
	if admin {
		return APIRunner{
			Address:     "localhost:19001",
			RPC:         "/admin-api/rpc",
			SwaggerPath: "application/api/spec/api-exported.yaml",
			IsAdmin:     true,
		}
	}
	return APIRunner{
		Address:     "localhost:19101",
		RPC:         "/api/rpc",
		SwaggerPath: "application/api/spec/api-exported.yaml",
		IsAdmin:     false,
	}
}

func (ar *APIRunner) String() string {
	res := fmt.Sprintf("Addr -> %s, RPC -> %s, IsAdmin -> %t, SwaggerPath -> %s\n", ar.Address, ar.RPC, ar.IsAdmin, ar.SwaggerPath)
	return res
}
