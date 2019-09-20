package critlog

import (
	"context"
	"github.com/insolar/insolar/insolar"
	"io"
	"time"
)

func NewBackpressureBuffer(output io.Writer, bufSize int, maxPenalty uint8, dropBufOnFatal bool) *BackpressureBuffer {
	if bufSize <= 1 {
		panic("illegal value")
	}

	return &BackpressureBuffer{
		output:         OutputHelper{output},
		dropBufOnFatal: dropBufOnFatal,
		maxPenalty:     maxPenalty,
		buffer:         make(chan bufEntry, bufSize),
	}
}

var _ insolar.LogLevelWriter = &BackpressureBuffer{}

type BackpressureBuffer struct {
	output OutputHelper
	fatal  FatalHelper

	buffer         chan bufEntry
	maxPenalty     uint8
	dropBufOnFatal bool
}

type bufEntry struct {
	lvl       insolar.LogLevel
	b         []byte
	flushMark bufferFlushMode
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
		if p.dropBufOnFatal {
			n, _ = p.doWrite(level, b, tillDepletion)
		} else {
			n, _ = p.output.DoWriteLevel(level, b)
		}
		return n, p.Close()

	case insolar.PanicLevel:
		n, err = p.doWrite(level, b, tillMark)
		if err != nil {
			_ = p.Flush()
			return n, err
		}
		return n, p.Flush()
	default:
		return p.doWrite(level, b, noFlush)
	}
}

type bufferFlushMode uint8

const (
	noFlush bufferFlushMode = iota
	tillMark
	tillDepletion
)

func (p *BackpressureBuffer) doWrite(level insolar.LogLevel, b []byte, flush bufferFlushMode) (int, error) {

	penalty := 0
	switch flush {
	case tillDepletion:
		// nothing
	case tillMark:
		p.buffer <- bufEntry{flushMark: tillMark}
	default:
		bufLen := len(p.buffer)
		if bufLen == 0 { // dirty check
			// direct write
			return p.output.DoWriteLevel(level, b)
		}
		penalty = int(p.maxPenalty+2) * len(p.buffer) / (1 + cap(p.buffer))
	}

	prevWasMark := false
	for i := 0; flush != noFlush || i < penalty; i++ {
		select {
		case be, ok := <-p.buffer:
			/*
				There is a chance that we will get a mark of someone else, but it is ok as long as
				the total count of flush writers and queued marks is equal.

				The full depletion writer must present the depletion mark before exiting.
			*/
			switch {
			case !ok:
				// break
			case be.flushMark == noFlush:
				prevWasMark = false
				_, _ = p.output.DoWriteLevel(be.lvl, be.b)
				continue
			case be.flushMark == tillDepletion:
				// return the mark and stop
				p.buffer <- be
			case flush == tillDepletion:
				p.buffer <- be
				if prevWasMark == true {
					time.Sleep(1 * time.Millisecond)
				} else {
					prevWasMark = true
				}
				continue
			case flush == tillMark:
				// break
			default:
				panic("illegal state")
			}
		default:
			if flush == tillDepletion {
				/* It will stay in the queue to signal other writers and the worker to stop */
				p.buffer <- bufEntry{flushMark: tillDepletion}
			}
		}
		return p.output.DoWriteLevel(level, b)
	}

	/*
		We paid our penalty and the queue didn't became empty.
		So we will leave our event for someone else to pick.
	*/
	p.buffer <- bufEntry{lvl: level, b: p.copyBytes(b)}
	return len(b), nil
}

func (p *BackpressureBuffer) worker(ctx context.Context) {
	prevWasMark := false
	for {
		select {
		case <-ctx.Done():
			return
		case be, ok := <-p.buffer:
			switch {
			case !ok:
				return
			case be.flushMark == noFlush:
				prevWasMark = false
				_, _ = p.output.DoWriteLevel(be.lvl, be.b)
			case be.flushMark == tillDepletion:
				// return the mark and stop
				p.buffer <- be
				return
			case be.flushMark == tillMark:
				/*
					Never take out the marks, otherwise the write will stuck.

					Presence of this mark also indicates that the queue is processed by the write,
					so this worker can hands off for a while.
				*/
				p.buffer <- be
				if prevWasMark == true {
					time.Sleep(10 * time.Millisecond)
				} else {
					prevWasMark = true
				}
			default:
				panic("illegal state")
			}
		}
	}
}

func (p *BackpressureBuffer) copyBytes(bytes []byte) []byte {
	var v []byte
	return append(v, bytes...)
}
