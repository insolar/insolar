//
// Copyright 2019 Insolar Technologies GmbH
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
//

package log

import (
	"errors"
	"io"
	"sync"
)

type ProxyWriter struct {
	mutex  sync.RWMutex
	target io.Writer
}

func (p *ProxyWriter) getTarget() io.Writer {
	p.mutex.RLock()
	t := p.target
	p.mutex.RUnlock()
	return t
}

func (p *ProxyWriter) setTarget(t io.Writer) {
	p.mutex.Lock()
	p.target = t
	p.mutex.Unlock()
}

func (p *ProxyWriter) Write(b []byte) (n int, err error) {
	return p.getTarget().Write(b)
}

func (p *ProxyWriter) Close() error {
	if c, ok := p.getTarget().(io.Closer); ok {
		return c.Close()
	}
	return errors.New("not supported: Close")
}

func (p *ProxyWriter) Flush() error {
	type flusher interface {
		Flush() error
	}
	if c, ok := p.getTarget().(flusher); ok {
		return c.Flush()
	}
	return errors.New("not supported: Flush")
}

func (p *ProxyWriter) Sync() error {
	type flusher interface {
		Sync() error
	}
	if c, ok := p.getTarget().(flusher); ok {
		return c.Sync()
	}
	return errors.New("not supported: Sync")
}
