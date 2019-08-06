/*
 *    Copyright 2019 Insolar Technologies
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package exporter

import (
	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/pulse"
	"github.com/insolar/insolar/ledger/heavy/executor"
	"github.com/pkg/errors"
)

type PulseServer struct {
	pulses    pulse.Calculator
	jetKeeper executor.JetKeeper
}

func NewPulseServer(pulses pulse.Calculator, jetKeeper executor.JetKeeper) *PulseServer {
	return &PulseServer{
		pulses:    pulses,
		jetKeeper: jetKeeper,
	}
}

func (p *PulseServer) Export(getPulses *GetPulses, stream PulseExporter_ExportServer) error {
	if getPulses.Count == 0 {
		return errors.New("count can't be 0")
	}

	read := uint32(0)
	if getPulses.PulseNumber == 0 {
		getPulses.PulseNumber = insolar.FirstPulseNumber
		err := stream.Send(&Pulse{
			PulseNumber: insolar.FirstPulseNumber,
		})
		if err != nil {
			return err
		}
		read++
	}
	currentPN := getPulses.PulseNumber
	for read < getPulses.Count {
		topPulse := p.jetKeeper.TopSyncPulse()
		if currentPN >= topPulse {
			return nil
		}

		pulse, err := p.pulses.Forwards(stream.Context(), currentPN, 1)
		if err != nil {
			return err
		}
		err = stream.Send(&Pulse{
			PulseNumber: pulse.PulseNumber,
		})
		if err != nil {
			return err
		}

		read++
		currentPN = pulse.PulseNumber
	}

	return nil
}
