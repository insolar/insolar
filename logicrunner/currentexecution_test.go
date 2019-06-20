package logicrunner

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/insolar/insolar/insolar/gen"
)

func TestTranscriptDequeue_Basic(t *testing.T) {
	a, r := assert.New(t), require.New(t)
	d := NewTranscriptDequeue()
	r.NotNil(d)

	d.Push(&Transcript{Nonce: 1}, &Transcript{Nonce: 2})

	tr := d.Pop()
	r.NotNil(tr)
	a.Equal(uint64(1), tr.Nonce)

	d.Prepend(&Transcript{Nonce: 3}, &Transcript{Nonce: 4})

	tr = d.Pop()
	r.NotNil(tr)
	a.Equal(uint64(3), tr.Nonce)

	tr = d.Pop()
	r.NotNil(tr)
	a.Equal(uint64(4), tr.Nonce)

	tr = d.Pop()
	r.NotNil(tr)
	a.Equal(uint64(2), tr.Nonce)

	a.Nil(d.Pop())
}

func TestTranscriptDequeue_Rotate(t *testing.T) {
	a, r := assert.New(t), require.New(t)
	d := NewTranscriptDequeue()
	r.NotNil(d)

	d.Push(&Transcript{Nonce: 1}, &Transcript{Nonce: 2})

	rotated := d.Rotate()
	r.Len(rotated, 2)

	a.Nil(d.Pop())

	rotated = d.Rotate()
	r.Len(rotated, 0)

	a.Nil(d.Pop())
}

func TestTranscriptDequeue_HasFromLedger(t *testing.T) {
	a, r := assert.New(t), require.New(t)
	d := NewTranscriptDequeue()
	r.NotNil(d)

	d.Prepend(&Transcript{Nonce: 3}, &Transcript{Nonce: 4})
	a.False(d.HasFromLedger() != nil)

	d.Push(&Transcript{FromLedger: true})
	a.True(d.HasFromLedger() != nil)

	d.Push(&Transcript{FromLedger: true})
	a.True(d.HasFromLedger() != nil)
}

func TestTranscriptDequeue_PopByReference(t *testing.T) {
	a, r := assert.New(t), require.New(t)
	d := NewTranscriptDequeue()
	r.NotNil(d)

	ref := gen.Reference()
	badRef := gen.Reference()

	d.Prepend(
		&Transcript{Nonce: 3, RequestRef: &badRef},
		&Transcript{Nonce: 4, RequestRef: &ref},
		&Transcript{Nonce: 5, RequestRef: &badRef},
	)

	tr := d.PopByReference(&ref)
	r.NotNil(tr)
	r.Equal(tr.RequestRef.Bytes(), ref.Bytes())

	tr = d.PopByReference(&ref)
	a.Nil(tr)

	a.Equal(d.Len(), 2)
}

func TestTranscriptDequeue_Take(t *testing.T) {
	a, r := assert.New(t), require.New(t)
	d := NewTranscriptDequeue()
	r.NotNil(d)

	for i := 0; i < 15; i++ {
		d.Push(&Transcript{Nonce: uint64(i)})
	}

	trs := d.Take(0)
	r.NotNil(d)
	a.Len(trs, 0)

	trs = d.Take(10)
	r.NotNil(d)
	a.Len(trs, 10)

	trs = d.Take(10)
	r.NotNil(d)
	a.Len(trs, 5)

	trs = d.Take(10)
	r.NotNil(d)
	a.Len(trs, 0)
}
