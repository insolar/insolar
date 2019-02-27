#!/bin/bash
set -e

export IP=`awk 'END{print $1}' /etc/hosts`

INSOLARD_BIN=/opt/bin/insolard
INSGORUND_BIN=/opt/bin/insgorund

GENERATE_BIN=/opt/bin/genconfig
FOREGO_BIN=/opt/bin/forego

export INSOLARD_TRANSPORT_LISTEN_PORT=${INSOLARD_TRANSPORT_LISTEN_PORT:-13831}
export INSOLARD_METRICS_LISTEN_PORT=${INSOLARD_METRICS_LISTEN_PORT:-8001}
export INSOLARD_RPC_LISTEN_PORT=${INSOLARD_RPC_LISTEN_PORT:-33001}
export INSOLARD_API_LISTEN_PORT=${INSOLARD_API_LISTEN_PORT:-19101}
export INSOLARD_LOG_LEVEL=${INSOLARD_LOG_LEVEL:-info}
export INSGORUND_LISTEN_PORT=${INSGORUND_LISTEN_PORT:-33002}

check_cert() {
    if   [ ! -f "/etc/insolar/cert.json" ]; then
        echo "Failed to find certificate, please provide cert.json file"
        return 1
    elif [ ! -f "/etc/insolar/keys.json" ]; then
        echo "Failed to find keys, please provide keys.json file"
        return 1
    fi
    return 0
}

check_genesis() {
    return 0
}

generate_config() {
    # initial loading, backup old configuration file
    if [ -f "/etc/insolar/insolard.yaml" -a ! -f "/etc/insolar/insolard.yaml.old" ]; then
        cp /etc/insolar/insolard.yaml /etc/insolar/insolard.yaml.old
    elif [ ! -f "/etc/insolar/insolard.yaml" ]; then
        cp /opt/tmp/insolard.yaml.default /etc/insolar/insolard.yaml
        touch /etc/insolar/insolard.yaml.old
    fi
    # generate (or regenerate configuration with right content)
    $GENERATE_BIN /etc/insolar/insolard.yaml
}

construct_insgorund_cmd() {
	local INSOLARD_RPC_LISTEN="$IP:$INSOLARD_RPC_LISTEN_PORT"

	local insgorund_cmd="$INSGORUND_BIN --listen $INSGORUND_LISTEN"
	insgorund_cmd="$insgorund_cmd --rpc $INSOLARD_RPC_ENDPOINT"
	insgorund_cmd="$insgorund_cmd --log-level=$INSOLARD_LOG_LEVEL"
	if [ ! -z "${INSOLARD_LOG_TO_FILE}" ]; then
	    insgorund_cmd="$insgorund_cmd 2> /var/log/insgorund.log"
	fi
	echo $insgorund_cmd
}

generate_procfile() {
    echo "insolard:  $(construct_insolar_cmd)"    > /opt/Procfile
    echo "insgorund: $(construct_insgorund_cmd)" >> /opt/Procfile
}

construct_insolar_cmd() {
    local insolar_cmd="$INSOLARD_BIN --config /etc/insolar/insolard.yaml"
    if [ ! -z "${INSOLARD_JAEGER_ENDPOINT}" ]; then
        insolar_cmd="$insolar_cmd --trace"
    fi
    if [ ! -z "${INSOLARD_LOG_TO_FILE}" ]; then
        insolar_cmd="$insolar_cmd 2> /var/log/insolard.log"
    fi
    echo $insolar_cmd
}

check_cert
check_genesis

# extreme hack
if [ -z "${INSGORUND_ENDPOINT}" ]; then
    export INSGORUND_LISTEN=127.0.0.1:33302
    export INSOLARD_RPC_ENDPOINT=127.0.0.1:33301
    generate_config
    generate_procfile
    cd /opt && cat Procfile && $FOREGO_BIN start
else
    generate_config
    $(construct_insolar_cmd)
fi

exit 0
