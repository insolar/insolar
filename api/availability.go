/*
 * Copyright 2019 Insolar Technologies GmbH
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/insolar/insolar/configuration"
	"github.com/insolar/insolar/instrumentation/inslogger"
)

// AvailabilityChecker component checks if insolar network can't process any new requests
type AvailabilityChecker interface {
	IsAvailable(context.Context) bool
}

const keeperRequestTimeout = 15 * time.Second
const checkAvailabilityPeriod = 5 * time.Second

type keeperResponse struct {
	Available bool `json:"available"`
}

type NetworkChecker struct {
	client    *http.Client
	enabled   bool
	keeperUrl string

	lock        *sync.RWMutex
	isAvailable bool
}

func NewNetworkChecker(cfg configuration.AvailabilityChecker) *NetworkChecker {
	return &NetworkChecker{
		client: &http.Client{
			Transport: &http.Transport{},
			Timeout:   keeperRequestTimeout,
		},
		enabled:     cfg.Enabled,
		keeperUrl:   cfg.KeeperUrl,
		lock:        &sync.RWMutex{},
		isAvailable: false,
	}
}

func (nc *NetworkChecker) Start(ctx context.Context) error {
	if !nc.enabled {
		nc.lock.Lock()
		defer nc.lock.Unlock()

		nc.isAvailable = true
		return nil
	}

	go func(ctx context.Context) {
		for range time.NewTicker(checkAvailabilityPeriod).C {
			nc.updateAvailability(ctx)
		}
	}(ctx)
	return nil
}

func (nc *NetworkChecker) updateAvailability(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	resp, err := nc.client.Get(nc.keeperUrl)
	defer func() {
		err := resp.Body.Close()
		logger.Error("Can't close body: ", err)
	}()

	nc.lock.Lock()
	defer nc.lock.Unlock()

	if err != nil {
		nc.isAvailable = false
		logger.Error("Can't get keeper status: ", err)
	}

	if resp != nil && resp.StatusCode != http.StatusOK {
		nc.isAvailable = false
		logger.Error("Can't get keeper status: no response or bad StatusCode")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		nc.isAvailable = false
		logger.Error("Can't get keeper status: Can't read body: ", err)
	}

	respObj := &keeperResponse{}
	err = json.Unmarshal(body, respObj)
	if err != nil {
		nc.isAvailable = false
		logger.Error("Can't get keeper status: Can't unmarshal body: ", err)
	}

	if !respObj.Available {
		logger.Error("Network is not available for request processing")
	}

	nc.isAvailable = respObj.Available
}

func (nc *NetworkChecker) IsAvailable(ctx context.Context) bool {
	nc.lock.RLock()
	defer nc.lock.RUnlock()
	return nc.isAvailable
}
