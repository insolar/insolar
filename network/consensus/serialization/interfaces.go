package serialization

import (
	"encoding/binary"
	"io"

	"github.com/insolar/insolar/network/consensus/common"
)

type SerializerTo interface {
	SerializeTo(writer io.Writer, signer common.DataSigner) (int64, error)
}

const (
	fieldBufSize  = 2048
	packetBufSize = 2048
)

var defaultByteOrder = binary.BigEndian
