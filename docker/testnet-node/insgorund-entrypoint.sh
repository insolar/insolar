#!/bin/bash
set -e

export IP=`awk 'END{print $1}' /etc/hosts`

INSGORUND_BIN=/opt/bin/insgorund

export INSOLARD_LOG_LEVEL=${INSOLARD_LOG_LEVEL:-info}
export INSOLARD_RPC_ENDPOINT=${INSOLARD_RPC_ENDPOINT:-}

INSGORUND_LISTEN=$IP:7777
INSGORUND_METRICS_LISTEN=$IP:8001
INSGORUND_CACHE_DIR="/tmp/insgorund"

mkdir -p $INSGORUND_CACHE_DIR

construct_insgorund_cmd() {
    local INSOLARD_RPC_LISTEN="$IP:$INSOLARD_RPC_LISTEN_PORT"

    local insgorund_cmd="$INSGORUND_BIN --listen $INSGORUND_LISTEN"
    insgorund_cmd="$insgorund_cmd --rpc $INSOLARD_RPC_ENDPOINT"
    insgorund_cmd="$insgorund_cmd --log-level=$INSOLARD_LOG_LEVEL"
    insgorund_cmd="$insgorund_cmd --metrics=$INSGORUND_METRICS_LISTEN"
    insgorund_cmd="$insgorund_cmd --directory=$INSGORUND_CACHE_DIR"
    echo $insgorund_cmd
}

echo "Running: $(construct_insgorund_cmd)"
$(construct_insgorund_cmd)

exit 0
