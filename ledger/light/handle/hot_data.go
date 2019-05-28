//
// Copyright 2019 Insolar Technologies GmbH
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
//

package handle

import (
	"context"
	"fmt"

	watermillMsg "github.com/ThreeDotsLabs/watermill/message"
	"github.com/insolar/insolar/insolar/flow"
	"github.com/insolar/insolar/insolar/message"
	"github.com/insolar/insolar/ledger/light/proc"
)

type HotData struct {
	dep       *proc.Dependencies
	wmmessage *watermillMsg.Message
	message   *message.HotData
}

func NewHotData(dep *proc.Dependencies, wmmessage *watermillMsg.Message, msg *message.HotData) *HotData {
	return &HotData{
		dep:       dep,
		wmmessage: wmmessage,
		message:   msg,
	}
}

func (s *HotData) Present(ctx context.Context, f flow.Flow) error {
	fmt.Println("start TypeHotRecords in Present")
	proc := proc.NewHotData(s.message, s.wmmessage)
	s.dep.HotData(proc)
	return f.Procedure(ctx, proc, false)
}
