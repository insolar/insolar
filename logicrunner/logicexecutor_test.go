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
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/testutils"
)

func TestLogicExecutor_ExecuteMethod(t *testing.T) {
	mc := minimock.NewController(t)
	defer mc.Finish()

	protoRef := gen.Reference()
	codeRef := gen.Reference()

	tests := []struct {
		name       string
		transcript *Transcript
		error      bool
		dc         artifacts.DescriptorsCache
		mm         MachinesManager
		res        artifacts.RequestResult
	}{
		{
			name: "success",
			transcript: &Transcript{
				ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
					ParentMock.Return(nil).
					MemoryMock.Return(nil).
					HeadRefMock.Return(nil).
					StateIDMock.Return(nil).
					IsPrototypeMock.Return(false).
					PrototypeMock.Return(nil, nil),
				Request: &record.IncomingRequest{
					Prototype: &protoRef,
				},
				LogicContext: &insolar.LogicCallContext{},
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
			res: &requestResult{
				sideEffectType: artifacts.RequestSideEffectAmend,
				memory:         []byte{1, 2, 3},
				result:         []byte{3, 2, 1},
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
				LogicContext: &insolar.LogicCallContext{},
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
			res: &requestResult{
				sideEffectType: artifacts.RequestSideEffectNone,
				result:         []byte{3, 2, 1},
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
				LogicContext: &insolar.LogicCallContext{},
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
			res: &requestResult{
				result: []byte{3, 2, 1},
			},
		},
		{
			name: "parent mismatch",
			transcript: &Transcript{
				ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).ParentMock.Return(nil).HeadRefMock.Return(nil),
				Request: &record.IncomingRequest{
					Prototype: &insolar.Reference{},
				},
				LogicContext: &insolar.LogicCallContext{},
			},
			mm: NewMachinesManagerMock(mc),
			dc: artifacts.NewDescriptorsCacheMock(mc).
				ByObjectDescriptorMock.
				Return(
					artifacts.NewObjectDescriptorMock(mc).HeadRefMock.Return(&protoRef),
					artifacts.NewCodeDescriptorMock(mc).
						RefMock.Return(&codeRef),
					nil,
				),
			error: true,
			res:   (artifacts.RequestResult)(nil),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := inslogger.TestContext(t)

			vs := &logicExecutor{MachinesManager: test.mm, DescriptorsCache: test.dc}
			res, err := vs.ExecuteMethod(ctx, test.transcript)
			if test.error {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, test.res, res)
		})
	}
}
