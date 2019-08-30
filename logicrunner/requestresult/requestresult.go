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

package requestresult

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

type RequestResult struct {
	SideEffectType     artifacts.RequestResultType // every
	RawResult          []byte                      // every
	RawObjectReference insolar.Reference           // every

	ParentReference insolar.Reference // activate
	ObjectImage     insolar.Reference // amend + activate
	ObjectStateID   insolar.ID        // amend + deactivate
	Memory          []byte            // amend + activate
}

func New(result []byte, objectRef insolar.Reference) *RequestResult {
	return &RequestResult{
		SideEffectType:     artifacts.RequestSideEffectNone,
		RawResult:          result,
		RawObjectReference: objectRef,
	}
}

func (s *RequestResult) Result() []byte {
	return s.RawResult
}

func (s *RequestResult) Activate() (insolar.Reference, insolar.Reference, []byte) {
	return s.ParentReference, s.ObjectImage, s.Memory
}

func (s *RequestResult) Amend() (insolar.ID, insolar.Reference, []byte) {
	return s.ObjectStateID, s.ObjectImage, s.Memory
}

func (s *RequestResult) Deactivate() insolar.ID {
	return s.ObjectStateID
}

func (s *RequestResult) SetActivate(parent, image insolar.Reference, memory []byte) {
	s.SideEffectType = artifacts.RequestSideEffectActivate

	s.ParentReference = parent
	s.ObjectImage = image
	s.Memory = memory
}

func (s *RequestResult) SetAmend(object artifacts.ObjectDescriptor, memory []byte) {
	s.SideEffectType = artifacts.RequestSideEffectAmend
	s.Memory = memory
	s.ObjectStateID = *object.StateID()

	prototype, _ := object.Prototype()
	s.ObjectImage = *prototype
}

func (s *RequestResult) SetDeactivate(object artifacts.ObjectDescriptor) {
	s.SideEffectType = artifacts.RequestSideEffectDeactivate
	s.ObjectStateID = *object.StateID()
}

func (s RequestResult) Type() artifacts.RequestResultType {
	return s.SideEffectType
}

func (s *RequestResult) ObjectReference() insolar.Reference {
	return s.RawObjectReference
}
