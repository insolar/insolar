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

package api

import (
	"context"

	"github.com/insolar/insolar/api/seedmanager"
	"github.com/insolar/insolar/core"
	"github.com/pkg/errors"
)

const (
	// REFERENCE is field for reference
	REFERENCE = "reference"
	// SEED is field to reference
	SEED = "seed"
)

// RequestHandler encapsulate processing of request
type RequestHandler struct {
	params              *Params
	contractRequester   core.ContractRequester
	rootDomainReference core.RecordRef
	seedManager         *seedmanager.SeedManager
	seedGenerator       seedmanager.SeedGenerator
	netCoordinator      core.NetworkCoordinator
}

// NewRequestHandler creates new query handler
func NewRequestHandler(params *Params, contractRequester core.ContractRequester, nc core.NetworkCoordinator, rootDomainReference core.RecordRef, smanager *seedmanager.SeedManager) *RequestHandler {
	return &RequestHandler{
		params:              params,
		contractRequester:   contractRequester,
		rootDomainReference: rootDomainReference,
		seedManager:         smanager,
		netCoordinator:      nc,
	}
}

// ProcessGetSeed processes get seed request
func (rh *RequestHandler) ProcessGetSeed(ctx context.Context) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	seed, err := rh.seedGenerator.Next()
	if err != nil {
		return nil, errors.Wrap(err, "[ ProcessGetSeed ]")
	}
	rh.seedManager.Add(*seed)

	result[SEED] = seed[:]

	return result, nil
}
