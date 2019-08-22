#!/usr/bin/env bash
# Helper for running monitoring stack for launchnet
#
# 1. without arguments just (re)creates all containers:
#
# ./scripts/monitor.sh
#
# 2. with arguments just wraps docker-compose command:
#
# ./scripts/monitor.sh -h
#
# examples:
#
# ./scripts/monitor.sh logs -f jaeger
# ./scripts/monitor.sh restart jaeger

# strict mode
set -euo pipefail
IFS=$'\n\t'

# configurable vars
INSOLAR_ARTIFACTS_DIR=${INSOLAR_ARTIFACTS_DIR:-".artifacts"}/
LAUNCHNET_BASE_DIR=${LAUNCHNET_BASE_DIR:-"${INSOLAR_ARTIFACTS_DIR}launchnet"}/

PROMETHEUS_IN_CONFIG=${LAUNCHNET_BASE_DIR}prometheus.yaml
CONFIG_DIR=scripts/monitor/prometheus/

JAEGER_WAIT_ATTEMPTS=20

set +x

export PROMETHEUS_CONFIG_DIR=../${LAUNCHNET_BASE_DIR}

if [[ $# -lt 1 ]]; then
# 1) if started without params  pretend to be clever and do what is expected:
# * prepare configuration
# * shutdown ans start all monitoring services
# * wait until Jaeger starts (we want to be sure tracer works before start benchmark)
    echo "pre-gen launchnet configs (required for prometheus config generation)"
    ./scripts/insolard/launchnet.sh -C
    mkdir -p ${CONFIG_DIR}

    echo "start monitoring stack"
    cd scripts/

    docker-compose down
    docker-compose up -d
    docker-compose ps

    echo "wait for jaeger..."
    set +e
    for (( i=1; i<=JAEGER_WAIT_ATTEMPTS; i++ )) do
        curl -s -f 'http://localhost:16687/'
        if [[ $? -eq 0 ]]; then
            echo "jaeger works"
            break
        fi
        if [[ $i -eq ${JAEGER_WAIT_ATTEMPTS} ]]; then
            echo "jaeger wait is failed"
            break
        fi
        sleep 2
    done
    echo "# Grafana: http://localhost:3000 admin:pass"
    echo "# Jaeger: http://localhost:16686"
    echo "# Prometheus: http://localhost:9090/targets"
    echo "# Kibana: http://localhost:5601 (starts slowly, be patient)"
    echo ""
    echo "enable collection of jaeger traces in launchnet by env variable:"
    echo "    INSOLAR_TRACER_JAEGER_AGENTENDPOINT=http://localhost:6831"
    exit
fi

# 2) with arguments just work as thin docker-compose wrapper

cd scripts/
docker-compose $@
