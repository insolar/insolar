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

package updateapproves

import (
	"github.com/insolar/insolar/application/contract/noderecord"
	"github.com/insolar/insolar/logicrunner/goplugin/foundation"
)

type ApproveResult byte

const (
	TypeUnknown ApproveResult = iota
	TypeAgree
	TypeDisagree
)

type UpdateApproves struct {
	foundation.BaseContract
	NodeRec   *noderecord.NodeRecord
	Result    ApproveResult
	Signature []byte
}

func New(nodeRec *noderecord.NodeRecord) *UpdateApproves {
	return &UpdateApproves{
		NodeRec:   nodeRec,
		Result:    TypeUnknown,
		Signature: nil,
	}
}

func (ua *UpdateApproves) Register(signature []byte) *foundation.Error {
	approveResult, err := ua.Verify(ua.NodeRec, signature)
	if err != nil {
		return &foundation.Error{S: err.Error()}
	}
	ua.Result = approveResult
	ua.Signature = signature
	return nil
}

func (ua *UpdateApproves) GetApproveResult() *ApproveResult {
	return &ua.Result
}

func (ua *UpdateApproves) Verify(nodeRec *noderecord.NodeRecord, signature []byte) (ApproveResult, error) {
	//
	//
	//
	return TypeUnknown, nil
}
