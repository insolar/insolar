// Copyright 2020 Insolar Network Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
