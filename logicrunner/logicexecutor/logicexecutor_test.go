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

package logicexecutor

import (
	"context"
	"testing"
	"time"

	"github.com/gojuno/minimock/v3"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/insolar/record"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/logicrunner/artifacts"
	"github.com/insolar/insolar/logicrunner/builtin/foundation"
	"github.com/insolar/insolar/logicrunner/common"
	"github.com/insolar/insolar/logicrunner/machinesmanager"
	"github.com/insolar/insolar/logicrunner/requestresult"
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
	res, err := le.Execute(ctx, &common.Transcript{Request: &record.IncomingRequest{CallType: 322}})
	require.Error(t, err)
	require.Nil(t, res)
}

func TestLogicExecutor_ExecuteMethod(t *testing.T) {
	objRef := gen.Reference()
	objRecordID := gen.ID()
	protoRef := gen.Reference()
	codeRef := gen.Reference()

	tests := []struct {
		name  string
		mocks func(ctx context.Context, t minimock.Tester) (LogicExecutor, *common.Transcript)
		error bool
		res   artifacts.RequestResult
	}{
		{
			name: "success",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
						ParentMock.Return(nil).
						MemoryMock.Return(nil).
						HeadRefMock.Return(&objRef).
						StateIDMock.Return(&objRecordID).
						PrototypeMock.Return(&protoRef, nil),
					Request: &record.IncomingRequest{
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						testutils.NewMachineLogicExecutorMock(mc).
							CallMethodMock.Return([]byte{1, 2, 3}, []byte{3, 2, 1}, nil),
						nil,
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByObjectDescriptorMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							RefMock.Return(&codeRef).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			res: &requestresult.RequestResult{
				RawObjectReference: objRef,
				ObjectImage:        protoRef,
				ObjectStateID:      objRecordID,
				SideEffectType:     artifacts.RequestSideEffectAmend,
				Memory:             []byte{1, 2, 3},
				RawResult:          []byte{3, 2, 1},
			},
		},
		{
			name: "success, no  Memory change",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
						ParentMock.Return(nil).
						MemoryMock.Return([]byte{1, 2, 3}).
						HeadRefMock.Return(&objRef),
					Request: &record.IncomingRequest{
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						testutils.NewMachineLogicExecutorMock(mc).
							CallMethodMock.Return([]byte{1, 2, 3}, []byte{3, 2, 1}, nil),
						nil,
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByObjectDescriptorMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							RefMock.Return(&codeRef).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			res: &requestresult.RequestResult{
				SideEffectType:     artifacts.RequestSideEffectNone,
				RawResult:          []byte{3, 2, 1},
				RawObjectReference: objRef,
			},
		},
		{
			name: "success, immutable call, no change",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
						ParentMock.Return(nil).
						MemoryMock.Return([]byte{1, 2, 3}).
						HeadRefMock.Return(&objRef),
					Request: &record.IncomingRequest{
						Prototype: &protoRef,
						Immutable: true,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						testutils.NewMachineLogicExecutorMock(mc).
							CallMethodMock.Return([]byte{1, 2, 3, 4, 5}, []byte{3, 2, 1}, nil),
						nil,
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByObjectDescriptorMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							RefMock.Return(&codeRef).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			res: &requestresult.RequestResult{
				SideEffectType:     artifacts.RequestSideEffectNone,
				RawResult:          []byte{3, 2, 1},
				RawObjectReference: objRef,
			},
		},
		{
			name: "success, deactivation",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
						ParentMock.Return(nil).
						MemoryMock.Return(nil).
						StateIDMock.Return(&objRecordID).
						HeadRefMock.Return(&objRef),
					Request: &record.IncomingRequest{
						Prototype: &protoRef,
					},
					Deactivate: true, // this a bit hacky
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						testutils.NewMachineLogicExecutorMock(mc).
							CallMethodMock.Return([]byte{1, 2, 3}, []byte{3, 2, 1}, nil),
						nil,
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByObjectDescriptorMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							RefMock.Return(&codeRef).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			res: &requestresult.RequestResult{
				SideEffectType:     artifacts.RequestSideEffectDeactivate,
				RawResult:          []byte{3, 2, 1},
				ObjectStateID:      objRecordID,
				RawObjectReference: objRef,
			},
		},
		{
			name: "parent mismatch",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).HeadRefMock.Return(&objRef),
					Request: &record.IncomingRequest{
						Prototype: insolar.NewEmptyReference(),
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByObjectDescriptorMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			error: false,
			res: &requestresult.RequestResult{
				RawResult: func() []byte {
					err := errors.New("proxy call error: try to call method of prototype as method of another prototype")
					errResBuf, err := foundation.MarshalMethodErrorResult(err)
					if err != nil {
						require.NoError(t, err)
					}
					return errResBuf
				}(),
				RawObjectReference: objRef,
			},
		},
		{
			name: "error, descriptors trouble",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc),
					Request: &record.IncomingRequest{
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByObjectDescriptorMock.
					Return(
						nil, nil, errors.New("some"),
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			error: true,
		},
		{
			name: "error, no such machine executor",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc),
					Request: &record.IncomingRequest{
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						nil, errors.New("some"),
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByObjectDescriptorMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			error: true,
		},
		{
			name: "error, execution failed",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
						ParentMock.Return(nil).
						MemoryMock.Return(nil).
						HeadRefMock.Return(&objRef),
					Request: &record.IncomingRequest{
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						testutils.NewMachineLogicExecutorMock(mc).
							CallMethodMock.Return(nil, nil, errors.New("some")),
						nil,
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByObjectDescriptorMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							RefMock.Return(&codeRef).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			error: true,
		},
		{
			name: "error, empty result",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
						ParentMock.Return(nil).
						MemoryMock.Return(nil).
						HeadRefMock.Return(&objRef),
					Request: &record.IncomingRequest{
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						testutils.NewMachineLogicExecutorMock(mc).
							CallMethodMock.Return(nil, []byte{}, nil),
						nil,
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByObjectDescriptorMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							RefMock.Return(&codeRef).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			error: true,
		},
		{
			name: "error, empty state",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					ObjectDescriptor: artifacts.NewObjectDescriptorMock(mc).
						ParentMock.Return(nil).
						MemoryMock.Return(nil).
						HeadRefMock.Return(&objRef),
					Request: &record.IncomingRequest{
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						testutils.NewMachineLogicExecutorMock(mc).
							CallMethodMock.Return([]byte{}, []byte{1, 2, 3}, nil),
						nil,
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByObjectDescriptorMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							RefMock.Return(&codeRef).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			ctx := inslogger.TestContext(t)

			le, tr := test.mocks(ctx, mc)
			// using Execute to increase coverage, calls should only go to ExecuteMethod
			res, err := le.Execute(ctx, tr)
			if test.error {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.res, res)
			}

			mc.Wait(1 * time.Minute)
			mc.Finish()
		})
	}
}

func TestLogicExecutor_ExecuteConstructor(t *testing.T) {
	protoRef := gen.Reference()
	callerRef := gen.Reference()
	codeRef := gen.Reference()
	baseRef := gen.Reference()
	objectRef := gen.Reference()

	tests := []struct {
		name  string
		mocks func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript)
		error bool
		res   *requestresult.RequestResult
	}{
		{
			name: "success",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					Request: &record.IncomingRequest{
						Object:    &objectRef,
						Base:      &baseRef,
						CallType:  record.CTSaveAsChild,
						Caller:    callerRef,
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						testutils.NewMachineLogicExecutorMock(mc).
							CallConstructorMock.Return([]byte{1, 2, 3}, []byte{3, 2, 1}, nil),
						nil,
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByPrototypeRefMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							RefMock.Return(&codeRef).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			res: &requestresult.RequestResult{
				SideEffectType:     artifacts.RequestSideEffectActivate,
				RawResult:          []byte{3, 2, 1},
				Memory:             []byte{1, 2, 3},
				ParentReference:    baseRef,
				ObjectImage:        protoRef,
				RawObjectReference: objectRef,
			},
		},
		{
			name: "empty state, logic error in constructor",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					Request: &record.IncomingRequest{
						Object:    &objectRef,
						Base:      &baseRef,
						CallType:  record.CTSaveAsChild,
						Caller:    callerRef,
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						testutils.NewMachineLogicExecutorMock(mc).
							CallConstructorMock.Return(nil, []byte{3, 2, 1}, nil),
						nil,
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByPrototypeRefMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							RefMock.Return(&codeRef).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			res: &requestresult.RequestResult{
				RawResult:          []byte{3, 2, 1},
				RawObjectReference: objectRef,
			},
		},
		{
			name: "error, executor problem",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					Request: &record.IncomingRequest{
						Object:    &objectRef,
						Base:      &baseRef,
						CallType:  record.CTSaveAsChild,
						Caller:    callerRef,
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						testutils.NewMachineLogicExecutorMock(mc).
							CallConstructorMock.Return(
							nil, nil, errors.New("some"),
						),
						nil,
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByPrototypeRefMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							RefMock.Return(&codeRef).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			error: true,
		},
		{
			name: "error, empty result",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					Request: &record.IncomingRequest{
						Object:    &objectRef,
						Base:      &baseRef,
						CallType:  record.CTSaveAsChild,
						Caller:    callerRef,
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(
						testutils.NewMachineLogicExecutorMock(mc).
							CallConstructorMock.Return(
							nil, []byte{}, nil,
						),
						nil,
					)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByPrototypeRefMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc).
							HeadRefMock.Return(&protoRef),
						artifacts.NewCodeDescriptorMock(mc).
							RefMock.Return(&codeRef).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			error: true,
		},
		{
			name: "error, no machine type executor",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					Request: &record.IncomingRequest{
						Object:    &objectRef,
						Base:      &baseRef,
						CallType:  record.CTSaveAsChild,
						Caller:    gen.Reference(),
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc).
					GetExecutorMock.
					Return(nil, errors.New("some"))
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByPrototypeRefMock.
					Return(
						artifacts.NewPrototypeDescriptorMock(mc),
						artifacts.NewCodeDescriptorMock(mc).
							MachineTypeMock.Return(insolar.MachineTypeBuiltin),
						nil,
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			error: true,
		},
		{
			name: "error, no descriptors",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					Request: &record.IncomingRequest{
						Object:    &objectRef,
						CallType:  record.CTSaveAsChild,
						Caller:    gen.Reference(),
						Prototype: &protoRef,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc)
				dc := artifacts.NewDescriptorsCacheMock(mc).
					ByPrototypeRefMock.
					Return(
						nil, nil, errors.New("some"),
					)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			error: true,
		},
		{
			name: "error, nil prototype",
			mocks: func(ctx context.Context, mc minimock.Tester) (LogicExecutor, *common.Transcript) {
				tr := &common.Transcript{
					Request: &record.IncomingRequest{
						Object:    &objectRef,
						CallType:  record.CTSaveAsChild,
						Caller:    gen.Reference(),
						Prototype: nil,
					},
				}
				mm := machinesmanager.NewMachinesManagerMock(mc)
				dc := artifacts.NewDescriptorsCacheMock(mc)
				return &logicExecutor{MachinesManager: mm, DescriptorsCache: dc}, tr
			},
			error: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			mc := minimock.NewController(t)
			ctx := inslogger.TestContext(t)

			le, tr := test.mocks(ctx, mc)
			// using Execute to increase coverage, calls should only go to ExecuteMethod
			res, err := le.Execute(ctx, tr)
			if test.error {
				require.Error(t, err)
				require.Nil(t, res)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.res, res)
			}

			mc.Wait(1 * time.Minute)
			mc.Finish()
		})
	}
}
