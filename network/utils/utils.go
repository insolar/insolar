//
// Modified BSD 3-Clause Clear License
//
// Copyright (c) 2019 Insolar Technologies GmbH
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without modification,
// are permitted (subject to the limitations in the disclaimer below) provided that
// the following conditions are met:
//  * Redistributions of source code must retain the above copyright notice, this list
//    of conditions and the following disclaimer.
//  * Redistributions in binary form must reproduce the above copyright notice, this list
//    of conditions and the following disclaimer in the documentation and/or other materials
//    provided with the distribution.
//  * Neither the name of Insolar Technologies GmbH nor the names of its contributors
//    may be used to endorse or promote products derived from this software without
//    specific prior written permission.
//
// NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED
// BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS
// AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES,
// INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY
// AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL
// THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS
// OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//
// Notwithstanding any other provisions of this license, it is prohibited to:
//    (a) use this software,
//
//    (b) prepare modifications and derivative works of this software,
//
//    (c) distribute this software (including without limitation in source code, binary or
//        object code form), and
//
//    (d) reproduce copies of this software
//
//    for any commercial purposes, and/or
//
//    for the purposes of making available this software to third parties as a service,
//    including, without limitation, any software-as-a-service, platform-as-a-service,
//    infrastructure-as-a-service or other similar online service, irrespective of
//    whether it competes with the products or services of Insolar Technologies GmbH.
//

package utils

import (
	"hash/crc32"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log"
)

func WaitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return true // completed normally
	case <-time.After(timeout):
		return false // timed out
	}
}

// GenerateShortID generate short ID for node without checking collisions
func GenerateShortID(ref insolar.Reference) insolar.ShortNodeID {
	return insolar.ShortNodeID(GenerateUintShortID(ref))
}

// GenerateShortID generate short ID for node without checking collisions
func GenerateUintShortID(ref insolar.Reference) uint32 {
	return crc32.ChecksumIEEE(ref[:])
}

func OriginIsDiscovery(cert insolar.Certificate) bool {
	return IsDiscovery(*cert.GetNodeRef(), cert)
}

func IsDiscovery(nodeID insolar.Reference, cert insolar.Certificate) bool {
	bNodes := cert.GetDiscoveryNodes()
	for _, discoveryNode := range bNodes {
		if nodeID.Equal(*discoveryNode.GetNodeRef()) {
			return true
		}
	}
	return false
}

func CloseVerbose(closer io.Closer) {
	err := closer.Close()
	if err != nil {
		log.Errorf("[ CloseVerbose ] Failed to close: %s", err.Error())
	}
}

// IsConnectionClosed checks err for connection closed, workaround for poll.ErrNetClosing https://github.com/golang/go/issues/4373
func IsConnectionClosed(err error) bool {
	if err == nil {
		return false
	}
	err = errors.Cause(err)
	return strings.Contains(err.Error(), "use of closed network connection")
}
