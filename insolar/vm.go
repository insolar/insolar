// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package insolar

import (
	"context"
)

// MachineType is a type of virtual machine
type MachineType int

// Real constants of MachineType
const (
	MachineTypeNotExist             = 0
	MachineTypeBuiltin  MachineType = iota + 1

	MachineTypesLastID
)

func (m MachineType) Equal(other MachineType) bool {
	return m == other
}

//go:generate minimock -i github.com/insolar/insolar/insolar.MachineLogicExecutor -o ../testutils -s _mock.go -g

// MachineLogicExecutor is an interface for implementers of one particular machine type
type MachineLogicExecutor interface {
	CallMethod(
		ctx context.Context, callContext *LogicCallContext,
		code Reference, data []byte,
		method string, args Arguments,
	) (
		newObjectState []byte, methodResults Arguments, err error,
	)
	CallConstructor(
		ctx context.Context, callContext *LogicCallContext,
		code Reference, name string, args Arguments,
	) (
		objectState []byte, result Arguments, err error,
	)
}

//go:generate minimock -i github.com/insolar/insolar/insolar.LogicRunner -o ../testutils -s _mock.go -g

// LogicRunner is an interface that should satisfy logic executor
type LogicRunner interface {
	LRI()
	OnPulse(context.Context, Pulse, Pulse) error
	AddUnwantedResponse(ctx context.Context, msg Payload) error
}

// CallMode indicates whether we execute or validate
type CallMode int

const (
	ExecuteCallMode CallMode = iota
	ValidateCallMode
)

func (m CallMode) String() string {
	switch m {
	case ExecuteCallMode:
		return "execute"
	case ValidateCallMode:
		return "validate"
	default:
		return "unknown"
	}
}

// LogicCallContext is a context of contract execution. Everything
// that is required to implement foundation functions. This struct
// shouldn't be used in core components.
type LogicCallContext struct {
	Mode CallMode // either "execution" or "validation"

	Request *Reference // reference of incoming request record

	Callee    *Reference // Contract that is called
	Parent    *Reference // Parent of the callee
	Prototype *Reference // Prototype (base class) of the callee
	Code      *Reference // Code reference of the callee

	Caller          *Reference // Contract that made the call
	CallerPrototype *Reference // Prototype (base class) of the caller

	TraceID string // trace mark for Jaeger and friends
	Pulse   Pulse  // prefetched pulse for call context
}

// ContractConstructor is a typedef for wrapper contract header
type ContractMethod func(oldState []byte, args []byte) (newState []byte, result []byte, err error)

// ContractMethods maps name to contract method
type ContractMethods map[string]ContractMethod

// ContractConstructor is a typedef of typical contract constructor
type ContractConstructor func(ref Reference, args []byte) (state []byte, result []byte, err error)

// ContractConstructors maps name to contract constructor
type ContractConstructors map[string]ContractConstructor

// ContractWrapper stores all needed about contract wrapper (it's methods/constructors)
type ContractWrapper struct {
	Methods      ContractMethods
	Constructors ContractConstructors
}

//go:generate stringer -type=PendingState

// PendingState is a state of execution for each object
type PendingState int

const (
	PendingUnknown PendingState = iota // PendingUnknown signalizes that we don't know about execution state
	NotPending                         // NotPending means that we know that this task is not executed by another VE
	InPending                          // InPending means that we know that method on object is executed by another VE
)

func (s PendingState) Equal(other PendingState) bool {
	return s == other
}
