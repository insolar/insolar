package bus

import (
	"github.com/insolar/insolar/insolar"
)

type Reply struct {
	Reply insolar.Reply
	Err   error
}

type Message struct {
	Parcel  insolar.Parcel
	ReplyTo chan Reply
}
