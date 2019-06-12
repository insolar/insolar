# how to start launchnet, bench and monitor

## run

    ./scripts/insolard/launchnet.sh -g

## monitor

    ./scripts/monitor.sh

## bench

    ./scripts/bench.sh -c=2 -r=40

## if you want to use jaeger in launchnet, add ENV param

	INSOLAR_TRACER_JAEGER_AGENTENDPOINT=<jaeger-addr>
