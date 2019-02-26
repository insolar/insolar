#!/bin/bash
set -e

export IP=`awk 'END{print $1}' /etc/hosts`

INSOLARD_BIN=/opt/bin/insolard
INSGORUND_BIN=/opt/bin/insgorund

GENERATE_BIN=/opt/bin/genconfig

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

construct_insolar_cmd() {
    local insolar_cmd="$INSOLARD_BIN --config /etc/insolar/insolard.yaml"
    if [ ! -z "${INSOLARD_JAEGER_ENDPOINT}" ]; then
        insolar_cmd="$insolar_cmd --trace"
    fi
    if [ ! -z "${INSOLARD_LOG_TO_FILE}"]; then
        insolar_cmd="$insolar_cmd 2> /var/log/insolard.log"
    fi
    echo $insolar_cmd
}

check_cert
generate_config
$(construct_insolar_cmd)

exit 0
