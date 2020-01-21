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

package critlog

import (
	"sync"

	"github.com/insolar/insolar/insolar"
)

var _ insolar.LoggerOutput = &ProxyLoggerOutput{}

type ProxyLoggerOutput struct {
	mutex  sync.RWMutex
	target insolar.LoggerOutput
}

func (p *ProxyLoggerOutput) GetTarget() insolar.LoggerOutput {
	p.mutex.RLock()
	t := p.target
	p.mutex.RUnlock()
	return t
}

func (p *ProxyLoggerOutput) SetTarget(t insolar.LoggerOutput) {
	for {
		if t == p {
			return
		}
		if tp, ok := t.(*ProxyLoggerOutput); ok {
			t = tp.GetTarget()
		} else {
			break
		}
	}

	p.mutex.Lock()
	p.target = t
	p.mutex.Unlock()
}

func (p *ProxyLoggerOutput) Write(b []byte) (n int, err error) {
	return p.GetTarget().Write(b)
}

func (p *ProxyLoggerOutput) Close() error {
	return p.GetTarget().Close()
}

func (p *ProxyLoggerOutput) LogLevelWrite(level insolar.LogLevel, b []byte) (int, error) {
	return p.GetTarget().LogLevelWrite(level, b)
}

func (p *ProxyLoggerOutput) Flush() error {
	return p.GetTarget().Flush()
}

func (p *ProxyLoggerOutput) LowLatencyWrite(level insolar.LogLevel, b []byte) (int, error) {
	return p.GetTarget().LowLatencyWrite(level, b)
}

func (p *ProxyLoggerOutput) IsLowLatencySupported() bool {
	return p.GetTarget().IsLowLatencySupported()
}
