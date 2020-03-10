#! /usr/bin/env bash

# Starts prometheus server for launchnet metrics collection.
# uses port 9990 (default 9090 mysteriously clash with default insolard metrics port)
#
# with argsument c (continue) removes previous collected data:
# ./insolar-scripts/prom/run.sh c

set -e
if [[ "$1" != "c" ]]; then
    echo "Clean Prometheus data dir"
    rm -rf insolar-scripts/prom/data/
fi
echo "Prometheus data dir size before start:" && du -sh insolar-scripts/prom/data/ || true
prometheus --config.file ./insolar-scripts/prom/server.yml --storage.tsdb.path=./insolar-scripts/prom/data --web.listen-address="0.0.0.0:9990"
echo "Prometheus data dir size after start:" && du -sh insolar-scripts/prom/data/
