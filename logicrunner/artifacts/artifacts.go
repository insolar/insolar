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

package artifacts

import (
	"context"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/payload"
	"github.com/insolar/insolar/insolar/record"
)

//go:generate minimock -i github.com/insolar/insolar/logicrunner/artifacts.Client -o ./ -s _mock.go -g

// Client is a high level storage interface.
type Client interface {
	// RegisterIncomingRequest creates an incoming request record in storage.
	RegisterIncomingRequest(ctx context.Context, request *record.IncomingRequest) (*payload.RequestInfo, error)
	// RegisterIncomingRequest creates an outgoing request record in storage.
	RegisterOutgoingRequest(ctx context.Context, request *record.OutgoingRequest) (*payload.RequestInfo, error)

	// RegisterResult saves VM method call result and side-effect
	RegisterResult(ctx context.Context, request insolar.Reference, result RequestResult) error

	// GetAbandonedRequest returns an incoming or outgoing request for an object.
	GetAbandonedRequest(ctx context.Context, objectRef, reqRef insolar.Reference) (record.Request, error)

	// GetPendings returns pending request IDs of an object.
	GetPendings(ctx context.Context, objectRef insolar.Reference) ([]insolar.Reference, error)

	// HasPendings returns true if object has unclosed requests.
	HasPendings(ctx context.Context, object insolar.Reference) (bool, error)

	// GetCode returns code from code record by provided reference according to provided machine preference.
	//
	// This method is used by VM to fetch code for execution.
	GetCode(ctx context.Context, ref insolar.Reference) (CodeDescriptor, error)

	// GetObject returns descriptor for provided state.
	//
	// If provided state is nil, the latest state will be returned (with deactivation check). Returned descriptor will
	// provide methods for fetching all related data.
	GetObject(ctx context.Context, head insolar.Reference) (ObjectDescriptor, error)

	// DeployCode creates new code record in storage.
	//
	// Code records are used to activate prototype.
	DeployCode(ctx context.Context, domain, request insolar.Reference, code []byte, machineType insolar.MachineType) (*insolar.ID, error)

	// ActivatePrototype creates activate object record in storage. Provided prototype reference will be used as objects prototype
	// memory as memory of created object. If memory is not provided, the prototype default memory will be used.
	//
	// Request reference will be this object's identifier and referred as "object head".
	ActivatePrototype(
		ctx context.Context,
		request, parent, code insolar.Reference,
		memory []byte,
	) error

	// State returns hash state for artifact manager.
	State() []byte

	// InjectCodeDescriptor injects code descriptor needed by builtin contracts
	InjectCodeDescriptor(insolar.Reference, CodeDescriptor)
	// InjectObjectDescriptor injects object descriptor needed by builtin contracts (to store prototypes)
	InjectObjectDescriptor(insolar.Reference, ObjectDescriptor)
	// InjectFinish finalizes all injects, all next injects will panic
	InjectFinish()
}

//go:generate minimock -i github.com/insolar/insolar/logicrunner/artifacts.CodeDescriptor -o ./ -s _mock.go -g

// CodeDescriptor represents meta info required to fetch all code data.
type CodeDescriptor interface {
	// Ref returns reference to represented code record.
	Ref() *insolar.Reference

	// MachineType returns code machine type for represented code.
	MachineType() insolar.MachineType

	// Code returns code data.
	Code() ([]byte, error)
}

//go:generate minimock -i github.com/insolar/insolar/logicrunner/artifacts.ObjectDescriptor -o ./ -s _mock.go -g

// ObjectDescriptor represents meta info required to fetch all object data.
type ObjectDescriptor interface {
	// HeadRef returns head reference to represented object record.
	HeadRef() *insolar.Reference

	// StateID returns reference to object state record.
	StateID() *insolar.ID

	// Memory fetches object memory from storage.
	Memory() []byte

	// IsPrototype determines if the object is a prototype.
	IsPrototype() bool

	// Code returns code reference.
	Code() (*insolar.Reference, error)

	// Prototype returns prototype reference.
	Prototype() (*insolar.Reference, error)

	// ChildPointer returns the latest child for this object.
	ChildPointer() *insolar.ID

	// Parent returns object's parent.
	Parent() *insolar.Reference
}

// RefIterator is used for iteration over affined children(parts) of container.
type RefIterator interface {
	Next() (*insolar.Reference, error)
	HasNext() bool
}

//go:generate minimock -i github.com/insolar/insolar/logicrunner/artifacts.DescriptorsCache -o ./ -s _mock.go -g

// DescriptorsCache provides convenient way to get prototype and code descriptors
// of objects without fetching them twice
type DescriptorsCache interface {
	ByPrototypeRef(ctx context.Context, protoRef insolar.Reference) (ObjectDescriptor, CodeDescriptor, error)
	ByObjectDescriptor(ctx context.Context, obj ObjectDescriptor) (ObjectDescriptor, CodeDescriptor, error)
	GetPrototype(ctx context.Context, ref insolar.Reference) (ObjectDescriptor, error)
	GetCode(ctx context.Context, ref insolar.Reference) (CodeDescriptor, error)
}

type RequestResultType uint8

const (
	RequestSideEffectNone RequestResultType = iota
	RequestSideEffectActivate
	RequestSideEffectAmend
	RequestSideEffectDeactivate
)

func (t RequestResultType) String() string {
	switch t {
	case RequestSideEffectNone:
		return "None"
	case RequestSideEffectActivate:
		return "Activate"
	case RequestSideEffectAmend:
		return "Amend"
	case RequestSideEffectDeactivate:
		return "Deactivate"
	default:
		return "Unknown"
	}
}

type RequestResult interface {
	Type() RequestResultType

	Activate() (insolar.Reference, insolar.Reference, []byte)
	Amend() (insolar.ID, insolar.Reference, []byte)
	Deactivate() insolar.ID

	Result() []byte
	ObjectReference() insolar.Reference
}
