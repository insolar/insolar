#!/bin/bash
set -e

export IP=`awk 'END{print $1}' /etc/hosts`

INSOLARD_BIN=/opt/bin/insolard
INSGORUND_BIN=/opt/bin/insgorund

GENERATE_BIN=/opt/bin/genconfig
FOREGO_BIN=/opt/bin/forego

export INSOLARD_LOG_LEVEL=${INSOLARD_LOG_LEVEL:-info}
export INSOLARD_RPC_ENDPOINT=${INSOLARD_RPC_ENDPOINT:-}

INSGORUND_LISTEN=$IP:7777
INSGORUND_METRICS_LISTEN=$IP:8001
INSGORUND_CACHE_DIR="/opt/insgorund"

construct_insgorund_cmd() {
    local INSOLARD_RPC_LISTEN="$IP:$INSOLARD_RPC_LISTEN_PORT"

    local insgorund_cmd="$INSGORUND_BIN --listen $INSGORUND_LISTEN"
    insgorund_cmd="$insgorund_cmd --rpc $INSOLARD_RPC_ENDPOINT"
    insgorund_cmd="$insgorund_cmd --log-level=$INSOLARD_LOG_LEVEL"
    insgorund_cmd="$insgorund_cmd --metrics=$INSGORUND_METRICS_LISTEN"
    insgorund_cmd=""
    if [ ! -z "${INSOLARD_LOG_TO_FILE}" ]; then
        insgorund_cmd="$insgorund_cmd 2> /var/log/insgorund.log"
    fi
    echo $insgorund_cmd
}

$(construct_insgorund_cmd)

exit 0
