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

package fakepulsar

import (
	"context"
	"encoding/binary"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/insolar/insolar/core"
)

// Fakepulsar needed when the network starts and can't receive a real pulse.

// onPulse is a callbaback for pulse recv.
type callbackOnPulse func(ctx context.Context, pulse core.Pulse)

// FakePulsar is a struct which uses at void network state.
type FakePulsar struct {
	onPulse        callbackOnPulse
	stop           chan bool
	timeoutMs      int32 // ms
	running        bool
	firstPulseTime int64
	pulseNum       int64
	mut            sync.Mutex
}

// NewFakePulsar creates and returns a new FakePulsar.
func NewFakePulsar(callback callbackOnPulse, timeoutMs int32) *FakePulsar {
	return &FakePulsar{
		onPulse:   callback,
		timeoutMs: timeoutMs,
		stop:      make(chan bool),
		running:   false,
	}
}

// GetFakePulse creates and returns a fake pulse.
func (fp *FakePulsar) GetFakePulse() *core.Pulse {
	return fp.newPulse()
}

// Start starts sending a fake pulse.
func (fp *FakePulsar) Start(ctx context.Context) {
	fp.running = true
	fp.firstPulseTime = time.Now().Unix()
	go func(fp *FakePulsar) {
		for {
			select {
			case <-time.After(time.Millisecond * time.Duration(fp.timeoutMs)):
				{
					if math.Abs(float64(time.Now().Second()-time.Unix(fp.firstPulseTime, 0).Second())) != float64(fp.timeoutMs/1000) {
						time.Sleep(time.Duration(math.Abs(float64(time.Now().Second() - time.Unix(fp.firstPulseTime, 0).Second()))))
					}
					fp.mut.Lock()
					defer fp.mut.Unlock()
					fp.pulseNum++
					fp.onPulse(ctx, *fp.GetFakePulse())
				}
			case <-fp.stop:
				return
			}
		}
	}(fp)
}

// Stop sending a fake pulse.
func (fp *FakePulsar) Stop(ctx context.Context) {
	if fp.running {
		fp.stop <- true
		close(fp.stop)
		fp.running = false
	}
}

func (fp *FakePulsar) newPulse() *core.Pulse {
	rand.Seed(fp.pulseNum)
	tmp := make([]byte, core.EntropySize)
	binary.BigEndian.PutUint64(tmp, rand.Uint64())
	var entropy core.Entropy
	copy(entropy[:], tmp[:core.EntropySize])
	return &core.Pulse{
		PulseNumber:     0,
		NextPulseNumber: 0,
		Entropy:         entropy,
	}
}

func (fp *FakePulsar) GetFirstPulseTime() int64 {
	return fp.firstPulseTime
}

func (fp *FakePulsar) GetPulseNum() int64 {
	return fp.pulseNum
}

func (fp *FakePulsar) SetPulseData(time, pulseNum int64) {
	fp.mut.Lock()
	defer fp.mut.Unlock()
	fp.firstPulseTime = time
	fp.pulseNum = pulseNum
}
