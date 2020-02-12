// Copyright 2020 Insolar Network Ltd.
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

package common

import (
	"github.com/tylerb/gls"
)

const glsSystemErrorKey = "systemError"

type SystemError interface {
	GetSystemError() error
	SetSystemError(err error)
}

type SystemErrorImpl struct{}

func NewSystemError() *SystemErrorImpl {
	return &SystemErrorImpl{}
}

func (h *SystemErrorImpl) GetSystemError() error {
	// SystemError means an error in the system (platform), not a particular contract.
	// For instance, timed out external call or failed deserialization means a SystemError.
	// In case of SystemError all following external calls during current method call return
	// an error and the result of the current method call is discarded (not registered).
	callContextInterface := gls.Get(glsSystemErrorKey)
	if callContextInterface == nil {
		return nil
	}
	return callContextInterface.(error)
}

func (h *SystemErrorImpl) SetSystemError(err error) {
	gls.Set(glsSystemErrorKey, err)
}
