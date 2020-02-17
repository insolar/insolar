// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package record

import (
	"github.com/insolar/insolar/insolar"

	"github.com/pkg/errors"
)

type Record interface {
	Marshal() (dAtA []byte, err error)
}

// StateID is a state of lifeline records.
type StateID int

const (
	// StateUndefined is used for special cases.
	StateUndefined = StateID(iota)
	// StateActivation means it's an activation record.
	StateActivation
	// StateAmend means it's an amend record.
	StateAmend
	// StateDeactivation means it's a deactivation record.
	StateDeactivation
)

func (s *StateID) Equal(other StateID) bool {
	return *s == other
}

// State is common object state record.
type State interface {
	Record
	// ID returns state id.
	ID() StateID
	// GetImage returns state code.
	GetImage() *insolar.Reference
	// GetIsPrototype returns state code.
	GetIsPrototype() bool
	// GetMemory returns state indexStorage.
	GetMemory() []byte
	// PrevStateID returns previous state id.
	PrevStateID() *insolar.ID
}

func (Activate) ID() StateID {
	return StateActivation
}

func (p Activate) GetImage() *insolar.Reference {
	return &p.Image
}

func (p Activate) GetIsPrototype() bool {
	return p.IsPrototype
}

func (p Activate) GetMemory() []byte {
	return p.Memory
}

func (Activate) PrevStateID() *insolar.ID {
	return nil
}

func (Amend) ID() StateID {
	return StateAmend
}

func (p Amend) GetImage() *insolar.Reference {
	return &p.Image
}

func (p Amend) GetIsPrototype() bool {
	return p.IsPrototype
}

func (p Amend) GetMemory() []byte {
	return p.Memory
}

func (p Amend) PrevStateID() *insolar.ID {
	return &p.PrevState
}

func (Deactivate) ID() StateID {
	return StateDeactivation
}

func (Deactivate) GetImage() *insolar.Reference {
	return nil
}

func (Deactivate) GetIsPrototype() bool {
	return false
}

func (Deactivate) GetMemory() []byte {
	return nil
}

func (p Deactivate) PrevStateID() *insolar.ID {
	return &p.PrevState
}

func (Genesis) PrevStateID() *insolar.ID {
	return nil
}

func (Genesis) ID() StateID {
	return StateActivation
}

func (Genesis) GetMemory() []byte {
	return nil
}

func (Genesis) GetImage() *insolar.Reference {
	return nil
}

func (Genesis) GetIsPrototype() bool {
	return false
}

//go:generate minimock -i github.com/insolar/insolar/insolar/record.Request -o ./ -s _mock.go -g

// Request is a common request interface.
type Request interface {
	Record
	// AffinityRef returns a pointer to the reference of the object the
	// Request is affine to. The result can be nil, e.g. in case of creating
	// a new object.
	AffinityRef() *insolar.Reference
	// ReasonRef returns a reference of the Request that caused the creating
	// of this Request.
	ReasonRef() insolar.Reference
	// ReasonAffinityRef returns a reference of an object reason request is
	// affine to.
	ReasonAffinityRef() insolar.Reference
	// GetCallType returns call type.
	GetCallType() CallType
	// IsAPIRequest tells is it API-request or not.
	IsAPIRequest() bool
	// IsCreationRequest checks a request-type.
	IsCreationRequest() bool
	// Validate validates request params and its combinations.
	Validate() error
	// IsTemporaryUploadCode tells us that that request is temporary hack
	// for uploading code.
	IsTemporaryUploadCode() bool
}

func (r *IncomingRequest) AffinityRef() *insolar.Reference {
	// IncomingRequests are affine to the Object on which the request
	// is going to be executed.
	// Exceptions are CTSaveAsMethod, we should
	// calculate hash of message, so call CalculateRequestAffinityRef.
	if r.IsCreationRequest() {
		return nil
	}
	return r.Object
}

func (r *IncomingRequest) ReasonRef() insolar.Reference {
	return r.Reason
}

func (r *IncomingRequest) ReasonAffinityRef() insolar.Reference {
	return r.Caller
}

func (r *IncomingRequest) IsAPIRequest() bool {
	return !r.APINode.IsEmpty()
}

func (r *IncomingRequest) IsCreationRequest() bool {
	return r.GetCallType() == CTSaveAsChild || r.GetCallType() == CTDeployPrototype
}

func (r *IncomingRequest) Validate() error {
	if r.ReasonRef().GetLocal().IsEmpty() {
		return errors.New("reason is empty")
	}
	// Incoming requests never should't be in detached state,
	// app code should check it and raise some kind of error.
	if r.IsAPIRequest() {
		return nil
	}
	if r.ReasonAffinityRef().IsEmpty() {
		return errors.New("reason object is not set on incoming request")
	}
	return nil
}

func (r *IncomingRequest) IsDetachedCall() bool {
	return r.ReturnMode == ReturnSaga
}

func (r *IncomingRequest) IsTemporaryUploadCode() bool {
	return r.GetCallType() == CTDeployPrototype
}

func (r *OutgoingRequest) AffinityRef() *insolar.Reference {
	// OutgoingRequests are affine to the Caller which created the Request.
	return &r.Caller
}

func (r *OutgoingRequest) ReasonRef() insolar.Reference {
	return r.Reason
}

func (r *OutgoingRequest) ReasonAffinityRef() insolar.Reference {
	return r.Caller
}

func (r *OutgoingRequest) IsAPIRequest() bool {
	return false
}

func (r *OutgoingRequest) IsCreationRequest() bool {
	return false
}

func (r *OutgoingRequest) IsDetached() bool {
	return r.ReturnMode == ReturnSaga
}

func (r *OutgoingRequest) Validate() error {
	if r.IsCreationRequest() {
		return errors.New("outgoing request cannot be creating request")
	}
	if r.ReasonRef().GetLocal().IsEmpty() {
		return errors.New("reason is empty")
	}

	return nil
}

func (r *OutgoingRequest) IsTemporaryUploadCode() bool {
	return false
}

func CalculateRequestAffinityRef(
	request Request,
	pulseNumber insolar.PulseNumber,
	scheme insolar.PlatformCryptographyScheme,
) *insolar.Reference {
	affinityRef := request.AffinityRef()
	if affinityRef == nil {
		virtualRecord := Wrap(request)
		hash := HashVirtual(scheme.ReferenceHasher(), virtualRecord)
		recID := insolar.NewID(pulseNumber, hash)
		affinityRef = insolar.NewReference(*recID)
	}
	return affinityRef
}

// ObjectIDFromRequest calculates object is from request.
func ObjectIDFromRequest(cs insolar.PlatformCryptographyScheme, request Request, requestID insolar.ID) (insolar.ID, error) {
	if !request.IsCreationRequest() {
		if request.AffinityRef() == nil {
			return insolar.ID{}, errors.New("affinity ref is empty")
		}
		return *request.AffinityRef().GetLocal(), nil
	}
	virtual := Wrap(request)
	buf, err := virtual.Marshal()
	if err != nil {
		return insolar.ID{}, err
	}
	hasher := cs.ReferenceHasher()
	_, err = hasher.Write(buf)
	if err != nil {
		return insolar.ID{}, errors.Wrap(err, "failed to calculate id")
	}
	return *insolar.NewID(requestID.Pulse(), hasher.Sum(nil)), nil
}
