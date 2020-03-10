#!/usr/bin/env bash
#
# Example of profiling 60 second profile on all insolard node
# (by default profiles 30 seconds):
#
# ./insolar-scripts/insolard/profile.sh [60]

set -e

# Changeable environment variables (parameters)
INSOLAR_ARTIFACTS_DIR=${INSOLAR_ARTIFACTS_DIR:-".artifacts"}/
LAUNCHNET_BASE_DIR=${LAUNCHNET_BASE_DIR:-"${INSOLAR_ARTIFACTS_DIR}launchnet"}/

PROF_TIME=${1:-"30"}
PROF_FILES_DIR=${LAUNCHNET_BASE_DIR}pprof/
WEB_PROFILE_PORT=8080

trap 'killall' INT TERM EXIT

killall() {
    trap '' INT TERM     # ignore INT and TERM while shutting down
    echo "Shutting down..."
    kill -TERM 0         # fixed order, send TERM not INT
    wait
    echo DONE
}

confs=${LAUNCHNET_BASE_DIR}"/*nodes/*/insolard.yaml"
prof_ports=$( grep listenaddress ${confs} |  grep -o ":\d\+" | grep -o "\d\+" | tr '\n' ' ' )
#echo "prof_ports: $prof_ports"

set -x
rm -rf ${PROF_FILES_DIR}
mkdir -p ${PROF_FILES_DIR}
{ set +x; } 2>/dev/null

echo "Fetching profile data from insolar nodes..."
for port in ${prof_ports}
do
    set -x
    curl "http://localhost:${port}/debug/pprof/profile?seconds=${PROF_TIME}" \
        --output ${PROF_FILES_DIR}prof_${port} \
        &> /dev/null &
    { set +x; } 2>/dev/null
done
wait

echo "Starting web servers with profile info..."

for port in ${prof_ports}
do
    echo "Start web profile server on localhost:${WEB_PROFILE_PORT}/ui"
    go tool pprof -http=:${WEB_PROFILE_PORT} ${PROF_FILES_DIR}prof_${port} &
    WEB_PROFILE_PORT=$((WEB_PROFILE_PORT + 1))
done
wait
