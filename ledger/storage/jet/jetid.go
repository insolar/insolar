package jet

import (
	"bytes"

	"github.com/insolar/insolar/core"
	"github.com/ugorji/go/codec"
)

// JetIDSet is an alias for map[JetID]struct{}
type JetIDSet map[JetID]struct{}

// JetID contains meta-params for jet
type JetID struct {
	ID    core.RecordID
	Depth int8
}

// Bytes serializes pulse.
func (j JetIDSet) Bytes() []byte {
	var buf bytes.Buffer
	enc := codec.NewEncoder(&buf, &codec.CborHandle{})
	enc.MustEncode(j)
	return buf.Bytes()
}
