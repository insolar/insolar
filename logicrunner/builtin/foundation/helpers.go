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

package foundation

import (
	"errors"

	"github.com/insolar/insolar/insolar"
)

// GetPulseNumber returns current pulse from context.
func GetPulseNumber() (insolar.PulseNumber, error) {
	req := GetLogicalContext().Request
	if req == nil {
		return insolar.PulseNumber(0), errors.New("request from LogicCallContext is nil, get pulse is failed")
	}
	return req.Record().Pulse(), nil
}

// GetRequestReference - Returns request reference from context.
func GetRequestReference() insolar.Reference {
	ctx := GetLogicalContext()
	if ctx.Request == nil {
		panic("context has no request set")
	}
	return *ctx.Request
}

// GetObject create proxy by address
// unimplemented
func GetObject(ref insolar.Reference) ProxyInterface {
	panic("not implemented")
}
