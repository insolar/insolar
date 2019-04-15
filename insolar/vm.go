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
	"time"
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
	Stop() error
}

//go:generate minimock -i github.com/insolar/insolar/insolar.LogicRunner -o ../testutils -s _mock.go

// LogicRunner is an interface that should satisfy logic executor
type LogicRunner interface {
	Execute(context.Context, Parcel) (res Reply, err error)
	HandleValidateCaseBindMessage(context.Context, Parcel) (res Reply, err error)
	HandleValidationResultsMessage(context.Context, Parcel) (res Reply, err error)
	HandleExecutorResultsMessage(context.Context, Parcel) (res Reply, err error)
	OnPulse(context.Context, Pulse) error
}

// LogicCallContext is a context of contract execution
type LogicCallContext struct {
	Mode            string     // either "execution" or "validation"
	Callee          *Reference // Contract that was called
	Request         *Reference // ref of request
	Prototype       *Reference // Image of the callee
	Code            *Reference // ref of contract code
	CallerPrototype *Reference // Image of the caller
	Parent          *Reference // Parent of the callee
	Caller          *Reference // Contract that made the call
	Time            time.Time  // Time when call was made
	Pulse           Pulse      // Number of the pulse
	Immutable       bool
	TraceID         string
}
