// +build never_run

package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/instrumentation/inslogger"
	"github.com/insolar/insolar/log"
)

func TestConsensusMain(t *testing.T) {

	startedAt := time.Now()

	ctx := context.Background()
	logger := inslogger.FromContext(ctx) //.WithCaller(false)
	//logger, _ = logger.WithLevelNumber(insolar.DebugLevel)
	//logger, _ = logger.WithFormat(insolar.TextFormat)
	logger = logger.Level(insolar.DebugLevel)
	log.SetGlobalLogger(logger)
	ctx = inslogger.SetLogger(ctx, log.GlobalLogger())
	_ = log.SetGlobalLevelFilter(insolar.DebugLevel)

	netStrategy := NewDelayNetStrategy(DelayStrategyConf{
		MinDelay:         10 * time.Millisecond,
		MaxDelay:         30 * time.Millisecond,
		Variance:         0.2,
		SpikeProbability: 0.1,
	})
	strategyFactory := &EmuRoundStrategyFactory{}

	nodes := NewEmuNodeIntros(generateNameList(0, 1, 3, 5)...)
	netBuilder := newEmuNetworkBuilder(ctx, netStrategy, strategyFactory)

	for i := range nodes {
		netBuilder.connectEmuNode(nodes, i)
	}

	netBuilder.StartNetwork(ctx)

	netBuilder.StartPulsar(10, 2, "pulsar0", nodes)

	// time.AfterFunc(time.Second, func() {
	//	netBuilder.network.DropHost("V0007")
	// })

	for {
		fmt.Println("===", time.Since(startedAt), "=================================================")
		time.Sleep(time.Second)
		if time.Since(startedAt) > time.Minute*30 {
			return
		}
	}
}
