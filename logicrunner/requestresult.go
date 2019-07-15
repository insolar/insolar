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

package logicrunner

import (
	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/artifacts"
)

type requestResult struct {
	sideEffectType  artifacts.RequestResultType // every
	result          []byte                      // every
	objectReference *insolar.Reference          // every

	asDelegate      bool               // activate
	parentReference *insolar.Reference // activate
	objectImage     *insolar.Reference // amend + activate
	objectStateID   *insolar.ID        // amend + deactivate
	memory          []byte             // amend + activate
}

func newRequestResult(result []byte) *requestResult {
	return &requestResult{result: result}
}

func (s *requestResult) Result() []byte {
	return s.result
}

func (s *requestResult) Activate() (*insolar.Reference, *insolar.Reference, bool, []byte) {
	return s.parentReference, s.objectImage, s.asDelegate, s.memory
}

func (s *requestResult) Amend() (*insolar.ID, *insolar.Reference, []byte) {
	return s.objectStateID, s.objectImage, s.memory
}

func (s *requestResult) Deactivate() *insolar.ID {
	return s.objectStateID
}

func (s *requestResult) SetActivate(parent, image *insolar.Reference, asDelegate bool, memory []byte) {
	s.sideEffectType = artifacts.RequestSideEffectActivate
	s.objectReference = parent

	s.asDelegate = asDelegate
	s.parentReference = parent
	s.objectImage = image
	s.memory = memory

}

func (s *requestResult) SetAmend(object artifacts.ObjectDescriptor, memory []byte) error {
	s.sideEffectType = artifacts.RequestSideEffectAmend
	s.memory = memory
	s.objectReference = object.HeadRef()
	s.objectStateID = object.StateID()

	if object.IsPrototype() {
		return errors.New("Can't update prototype")
	}

	prototype, err := object.Prototype()
	if err != nil {
		return errors.Wrap(err, "Failed to obtain prototype/code of object")
	}

	s.objectImage = prototype
	return nil
}

func (s *requestResult) SetDeactivate(object artifacts.ObjectDescriptor) {
	s.sideEffectType = artifacts.RequestSideEffectDeactivate
	s.objectReference = object.HeadRef()
	s.objectStateID = object.StateID()
}

func (s *requestResult) SetNone(objectRef *insolar.Reference) {
	s.sideEffectType = artifacts.RequestSideEffectNone
	s.objectReference = objectRef
}

func (s requestResult) Type() artifacts.RequestResultType {
	return s.sideEffectType
}

func (s *requestResult) ObjectReference() *insolar.Reference {
	return s.objectReference
}
