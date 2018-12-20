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
	result := crc32.ChecksumIEEE(ref[:])
	return core.ShortNodeID(result)
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
