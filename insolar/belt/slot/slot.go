package slot

import (
	"github.com/insolar/insolar/insolar/belt"
)

type Slot struct {
}

func (s *Slot) Add(belt.FlowController) (belt.ID, error) {
	panic("implement me")
}

func (s *Slot) Remove(belt.ID) error {
	panic("implement me")
}

func NewSlot() *Slot {
	return &Slot{}
}
