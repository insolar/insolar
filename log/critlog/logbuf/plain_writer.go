// +build ignore

package logbuf

import "io"

var _ io.Writer = &PlainBuffer{}
var _ io.WriterTo = &PlainBuffer{}

type MissedFunc func(missed uint32)

func NewPlainBuffer(buffer PagedBuffer, idleFn func(), missedFn MissedFunc) PlainBuffer {
	return PlainBuffer{
		buffer:   buffer,
		idleFn:   idleFn,
		missedFn: missedFn,
	}
}

type PlainBuffer struct {
	buffer   PagedBuffer
	idleFn   func()
	missedFn MissedFunc
}

func (p *PlainBuffer) WriteTo(w io.Writer) (int64, error) {
	pg := p.buffer.FlushPages()
	pg.StartAccess(p.idleFn)
	defer pg.StopAccess()

	c := pg.Count()
	totalN := int64(0)
	for i := 0; i < c; i++ {
		writes, buf := pg.Page(i)
		if buf == nil {
			if writes > 0 && p.missedFn != nil {
				p.missedFn(writes)
			}
			continue
		}
		nc, err := w.Write(buf)
		totalN += int64(nc)
		if err != nil {
			return totalN, err
		}
	}
	return totalN, nil
}

func (p *PlainBuffer) Write(b []byte) (n int, err error) {
	if len(b) == 0 {
		return 0, nil
	}
	pg, buf := p.buffer.allocateBuffer(uint32(len(b)))
	copy(buf, b)
	pg.stopAccess()
	return len(b), nil
}
