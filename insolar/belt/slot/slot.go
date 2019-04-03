package slot

import (
	"sync"

	"github.com/insolar/insolar/insolar/belt"
)

type Slot struct {
	lock        sync.RWMutex
	idCounter   belt.ID
	controllers map[belt.ID]belt.FlowController
}

func NewSlot() *Slot {
	return &Slot{
		controllers: map[belt.ID]belt.FlowController{},
	}
}

func (s *Slot) Add(c belt.FlowController) (belt.ID, error) {
	s.lock.Lock()
	id := s.idCounter
	s.idCounter++
	s.controllers[id] = c
	s.lock.Unlock()
	return id, nil
}

func (s *Slot) Remove(id belt.ID) error {
	s.lock.Lock()
	delete(s.controllers, id)
	s.lock.Unlock()
	return nil
}
