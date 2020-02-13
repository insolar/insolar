// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/insolar/blob/master/LICENSE.md.

package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
	"sync/atomic"
	"time"

	"go.opencensus.io/stats"

	"github.com/insolar/insolar/instrumentation/inslogger"
)

type Keeper struct {
	config      KeeperConfig
	configLock  sync.RWMutex
	isAvailable uint32
}

func NewKeeper(cfg KeeperConfig) *Keeper {
	return &Keeper{
		config:      cfg,
		isAvailable: 0,
	}
}

func (k *Keeper) SetConfig(cfg KeeperConfig) {
	k.configLock.Lock()
	defer k.configLock.Unlock()
	k.config = cfg
}

func (k *Keeper) Config() KeeperConfig {
	k.configLock.RLock()
	defer k.configLock.RUnlock()
	return k.config
}

func (k *Keeper) Run(ctx context.Context) {
	go k.startChecker(ctx)
	go k.startServer(ctx)
}

func (k *Keeper) startChecker(ctx context.Context) {
	ticker := time.NewTicker(k.Config().PollPeriod)
	for range ticker.C {
		k.checkMetrics(ctx)
	}
}

func (k *Keeper) checkMetrics(ctx context.Context) {
	queries := k.Config().Queries
	var isOK uint32 = 1
	if !k.Config().FakeTrue {
		for _, q := range queries {
			if !k.checkMetric(ctx, q) {
				isOK = 0
			}
		}
	}
	atomic.StoreUint32(&k.isAvailable, isOK)
	stats.Record(ctx, IsAvailable.M(int64(isOK)))
}

func (k *Keeper) checkMetric(ctx context.Context, query string) bool {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("Checking metric: %s", query)

	metricURL := k.Config().QueryURL + url.QueryEscape(query)
	resp, err := http.Get(metricURL)
	if err != nil {
		logger.Errorf("Error while getting <<%s>> query: %s", query, err.Error())
		return false
	}
	defer func() {
		if resp != nil && resp.Body != nil {
			err = resp.Body.Close()
			if err != nil {
				logger.Warnf("Failed to close response body: %s", err.Error())
			}
		}
	}()

	promRsp := PromRsp{}
	err = json.NewDecoder(resp.Body).Decode(&promRsp)
	if err != nil {
		logger.Errorf("Metric <<%s>>. Failed to parse Prometheus response: %s", query, err.Error())
		return false
	}
	if promRsp.Status != "success" {
		logger.Errorf("Metric <<%s>>. Bad response from Prometheus: %s: %s", query, promRsp.ErrorType, promRsp.Error)
		return false
	}
	for _, res := range promRsp.Data.Result {
		ts, ok := res.Value[0].(float64)
		if !ok {
			logger.Errorf("Metric <<%s>> failed on instance %s. Cannot parse timestamp %v", query, res.Metric.Instance, res.Value[0])
			return false
		}
		metricLag := time.Since(time.Unix(int64(ts), 0))
		logger.Debugf("Metric <<%s>> on instance %s lag: %s", query, res.Metric.Instance, metricLag)
		if metricLag > k.Config().MaxMetricLag {
			logger.Errorf("Metric <<%s>> failed on instance %s. Last data point was too long ago: %s", query, res.Metric.Instance, res.Value[0], metricLag)
			return false
		}

		if res.Value[1] != "1" {
			logger.Infof("Metric <<%s>> failed on instance %s. Value is %s", query, res.Metric.Instance, res.Value[1])
			return false
		}
	}

	return true
}

func (k *Keeper) startServer(ctx context.Context) {
	logger := inslogger.FromContext(ctx)
	http.HandleFunc("/check", func(writer http.ResponseWriter, request *http.Request) {
		response := KeeperRsp{
			Available: atomic.LoadUint32(&k.isAvailable) > 0,
		}
		writer.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(writer).Encode(response)
		if err != nil {
			logger.Errorf("Failed to encode response: %s", err.Error())
		}
	})
	logger.Infof("Starting http server on %s", k.Config().ListenAddress)
	logger.Fatal(http.ListenAndServe(k.Config().ListenAddress, nil))
}
