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

package api

import (
	"net/http"

	"github.com/insolar/insolar/api/requester"
	"github.com/insolar/rpc/v2"
)

// AdminContractService is a service that provides API for working with smart contracts.
type AdminContractService struct {
	runner         *Runner
	allowedMethods map[string]bool
}

// NewAdminContractService creates new AdminContract service instance.
func NewAdminContractService(runner *Runner) *AdminContractService {
	methods := map[string]bool{
		"migration.getInfo":          true,
		"migration.deactivateDaemon": true,
		"migration.activateDaemon":   true,
		"migration.checkDaemon":      true,
		"migration.addAddresses":     true,
		"deposit.migration":          true,
		"member.getBalance":          true,
		"contract.registerNode":      true,
		"contract.getNodeRef":        true,
		"CreateHelloWorld":           true,
		"Greet":                      true,
		"Count":                      true,
		"CreateChild":                true,
		"ReturnObj":                  true,
		"PulseNumber":                true,
	}
	return &AdminContractService{runner: runner, allowedMethods: methods}
}

func (cs *AdminContractService) Call(req *http.Request, args *requester.Params, requestBody *rpc.RequestBody, result *requester.ContractResult) error {
	return wrapCall(cs.runner, cs.allowedMethods, req, args, requestBody, result)
}
