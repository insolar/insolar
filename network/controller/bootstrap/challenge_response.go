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

package bootstrap

import (
	"encoding/gob"

	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

type ChallengeResponseController struct {
}

type Nonce []byte
type SignedNonce []byte
type SessionID uint64

// Node                           Discovery Node
//  | ------ ChallengeRequest -----> |
//  | <-- SignedChallengeResponse -- |
//  | --- SignedChallengeRequest --> |
//  | <----- ChallengeResponse ----- |
// ------------------------------------

type SessionHeader struct {
	Success   bool
	SessionID SessionID
	Error     string
}

type ChallengeRequest struct {
	Certificate core.Certificate
	Nonce       Nonce
}

type SignedChallengeResponse struct {
	SessionHeader SessionHeader

	SignedNonce       SignedNonce
	XorDiscoveryNonce Nonce
	DiscoveryNonce    Nonce
}

type SignedChallengeRequest struct {
	SessionHeader SessionHeader

	SignedDiscoveryNonce SignedNonce
	XorNonce             Nonce
}

type ChallengeResponse struct {
	SessionHeader SessionHeader

	CurrentPulse core.Pulse
	State        core.NetworkState
}

func init() {
	gob.Register(&ChallengeRequest{})
	gob.Register(&SignedChallengeResponse{})
	gob.Register(&SignedChallengeRequest{})
	gob.Register(&ChallengeResponse{})
}

func (crc *ChallengeResponseController) SendRequest(request *ChallengeRequest) (*ChallengeResponse, error) {
	// TODO: double challenge response
	return nil, errors.New("not implemented")
}
