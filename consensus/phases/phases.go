/*
 *    Copyright 2018 Insolar
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

package phases

import (
	"github.com/pkg/errors"
)

// FirstPhase is a first phase.
type FirstPhase struct {
	next *SecondPhase
}

func (fp *FirstPhase) HandlePulse(localClaims []ReferendumClaim, data *PulseData) error {
	result, claims, err := fp.getPulseProof(data)
	if err != nil {
		return errors.Wrap(err, "Failed to get a pulse proof")
	}
	return fp.next.Calculate(result, claims)
}

func (fp *FirstPhase) getPulseProof(pulse *PulseData) ([]NodePulseProof, []ReferendumClaim, error) {
	return nil, nil, nil
}

// SecondPhase is a second phase.
type SecondPhase struct {
}

func (sp *SecondPhase) Calculate(proof []NodePulseProof, claims []ReferendumClaim) error {
	return nil
}
