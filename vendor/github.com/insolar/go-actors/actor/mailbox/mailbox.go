package mailbox

import (
	"github.com/insolar/go-actors/actor"
	"github.com/insolar/go-actors/actor/errors"
	"sync"
)

// Mailbox is a queue of messages sent to a given actor.
type Mailbox interface {
	// Enqueue places a new message to the Mailbox. If Mailbox is full
	// method returns errors.MailboxFull. Mailboxes are unbounded unless
	// SetLimit was called with a value of `limit` > 0.
	Enqueue(message actor.Message) error

	// Enqueue places a new message in the front of the messages queue. The
	// message will be Dequeued before any messages placed using Enqueue.
	EnqueueFront(message actor.Message) error

	// Dequeue gets a next message from the Mailbox. If Mailbox is empty,
	// Dequeue blocks until someone places a message to the Mailbox. There
	// should be only one goroutine that uses this method.
	Dequeue() actor.Message

	// Stash places a message to the separate queue of delayed messages.
	// This method should be called from the same single goroutine that calls Dequeue.
	Stash(message actor.Message)

	// Unstash places all previously Stash'ed messages in the front of the message queue.
	// The original order of Stash'ed messages is preserved. This method should be called
	// from the same single goroutine that calls Dequeue.
	Unstash()

	// SetLimit sets the maximum capacity of the Mailbox. This limit doesn't affect
	// stashed and priority messages - the amount of these messages is always unlimited.
	SetLimit(limit int)
}

// mailbox implements actor.Mailbox interface.
type mailbox struct {
	lock             sync.Mutex
	dequeueIsWaiting bool
	notifyDequeue    chan struct{}
	limit            int
	regularQueue     queue
	priorityQueue    queue
	stashQueue       queue
}

// New creates a new Mailbox
func New() Mailbox {
	mbox := &mailbox{}
	mbox.regularQueue.init()
	mbox.priorityQueue.init()
	mbox.stashQueue.init()
	mbox.notifyDequeue = make(chan struct{}, 1)
	return mbox
}

func (mb *mailbox) SetLimit(limit int) {
	if limit > 0 {
		mb.limit = limit
	}
}

func (mb *mailbox) Enqueue(message actor.Message) error {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	if mb.limit != 0 && mb.regularQueue.length() == mb.limit {
		return errors.MailboxFull
	}

	mb.regularQueue.enqueue(message)
	if mb.dequeueIsWaiting {
		mb.dequeueIsWaiting = false
		mb.notifyDequeue <- struct{}{}
	}
	return nil
}

func (mb *mailbox) EnqueueFront(message actor.Message) error {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.priorityQueue.enqueue(message)
	if mb.dequeueIsWaiting {
		mb.dequeueIsWaiting = false
		mb.notifyDequeue <- struct{}{}
	}
	return nil
}

func (mb *mailbox) Stash(message actor.Message) {
	mb.lock.Lock()
	defer mb.lock.Unlock()
	mb.stashQueue.enqueue(message)
}

func (mb *mailbox) Unstash() {
	mb.lock.Lock()
	defer mb.lock.Unlock()

	mb.regularQueue.moveFromQueue(&mb.stashQueue)

	// There is no need to check mb.dequeueIsWaiting here because
	// Unstash and Dequeue are called from the same goroutine.
	// Since we are here we know there is no another goroutine
	// that is waiting in Dequeue.
}

func (mb *mailbox) Dequeue() actor.Message {
	var msg actor.Message
	for {
		mb.lock.Lock()

		if !mb.priorityQueue.empty() {
			msg = mb.priorityQueue.dequeue()
			break
		}

		if !mb.regularQueue.empty() {
			msg = mb.regularQueue.dequeue()
			break
		}

		// all queues are empty
		mb.dequeueIsWaiting = true
		mb.lock.Unlock()
		<-mb.notifyDequeue
	}
	mb.lock.Unlock()
	return msg
}
