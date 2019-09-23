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

package logadapter

import (
	"context"
	"github.com/insolar/insolar/insolar"
	"time"
)

type MetricsHelper struct {
}

type DurationReportFunc func(d time.Duration)

func (p *MetricsHelper) OnNewEvent(ctx context.Context, level insolar.LogLevel) {
	if p == nil {
		return
	}
	//stats.Record(m.metrics, statLogCalls.M(1))
}

func (p *MetricsHelper) OnFilteredEvent(ctx context.Context, level insolar.LogLevel) {
	if p == nil {
		return
	}
	//stats.Record(m.metrics, statLogWrites.M(1))
}

func (p *MetricsHelper) OnWriteDuration(d time.Duration) {
	if p == nil {
		return
	}

}

func (p *MetricsHelper) GetOnWriteDurationReport() DurationReportFunc {
	if p == nil {
		return nil
	}
	return p.OnWriteDuration
}
