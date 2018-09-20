/*
 *    Copyright 2018 Insolar
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
	"log"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/eventbus/event"
	"github.com/insolar/insolar/eventbus/reaction"
)

// HandleEvent performs event processing.
func (am *LedgerArtifactManager) HandleEvent(e core.Event) (core.Reaction, error) {

	machinePref := []core.MachineType{
		core.MachineTypeBuiltin,
		core.MachineTypeGoPlugin,
	}

	switch m := e.(type) {

	case *event.ActivateObjDelegate:
		ref, err := am.ActivateObjDelegate(
			// core.RecordRef{}, core.RecordRef{}, m.Class, m.Into, m.Body,
			m.Domain, m.Request, m.Class, m.Parent, m.Memory,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't save new object")
		}
		return &reaction.CommonReaction{Data: []byte(ref.String())}, nil

	case *event.ActivateObj:
		ref, err := am.ActivateObj(
			// core.RecordRef{}, core.RecordRef{}, m.Class, m.Into, m.Body,
			m.Domain, m.Request, m.Class, m.Parent, m.Memory,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't save new object")
		}
		return &reaction.CommonReaction{Data: []byte(ref.String())}, nil

	case *event.UpdateObj:
		_, err := am.UpdateObj(
			// core.RecordRef{}, core.RecordRef{}, m.Object, m.Body,
			m.Domain, m.Request, m.Obj, m.Memory,
		)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't update object")
		}
		return &reaction.CommonReaction{}, nil

	case *event.GetLatestObj:
		objDesc, err := am.GetLatestObj(m.Head)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get object")
		}

		classDesc, err := objDesc.ClassDescriptor(nil)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get object's class")
		}

		data, err := objDesc.Memory()
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get object's data")
		}

		codeDesc, err := classDesc.CodeDescriptor(machinePref)
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get object's code descriptor")
		}

		mt, err := codeDesc.MachineType()
		if err != nil {
			return nil, errors.Wrap(err, "couldn't get machine type")
		}

		return &reaction.ObjectBodyReaction{
			Body:        data,
			Code:        *codeDesc.Ref(),
			Class:       *classDesc.HeadRef(),
			MachineType: mt,
		}, nil
	default:
		log.Fatalf("Unknown event type %T", e)
	}
	return nil, errors.New("unknown event type")
}
