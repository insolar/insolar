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

package logicexecutor

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

// ensure artifacts interface
var _ artifacts.RequestResult = &RequestResult{}

type RequestResult struct {
	SideEffectType     artifacts.RequestResultType // every
	ResultData         []byte                      // every
	ObjectReferenceRef insolar.Reference           // every

	parentReference insolar.Reference // activate
	objectImage     insolar.Reference // amend + activate
	objectStateID   insolar.ID        // amend + deactivate
	Memory          []byte            // amend + activate
}

func NewRequestResult(result []byte, objectRef insolar.Reference) *RequestResult {
	return &RequestResult{
		SideEffectType:     artifacts.RequestSideEffectNone,
		ResultData:         result,
		ObjectReferenceRef: objectRef,
	}
}

func (s *RequestResult) Result() []byte {
	return s.ResultData
}

func (s *RequestResult) Activate() (insolar.Reference, insolar.Reference, []byte) {
	return s.parentReference, s.objectImage, s.Memory
}

func (s *RequestResult) Amend() (insolar.ID, insolar.Reference, []byte) {
	return s.objectStateID, s.objectImage, s.Memory
}

func (s *RequestResult) Deactivate() insolar.ID {
	return s.objectStateID
}

func (s *RequestResult) SetActivate(parent, image insolar.Reference, memory []byte) {
	s.SideEffectType = artifacts.RequestSideEffectActivate

	s.parentReference = parent
	s.objectImage = image
	s.Memory = memory
}

func (s *RequestResult) SetAmend(object artifacts.ObjectDescriptor, memory []byte) {
	s.SideEffectType = artifacts.RequestSideEffectAmend
	s.Memory = memory
	s.objectStateID = *object.StateID()

	prototype, _ := object.Prototype()
	s.objectImage = *prototype
}

func (s *RequestResult) SetDeactivate(object artifacts.ObjectDescriptor) {
	s.SideEffectType = artifacts.RequestSideEffectDeactivate
	s.objectStateID = *object.StateID()
}

func (s RequestResult) Type() artifacts.RequestResultType {
	return s.SideEffectType
}

func (s *RequestResult) ObjectReference() insolar.Reference {
	return s.ObjectReferenceRef
}
