package mailbox

import "github.com/insolar/go-actors/actor"

type queue_item struct {
	next *queue_item
	prev *queue_item
	payload actor.Message
}

type queue struct {
	len int
	head queue_item
}

func (q *queue) init() {
	q.head.next = &q.head
	q.head.prev = &q.head
	q.head.payload = nil
	q.len = 0
}

func (q *queue) empty() bool {
	return q.head.next == &q.head
}

func (q *queue) length() int {
	return q.len
}

func (q *queue) enqueue(payload actor.Message) {
	new :=  &queue_item{
		payload: payload,
	}

	new.next = q.head.next
	new.prev = &q.head
	q.head.next.prev = new
	q.head.next = new
	q.len++
}

// empty() should be checked before this call
func (q *queue) dequeue() actor.Message {
	msg := q.head.prev.payload
	old := q.head.prev

	old.prev.next = &q.head
	q.head.prev = old.prev
	old.next = nil
	old.prev = nil
	q.len--
	return msg
}

func (q1 *queue) moveFromQueue(q2 *queue) {
	if q2.empty() {
		return
	}

	q2.head.next.prev = q1.head.prev
	q1.head.prev.next = q2.head.next
	q2.head.prev.next = &q1.head
	q1.head.prev = q2.head.prev

	q1.len += q2.len
	q2.init()
}
