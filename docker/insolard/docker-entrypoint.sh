#!/bin/bash
set -e

export IP=`awk 'END{print $1}' /etc/hosts`

INSOLARD_BIN=/opt/bin/insolard
INSGORUND_BIN=/opt/bin/insgorund

FOREGO_BIN=/opt/bin/forego
GENERATE_BIN=/opt/bin/genconfig

export INSOLARD_TRANSPORT_LISTEN_PORT=${INSOLARD_TRANSPORT_LISTEN_PORT:-13831}
export INSOLARD_METRICS_LISTEN_PORT=${INSOLARD_METRICS_LISTEN_PORT:-8001}
export INSOLARD_RPC_LISTEN_PORT=${INSOLARD_RPC_LISTEN_PORT:-33001}
export INSOLARD_API_LISTEN_PORT=${INSOLARD_API_LISTEN_PORT:-19101}
export INSOLARD_LOG_LEVEL=${INSOLARD_LOG_LEVEL:-info}
export INSGORUND_LISTEN_PORT=${INSGORUND_LISTEN_PORT:-33002}

export INSOLARD_ROLE=${INSOLARD_ROLE:-insolard+insgorund}
INSOLARD_ROLE=`echo "${INSOLARD_ROLE}" | awk '{print tolower($0)}'`

check_cert() {
    if   [ ! -f "/opt/config/cert.json" ]; then
        echo "Failed to find certificate, please provide cert.json file"
        return 1
    elif [ ! -f "/opt/config/keys.json" ]; then
        echo "Failed to find keys, please provide keys.json file"
        return 1
    fi
    return 0
}

generate_config() {
    # initial loading, backup old configuration file
    if [ -f "/opt/config/insolard.yaml" -a ! -f "/opt/config/insolard.yaml.old" ]; then
        cp /opt/config/insolard.yaml /opt/config/insolard.yaml.old
    elif [ ! -f "/opt/config/insolard.yaml" ]; then
        cp /opt/tmp/insolard.yaml.default /opt/config/insolard.yaml
        touch /opt/config/insolard.yaml.old
    fi
    # generate (or regenerate configuration with right content)
    $GENERATE_BIN /opt/config/insolard.yaml
}

construct_insolar_cmd() {
    local insolar_cmd="$INSOLARD_BIN --config /opt/config/insolard.yaml"
    if [ ! -z "${INSOLARD_JAEGER_ENDPOINT}" ]; then
        insolar_cmd="$insolar_cmd --trace"
    fi
    if [ ! -z "${INSOLARD_LOG_TO_FILE}"]; then
        insolar_cmd="$insolar_cmd 2> /var/log/insolard.log"
    fi
    echo $insolar_cmd
}

construct_insgorund_cmd() {
    local INSOLARD_RPC_LISTEN="$IP:$INSOLARD_RPC_LISTEN_PORT"

    local insgorund_cmd="$INSGORUND_BIN --listen $INSGORUND_LISTEN"
    insgorund_cmd="$insgorund_cmd --rpc $INSOLARD_RPC_ENDPOINT"
    insgorund_cmd="$insgorund_cmd --log-level=$INSOLARD_LOG_LEVEL"
    if [ ! -z "${INSOLARD_LOG_TO_FILE}"]; then
        insgorund_cmd="$insgorund_cmd 2> /opt/log/insgorund.log"
    fi
    echo $insgorund_cmd
}

generate_procfile() {
    echo "insolard:  $(construct_insolar_cmd)"    > /opt/Procfile
    echo "insgorund: $(construct_insgorund_cmd)" >> /opt/Procfile
}

if [ "$INSOLARD_ROLE" = "insolard+insgorund" ]; then
    export INSGORUND_LISTEN="127.0.0.1:${INSGORUND_LISTEN_PORT}"
    export INSOLARD_RPC_ENDPOINT="${IP}:${INSOLARD_RPC_LISTEN_PORT}"
    check_cert
    generate_config
    generate_procfile
    cat /opt/Procfile
    $FOREGO_BIN start
elif [ "${INSOLARD_ROLE}" = "insolard" ]; then
    check_cert
    generate_config
    $(construct_insolar_cmd)
elif [ "${INSOLARD_ROLE}" = "insgorund" ]; then
    if [ -z "${INSOLARD_RPC_ENDPOINT}" ]; then
        echo "Insgorund role expects INSOLARD_RPC_ENDPOINT variable"
        exit 1
    fi
    export INSGORUND_LISTEN="${IP}:${INSGORUND_LISTEN_PORT}"

    $(construct_insgorund_cmd)
else
    echo "Failed to understand container role: '${INSOLARD_ROLE}', expected insolard/insgorund/insolard+insgorund"
    exit 1
fi

exit 0
