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

package critlog

import (
	"errors"
	"github.com/insolar/insolar/insolar"
	"io"
	"sync/atomic"
)

var _ insolar.LogLevelWriter = &FlushBypass{}

func NewFlushBypass(output io.Writer) FlushBypass {
	return FlushBypass{Writer: output}
}

type FlushBypass struct {
	io.Writer
	closed uint32
}

func (p *FlushBypass) DoLevelWrite(level insolar.LogLevel, b []byte) (n int, err error) {
	if lw, ok := p.Writer.(insolar.LogLevelWriter); ok {
		return lw.LogLevelWrite(level, b)
	}
	return p.Writer.Write(b)
}

func (p *FlushBypass) DoWrite(b []byte) (n int, err error) {
	return p.Writer.Write(b)
}

func (p *FlushBypass) LogLevelWrite(level insolar.LogLevel, b []byte) (int, error) {
	if p.IsClosed() {
		return 0, p.ClosedError()
	}
	return p.DoLevelWrite(level, b)
}

func (p *FlushBypass) ClosedError() error {
	return errors.New("closed")
}

func (p *FlushBypass) Write(b []byte) (int, error) {
	if p.IsClosed() {
		return 0, p.ClosedError()
	}
	return p.DoLevelWrite(insolar.NoLevel, b)
}

func (p *FlushBypass) SetClosed() bool {
	return atomic.CompareAndSwapUint32(&p.closed, 0, 1)
}

func (p *FlushBypass) IsClosed() bool {
	return atomic.LoadUint32(&p.closed) != 0
}

func (p *FlushBypass) DoClose() (bool, error) {
	if f, ok := p.Writer.(io.Closer); ok {
		return true, f.Close()
	}
	return false, errors.New("unsupported: Close")
}

func (p *FlushBypass) DoFlush() (bool, error) {
	type flusher interface {
		Flush() error
	}
	if f, ok := p.Writer.(flusher); ok {
		return true, f.Flush()
	}
	return false, errors.New("unsupported: Flush")
}

func (p *FlushBypass) DoSync() (bool, error) {
	type syncer interface {
		Sync() error
	}
	if f, ok := p.Writer.(syncer); ok {
		return true, f.Sync()
	}
	return false, errors.New("unsupported: Sync")
}

func (p *FlushBypass) DoFlushOrSync() (bool, error) {
	if ok, err := p.DoFlush(); ok && err == nil {
		return true, nil
	}

	if ok, err := p.DoSync(); ok {
		return err == nil, err
	}

	return false, errors.New("unsupported: Flush")
}

func (p *FlushBypass) Close() error {
	if p.SetClosed() {
		_, _ = p.DoClose()
		return nil
	}
	return p.ClosedError()
}

func (p *FlushBypass) Flush() error {
	if p.IsClosed() {
		return p.ClosedError()
	}
	_, err := p.DoFlush()
	return err
}

func (p *FlushBypass) Sync() error {
	if p.IsClosed() {
		return p.ClosedError()
	}
	_, err := p.DoSync()
	return err
}

func (p *FlushBypass) FlushOrSync() error {
	if p.IsClosed() {
		return p.ClosedError()
	}
	_, err := p.DoFlushOrSync()
	return err
}
