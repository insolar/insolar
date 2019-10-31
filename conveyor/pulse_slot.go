package conveyor

import (
	"github.com/insolar/insolar/pulse"
)

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

func (p *PulseSlot) isAcceptedFutureOrPresent(pn pulse.Number) (isFuture, isAccepted bool) {
	presentPN, futurePN := p.pulseManager.GetPresentPulse()
	if ps, ok := p._isAcceptedPresent(presentPN, pn); ps != Future {
		return false, ok
	}
	return true, pn >= futurePN
}

func (p *PulseSlot) _isAcceptedPresent(presentPN, pn pulse.Number) (PulseSlotState, bool) {
	var isProhibited bool

	switch pr, ps := p.pulseData.PulseRange(); {
	case ps != Present:
		return ps, false
	case pn == presentPN:
		return ps, true
	case pn < pr.LeftBoundNumber():
		// pn belongs to Past or Antique for sure
		return Past, false
	case pr.IsSingular():
		return ps, false
	case !pr.EnumNonArticulatedNumbers(func(n pulse.Number, prevDelta, nextDelta uint16) bool {
		switch {
		case n == pn:
		case pn.IsEqOrOut(n, prevDelta, nextDelta):
			return pn < n // stop search, as EnumNumbers from smaller to higher pulses
		default:
			// this number is explicitly prohibited by a known pulse data
			isProhibited = true
			// and stop now
		}
		return true
	}):
		// we've seen neither positive nor negative match
		if pr.IsArticulated() {
			// Present range is articulated, so anything that is not wrong - can be valid as present
			return ps, true
		}
		fallthrough
	case isProhibited: // Present range is articulated, so anything that is not wrong - can be valid
		return ps, false
	default:
		// we found a match in a range of the present slot
		return ps, true
	}
}

func (p *PulseSlot) HasPulseData(pn pulse.Number) bool {
	if pr, _ := p.pulseData.PulseRange(); pr != nil {
		return true
	}
	return p.pulseManager.HasPulseData(pn)
}
