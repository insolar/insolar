package critlog

import (
	"context"
	"errors"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/args"
	"io"
	"math"
	"runtime"
	"sync/atomic"
	"time"
)

type BackpressureBufferFlags uint8

const (
	BufferDropOnFatal BackpressureBufferFlags = 1 << iota
	BufferWriteDelayFairness
	BufferDirectForRegular
	BufferTrackWriteDuration
	BufferReuse
)

func NewBackpressureBuffer(output io.Writer, bufSize int, extraPenalty uint8, maxParWrites uint8,
	flags BackpressureBufferFlags, missFn MissedEventFunc,
) *BackpressureBuffer {

	if bufSize <= 1 {
		panic("illegal value")
	}

	return &BackpressureBuffer{
		output:       NewFlushBypass(output),
		flags:        flags,
		extraPenalty: extraPenalty,
		maxParWrites: maxParWrites,
		missFn:       missFn,
		buffer:       make(chan bufEntry, bufSize),
	}
}

type MissedEventFunc func(missed int) (insolar.LogLevel, []byte)

var _ insolar.LogLevelWriter = &BackpressureBuffer{}

type BackpressureBuffer struct {
	output FlushBypass
	fatal  FatalHelper
	missFn MissedEventFunc

	buffer chan bufEntry

	writeSeq      uint32 // atomic
	pendingWrites uint32 // atomic
	missCount     uint32 // atomic
	avgDelayMicro uint32 // atomic

	maxParWrites uint8
	extraPenalty uint8
	flags        BackpressureBufferFlags
}

type bufEntry struct {
	lvl       insolar.LogLevel
	b         []byte
	start     int64
	flushMark bufferFlushMode
}

/* The buffer requires a worker to scrap the buffer. Multiple workers are ok, but aren't necessary. */
func (p *BackpressureBuffer) StartWorker(ctx context.Context) *BackpressureBuffer {
	go p.worker(ctx)
	return p
}

const internalOpLevel = insolar.LogLevel(255)

func (p *BackpressureBuffer) Close() error {
	if p.fatal.IsFatal() {
		if p.output.SetClosed() {
			_, _ = p.output.DoClose()
		}
		p.fatal.LockFatal()
		return nil
	}

	if !p.output.SetClosed() {
		return errors.New("closed")
	}

	_, _ = p.flushWrite(internalOpLevel, nil, tillDepletion, 0)
	_, _ = p.output.DoClose()
	return nil
}

func (p *BackpressureBuffer) Flush() error {
	_, err := p.flushWrite(internalOpLevel, nil, tillFlushMark, 0)
	_, _ = p.output.DoFlushOrSync()
	return err
}

func (p *BackpressureBuffer) Write(b []byte) (n int, err error) {
	return p.LogLevelWrite(insolar.NoLevel, b)
}

func (p *BackpressureBuffer) IsLowLatencySupported() bool {
	return true
}

func (p *BackpressureBuffer) GetBareOutput() io.Writer {
	return p.output.Writer
}

func (p *BackpressureBuffer) LowLatencyWrite(level insolar.LogLevel, b []byte) (n int, err error) {
	if p.fatal.IsFatal() {
		return p.fatal.PostFatalWrite(level, b)
	}

	if level == insolar.FatalLevel {
		return p.writeFatal(level, b)
	}

	be := p.newQueueEntry(level, b)

	for i := 0; ; i++ {
		select {
		case p.buffer <- be:
		default:
			if i < 5 {
				runtime.Gosched()
				continue
			}
			atomic.AddUint32(&p.missCount, 1)
		}
		break
	}
	return len(b), nil
}

func (p *BackpressureBuffer) newQueueEntry(level insolar.LogLevel, b []byte) bufEntry {
	if p.flags&BufferReuse != 0 {
		return bufEntry{lvl: level, b: b}
	}
	var v []byte
	return bufEntry{lvl: level, b: append(v, b...)}
}

func (p *BackpressureBuffer) LogLevelWrite(level insolar.LogLevel, b []byte) (n int, err error) {
	if p.fatal.IsFatal() {
		return p.fatal.PostFatalWrite(level, b)
	}

	startNano := time.Now().UnixNano()

	switch level {
	case insolar.FatalLevel:
		return p.writeFatal(level, b)

	case insolar.PanicLevel:
		n, err = p.flushWrite(level, b, tillFlushMark, 0)
		if err == nil {
			_, _ = p.output.DoFlushOrSync()
		}
		return n, err

	default:
		if p.flags&BufferDirectForRegular != 0 {
			return p.flushWrite(level, b, noFlush, startNano)
		}
		return p.checkWrite(level, b, startNano)
	}
}

func (p *BackpressureBuffer) writeFatal(level insolar.LogLevel, b []byte) (n int, err error) {
	if !p.fatal.SetFatal() {
		return p.fatal.PostFatalWrite(level, b)
	}
	if p.flags&BufferDropOnFatal != 0 {
		n, _ = p.flushWrite(level, b, tillDepletion, 0)
	} else {
		n, _ = p.directWrite(level, b, 0)
	}
	p.writeMissedCount(p.getMissedCount() + len(p.buffer))

	if ok, err := p.output.DoFlushOrSync(); ok && err == nil {
		return n, nil
	}

	return n, p.output.Close()
}

func (p *BackpressureBuffer) checkWrite(level insolar.LogLevel, b []byte, startNano int64) (int, error) {
	writeId := atomic.AddUint32(&p.writeSeq, 1)

	for i := uint8(0); ; i++ {
		pendingWrites := atomic.LoadUint32(&p.pendingWrites)

		if pendingWrites >= uint32(p.maxParWrites) || !p.drawStraw(writeId, pendingWrites) {
			return p.fairQueueWrite(level, b, startNano)
		}

		if atomic.CompareAndSwapUint32(&p.pendingWrites, pendingWrites, pendingWrites+1) {
			break
		}

		if i >= (1+p.maxParWrites)<<1 {
			// too many retries
			return p.fairQueueWrite(level, b, startNano)
		}
		runtime.Gosched()
	}

	defer atomic.AddUint32(&p.pendingWrites, ^uint32(0)) // -1
	return p.flushWrite(level, b, noFlush, startNano)
}

type bufferFlushMode uint8

const (
	noFlush bufferFlushMode = iota
	tillFlushMark
	tillDepletion
)

func (p *BackpressureBuffer) fairQueueWrite(level insolar.LogLevel, b []byte, startNano int64) (int, error) {
	n, err := p.queueWrite(level, b)
	if startNano != 0 && p.flags&BufferWriteDelayFairness != 0 {
		waitNano := int64(p.GetWriteDuration()) - (time.Now().UnixNano() - startNano)
		if waitNano > 0 {
			time.Sleep(time.Duration(waitNano))
		}
	}
	return n, err
}

func (p *BackpressureBuffer) queueWrite(level insolar.LogLevel, b []byte) (int, error) {
	p.buffer <- p.newQueueEntry(level, b)
	return len(b), nil
}

func (p *BackpressureBuffer) directWrite(level insolar.LogLevel, b []byte, startNano int64) (int, error) {
	if level == internalOpLevel && b == nil {
		return 0, nil
	}
	n, err := p.output.DoLevelWrite(level, b)
	if startNano > 0 && p.flags&BufferTrackWriteDuration != 0 {
		writeDuration := time.Now().UnixNano() - startNano
		p.ApplyWriteDuration(time.Duration(writeDuration))
	}
	return n, err
}

func (p *BackpressureBuffer) flushWrite(level insolar.LogLevel, b []byte, flush bufferFlushMode, startNano int64) (int, error) {

	penalty := 1 // every worker has to write at least +1 event from the queue
	switch flush {
	case tillDepletion:
		// nothing
	case tillFlushMark:
		p.buffer <- bufEntry{flushMark: tillFlushMark}
	default:
		bufLen := len(p.buffer)
		if bufLen == 0 { // dirty check
			// direct write
			return p.directWrite(level, b, startNano)
		}
		// extra penalty is added proportionally to queue occupation
		penalty += int(p.extraPenalty+1) * len(p.buffer) / (1 + cap(p.buffer))
	}

	hasDepletionMark := false
	prevWasFlushMark := false
	for i := 0; flush != noFlush || i <= penalty; i++ {
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
				prevWasFlushMark = false
				_, _ = p.directWrite(be.lvl, be.b, startNano)
				continue
			case be.flushMark == tillDepletion:
				// return the mark and stop
				hasDepletionMark = true
				p.buffer <- be
				// break
			case be.flushMark == tillFlushMark:
				prevWasFlushMark = true
				switch flush {
				case tillDepletion:
					/* we don't need it - put it back for another worker */
					p.buffer <- be
					if prevWasFlushMark == true {
						time.Sleep(1 * time.Millisecond)
					} else {
						prevWasFlushMark = true
					}
					continue
				case noFlush:
					/* we don't need it - put it back for another worker */
					p.buffer <- be
					continue
				}
				// break
			default:
				panic("illegal state")
			}
		default:
			if flush == tillDepletion && !hasDepletionMark {
				/* It will stay in the queue to signal other writers and the worker to stop */
				p.buffer <- bufEntry{flushMark: tillDepletion}
			}
		}
		return p.directWrite(level, b, startNano)
	}

	/*
		We paid our penalty and the queue didn't became empty.
		Lets leave our event for someone else to pick.
	*/
	return p.queueWrite(level, b)
}

func (p *BackpressureBuffer) worker(ctx context.Context) {

	atomic.AddUint32(&p.pendingWrites, 1)
	defer atomic.AddUint32(&p.pendingWrites, ^uint32(0)) // -1

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
				_, _ = p.directWrite(be.lvl, be.b, be.start)
			case be.flushMark == tillDepletion:
				// return the mark and stop
				p.buffer <- be
				return
			case be.flushMark == tillFlushMark:
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
		p.getAndWriteMissed()
	}
}

func (p *BackpressureBuffer) drawStraw(writeId uint32, writersInQueue uint32) bool {
	return writersInQueue == 0 || (writeId%args.Prime(int(writersInQueue-1))) == 0
}

func (p *BackpressureBuffer) getMissedCount() int {
	return int(atomic.SwapUint32(&p.missCount, 0))
}

func (p *BackpressureBuffer) getAndWriteMissed() {
	if p.missFn == nil || p.output.IsClosed() || p.fatal.IsFatal() {
		return
	}
	p.writeMissedCount(p.getMissedCount())
}

func (p *BackpressureBuffer) writeMissedCount(missedCount int) {
	if p.missFn == nil || missedCount == 0 {
		return
	}
	lvl, missMsg := p.missFn(missedCount)
	if lvl == insolar.NoLevel || len(missMsg) == 0 {
		return
	}
	_, _ = p.output.DoLevelWrite(lvl, missMsg)
}

func (p *BackpressureBuffer) GetWriteDuration() time.Duration {
	return time.Duration(atomic.LoadUint32(&p.avgDelayMicro)) * time.Microsecond
}

func (p *BackpressureBuffer) ApplyWriteDuration(d time.Duration) {
	for {
		v := atomic.LoadUint32(&p.avgDelayMicro)

		vv := uint64(d / time.Microsecond)
		switch {
		case vv == 0:
			vv = 1
		case vv > math.MaxUint32:
			vv = math.MaxUint32
		}

		if v != 0 {
			vv = (vv + uint64(v)) >> 1
		}

		if atomic.CompareAndSwapUint32(&p.avgDelayMicro, v, uint32(vv)) {
			return
		}
	}
}
