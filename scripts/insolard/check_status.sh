#!/usr/bin/env bash
set -o pipefail

confs=$( dirname $0 )"/configs/generated_configs/discoverynodes"
api_ports=$( grep -A 1 "apirunner" $confs/insolar_*.yaml  | grep address  |  grep -o ":\d\+" | grep -o "\d\+" | tr '\n' ' ' )
grep_exit=$?
if [[ $grep_exit -ne 0 ]]; then
    exit 1
fi

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
