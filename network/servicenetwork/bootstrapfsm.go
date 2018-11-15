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

package servicenetwork

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
)

type state int

const (
	unknown state = iota
	collectingDiscoveryConsensus
	consensusOk
)

type bootstrapFSM struct {
	current state
}

func newBootstrapFSM() *bootstrapFSM {
	return &bootstrapFSM{}
}

func (bFSM *bootstrapFSM) bypassDiscoveryNodes() {

}

func (bFSM *bootstrapFSM) hasBootstrapConsensus(ctx context.Context) (bool, error) {
	switch bFSM.current {
	case unknown:
	case collectingDiscoveryConsensus:

	case consensusOk:
		return true, nil
	default:
		return false, errors.New("[ hasBootstrapConsensus ] Bad state: " + strconv.Itoa(int(bFSM.current)))
	}

	return false, nil
}
