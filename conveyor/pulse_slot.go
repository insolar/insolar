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
	_, ps := p.pulseData.PulseRange()
	return ps
}

func (p *PulseSlot) PulseData() pulse.Data {
	pd, _ := p.pulseData.PulseData()
	if pd.IsEmpty() {
		// possible incorrect injection for SM in the Antique slot
		panic("illegal state - not initialized")
	}
	return pd
}

func (p *PulseSlot) IsAcceptedFutureOrPresent(pn pulse.Number) (isFuture, isAccepted bool) {
	presentPN, futurePN := p.pulseManager.GetPresentPulse()
	if p.State() == Future {
		return true, pn >= futurePN
	}
	return false, pn == presentPN // TODO consider a set of pulse numbers when there is a skipped pulse
}

func (p *PulseSlot) HasPulseData(pn pulse.Number) bool {
	if pr, _ := p.pulseData.PulseRange(); pr != nil {
		return true
	}
	return p.pulseManager.HasPulseData(pn)
}
