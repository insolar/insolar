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
	"sync"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/logicrunner/transcript"
)

type Element struct {
	prev  *Element
	next  *Element
	value *transcript.Transcript
}

// TODO: probably it's better to rewrite it using linked list
type TranscriptDequeue struct {
	lock   sync.Locker
	first  *Element
	last   *Element
	length int
}

func NewTranscriptDequeue() *TranscriptDequeue {
	return &TranscriptDequeue{
		lock:   &sync.Mutex{},
		first:  nil,
		last:   nil,
		length: 0,
	}
}

func (d *TranscriptDequeue) pushOne(el *transcript.Transcript) {
	newElement := &Element{value: el}
	lastElement := d.last

	if lastElement != nil {
		newElement.prev = lastElement
		lastElement.next = newElement
		d.last = newElement
	} else {
		d.first, d.last = newElement, newElement
	}

	d.length++
}

func (d *TranscriptDequeue) Push(els ...*transcript.Transcript) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for _, el := range els {
		d.pushOne(el)
	}
}

func (d *TranscriptDequeue) prependOne(el *transcript.Transcript) {
	newElement := &Element{value: el}
	firstElement := d.first

	if firstElement != nil {
		newElement.next = firstElement
		firstElement.prev = newElement
		d.first = newElement
	} else {
		d.first, d.last = newElement, newElement
	}

	d.length++
}

func (d *TranscriptDequeue) Prepend(els ...*transcript.Transcript) {
	d.lock.Lock()
	defer d.lock.Unlock()

	for i := len(els) - 1; i >= 0; i-- {
		d.prependOne(els[i])
	}
}

func (d *TranscriptDequeue) Pop() *transcript.Transcript {
	elements := d.Take(1)
	if len(elements) == 0 {
		return nil
	}
	return elements[0]
}

func (d *TranscriptDequeue) Has(ref insolar.Reference) bool {
	d.lock.Lock()
	defer d.lock.Unlock()

	for elem := d.first; elem != nil; elem = elem.next {
		if elem.value.RequestRef.Compare(ref) == 0 {
			return true
		}
	}
	return false
}

func (d *TranscriptDequeue) PopByReference(ref insolar.Reference) *transcript.Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	for elem := d.first; elem != nil; elem = elem.next {
		if elem.value.RequestRef.Compare(ref) == 0 {
			if elem.prev != nil {
				elem.prev.next = elem.next
			} else {
				d.first = elem.next
			}
			if elem.next != nil {
				elem.next.prev = elem.prev
			} else {
				d.last = elem.prev
			}

			d.length--

			return elem.value
		}
	}

	return nil
}

func (d *TranscriptDequeue) HasFromLedger() *transcript.Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	for elem := d.first; elem != nil; elem = elem.next {
		if elem.value.FromLedger {
			return elem.value
		}
	}
	return nil
}

func (d *TranscriptDequeue) commonPeek(count int) (*Element, []*transcript.Transcript) {
	if d.length < count {
		count = d.length
	}

	rv := make([]*transcript.Transcript, count)

	var lastElement *Element
	for i := 0; i < count; i++ {
		if lastElement == nil {
			lastElement = d.first
		} else {
			lastElement = lastElement.next
		}
		rv[i] = lastElement.value
	}

	return lastElement, rv
}

func (d *TranscriptDequeue) take(count int) []*transcript.Transcript {
	lastElement, rv := d.commonPeek(count)
	if lastElement != nil {
		if lastElement.next == nil {
			d.first, d.last = nil, nil
		} else {
			lastElement.next.prev, d.first = nil, lastElement.next
			lastElement.next = nil
		}

		d.length -= len(rv)
	}

	return rv
}

func (d *TranscriptDequeue) Peek(count int) []*transcript.Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	_, rv := d.commonPeek(count)
	return rv
}

func (d *TranscriptDequeue) Take(count int) []*transcript.Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	return d.take(count)
}

func (d *TranscriptDequeue) Rotate() []*transcript.Transcript {
	d.lock.Lock()
	defer d.lock.Unlock()

	return d.take(d.length)
}

func (d *TranscriptDequeue) Length() int {
	d.lock.Lock()
	defer d.lock.Unlock()

	return d.length
}
