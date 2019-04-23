#! /usr/bin/env bash
set -o pipefail

# Changeable environment variables (parameters)
INSOLAR_ARTIFACTS_DIR=${INSOLAR_ARTIFACTS_DIR:-".artifacts"}/
LAUNCHNET_BASE_DIR=${LAUNCHNET_BASE_DIR:-"${INSOLAR_ARTIFACTS_DIR}launchnet"}/

# check is discovery nodes ready

confs=${LAUNCHNET_BASE_DIR}discoverynodes/
#(>&2 echo "grep -A 1 'apirunner' ${confs}*/insolard.yaml  | grep address  |  grep -o ':\d\+' | grep -o '\d\+' | tr '\n' ' '")
api_ports=$( grep -A 1 "apirunner" ${confs}*/insolard.yaml  | grep address  |  grep -o ":\d\+" | grep -o "\d\+" | tr '\n' ' ' )
grep_exit=$?
if [[ $grep_exit -ne 0 ]]; then
    exit 1
fi
#(>&2 echo "api_ports=$api_ports")

exit_code=0
for port in $api_ports
do
    echo "Port $port. Status: "
    status=$(curl --header "Content-Type:application/json"  --data '{ "jsonrpc": "2.0", "method": "status.Get","id": "" }'  "localhost:$port/api/rpc" 2>/dev/null |  python -m json.tool | grep -i NodeState)
    echo "    $status"
    echo
    if [[ $status != *"NodeReady"* ]]; then
        exit_code=1
    fi
done
exit $exit_code
