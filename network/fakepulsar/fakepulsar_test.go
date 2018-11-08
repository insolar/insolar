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
	"testing"
	"time"

	"github.com/insolar/insolar/core"
	"github.com/stretchr/testify/assert"
)

func onPulse(pulse core.Pulse) {
}

func TestGetFakePulse(t *testing.T) {
	pulsar := NewFakePulsar(onPulse, 1000)
	pulse := pulsar.GetFakePulse()
	assert.NotNil(t, pulse)
}

func TestFakePulsar_Start(t *testing.T) {
	pulsar := NewFakePulsar(onPulse, 1000)
	pulsar.Start()
	time.Sleep(time.Millisecond * 1100)
	pulsar.Stop()
}
