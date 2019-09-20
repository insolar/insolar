package critlog

import "github.com/insolar/insolar/insolar"

//var _ insolar.LogLevelWriter = &BackpressureBuffer{}

type BackpressureBuffer struct {
	output OutputHelper
	fatal  FatalHelper

	dropBufOnFatal bool

	buffer  chan bufEntry
	status  uint32 // atomic 0 = working, 1 = fatal, 2 = closed
	penalty int
}

type bufEntry struct {
	lvl insolar.LogLevel
	b   []byte
}

func (p *BackpressureBuffer) Close() error {
	return p.output.DoClose()
}

func (p *BackpressureBuffer) Flush() error {
	_ = p.output.DoFlush()
	return nil
}

func (p *BackpressureBuffer) Write(b []byte) (n int, err error) {
	return p.LogLevelWrite(insolar.NoLevel, b)
}

func (p *BackpressureBuffer) LogLevelWrite(level insolar.LogLevel, b []byte) (n int, err error) {
	if p.fatal.IsFatal() {
		return p.fatal.PostFatalWrite(level, b)
	}

	switch level {
	case insolar.FatalLevel:
		if !p.fatal.SetFatal() {
			return p.fatal.PostFatalWrite(level, b)
		}
		if !p.dropBufOnFatal {
			// TODO send the buffer out
		}
		n, _ = p.output.DoWriteLevel(level, b)
		return n, p.Close()

	case insolar.PanicLevel:
		n, err = p.doWrite(level, b, true)
		if err != nil {
			_ = p.Flush()
			return n, err
		}
		return n, p.Flush()
	default:
		return p.doWrite(level, b, false)
	}

}

func (p *BackpressureBuffer) doWrite(level insolar.LogLevel, b []byte, tillDepletion bool) (int, error) {
	if len(p.buffer) == 0 { // dirty check
		// direct write
		return p.output.DoWriteLevel(level, b)
	}

	for i := 0; tillDepletion || i <= p.penalty; i++ {
		select {
		case be, ok := <-p.buffer:
			if !ok {
				return p.output.DoWriteLevel(level, b)
			}
			_, _ = p.output.DoWriteLevel(be.lvl, be.b)
			continue
		default:
			if len(p.buffer) == 0 { // dirty check
				// direct write
				return p.output.DoWriteLevel(level, b)
			}
		}
		break
	}

	//defer func() {
	//	recovered :=
	//}()

	p.buffer <- bufEntry{level, p.copyBytes(b)}
	return len(b), nil
}

func (p *BackpressureBuffer) write(level insolar.LogLevel, bytes []byte) (int, error) {
	panic("not implemented")
}

func (p *BackpressureBuffer) copyBytes(bytes []byte) []byte {
	var v []byte
	return append(v, bytes...)
}
