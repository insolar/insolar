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

func (p *PulseSlot) IsAcceptedFutureOrPresent(pn pulse.Number) (isFuture, isAccepted bool) {
	presentPN, futurePN := p.pulseManager.GetPresentPulse()
	if p.State() == Future {
		return true, pn >= futurePN
	}
	return false, pn == presentPN // TODO consider a set of pulse numbers when there is a skipped pulse
}

func (p *PulseSlot) HasPulseData(pn pulse.Number) bool {
	pd := p.pulseData.PulseData()
	if pn == pd.PulseNumber && pd.IsValidPulsarData() {
		return true
	}
	return p.pulseManager.HasPulseData(pn)
}
