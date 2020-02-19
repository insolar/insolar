// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

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
