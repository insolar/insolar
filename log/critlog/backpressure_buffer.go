package critlog

import (
	"context"
	"errors"
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/network/consensus/common/args"
	"io"
	"runtime"
	"sync/atomic"
	"time"
)

type BackpressureBufferFlags uint8

const (
	// Buffer content will not be flushed on fatal, instead a "missing X" message will be added
	BufferDropOnFatal BackpressureBufferFlags = 1 << iota
	// Buffer may apply additional delay to writes done into a queue to equalize timings
	// This mode requires either BufferTrackWriteDuration flag or use of SetAvgWriteDuration
	// This flag has no effect when BufferBypassForRegular is set
	BufferWriteDelayFairness
	// With this flags the buffer will update GetAvgWriteDuration with every regular write
	BufferTrackWriteDuration
	// Regular (not-lowLatency) writes will go directly to output, ignoring queue and parallel write limits
	BufferBypassForRegular
	// When a worker is started, but all links to BackpressureBuffer were lost, then the worker will be stopped.
	// And with this flag present, the buffer (and its underlying) will also be closed.
	BufferCloseOnStop
	// USE WITH CAUTION! This flag enables to use argument of Write([]byte) outside of the call.
	// This is AGAINST existing conventions and MUST ONLY be used when a writer's code is proprietary and never reuses the argument.
	BufferReuse
)

func NewBackpressureBuffer(output io.Writer, bufSize int, extraPenalty uint8, maxParWrites uint8,
	flags BackpressureBufferFlags, missFn MissedEventFunc,
) *BackpressureBuffer {

	if bufSize <= 1 {
		panic("illegal value")
	}

	if maxParWrites == 1 {
		// fairness is slower with single writer
		flags &^= BufferWriteDelayFairness
	}

	return &BackpressureBuffer{&internalBackpressureBuffer{
		output:       NewFlushBypass(output),
		flags:        flags,
		extraPenalty: extraPenalty,
		maxParWrites: maxParWrites,
		missFn:       missFn,
		buffer:       make(chan bufEntry, bufSize),
	}, 0}
}

type MissedEventFunc func(missed int) (insolar.LogLevel, []byte)

var _ insolar.LogLevelWriter = &internalBackpressureBuffer{}

/*
Provides weak-reference behavior to enable auto-stop of workers
*/
type BackpressureBuffer struct {
	*internalBackpressureBuffer
	status uint32 // atomic
}

type internalBackpressureBuffer struct {
	output FlushBypass
	fatal  FatalHelper
	missFn MissedEventFunc

	buffer chan bufEntry

	writeSeq      uint32        // atomic
	pendingWrites uint32        // atomic
	missCount     uint32        // atomic
	avgDelayNano  time.Duration // atomic

	maxParWrites uint8
	extraPenalty uint8
	flags        BackpressureBufferFlags
}

type bufEntry struct {
	lvl       insolar.LogLevel
	b         []byte
	start     int64
	flushMark bufferMark
}

/*
The buffer requires a worker to scrap the buffer. Multiple workers are ok, but aren't necessary.
Start of the worker will also attach a finalizer to the buffer.
*/
func (p *BackpressureBuffer) StartWorker(ctx context.Context) *BackpressureBuffer {
	go p.worker(ctx)

	if atomic.AddUint32(&p.status, 1) != 1 {
		return p
	}

	if p.flags&BufferCloseOnStop != 0 {
		runtime.SetFinalizer(p, func(pp *BackpressureBuffer) {
			_ = pp.Close()
		})
	} else {
		runtime.SetFinalizer(p, func(pp *BackpressureBuffer) {
			pp.closeWorker()
		})
	}

	return p
}

const internalOpLevel = insolar.LogLevel(255)

func (p *internalBackpressureBuffer) Close() error {
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

	_, _ = p.flushTillDepletion(internalOpLevel, nil, 0)
	_, _ = p.output.DoClose()
	return nil
}

func (p *internalBackpressureBuffer) closeWorker() {
	if p.fatal.IsFatal() || p.output.IsClosed() {
		return
	}
	_, _ = p.flushTillDepletion(internalOpLevel, nil, 0)
}

// NB! Flush() may NOT be able to clean up whole buffer when there are too many pending writers
func (p *internalBackpressureBuffer) Flush() error {
	if p.fatal.IsFatal() || p.output.IsClosed() {
		return nil
	}

	_, err := p.flushTillMark(internalOpLevel, nil, 0)
	_, _ = p.output.DoFlushOrSync()
	return err
}

func (p *internalBackpressureBuffer) Write(b []byte) (n int, err error) {
	return p.LogLevelWrite(insolar.NoLevel, b)
}

func (p *internalBackpressureBuffer) IsLowLatencySupported() bool {
	return true
}

func (p *internalBackpressureBuffer) GetBareOutput() io.Writer {
	return p.output.Writer
}

func (p *internalBackpressureBuffer) LowLatencyWrite(level insolar.LogLevel, b []byte) (n int, err error) {
	if p.fatal.IsFatal() {
		return p.fatal.PostFatalWrite(level, b)
	}

	if level == insolar.FatalLevel {
		return p.writeFatal(level, b)
	}

	be := p.newQueueEntry(level, b, time.Now().UnixNano())

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

func (p *internalBackpressureBuffer) newQueueEntry(level insolar.LogLevel, b []byte, startNano int64) bufEntry {
	if p.flags&BufferReuse != 0 {
		return bufEntry{lvl: level, b: b, start: startNano}
	}
	var v []byte
	return bufEntry{lvl: level, b: append(v, b...), start: startNano}
}

func (p *internalBackpressureBuffer) LogLevelWrite(level insolar.LogLevel, b []byte) (n int, err error) {
	if p.fatal.IsFatal() {
		return p.fatal.PostFatalWrite(level, b)
	}

	startNano := time.Now().UnixNano()

	switch level {
	case insolar.FatalLevel:
		return p.writeFatal(level, b)

	case insolar.PanicLevel:
		n, err = p.flushTillMark(level, b, 0)
		if err == nil {
			_, _ = p.output.DoFlushOrSync()
		}
		return n, err

	default:
		if p.flags&BufferBypassForRegular != 0 {
			return p.flushWrite(level, b, startNano)
		}
		return p.checkWrite(level, b, startNano)
	}
}

func (p *internalBackpressureBuffer) writeFatal(level insolar.LogLevel, b []byte) (n int, err error) {
	if !p.fatal.SetFatal() {
		return p.fatal.PostFatalWrite(level, b)
	}
	if p.flags&BufferDropOnFatal != 0 {
		n, _ = p.flushTillDepletion(level, b, 0)
	} else {
		n, _ = p.directWrite(level, b, 0)
	}
	p.writeMissedCount(p.getAndResetMissedCount() + len(p.buffer))

	if ok, err := p.output.DoFlushOrSync(); ok && err == nil {
		return n, nil
	}

	return n, p.output.Close()
}

func (p *internalBackpressureBuffer) checkWrite(level insolar.LogLevel, b []byte, startNano int64) (int, error) {
	writeSeq := atomic.AddUint32(&p.writeSeq, 1)

	for i := uint8(0); ; i++ {
		pendingWrites := atomic.LoadUint32(&p.pendingWrites)

		if pendingWrites >= uint32(p.maxParWrites) || !p.drawStraw(writeSeq, pendingWrites) {
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
	return p.flushWrite(level, b, startNano)
}

type bufferMark uint8

const (
	noMark bufferMark = iota
	flushMark
	depletionMark
)

func (p *internalBackpressureBuffer) fairQueueWrite(level insolar.LogLevel, b []byte, startNano int64) (int, error) {
	waitNano := int64(p.GetAvgWriteDuration())
	n, err := p.queueWrite(level, b, startNano)

	if n > 0 && err == nil && startNano > 0 && p.flags&BufferWriteDelayFairness != 0 {
		waitNano -= time.Now().UnixNano() - startNano
		if waitNano > 0 {
			time.Sleep(time.Duration(waitNano))
		}
	}

	return n, err
}

func (p *internalBackpressureBuffer) queueWrite(level insolar.LogLevel, b []byte, startName int64) (int, error) {
	p.buffer <- p.newQueueEntry(level, b, startName)
	return len(b), nil
}

func (p *internalBackpressureBuffer) directWrite(level insolar.LogLevel, b []byte, startNano int64) (int, error) {
	n, err := p.output.DoLevelWrite(level, b)

	if n > 0 && err == nil && startNano > 0 && p.flags&BufferTrackWriteDuration != 0 {
		writeDuration := time.Now().UnixNano() - startNano
		p.updateWriteDuration(time.Duration(writeDuration))
	}
	return n, err
}

func (p *internalBackpressureBuffer) flushWrite(level insolar.LogLevel, b []byte, startNano int64) (int, error) {

	bufLen := len(p.buffer)
	if bufLen == 0 { // dirty check
		return p.directWrite(level, b, startNano)
	}

	// every worker has to write at least +1 event from the queue
	// extra penalty is added proportionally to queue occupation
	penalty := 1 + int(p.extraPenalty+1)*len(p.buffer)/(1+cap(p.buffer))

	hasPrevFlushMark := false
	for i := 0; i <= penalty; i++ {
		isContinue, markType, err := p.writeOneFromQueue(noMark)

		switch {
		case err != nil:
			return 0, err
		case !isContinue:
			return p.directWrite(level, b, startNano)
		case markType == noMark:
			hasPrevFlushMark = false
		case markType != flushMark:
			panic("illegal state")
		case hasPrevFlushMark:
			time.Sleep(1 * time.Millisecond)
		default:
			hasPrevFlushMark = true
		}
	}
	/*
		We paid our penalty and the queue didn't became empty.
		Lets leave our event for someone else to pick.
	*/
	return p.queueWrite(level, b, startNano)
}

func (p *internalBackpressureBuffer) flushTillMark(level insolar.LogLevel, b []byte, startNano int64) (int, error) {

	hasPrevFlushMark := false
	hasAddedFlushMark := false

	markEntry := bufEntry{flushMark: flushMark}

outer:
	for maxFlushCount := len(p.buffer); maxFlushCount >= 0; maxFlushCount-- {
		select {
		case p.buffer <- markEntry:
			hasAddedFlushMark = true
			break outer
		default:
			isContinue, markType, err := p.writeOneFromQueue(noMark)
			switch {
			case err != nil:
				return 0, err
			case !isContinue:
				break outer
			case markType == noMark:
				hasPrevFlushMark = false
			case markType != flushMark:
				panic("illegal state")
			case hasPrevFlushMark:
				time.Sleep(1 * time.Millisecond)
			default:
				hasPrevFlushMark = true
			}
		}
	}

	if hasAddedFlushMark {
		// clean up till the mark
		for {
			isContinue, markType, err := p.writeOneFromQueue(flushMark)
			switch {
			case err != nil:
				return 0, err
			case isContinue:
				continue
			case markType == noMark:
				// buffer is empty ... another worker pulled the marker out
				time.Sleep(1 * time.Millisecond)
				continue
			}
			break
		}
	}

	p.getAndWriteMissed()
	if level == internalOpLevel && b == nil {
		return 0, nil
	}
	return p.directWrite(level, b, startNano)
}

func (p *internalBackpressureBuffer) flushTillDepletion(level insolar.LogLevel, b []byte, startNano int64) (int, error) {

	prevWasFlushMark := false
	markEntry := bufEntry{flushMark: depletionMark}

outer:
	for {
		isContinue, markType, err := p.writeOneFromQueue(depletionMark)
		switch {
		case err != nil:
			return 0, err
		case isContinue:
			// buffer is not empty, continue
		case markType != depletionMark:
			select {
			case p.buffer <- markEntry:
				//
			default:
				// buffer should be empty by now .. but lets try again
				continue
			}
			fallthrough
		default:
			break outer
		}

		switch {
		case markType == noMark:
			prevWasFlushMark = false
		case markType != flushMark:
			panic("illegal state")
		case prevWasFlushMark:
			// let flusher(s) do the job first
			time.Sleep(1 * time.Millisecond)
		default:
			prevWasFlushMark = true
		}
	}

	p.getAndWriteMissed()

	var n int
	var err error
	if level != internalOpLevel || b != nil {
		n, err = p.directWrite(level, b, startNano)
	}

	for atomic.LoadUint32(&p.pendingWrites) > 0 {
		time.Sleep(time.Millisecond)
	}
	return n, err
}

func (p *internalBackpressureBuffer) writeOneFromQueue(flush bufferMark) (bool, bufferMark, error) {
	select {
	case be, ok := <-p.buffer:
		/*
			There is a chance that we will get a mark of someone else, but it is ok as long as
			the total count of flush writers and queued marks is equal.

			The full depletion writer must present the depletion mark before exiting.
		*/
		switch {
		case !ok:
			return false, 0, errors.New("buffer channel is closed")

		case be.flushMark == noMark:
			_, _ = p.directWrite(be.lvl, be.b, be.start)
			return true, noMark, nil

		case be.flushMark == depletionMark:
			// restore the mark and stop
			p.buffer <- be
			return false, depletionMark, nil

		case be.flushMark == flushMark:
			switch flush {
			case flushMark:
				// consume the mark and stop
				return false, flushMark, nil
			case depletionMark, noMark:
				/* we don't need flushMark - put it back for another worker */
				p.buffer <- be
			default:
				panic("illegal state")
			}
			return true, flushMark, nil

		default:
			panic("illegal state")
		}
	default:
		return false, noMark, nil
	}
}

func (p *internalBackpressureBuffer) worker(ctx context.Context) {

	atomic.AddUint32(&p.pendingWrites, 1)
	defer atomic.AddUint32(&p.pendingWrites, ^uint32(0)) // -1

	prevWasMark := false
	for {
		p.getAndWriteMissed()

		select {
		case <-ctx.Done():
			return
		case be, ok := <-p.buffer:
			switch {
			case !ok:
				return
			case be.flushMark == noMark:
				prevWasMark = false
				_, _ = p.directWrite(be.lvl, be.b, be.start)
			case be.flushMark == depletionMark:
				// make sure to clean up the queue
				select {
				case be2, ok := <-p.buffer:
					if !ok {
						return
					}
					p.buffer <- be2
					if be2.flushMark != depletionMark {
						p.buffer <- be
					}
					continue
				default:
				}
				// return the mark and stop
				p.buffer <- be
				return
			case be.flushMark == flushMark:
				/*
					Never take out the marks, otherwise the write will stuck.

					Presence of this mark also indicates that the queue is processed by the write,
					so this worker can hands off for a while.
				*/
				p.buffer <- be
				if prevWasMark {
					time.Sleep(1 * time.Millisecond)
				} else {
					prevWasMark = true
				}
			default:
				panic("illegal state")
			}
		}
	}
}

func (p *internalBackpressureBuffer) drawStraw(writerSeq uint32, writersInQueue uint32) bool {
	return writersInQueue == 0 || (writerSeq%args.Prime(int(writersInQueue-1))) == 0
}

func (p *internalBackpressureBuffer) getAndResetMissedCount() int {
	return int(atomic.SwapUint32(&p.missCount, 0))
}

func (p *internalBackpressureBuffer) getAndWriteMissed() {
	if p.missFn == nil || p.output.IsClosed() || p.fatal.IsFatal() {
		return
	}
	p.writeMissedCount(p.getAndResetMissedCount())
}

func (p *internalBackpressureBuffer) writeMissedCount(missedCount int) {
	if p.missFn == nil || missedCount == 0 {
		return
	}
	lvl, missMsg := p.missFn(missedCount)
	if lvl == insolar.NoLevel || len(missMsg) == 0 {
		return
	}
	_, _ = p.output.DoLevelWrite(lvl, missMsg)
}

func (p *internalBackpressureBuffer) GetAvgWriteDuration() time.Duration {
	return time.Duration(atomic.LoadInt64((*int64)(&p.avgDelayNano)))
}

func (p *internalBackpressureBuffer) SetAvgWriteDuration(d time.Duration) {
	if d < 0 {
		d = 0
	}
	atomic.StoreInt64((*int64)(&p.avgDelayNano), int64(d))
}

func (p *internalBackpressureBuffer) updateWriteDuration(d time.Duration) {
	if d <= 0 {
		return
	}
	for {
		v := p.GetAvgWriteDuration()
		vv := d
		if v > 0 {
			vv = (vv + 3*v) >> 2
		}

		if atomic.CompareAndSwapInt64((*int64)(&p.avgDelayNano), int64(v), int64(vv)) {
			return
		}
	}
}