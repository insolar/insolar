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

package foundation

import (
	"encoding/binary"
	"errors"
	"math/rand"
	"strings"

	"github.com/insolar/insolar/insolar"
)

// GetPulseNumber returns current pulse from context.
func GetPulseNumber() (insolar.PulseNumber, error) {
	req := GetLogicalContext().Request
	if req == nil {
		return insolar.PulseNumber(0), errors.New("request from LogicCallContext is nil, get pulse is failed")
	}
	return req.GetLocal().Pulse(), nil
}

// GetRequestReference - Returns request reference from context.
func GetRequestReference() (*insolar.Reference, error) {
	ctx := GetLogicalContext()
	if ctx.Request == nil {
		return nil, errors.New("request from LogicCallContext is nil, get pulse is failed")
	}
	return ctx.Request, nil
}

// NewSource returns source initialized with entropy from pulse.
func NewSource() rand.Source {
	randNum := binary.LittleEndian.Uint64(GetLogicalContext().Pulse.Entropy[:])
	return rand.NewSource(int64(randNum))
}

// GetObject creates proxy by address
// unimplemented
func GetObject(ref insolar.Reference) ProxyInterface {
	panic("not implemented")
}

// TrimPublicKey trims public key
func TrimPublicKey(publicKey string) string {
	return strings.Join(strings.Split(strings.TrimSpace(between(publicKey, "KEY-----", "-----END")), "\n"), "")
}

// TrimAddress trims address
func TrimAddress(address string) string {
	return strings.ToLower(strings.Join(strings.Split(strings.TrimSpace(address), "\n"), ""))
}

func between(value string, a string, b string) string {
	// Get substring between two strings.
	pos := strings.Index(value, a)
	if pos == -1 {
		return value
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return value
	}
	posFirst := pos + len(a)
	if posFirst >= posLast {
		return value
	}
	return value[posFirst:posLast]
}
