package mimic

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/longbits"
	"github.com/insolar/insolar/pulse"
)

type PulseGenerator struct {
	delta     uint16
	pulseList []pulse.Data
}

const nanosecondsInSecond = int64(time.Second / time.Nanosecond)

func NewPulseData(p insolar.Pulse) pulse.Data {
	data := pulse.NewPulsarData(
		p.PulseNumber,
		uint16(p.NextPulseNumber-p.PulseNumber),
		uint16(p.PulseNumber-p.PrevPulseNumber),
		longbits.NewBits512FromBytes(p.Entropy[:]).FoldToBits256(),
	)
	data.Timestamp = uint32(p.PulseTimestamp / nanosecondsInSecond)
	data.PulseEpoch = p.EpochPulseNumber
	return data
}

func NewPulse(pulseData pulse.Data) insolar.Pulse {
	var prev insolar.PulseNumber
	if !pulseData.IsFirstPulse() {
		prev = pulseData.PrevPulseNumber()
	} else {
		prev = pulseData.PulseNumber
	}

	entropy := insolar.Entropy{}
	bs := pulseData.PulseEntropy.AsBytes()
	copy(entropy[:], bs)
	copy(entropy[pulseData.PulseEntropy.FixedByteSize():], bs)

	return insolar.Pulse{
		PulseNumber:      pulseData.PulseNumber,
		NextPulseNumber:  pulseData.NextPulseNumber(),
		PrevPulseNumber:  prev,
		PulseTimestamp:   int64(pulseData.Timestamp) * nanosecondsInSecond,
		EpochPulseNumber: pulseData.PulseEpoch,
		Entropy:          entropy,
	}
}

func NewPulseGenerator(delta uint16) *PulseGenerator {
	g := &PulseGenerator{delta: delta}
	g.appendPulse(NewPulseData(*insolar.GenesisPulse))
	return g
}

func (g PulseGenerator) GetLastPulse() pulse.Data {
	return g.pulseList[len(g.pulseList)-1]
}

func (g PulseGenerator) GetLastPulseAsPulse() insolar.Pulse {
	return NewPulse(g.pulseList[len(g.pulseList)-1])
}

func (g *PulseGenerator) appendPulse(data pulse.Data) {
	g.pulseList = append(g.pulseList, data)
}

func generateEntropy() longbits.Bits256 {
	var entropy longbits.Bits256

	actualLength, err := rand.Read(entropy[:])
	if err != nil {
		panic(fmt.Sprintf("failed to read %d random bytes: %s", len(entropy), err.Error()))
	} else if actualLength != len(entropy) {
		panic(fmt.Sprintf("unreachable, %d != %d", len(entropy), actualLength))
	}

	return entropy
}

func (g PulseGenerator) generateNextPulse() pulse.Data {
	var (
		nextPulse pulse.Data
		prevPulse = g.GetLastPulse()
	)

	if !prevPulse.IsEmpty() || prevPulse.PulseNumber == pulse.MinTimePulse {
		nextPulse = prevPulse.CreateNextPulsarPulse(g.delta, generateEntropy)
	} else {
		nextPulse = pulse.NewFirstPulsarData(g.delta, generateEntropy())
	}

	return nextPulse
}

func (g *PulseGenerator) Generate() pulse.Data {
	newPulse := g.generateNextPulse()
	g.appendPulse(newPulse)
	return newPulse
}
