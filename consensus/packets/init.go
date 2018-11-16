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

package packets

import (
	"github.com/insolar/insolar/log"
)

// init packets and claims size variables
func init() {
	sizeOf := func(s Serializer) uint16 {
		data, err := s.Serialize()
		if err != nil {
			log.Fatalln("Failed to init packets package: ", err.Error())
		}
		return uint16(len(data))
	}

	phase1PacketSizeForClaims = phase1PacketMaxSize - int(sizeOf(&Phase1Packet{}))

	claimSizeMap = make(map[ClaimType]uint16, 5)
	claimSizeMap[TypeNodeJoinClaim] = sizeOf(&NodeJoinClaim{})
	claimSizeMap[TypeCapabilityPollingAndActivation] = sizeOf(&CapabilityPoolingAndActivation{})
	claimSizeMap[TypeNodeViolationBlame] = sizeOf(&NodeViolationBlame{})
	claimSizeMap[TypeNodeBroadcast] = sizeOf(&NodeBroadcast{})
	claimSizeMap[TypeNodeLeaveClaim] = sizeOf(&NodeLeaveClaim{})
}
