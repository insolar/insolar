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

package updatepackage

import (
	"github.com/insolar/insolar/application/contract/updateapproves"
	approves "github.com/insolar/insolar/application/proxy/updateapproves"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
	"github.com/insolar/insolar/updater/request"
)

type UpdatePackage struct {
	foundation.BaseContract
	UpdateVersion       *request.Version
	consensusCountNodes int
}

func (up *UpdatePackage) AddResult(result updateapproves.ApproveResult, node *core.RecordRef, sig []byte) core.RecordRef {
	ah := approves.New(node, result, sig)
	a := ah.AsChild(up.GetReference())
	return a.GetReference()
}

func (up *UpdatePackage) IsConsensus() bool {
	totalCount := up.getTotalAgrees()
	if up.consensusCountNodes <= totalCount {
		return true
	}
	return false
}

func (up *UpdatePackage) getTotalAgrees() int {
	crefs, err := up.GetChildrenTyped(approves.GetClass())
	if err != nil {
		panic(err)
	}
	totalCount := 0
	for _, cref := range crefs {
		obj := approves.GetObject(cref)
		if obj.GetApproveResult() == updateapproves.TypeAgree {
			totalCount++
		}
	}
	return totalCount
}

func New(updateVer *request.Version, consensusCountNodes int) *UpdatePackage {
	return &UpdatePackage{
		UpdateVersion:       updateVer,
		consensusCountNodes: consensusCountNodes,
	}
}
