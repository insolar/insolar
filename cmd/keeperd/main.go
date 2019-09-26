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

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/viper"

	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
)

var configPath string

func main() {
	var rootCmd = &cobra.Command{
		Use: "keeperd --config=<path to config>",
		Run: rootCommand,
	}
	rootCmd.Flags().StringVarP(&configPath, "config", "c", "", "path to config file")
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Wrong input params:", err.Error())
	}
}

func rootCommand(cmd *cobra.Command, args []string) {
	jww.SetStdoutThreshold(jww.LevelInfo)
	var err error

	vp := viper.New()
	cfg := NewConfig()
	if len(configPath) != 0 {
		vp.SetConfigFile(configPath)
	}
	err = vp.ReadInConfig()
	if err != nil {
		log.Warn("failed to load configuration from file: ", err.Error())
	}
	err = vp.Unmarshal(&cfg)
	if err != nil {
		log.Warn("failed to load configuration from file: ", err.Error())
	}

	ctx := context.Background()
	ctx, logger := inslogger.InitNodeLogger(ctx, cfg.Log, "main", "", "keeperd")

	keeper := NewKeeper(cfg.Keeper)
	keeper.Run(ctx)

	vp.WatchConfig()
	vp.OnConfigChange(func(e fsnotify.Event) {
		logger.Info("Reloading config file")
		cfg := NewConfig()
		err := vp.Unmarshal(&cfg)
		if err != nil {
			logger.Errorf("Failed to reload config: %s", err.Error())
			return
		}
		keeper.SetConfig(cfg.Keeper)
	})

	var gracefulStop = make(chan os.Signal, 1)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)

	<-gracefulStop
}

type Keeper struct {
	config      KeeperConfig
	isAvailable uint32
}

func NewKeeper(cfg KeeperConfig) *Keeper {
	return &Keeper{
		config:      cfg,
		isAvailable: 0,
	}
}

func (k *Keeper) SetConfig(cfg KeeperConfig) {
	k.config = cfg
}

func (k *Keeper) Run(ctx context.Context) {
	go k.startChecker(ctx)
	go k.startServer(ctx)
}

func (k *Keeper) startChecker(ctx context.Context) {
	ticker := time.NewTicker(k.config.PollPeriod)
	for range ticker.C {
		k.checkMetrics(ctx)
	}
}

func (k *Keeper) checkMetrics(ctx context.Context) {
	queries := k.config.Queries
	var isOK uint32 = 1
	if !k.config.FakeTrue {
		for _, q := range queries {
			if !k.checkMetric(ctx, q) {
				isOK = 0
			}
		}
	}
	atomic.StoreUint32(&k.isAvailable, isOK)
}

func (k *Keeper) checkMetric(ctx context.Context, query string) bool {
	logger := inslogger.FromContext(ctx)
	logger.Debugf("Checking metric: %s", query)

	metricURL := k.config.QueryURL + url.QueryEscape(query)
	resp, err := http.Get(metricURL)
	if err != nil {
		logger.Errorf("Error while getting <<%s>> query: %s", query, err.Error())
		return false
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			logger.Warnf("Failed to close response body: %s", err.Error())
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
		if metricLag > k.config.MaxMetricLag {
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
	logger.Infof("Starting http server on %s", k.config.ListenAddress)
	logger.Fatal(http.ListenAndServe(k.config.ListenAddress, nil))
}
