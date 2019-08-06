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

package transcriptdequeue

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/insolar/insolar/insolar/gen"
	"github.com/insolar/insolar/logicrunner/transcript"
)

type TranscriptDequeueSuite struct{ suite.Suite }

func TestTranscriptDequeue(t *testing.T) { suite.Run(t, new(TranscriptDequeueSuite)) }

func (s *TranscriptDequeueSuite) TestBasic() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	// [] + [1, 2]
	d.Push(&transcript.Transcript{Nonce: 1}, &transcript.Transcript{Nonce: 2})

	// 1, [2]
	tr := d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(1), tr.Nonce)

	// [3, 4] + [2]
	d.Prepend(&transcript.Transcript{Nonce: 3}, &transcript.Transcript{Nonce: 4})

	// 3, [4, 2]
	tr = d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(3), tr.Nonce)

	// 4, [2]
	tr = d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(4), tr.Nonce)

	// 2, []
	tr = d.Pop()
	s.Require().NotNil(tr)
	s.Equal(uint64(2), tr.Nonce)

	// nil, []
	s.Nil(d.Pop())
}

func (s *TranscriptDequeueSuite) TestRotate() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	d.Push(&transcript.Transcript{Nonce: 1}, &transcript.Transcript{Nonce: 2})

	rotated := d.Rotate()
	s.Require().Len(rotated, 2)

	s.Nil(d.Pop())

	rotated = d.Rotate()
	s.Require().Len(rotated, 0)

	s.Nil(d.Pop())
}

func (s *TranscriptDequeueSuite) TestHasFromLedger() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	d.Prepend(&transcript.Transcript{Nonce: 3}, &transcript.Transcript{Nonce: 4})
	s.False(d.HasFromLedger() != nil)

	d.Push(&transcript.Transcript{FromLedger: true})
	s.True(d.HasFromLedger() != nil)

	d.Push(&transcript.Transcript{FromLedger: true})
	s.True(d.HasFromLedger() != nil)
}

func (s *TranscriptDequeueSuite) TestPopByReference() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	ref1, ref2, ref3 := gen.Reference(), gen.Reference(), gen.Reference()

	d.Prepend(
		&transcript.Transcript{Nonce: 3, RequestRef: ref1},
		&transcript.Transcript{Nonce: 4, RequestRef: ref2},
		&transcript.Transcript{Nonce: 5, RequestRef: ref3},
	)

	tr := d.PopByReference(ref2)
	s.Require().NotNil(tr)
	s.Require().Equal(tr.RequestRef.Bytes(), ref2.Bytes())

	tr = d.PopByReference(ref2)
	s.Nil(tr)

	s.Nil(d.first.prev)
	s.Equal(d.first.next, d.last)
	s.Nil(d.last.next)
	s.Equal(d.last.prev, d.first)
	s.Equal(d.first.value.Nonce, uint64(3))
	s.Equal(d.last.value.Nonce, uint64(5))

	s.Equal(d.Length(), 2)
}

func (s *TranscriptDequeueSuite) TestPopByReferenceHead() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	ref1, ref2, ref3 := gen.Reference(), gen.Reference(), gen.Reference()
	el1 := &transcript.Transcript{Nonce: 3, RequestRef: ref1}
	el2 := &transcript.Transcript{Nonce: 4, RequestRef: ref2}
	el3 := &transcript.Transcript{Nonce: 5, RequestRef: ref3}
	d.Prepend(el1, el2, el3)

	tr := d.PopByReference(ref1)
	s.Require().NotNil(tr)
	s.Require().Equal(tr.RequestRef.Bytes(), ref1.Bytes())

	s.Nil(d.first.prev)
	s.Equal(d.first.next, d.last)
	s.Nil(d.last.next)
	s.Equal(d.last.prev, d.first)
	s.Equal(d.first.value.Nonce, uint64(4))
	s.Equal(d.last.value.Nonce, uint64(5))

	s.Equal(d.Length(), 2)
}

func (s *TranscriptDequeueSuite) TestPopByReferenceTail() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	ref1, ref2, ref3 := gen.Reference(), gen.Reference(), gen.Reference()
	el1 := &transcript.Transcript{Nonce: 3, RequestRef: ref1}
	el2 := &transcript.Transcript{Nonce: 4, RequestRef: ref2}
	el3 := &transcript.Transcript{Nonce: 5, RequestRef: ref3}
	d.Prepend(el1, el2, el3)

	tr := d.PopByReference(ref3)
	s.Require().NotNil(tr)
	s.Require().Equal(tr.RequestRef.Bytes(), ref3.Bytes())

	s.Nil(d.first.prev)
	s.Equal(d.first.next, d.last)
	s.Nil(d.last.next)
	s.Equal(d.last.prev, d.first)
	s.Equal(d.first.value.Nonce, uint64(3))
	s.Equal(d.last.value.Nonce, uint64(4))

	s.Equal(d.Length(), 2)
}

func (s *TranscriptDequeueSuite) TestPopByReferenceOneElement() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	ref1 := gen.Reference()

	d.Prepend(
		&transcript.Transcript{Nonce: 3, RequestRef: ref1},
	)

	tr := d.PopByReference(ref1)
	s.Require().NotNil(tr)
	s.Require().Equal(tr.RequestRef.Bytes(), ref1.Bytes())

	s.Nil(d.first)
	s.Nil(d.last)
	s.Equal(d.Length(), 0)
}

func (s *TranscriptDequeueSuite) TestTake() {
	d := NewTranscriptDequeue()
	s.Require().NotNil(d)

	for i := 0; i < 15; i++ {
		d.Push(&transcript.Transcript{Nonce: uint64(i)})
	}

	trs := d.Take(0)
	s.Require().NotNil(d)
	s.Len(trs, 0)

	trs = d.Take(10)
	s.Require().NotNil(d)
	s.Len(trs, 10)

	trs = d.Take(10)
	s.Require().NotNil(d)
	s.Len(trs, 5)

	trs = d.Take(10)
	s.Require().NotNil(d)
	s.Len(trs, 0)
}
