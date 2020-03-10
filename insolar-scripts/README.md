# how to start launchnet and monitor

To run launchnet, you must provide several files:

`/scripts/insolard/bootstrap_template.yaml` - template for bootstrap config

`/scripts/insolard/generate_initial_data.sh` - script, which creates init data for application (members keys, migration addresses, etc.)

## run

    ./insolar-scripts/insolard/launchnet.sh -g

## monitor

    ./insolar-scripts/monitor.sh

## if you want to use jaeger in launchnet, add ENV param

	INSOLAR_TRACER_JAEGER_AGENTENDPOINT=<jaeger-addr>
