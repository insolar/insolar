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

// Package delegationtoken is about an authorization token that allows a node to perform
// actions it can not normally perform during this pulse
package delegationtoken

import (
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type PendingExecution struct {
	Signature []byte
}

func (t *PendingExecution) Type() core.DelegationTokenType {
	return core.DTTypePendingExecution
}

func (t *PendingExecution) Verify(parcel core.Parcel) (bool, error) {
	switch mt := parcel.Message().Type(); mt {

	//TODO: stab should start verification
	case core.TypeCallMethod:
		return false, nil
	default:
		return false, errors.Errorf("Message of type %s can't be delegated with %s token", t.Type(), mt)
	}
}

// GetObjectRedirect is a redirect token for the GetObject method
type GetObjectRedirect struct {
	Signature []byte
}

// Type implementation of Token interface.
func (t *GetObjectRedirect) Type() core.DelegationTokenType {
	return core.DTTypeGetObjectRedirect
}

// Verify implementation of Token interface.
func (t *GetObjectRedirect) Verify(parcel core.Parcel) (bool, error) {
	panic("")
}

// GetChildrenRedirect is a redirect token for the GetObject method
type GetChildrenRedirect struct {
	Signature []byte
}

// Type implementation of Token interface.
func (t *GetChildrenRedirect) Type() core.DelegationTokenType {
	return core.DTTypeGetChildrenRedirect
}

// Verify implementation of Token interface.
func (t *GetChildrenRedirect) Verify(parcel core.Parcel) (bool, error) {
	panic("")
}

// GetCodeRedirect is a redirect token for the GetObject method
type GetCodeRedirect struct {
	Signature []byte
}

// Type implementation of Token interface.
func (t *GetCodeRedirect) Type() core.DelegationTokenType {
	return core.DTTypeGetCodeRedirect
}

// Verify implementation of Token interface.
func (t *GetCodeRedirect) Verify(parcel core.Parcel) (bool, error) {
	panic("")
}
