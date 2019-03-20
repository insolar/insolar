/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package artifactmanager

import (
	"context"

	"github.com/insolar/insolar/conveyor/adapter"
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/core/reply"
	"github.com/insolar/insolar/ledger/storage"
	"github.com/insolar/insolar/ledger/storage/object"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
)

// NewGetCodeAdapter creates new instance of adapter for getting code
func NewGetCodeAdapter() adapter.PulseConveyorAdapterTaskSink {
	return adapter.NewAdapterWithQueue(NewCodeGetter())
}

// GetCodeTask is task for adapter for getting code
type GetCodeTask struct {
	Parcel core.Parcel
}

// GetCodeResp is response for adapter for getting code
type GetCodeResp struct {
	Parcel core.Reply
	Err    error
}

// CodeGetter is worker for adapter for getting code
type CodeGetter struct {
	ObjectStorage          storage.ObjectStorage       `inject:""`
	JetCoordinator         core.JetCoordinator         `inject:""`
	DelegationTokenFactory core.DelegationTokenFactory `inject:""`
}

// NewCodeGetter returns new instance of worker which get code
func NewCodeGetter() adapter.Processor {
	return &CodeGetter{}
}

// Process implements Processor interface
func (rs *CodeGetter) Process(adapterID uint32, task adapter.AdapterTask, cancelInfo adapter.CancelInfo) adapter.Events {
	payload, ok := task.TaskPayload.(GetCodeTask)
	var msg GetCodeResp
	if !ok {
		msg.Err = errors.Errorf("[ LedgerCommunicator.Process ] Incorrect payload type: %T", task.TaskPayload)
		return adapter.Events{RespPayload: msg}
	}

	done := make(chan GetCodeResp, 1)
	go func(payload GetCodeTask) {
		ctx := context.Background()
		parcel, err := rs.handleGetCode(ctx, payload.Parcel)
		done <- GetCodeResp{parcel, err}
	}(payload)

	var flushed bool

	select {
	case <-cancelInfo.Cancel():
		log.Info("[ LedgerCommunicator.Process ] Cancel. Return Nil as Response")
	case <-cancelInfo.Flush():
		log.Info("[ LedgerCommunicator.Process ] Flush. DON'T Return Response")
		flushed = true
	case resp := <-done:
		log.Info("[ LedgerCommunicator.Process ] Process was dome successfully")
		msg = resp
	}

	return adapter.Events{RespPayload: msg, Flushed: flushed}
}

func (rs *CodeGetter) handleGetCode(ctx context.Context, parcel core.Parcel) (core.Reply, error) {
	msg := parcel.Message().(*message.GetCode)
	jetID := *core.NewJetID(0, nil)

	codeRec, err := rs.getCode(ctx, msg.Code.Record())
	if err == core.ErrNotFound {
		// We don't have code record. Must be on another node.
		node, err := rs.JetCoordinator.NodeForJet(ctx, core.RecordID(jetID), parcel.Pulse(), msg.Code.Record().Pulse())
		if err != nil {
			return nil, err
		}
		return reply.NewGetCodeRedirect(rs.DelegationTokenFactory, parcel, node)
	}
	if err != nil {
		return nil, err
	}
	code, err := rs.ObjectStorage.GetBlob(ctx, core.RecordID(jetID), codeRec.Code)
	if err != nil {
		return nil, err
	}

	rep := reply.Code{
		Code:        code,
		MachineType: codeRec.MachineType,
	}

	return &rep, nil
}

func (rs *CodeGetter) getCode(ctx context.Context, id *core.RecordID) (*object.CodeRecord, error) {
	jetID := *core.NewJetID(0, nil)

	rec, err := rs.ObjectStorage.GetRecord(ctx, core.RecordID(jetID), id)
	if err != nil {
		return nil, err
	}
	codeRec, ok := rec.(*object.CodeRecord)
	if !ok {
		return nil, errors.Wrap(ErrInvalidRef, "failed to retrieve code record")
	}

	return codeRec, nil
}
