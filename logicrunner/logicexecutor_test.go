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
	"testing"

	"github.com/gojuno/minimock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/testutils"
)

func TestLogicExecutor_New(t *testing.T) {
	le := NewLogicExecutor()
	require.NotNil(t, le)
}

func TestLogicExecutor_Execute(t *testing.T) {
	// testing error case only
	// success tested in ExecuteMethod and ExecuteConstructor
	ctx := inslogger.TestContext(t)
	le := &logicExecutor{}
	res, err := le.Execute(ctx, &Transcript{Request: &record.IncomingRequest{CallType: 322}})
	require.Error(t, err)
	require.Nil(t, res)
}

func TestLogicExecutor_ExecuteMethod(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	ctx := inslogger.TestContext(t)

	protoRef := gen.Reference()
	codeRef := gen.Reference()

	tests := []struct {
		name       string
		transcript *Transcript
		error      bool
		dc         artifacts.DescriptorsCache
		mm         MachinesManager
		res        *RequestResult
	}{
		{
			name: "success",
			transcript: &Transcript{
				ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
					ParentMock.Return(nil).
					MemoryMock.Return(nil).
					HeadRefMock.Return(nil),
				Request: &record.IncomingRequest{
					Prototype: &protoRef,
				},
			},
			mm: NewMachinesManagerMock(mc).
				GetExecutorMock.
				Return(
					testutils.NewMachineLogicExecutorMock(mc).
						CallMethodMock.Return([]byte{1, 2, 3}, []byte{3, 2, 1}, nil),
					nil,
				),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByObjectDescriptorMock.
				Return(
					artifacts.NewObjectDescriptorMock(mc).
						HeadRefMock.Return(&protoRef),
					artifacts.NewCodeDescriptorMock(mc).
						RefMock.Return(&codeRef).
						MachineTypeMock.Return(insolar.MachineTypeBuiltin),
					nil,
				),
			res: &RequestResult{
				NewMemory: []byte{1, 2, 3},
				Result:    []byte{3, 2, 1},
			},
		},
		{
			name: "success, no memory change",
			transcript: &Transcript{
				ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
					ParentMock.Return(nil).
					MemoryMock.Return([]byte{1, 2, 3}).
					HeadRefMock.Return(nil),
				Request: &record.IncomingRequest{
					Prototype: &protoRef,
				},
			},
			mm: NewMachinesManagerMock(mc).
				GetExecutorMock.
				Return(
					testutils.NewMachineLogicExecutorMock(mc).
						CallMethodMock.Return([]byte{1, 2, 3}, []byte{3, 2, 1}, nil),
					nil,
				),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByObjectDescriptorMock.
				Return(
					artifacts.NewObjectDescriptorMock(mc).
						HeadRefMock.Return(&protoRef),
					artifacts.NewCodeDescriptorMock(mc).
						RefMock.Return(&codeRef).
						MachineTypeMock.Return(insolar.MachineTypeBuiltin),
					nil,
				),
			res: &RequestResult{
				Result: []byte{3, 2, 1},
			},
		},
		{
			name: "success, immutable call, no change",
			transcript: &Transcript{
				ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
					ParentMock.Return(nil).
					MemoryMock.Return([]byte{1, 2, 3}).
					HeadRefMock.Return(nil),
				Request: &record.IncomingRequest{
					Prototype: &protoRef,
					Immutable: true,
				},
			},
			mm: NewMachinesManagerMock(mc).
				GetExecutorMock.
				Return(
					testutils.NewMachineLogicExecutorMock(mc).
						CallMethodMock.Return([]byte{1, 2, 3, 4, 5}, []byte{3, 2, 1}, nil),
					nil,
				),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByObjectDescriptorMock.
				Return(
					artifacts.NewObjectDescriptorMock(mc).
						HeadRefMock.Return(&protoRef),
					artifacts.NewCodeDescriptorMock(mc).
						RefMock.Return(&codeRef).
						MachineTypeMock.Return(insolar.MachineTypeBuiltin),
					nil,
				),
			res: &RequestResult{
				Result: []byte{3, 2, 1},
			},
		},
		{
			name: "success, deactivation",
			transcript: &Transcript{
				ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
					ParentMock.Return(nil).
					MemoryMock.Return(nil).
					HeadRefMock.Return(nil),
				Request: &record.IncomingRequest{
					Prototype: &protoRef,
				},
				Deactivate:   true, // this a bit hacky
			},
			mm: NewMachinesManagerMock(mc).
				GetExecutorMock.
				Return(
					testutils.NewMachineLogicExecutorMock(mc).
						CallMethodMock.Return([]byte{1, 2, 3}, []byte{3, 2, 1}, nil),
					nil,
				),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByObjectDescriptorMock.
				Return(
					artifacts.NewObjectDescriptorMock(mc).
						HeadRefMock.Return(&protoRef),
					artifacts.NewCodeDescriptorMock(mc).
						RefMock.Return(&codeRef).
						MachineTypeMock.Return(insolar.MachineTypeBuiltin),
					nil,
				),
			res: &RequestResult{
				Deactivation: true,
				Result:       []byte{3, 2, 1},
			},
		},
		{
			name: "parent mismatch",
			transcript: &Transcript{
				ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc),
				Request: &record.IncomingRequest{
					Prototype: &insolar.Reference{},
				},
			},
			mm: NewMachinesManagerMock(mc),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByObjectDescriptorMock.
				Return(
					artifacts.NewObjectDescriptorMock(mc).HeadRefMock.Return(&protoRef),
					artifacts.NewCodeDescriptorMock(mc),
					nil,
				),
			error: true,
		},
		{
			name: "error, descriptors trouble",
			transcript: &Transcript{
				ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc),
				Request: &record.IncomingRequest{
					Prototype: &protoRef,
				},
			},
			mm: NewMachinesManagerMock(mc),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByObjectDescriptorMock.
				Return(
					nil, nil, errors.New("some"),
				),
			error: true,
		},
		{
			name: "error, no such machine executor",
			transcript: &Transcript{
				ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc),
				Request: &record.IncomingRequest{
					Prototype: &protoRef,
				},
			},
			mm: NewMachinesManagerMock(mc).
				GetExecutorMock.
				Return(
					nil, errors.New("some"),
				),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByObjectDescriptorMock.
				Return(
					artifacts.NewObjectDescriptorMock(mc).
						HeadRefMock.Return(&protoRef),
					artifacts.NewCodeDescriptorMock(mc).
						MachineTypeMock.Return(insolar.MachineTypeBuiltin),
					nil,
				),
			error: true,
		},
		{
			name: "error, execution failed",
			transcript: &Transcript{
				ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
					ParentMock.Return(nil).
					MemoryMock.Return(nil).
					HeadRefMock.Return(nil),
				Request: &record.IncomingRequest{
					Prototype: &protoRef,
				},
			},
			mm: NewMachinesManagerMock(mc).
				GetExecutorMock.
				Return(
					testutils.NewMachineLogicExecutorMock(mc).
						CallMethodMock.Return(nil, nil, errors.New("some")),
					nil,
				),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByObjectDescriptorMock.
				Return(
					artifacts.NewObjectDescriptorMock(mc).
						HeadRefMock.Return(&protoRef),
					artifacts.NewCodeDescriptorMock(mc).
						RefMock.Return(&codeRef).
						MachineTypeMock.Return(insolar.MachineTypeBuiltin),
					nil,
				),
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vs := &logicExecutor{MachinesManager: test.mm, DescriptorsCache: test.dc}
			// using Execute to increase coverage, calls should only go to ExecuteMethod
			res, err := vs.Execute(ctx, test.transcript)
			if test.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.res, res)
		})
	}
}

func TestLogicExecutor_ExecuteConstructor(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()
	ctx := inslogger.TestContext(t)

	protoRef := gen.Reference()
	codeRef := gen.Reference()

	tests := []struct {
		name       string
		transcript *Transcript
		error      bool
		dc         artifacts.DescriptorsCache
		mm         MachinesManager
		res        *RequestResult
	}{
		{
			name: "success",
			transcript: &Transcript{
				Request: &record.IncomingRequest{
					CallType:  record.CTSaveAsChild,
					Caller:    gen.Reference(),
					Prototype: &protoRef,
				},
			},
			mm: NewMachinesManagerMock(mc).
				GetExecutorMock.
				Return(
					testutils.NewMachineLogicExecutorMock(mc).
						CallConstructorMock.Return([]byte{1, 2, 3}, nil),
					nil,
				),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByPrototypeRefMock.
				Return(
					artifacts.NewObjectDescriptorMock(mc).
						HeadRefMock.Return(&protoRef),
					artifacts.NewCodeDescriptorMock(mc).
						RefMock.Return(&codeRef).
						MachineTypeMock.Return(insolar.MachineTypeBuiltin),
					nil,
				),
			res: &RequestResult{
				Activation: true,
				NewMemory:  []byte{1, 2, 3},
			},
		},
		{
			name: "error, executor problem",
			transcript: &Transcript{
				Request: &record.IncomingRequest{
					CallType:  record.CTSaveAsChild,
					Caller:    gen.Reference(),
					Prototype: &protoRef,
				},
			},
			mm: NewMachinesManagerMock(mc).
				GetExecutorMock.
				Return(
					testutils.NewMachineLogicExecutorMock(mc).
						CallConstructorMock.Return(nil, errors.New("some")),
					nil,
				),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByPrototypeRefMock.
				Return(
					artifacts.NewObjectDescriptorMock(mc).
						HeadRefMock.Return(&protoRef),
					artifacts.NewCodeDescriptorMock(mc).
						RefMock.Return(&codeRef).
						MachineTypeMock.Return(insolar.MachineTypeBuiltin),
					nil,
				),
			error: true,
		},
		{
			name: "error, no machine type executor",
			transcript: &Transcript{
				Request: &record.IncomingRequest{
					CallType:  record.CTSaveAsChild,
					Caller:    gen.Reference(),
					Prototype: &protoRef,
				},
			},
			mm: NewMachinesManagerMock(mc).
				GetExecutorMock.
				Return(nil, errors.New("some")),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByPrototypeRefMock.
				Return(
					artifacts.NewObjectDescriptorMock(mc),
					artifacts.NewCodeDescriptorMock(mc).
						MachineTypeMock.Return(insolar.MachineTypeBuiltin),
					nil,
				),
			error: true,
		},
		{
			name: "error, no descriptors",
			transcript: &Transcript{
				Request: &record.IncomingRequest{
					CallType:  record.CTSaveAsChild,
					Caller:    gen.Reference(),
					Prototype: &protoRef,
				},
			},
			mm: NewMachinesManagerMock(mc),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByPrototypeRefMock.
				Return(
					nil, nil, errors.New("some"),
				),
			error: true,
		},
		{
			name: "error, nil prototype",
			transcript: &Transcript{
				Request: &record.IncomingRequest{
					CallType:  record.CTSaveAsChild,
					Caller:    gen.Reference(),
					Prototype: nil,
				},
			},
			error: true,
		},
		{
			name: "error, empty caller",
			transcript: &Transcript{
				Request: &record.IncomingRequest{
					CallType:  record.CTSaveAsChild,
				},
			},
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			vs := &logicExecutor{MachinesManager: test.mm, DescriptorsCache: test.dc}
			res, err := vs.Execute(ctx, test.transcript)
			if test.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.res, res)
		})
	}
}
