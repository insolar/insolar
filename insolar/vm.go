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
	MachineTypeGoPlugin

	MachineTypesLastID
)

func (m MachineType) Equal(other MachineType) bool {
	return m == other
}

//go:generate minimock -i github.com/insolar/insolar/insolar.MachineLogicExecutor -o ../testutils -s _mock.go

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
		objectState []byte, err error,
	)
}

//go:generate minimock -i github.com/insolar/insolar/insolar.LogicRunner -o ../testutils -s _mock.go

// LogicRunner is an interface that should satisfy logic executor
type LogicRunner interface {
	LRI()
	OnPulse(context.Context, Pulse) error
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

	TraceID string // trace mark for Jaegar and friends
}

// ContractConstructor is a typedef for wrapper contract header
type ContractMethod func([]byte, []byte) ([]byte, []byte, error)

// ContractMethods maps name to contract method
type ContractMethods map[string]ContractMethod

// ContractConstructor is a typedef of typical contract constructor
type ContractConstructor func([]byte) ([]byte, error)

// ContractConstructors maps name to contract constructor
type ContractConstructors map[string]ContractConstructor

// ContractWrapper stores all needed about contract wrapper (it's methods/constructors)
type ContractWrapper struct {
	GetCode      ContractMethod
	GetPrototype ContractMethod

	Methods      ContractMethods
	Constructors ContractConstructors
}

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
