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

// Package delegationtoken is about an authorization token that allows a node to perform
// actions it can not normally perform during this pulse
package delegationtoken

import (
	"encoding/gob"

	"github.com/insolar/insolar/insolar"
	"github.com/pkg/errors"
)

type PendingExecutionToken struct {
	Signature []byte
}

func (t *PendingExecutionToken) Type() insolar.DelegationTokenType {
	return insolar.DTTypePendingExecution
}

func (t *PendingExecutionToken) Verify(parcel insolar.Parcel) (bool, error) {
	switch mt := parcel.Message().Type(); mt {
	default:
		return false, errors.Errorf("Message of type %s can't be delegated with %s token", t.Type(), mt)
	}
}

// GetObjectRedirectToken is a redirect token for the GetObject method
type GetObjectRedirectToken struct {
	Signature []byte
}

// Type implementation of Token interface.
func (t *GetObjectRedirectToken) Type() insolar.DelegationTokenType {
	return insolar.DTTypeGetObjectRedirect
}

// Verify implementation of Token interface.
func (t *GetObjectRedirectToken) Verify(parcel insolar.Parcel) (bool, error) {
	panic("implement me")
}

// GetCodeRedirectToken is a redirect token for the GetObject method
type GetCodeRedirectToken struct {
	Signature []byte
}

// Type implementation of Token interface.
func (t *GetCodeRedirectToken) Type() insolar.DelegationTokenType {
	return insolar.DTTypeGetCodeRedirect
}

// Verify implementation of Token interface.
func (t *GetCodeRedirectToken) Verify(parcel insolar.Parcel) (bool, error) {
	panic("implement me")
}

func init() {
	gob.Register(&PendingExecutionToken{})
	gob.Register(&GetObjectRedirectToken{})
	gob.Register(&GetCodeRedirectToken{})
}
