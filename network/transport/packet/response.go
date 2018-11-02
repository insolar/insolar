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

package packet

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/network"
	"github.com/insolar/insolar/network/transport/host"
	"github.com/insolar/insolar/network/transport/relay"
)

// CheckNodePrivState - state of check node privileges request.
type CheckNodePrivState int

const (
	// Error - some error, see error string.
	Error = CheckNodePrivState(iota + 1)
	// Confirmed - state confirmed.
	Confirmed
	// Declined - state declined.
	Declined
)

// ResponseCheckNodePriv is data for check node privileges response.
type ResponseCheckNodePriv struct {
	State CheckNodePrivState
	Error string
}

// ResponseDataFindHost is data for FindHost response.
type ResponseDataFindHost struct {
	Closest []*host.Host
}

// ResponseDataFindValue is data for FindValue response.
type ResponseDataFindValue struct {
	Closest []*host.Host
	Value   []byte
}

// ResponseDataStore is data for Store response.
type ResponseDataStore struct {
	Success bool
}

// ResponseDataRPC is data for RPC response.
type ResponseDataRPC struct {
	Success bool
	Result  []byte
	Error   string
}

// ResponseCascadeSend is the response data of a cascade sending call
type ResponseCascadeSend struct {
	Success bool
	Error   string
}

// ResponsePulse is the response for a new pulse from a pulsar
type ResponsePulse struct {
	Success bool
	Error   string
}

// ResponseGetRandomHosts is the response containing random hosts of the DHT network
type ResponseGetRandomHosts struct {
	Hosts []host.Host
	Error string
}

// ResponseRelay is data for relay request response.
type ResponseRelay struct {
	State relay.State
}

// ResponseAuthentication is data for authentication request response.
type ResponseAuthentication struct {
	Success       bool
	AuthUniqueKey []byte
}

// ResponseCheckOrigin is data for check originality request response.
type ResponseCheckOrigin struct {
	AuthUniqueKey []byte
}

// ResponseObtainIP is data for get a IP of requesting host.
type ResponseObtainIP struct {
	IP string
}

// ResponseRelayOwnership is data to response to relay ownership request.
type ResponseRelayOwnership struct {
	Accepted bool
}

// ResponseKnownOuterHosts is data to answer if origin host know more outer hosts.
type ResponseKnownOuterHosts struct {
	ID         string // 	id of host in which more known outer hosts
	OuterHosts int    // number of known outer hosts
}

// ResponseGetNonce is data to answer to authorization request.
type ResponseGetNonce struct {
	Nonce []byte
	Error string
}

type ResponseCheckSignedNonce struct {
	Error       string
	ActiveNodes []core.Node
}

// RequestExchangeUnsyncLists is request to exchange unsync lists during consensus
type ResponseExchangeUnsyncLists struct {
	UnsyncList []core.Node
	Error      string
}

// RequestExchangeUnsyncHash is request to exchange hash of merged unsync lists during consensus
type ResponseExchangeUnsyncHash struct {
	UnsyncHash []*network.NodeUnsyncHash
	Error      string
}

// ResponseDisconnect id data to answer to disconnected node.
type ResponseDisconnect struct {
	Disconnected bool
	Error        error
}
