#! /usr/bin/env bash

# Starts prometheus server for launchnet metrics collection.
# uses port 9990 (default 9090 mysteriously clash with default insolard metrics port)
#
# with argsument c (continue) removes previous collected data:
# ./scripts/prom/run.sh c

set -e
if [[ "$1" != "c" ]]; then
    echo "Clean Prometheus data dir"
    rm -rf scripts/prom/data/
fi
echo "Prometheus data dir size before start:" && du -sh scripts/prom/data/ || true
prometheus --config.file ./scripts/prom/server.yml --storage.tsdb.path=./scripts/prom/data --web.listen-address="0.0.0.0:9990"
echo "Prometheus data dir size after start:" && du -sh scripts/prom/data/
