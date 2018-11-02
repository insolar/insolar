/*
 *    Copyright 2018 Insolar
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package auth

import (
	"github.com/insolar/insolar/network/controller/common"
	"github.com/insolar/insolar/network/hostnetwork"
	"github.com/pkg/errors"
)

//
type AuthorizationController struct {
	bootstrapController common.BootstrapController
	signer              *Signer
	transport           hostnetwork.InternalTransport
}

func (ac *AuthorizationController) Authorize() error {
	hosts := ac.bootstrapController.GetBootstrapHosts()
	if len(hosts) == 0 {
		return errors.New("Empty list of bootstrap hosts")
	}
	return nil
}

func NewAuthorizationController(bootstrapController common.BootstrapController,
	transport hostnetwork.InternalTransport) *AuthorizationController {
	return &AuthorizationController{bootstrapController: bootstrapController, transport: transport}
}
