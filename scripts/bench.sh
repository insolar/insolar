#! /usr/bin/env bash
# Usage:
# ./scripts/bench.sh -c 2 -r=40

# configurable vars
INSOLAR_ARTIFACTS_DIR=${INSOLAR_ARTIFACTS_DIR:-".artifacts"}/
LAUNCHNET_BASE_DIR=${LAUNCHNET_BASE_DIR:-"${INSOLAR_ARTIFACTS_DIR}launchnet"}/

echo "wait for launchnet..."
for i in {1..10}; do
    echo "Attempt $i"
    scripts/insolard/check_status.sh > /dev/null
    if [ $? -eq 0 ]; then
        break
    fi
    sleep 3
done

echo "Start bench"
sleep 1
while :; do
    set -x
    ./bin/benchmark -k=${LAUNCHNET_BASE_DIR}configs/root_member_keys.json $@
    { set +x; } 2>/dev/null
    sleep 5
done
