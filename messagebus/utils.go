package messagebus

import (
	"github.com/insolar/insolar/core"
	"github.com/insolar/insolar/core/message"
	"github.com/insolar/insolar/cryptohelpers/hash"
)

func GetMessageHash(msg core.SignedMessage) []byte {
	return hash.SHA3Bytes256(message.SignedToBytes(msg))
}
