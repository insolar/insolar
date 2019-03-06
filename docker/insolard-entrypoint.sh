#!/bin/bash
set -e

export IP=`awk 'END{print $1}' /etc/hosts`

INSOLARD_BIN=/opt/bin/insolard
GENERATE_BIN=/opt/bin/genconfig

export INSOLARD_TRANSPORT_LISTEN_PORT=${INSOLARD_TRANSPORT_LISTEN_PORT:-7900}
export INSOLARD_LOG_LEVEL=${INSOLARD_LOG_LEVEL:-info}

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

construct_insolar_cmd() {
    local insolar_cmd="$INSOLARD_BIN --config /etc/insolar/insolard.yaml"
    if [ ! -z "${INSOLARD_JAEGER_ENDPOINT}" ]; then
        insolar_cmd="$insolar_cmd --trace"
    fi
    echo $insolar_cmd
}

check_cert
check_genesis
generate_config
$(construct_insolar_cmd)

exit 0
