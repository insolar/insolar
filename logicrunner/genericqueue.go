package logicrunner

type GenericQueue interface {
	Empty() bool
	Enqueue(payload interface{})
	Dequeue() interface{} // check Empty() before calling Dequeue()
}

type GenericQueueItem struct {
	next    *GenericQueueItem
	prev    *GenericQueueItem
	payload interface{}
}

func NewGenericQueue() GenericQueue {
	q := &GenericQueueItem{}
	q.init()
	return q
}

func (q *GenericQueueItem) init() {
	q.next = q
	q.prev = q
	q.payload = nil
}

func (q *GenericQueueItem) Empty() bool {
	return q.next == q
}

func (q *GenericQueueItem) Enqueue(payload interface{}) {
	item := &GenericQueueItem{
		payload: payload,
	}

	item.next = q.next
	item.prev = q
	q.next.prev = item
	q.next = item
}

// Dequeue() dequeues an item.
// Please note that Empty() should be checked before this call
func (q *GenericQueueItem) Dequeue() interface{} {
	msg := q.prev.payload
	old := q.prev

	old.prev.next = q
	q.prev = old.prev
	old.next = nil
	old.prev = nil
	return msg
}
