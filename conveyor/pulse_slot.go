package conveyor

import "github.com/insolar/insolar/pulse"

type PulseSlotState uint8

const (
	_ PulseSlotState = iota
	Future
	Present
	Past
	Antique // non-individual past
)

type PulseSlot struct {
	pulseManager *PulseDataManager
	pulseData    pulseDataHolder
}

func (p *PulseSlot) State() PulseSlotState {
	return p.pulseData.State()
}

func (p *PulseSlot) PulseData() pulse.Data {
	return p.pulseData.PulseData()
}

func (p *PulseSlot) IsFuture(number pulse.Number) bool {

}

func (p *PulseSlot) IsAccepted(number pulse.Number) bool {

}

func (p *PulseSlot) HasPulseData(number pulse.Number) bool {

}
