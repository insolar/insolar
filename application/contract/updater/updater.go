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
package updater

import (
	approves "github.com/insolar/insolar/application/contract/updateapproves"
	proxyUP "github.com/insolar/insolar/application/proxy/updatepackage"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/updater/request"
)

type Updater struct {
	foundation.BaseContract
	CurrentVersion *request.Version
}

func (u *Updater) RegisterNewUpdate(ver string, consensusCountNodes int) core.RecordRef {
	up := proxyUP.New(request.NewVersion(ver), map[core.RecordRef]*approves.UpdateApproves{}, consensusCountNodes, 0)
	pk := up.AsChild(u.GetReference())
	return pk.GetReference()
}

func New(curVer *request.Version) *Updater {
	return &Updater{
		CurrentVersion: curVer,
	}
}
