// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package critlog

import (
	"io"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/log/logoutput"
)

func NewFatalDirectWriter(output *logoutput.Adapter) *FatalDirectWriter {
	if output == nil {
		panic("illegal value")
	}

	return &FatalDirectWriter{
		output: output,
	}
}

var _ insolar.LogLevelWriter = &FatalDirectWriter{}
var _ io.WriteCloser = &FatalDirectWriter{}

type FatalDirectWriter struct {
	output *logoutput.Adapter
}

func (p *FatalDirectWriter) Close() error {
	return p.output.Close()
}

func (p *FatalDirectWriter) Flush() error {
	return p.output.Flush()
}

func (p *FatalDirectWriter) Write(b []byte) (n int, err error) {
	return p.output.Write(b)
}

func (p *FatalDirectWriter) LowLatencyWrite(level insolar.LogLevel, b []byte) (int, error) {
	return p.LogLevelWrite(level, b)
}

func (p *FatalDirectWriter) IsLowLatencySupported() bool {
	return false
}

func (p *FatalDirectWriter) LogLevelWrite(level insolar.LogLevel, b []byte) (n int, err error) {
	switch level {
	case insolar.FatalLevel:
		if !p.output.SetFatal() {
			break
		}
		n, _ = p.output.DirectLevelWrite(level, b)
		_ = p.output.DirectFlushFatal()
		return n, nil

	case insolar.PanicLevel:
		n, err = p.output.LogLevelWrite(level, b)
		_ = p.output.Flush()
		return n, err
	}
	return p.output.LogLevelWrite(level, b)
}
