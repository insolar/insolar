/*
 * The Clear BSD License
 *
 * Copyright (c) 2019 Insolar Technologies
 *
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without modification, are permitted (subject to the limitations in the disclaimer below) provided that the following conditions are met:
 *
 *  Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.
 *  Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.
 *  Neither the name of Insolar Technologies nor the names of its contributors may be used to endorse or promote products derived from this software without specific prior written permission.
 *
 * NO EXPRESS OR IMPLIED LICENSES TO ANY PARTY'S PATENT RIGHTS ARE GRANTED BY THIS LICENSE. THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 */

package utils

import (
	"hash/crc32"
	"io"
	"sync"
	"sync/atomic"
	"time"

	"github.com/insolar/insolar/core"
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

// AtomicLoadAndIncrementUint64 performs CAS loop, increments counter and returns old value.
func AtomicLoadAndIncrementUint64(addr *uint64) uint64 {
	for {
		val := atomic.LoadUint64(addr)
		if atomic.CompareAndSwapUint64(addr, val, val+1) {
			return val
		}
	}
}

// GenerateShortID generate short ID for node without checking collisions
func GenerateShortID(ref core.RecordRef) core.ShortNodeID {
	return core.ShortNodeID(GenerateUintShortID(ref))
}

// GenerateShortID generate short ID for node without checking collisions
func GenerateUintShortID(ref core.RecordRef) uint32 {
	return crc32.ChecksumIEEE(ref[:])
}

func OriginIsDiscovery(cert core.Certificate) bool {
	bNodes := cert.GetDiscoveryNodes()
	for _, discoveryNode := range bNodes {
		if cert.GetNodeRef().Equal(*discoveryNode.GetNodeRef()) {
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

// FindDiscoveriesInNodeList returns only discovery nodes from active node list
func FindDiscoveriesInNodeList(nodes []core.Node, cert core.Certificate) []core.Node {
	discovery := cert.GetDiscoveryNodes()
	result := make([]core.Node, 0)

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
