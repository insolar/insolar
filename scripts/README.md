# how to start launchnet, bench and monitor

## run

    ./insolar-scripts/insolard/launchnet.sh -g

## monitor

    ./insolar-scripts/monitor.sh

## bench

    ./scripts/bench.sh -c=2 -r=40

## if you want to use jaeger in launchnet, add ENV param

	INSOLAR_TRACER_JAEGER_AGENTENDPOINT=<jaeger-addr>
