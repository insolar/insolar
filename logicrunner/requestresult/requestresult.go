// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
