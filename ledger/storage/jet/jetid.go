package jet

import (
	"bytes"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

// NewID generates RecordID for jet.
func NewID(depth uint8, prefix []byte) *core.RecordID {
	var id core.RecordID
	copy(id[:core.PulseNumberSize], core.PulseNumberJet.Bytes())
	id[core.PulseNumberSize] = depth
	copy(id[core.PulseNumberSize+1:], prefix)
	return &id
}

// IDSet is an alias for map[ID]struct{}
type IDSet map[core.RecordID]struct{}

// Has checks if passed id is in IDSet set.
func (j IDSet) Has(id core.RecordID) bool {
	_, ok := j[id]
	return ok
}

// Bytes serializes pulse.
func (j IDSet) Bytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	enc.MustEncode(j)
	return buf.Bytes()
}

// Jet extracts depth and prefix from jet id.
func Jet(id core.RecordID) (uint8, []byte) {
	if id.Pulse() != core.PulseNumberJet {
		panic("provided id in not a jet id")
	}
	return id[core.PulseNumberSize], id[core.PulseNumberSize+1:]
}
