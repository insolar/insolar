// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package logoutput

import (
	"errors"
	"io"
	"sync/atomic"

	"github.com/insolar/insolar/insolar"
)

type LogFlushFunc func() error

func NewAdapter(output io.Writer, protectedClose bool, flushFn, fatalFlushFn LogFlushFunc) *Adapter {
	flags := adapterState(0)
	if protectedClose {
		flags |= adapterProtectClose
	}

	if w, ok := output.(insolar.LogLevelWriter); ok {
		return &Adapter{output: w, flushFn: flushFn, state: uint32(flags)}
	}

	return &Adapter{output: writerAdapter{output}, flushFn: flushFn, fatalFlushFn: fatalFlushFn, state: uint32(flags)}
}

var errClosed = errors.New("closed")

type Adapter struct {
	output       insolar.LogLevelWriter
	flushFn      LogFlushFunc
	fatalFlushFn LogFlushFunc
	state        uint32 // atomic
}

type adapterState uint32

const (
	adapterClosed adapterState = 1 << iota
	adapterFatal

	adapterProtectClose
	adapterPanicOnFatal // for test usage
)

func (p *Adapter) getState() adapterState {
	return adapterState(atomic.LoadUint32(&p.state))
}

func (p *Adapter) setState(flags adapterState) adapterState {
	for {
		s := p.getState()
		if s&flags == flags {
			return s
		}
		if atomic.CompareAndSwapUint32(&p.state, uint32(s), uint32(s|flags)) {
			return s
		}
	}
}

func (p *Adapter) applyState() (ok bool, err error) {
	prev := p.getState()
	if prev&adapterFatal != 0 {
		p.LockFatal()
	}
	if prev&adapterClosed != 0 {
		return false, errClosed
	}
	return prev&adapterFatal == 0, nil
}

func (p *Adapter) DirectFlushFatal() error {
	prev := p.setState(adapterFatal | adapterClosed)
	if prev&adapterClosed != 0 {
		return errClosed
	}

	err := p.output.Flush()
	if p.fatalFlushFn != nil {
		err = p.fatalFlushFn()
	}

	if err != nil && prev&adapterProtectClose == 0 {
		err = p.output.Close()
	}
	return err
}

func (p *Adapter) _directClose(prev adapterState) error {
	if prev&adapterClosed != 0 {
		return errClosed
	}
	if prev&adapterProtectClose != 0 {
		return nil
	}
	return p.output.Close()
}

func (p *Adapter) DirectClose() error {
	prev := p.setState(adapterClosed)
	return p._directClose(prev &^ adapterClosed)
}

func (p *Adapter) Close() error {
	prev := p.setState(adapterClosed)
	if prev&adapterFatal != 0 {
		defer p.LockFatal()
	}
	return p._directClose(prev)
}

func (p *Adapter) Flush() error {
	if ok, err := p.applyState(); !ok {
		return err
	}
	if p.flushFn != nil {
		_ = p.output.Flush()
		return p.flushFn()
	}
	return p.output.Flush()
}

func (p *Adapter) Write(b []byte) (int, error) {
	if ok, err := p.applyState(); !ok {
		return 0, err
	}
	return p.output.Write(b)
}

func (p *Adapter) LogLevelWrite(level insolar.LogLevel, b []byte) (int, error) {
	if ok, err := p.applyState(); !ok {
		return 0, err
	}
	return p.output.LogLevelWrite(level, b)
}

func (p *Adapter) DirectLevelWrite(level insolar.LogLevel, b []byte) (int, error) {
	return p.output.LogLevelWrite(level, b)
}

func (p *Adapter) SetClosed() bool {
	return p.setState(adapterClosed)&adapterClosed == 0
}

func (p *Adapter) IsClosed() bool {
	return p.getState()&adapterClosed != 0
}

func (p *Adapter) SetFatal() bool {
	return p.setState(adapterFatal)&adapterFatal == 0
}

func (p *Adapter) IsFatal() bool {
	return p.getState()&adapterFatal != 0
}

func (p *Adapter) LockFatal() {
	if p.getState()&adapterPanicOnFatal != 0 {
		panic("fatal lock")
	}
	select {} // lock it down forever
}

/* =============================  */

var _ insolar.LogLevelWriter = &writerAdapter{}

type writerAdapter struct {
	output io.Writer
}

func (p writerAdapter) Close() error {
	if f, ok := p.output.(io.Closer); ok {
		return f.Close()
	}
	return errors.New("unsupported: Close")
}

func (p writerAdapter) Flush() error {
	type flusher interface {
		Flush() error
	}
	type syncer interface {
		Sync() error
	}

	switch ww := p.output.(type) {
	case flusher:
		return ww.Flush()
	case syncer:
		return ww.Sync()
	}

	return errors.New("unsupported: Flush")
}

func (p writerAdapter) Write(b []byte) (int, error) {
	return p.output.Write(b)
}

func (p writerAdapter) LogLevelWrite(_ insolar.LogLevel, b []byte) (int, error) {
	return p.output.Write(b)
}
