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

package network

import (
	"bytes"
	"context"
	"github.com/insolar/insolar/network/node"
	"io"
	"math"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
	"github.com/pkg/errors"
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

// CheckShortIDCollision returns true if nodes contains node with such ShortID
func CheckShortIDCollision(nodes []insolar.NetworkNode, id insolar.ShortNodeID) bool {
	for _, n := range nodes {
		if id == n.ShortID() {
			return true
		}
	}

	return false
}

// GenerateUniqueShortID correct ShortID of the node so it does not conflict with existing active node list
func GenerateUniqueShortID(nodes []insolar.NetworkNode, nodeID insolar.Reference) insolar.ShortNodeID {
	shortID := insolar.ShortNodeID(node.GenerateUintShortID(nodeID))
	if !CheckShortIDCollision(nodes, shortID) {
		return shortID
	}
	return regenerateShortID(nodes, shortID)
}

func regenerateShortID(nodes []insolar.NetworkNode, shortID insolar.ShortNodeID) insolar.ShortNodeID {
	shortIDs := make([]insolar.ShortNodeID, len(nodes))
	for i, activeNode := range nodes {
		shortIDs[i] = activeNode.ShortID()
	}
	sort.Slice(shortIDs, func(i, j int) bool {
		return shortIDs[i] < shortIDs[j]
	})
	return generateNonConflictingID(shortIDs, shortID)
}

func generateNonConflictingID(sortedSlice []insolar.ShortNodeID, conflictingID insolar.ShortNodeID) insolar.ShortNodeID {
	index := sort.Search(len(sortedSlice), func(i int) bool {
		return sortedSlice[i] >= conflictingID
	})
	result := conflictingID
	repeated := false
	for {
		if result == math.MaxUint32 {
			if !repeated {
				repeated = true
				result = 0
				index = 0
			} else {
				panic("[ generateNonConflictingID ] shortID overflow twice")
			}
		}
		index++
		result++
		if index >= len(sortedSlice) || result != sortedSlice[index] {
			return result
		}
	}
}

// ExcludeOrigin returns DiscoveryNode slice without Origin
func ExcludeOrigin(discoveryNodes []insolar.DiscoveryNode, origin insolar.Reference) []insolar.DiscoveryNode {
	for i, discoveryNode := range discoveryNodes {
		if origin.Equal(*discoveryNode.GetNodeRef()) {
			return append(discoveryNodes[:i], discoveryNodes[i+1:]...)
		}
	}
	return discoveryNodes
}

// FindDiscoveryByRef tries to find discovery node in Certificate by reference
func FindDiscoveryByRef(cert insolar.Certificate, ref insolar.Reference) insolar.DiscoveryNode {
	bNodes := cert.GetDiscoveryNodes()
	for _, discoveryNode := range bNodes {
		if ref.Equal(*discoveryNode.GetNodeRef()) {
			return discoveryNode
		}
	}
	return nil
}

func OriginIsDiscovery(cert insolar.Certificate) bool {
	return IsDiscovery(*cert.GetNodeRef(), cert)
}

func IsDiscovery(nodeID insolar.Reference, cert insolar.Certificate) bool {
	return FindDiscoveryByRef(cert, nodeID) != nil
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

// FindDiscoveriesInNodeList returns only discovery nodes from active node list
func FindDiscoveriesInNodeList(nodes []insolar.NetworkNode, cert insolar.Certificate) []insolar.NetworkNode {
	discovery := cert.GetDiscoveryNodes()
	result := make([]insolar.NetworkNode, 0)

	for _, d := range discovery {
		for _, n := range nodes {
			if d.GetNodeRef().Equal(n.ID()) {
				result = append(result, n)
				break
			}
		}
	}

	return result
}

func IsClosedPipe(err error) bool {
	if err == nil {
		return false
	}
	err = errors.Cause(err)
	return strings.Contains(err.Error(), "read/write on closed pipe")
}

func NewPulseContext(ctx context.Context, pulseNumber uint32) context.Context {
	insTraceID := "pulse_" + strconv.FormatUint(uint64(pulseNumber), 10)
	ctx = inslogger.ContextWithTrace(ctx, insTraceID)
	return ctx
}

type CapturingReader struct {
	io.Reader
	buffer bytes.Buffer
}

func NewCapturingReader(reader io.Reader) *CapturingReader {
	return &CapturingReader{Reader: reader}
}

func (r *CapturingReader) Read(p []byte) (int, error) {
	n, err := r.Reader.Read(p)
	r.buffer.Write(p)
	return n, err
}

func (r *CapturingReader) Captured() []byte {
	return r.buffer.Bytes()
}
