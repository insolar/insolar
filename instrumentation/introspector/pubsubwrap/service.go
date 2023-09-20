package pubsubwrap

import (
	"github.com/insolar/insolar/instrumentation/introspector/introproto"
)

// PublisherService implements introproto.PublisherServer.
type PublisherService struct {
	*MessageLockerByType
	*MessageStatByType
}

// programming and compile time check
var _ introproto.PublisherServer = PublisherService{}

// NewPublisherService creates PublisherService.
func NewPublisherService(ml *MessageLockerByType, ms *MessageStatByType) PublisherService {
	return PublisherService{
		MessageLockerByType: ml,
		MessageStatByType:   ms,
	}
}
