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

package message

import (
	"github.com/insolar/insolar/insolar"
)

type ValidationResults struct {
	Caller           insolar.Reference
	RecordRef        insolar.Reference
	PassedStepsCount int
	Error            string
}

// AllowedSenderObjectAndRole implements interface method
func (vr *ValidationResults) AllowedSenderObjectAndRole() (*insolar.Reference, insolar.DynamicRole) {
	return &vr.RecordRef, insolar.DynamicRoleVirtualValidator
}

// DefaultRole returns role for this event
func (*ValidationResults) DefaultRole() insolar.DynamicRole {
	return insolar.DynamicRoleVirtualExecutor
}

// DefaultTarget returns of target of this event.
func (vr *ValidationResults) DefaultTarget() *insolar.Reference {
	return &vr.RecordRef
}

func (vr *ValidationResults) Type() insolar.MessageType {
	return insolar.TypeValidationResults
}

// TODO change after changing pulsar
func (vr *ValidationResults) GetCaller() *insolar.Reference {
	return &vr.Caller // TODO actually it's not right. There is no caller.
}

func (vr *ValidationResults) GetReference() insolar.Reference {
	return vr.RecordRef
}
