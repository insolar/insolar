package system

import (
	"github.com/insolar/go-actors/actor"
	"github.com/insolar/go-actors/actor/errors"
	"github.com/insolar/go-actors/actor/mailbox"
	"sync"
)

// system implements actor.System
type system struct {
	wg        sync.WaitGroup
	lock      sync.Mutex
	lastPid   actor.Pid
	mailboxes map[actor.Pid]mailbox.Mailbox
}

// New creates a new System
func New() actor.System {
	return &system{
		mailboxes: make(map[actor.Pid]mailbox.Mailbox),
	}
}

func actorLoop(actor actor.Actor, mailbox mailbox.Mailbox) {
	var err error
	for {
		prevState := actor
		message := mailbox.Dequeue()
		actor, err = prevState.Receive(message)

		if err == errors.Stash {
			mailbox.Stash(message)
		} else if actor != prevState {
			mailbox.Unstash()
		}

		if err != nil && err != errors.Stash {
			break
		}
	}
}

func (s *system) Spawn(constructor actor.Constructor) actor.Pid {
	mbox := mailbox.New()

	s.lock.Lock()
	s.lastPid++ // check for wraparound?
	pid := s.lastPid
	s.mailboxes[pid] = mbox
	s.lock.Unlock()

	s.wg.Add(1)

	// We have to create an actor synchronously. Otherwise it's
	// possible that someone will call Send before the Mailbox
	// will be fully initialized, i.e. including SetLimit.
	act, limit := constructor(s, pid)
	mbox.SetLimit(limit)

	go func() {
		actorLoop(act, mbox)

		s.lock.Lock()
		delete(s.mailboxes, pid)
		s.lock.Unlock()

		s.wg.Done()
	}()

	return pid
}

func (s *system) Send(pid actor.Pid, message actor.Message) error {
	s.lock.Lock()
	mbox, ok := s.mailboxes[pid]
	s.lock.Unlock()

	if !ok {
		return errors.InvalidPid
	}

	return mbox.Enqueue(message)
}

func (s *system) SendPriority(pid actor.Pid, message actor.Message) error {
	s.lock.Lock()
	mbox, ok := s.mailboxes[pid]
	s.lock.Unlock()

	if !ok {
		return errors.InvalidPid
	}

	return mbox.EnqueueFront(message)
}

func (s *system) AwaitTermination() {
	s.wg.Wait()
}
