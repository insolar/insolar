package jet

import (
	"bytes"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

// IDSet is an alias for map[ID]struct{}
type IDSet map[core.RecordID]struct{}

// Bytes serializes pulse.
func (j IDSet) Bytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	enc.MustEncode(j)
	return buf.Bytes()
}
