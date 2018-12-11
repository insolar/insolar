package jet

import (
	"bytes"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

// IDSet is an alias for map[ID]struct{}
type IDSet map[ID]struct{}

// ID contains meta-params for jet
type ID struct {
	ID    core.RecordID
	Depth int8
}

// Bytes serializes pulse.
func (j IDSet) Bytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	enc.MustEncode(j)
	return buf.Bytes()
}
