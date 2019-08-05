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
	"github.com/pkg/errors"
)

type PulseServer struct {
	pulses pulse.Calculator
}

func (p *PulseServer) Export(getPulses *GetPulses, stream PulseExporter_ExportServer) error {
	if getPulses.Count == 0 {
		return errors.New("count can't be 0")
	}

	if getPulses.PulseNumber == 0 {
		getPulses.PulseNumber = insolar.FirstPulseNumber
		err := stream.Send(&Pulse{
			PulseNumber: insolar.FirstPulseNumber,
		})
		if err != nil {
			return err
		}
	}

	currentPN := getPulses.PulseNumber
	read := uint32(0)
	for pn, err := p.pulses.Forwards(stream.Context(), currentPN, 1); err == nil && read <= getPulses.Count; {
		err := stream.Send(&Pulse{
			PulseNumber: pn.PulseNumber,
		})
		if err != nil {
			return err
		}
		read++
	}

	return nil
}
