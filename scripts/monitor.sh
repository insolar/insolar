#! /usr/bin/env bash
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

# configurable vars
INSOLAR_ARTIFACTS_DIR=${INSOLAR_ARTIFACTS_DIR:-".artifacts"}/
LAUNCHNET_BASE_DIR=${LAUNCHNET_BASE_DIR:-"${INSOLAR_ARTIFACTS_DIR}launchnet"}/

PROMETHEUS_IN_CONFIG=${LAUNCHNET_BASE_DIR}prometheus.yaml
CONFIG_DIR=scripts/monitor/prometheus/

set -e


if [ $# -lt 1 ]; then
    echo "pregen configs (required for prometheus config generation)"
    ./scripts/insolard/launchnet.sh -C
    mkdir -p ${CONFIG_DIR}
    cp ${PROMETHEUS_IN_CONFIG} ${CONFIG_DIR}

    echo "start monitoring stack"
    cd scripts/

    docker-compose down
    docker-compose up -d
    docker-compose ps

    echo "wait for jaeger..."
    set +e
    for i in {1..10}; do
        curl -s -f 'http://localhost:16687/'
        if [ $? -eq 0 ]; then
            break
        fi
        sleep 2
    done
    echo "# Grafana: http://localhost:3000 admin:pass"
    echo "# Jaeger: http://localhost:16686"
    echo "# Prometheus: http://localhost:9090/targets"
    echo "# Kibana: http://localhost:5601 (starts slowly, be patient)"
    exit
fi

cd scripts/
docker-compose $@
