#!/usr/bin/env bash

confs=$( dirname $0 )"/configs/generated_configs/"
api_ports=$( grep -A 1 "apirunner" $confs/*  | grep address  |  grep -o ":\d\+" | grep -o "\d\+" | tr '\n' ' ' )

for port in $api_ports
do
    echo "Port $port. Status: "
    data=$(curl --header "Content-Type:application/json"  --data '{ "jsonrpc": "2.0", "method": "status.Get","id": "" }'  "localhost:$port/api/rpc" 2>/dev/null |  python -m json.tool )
    echo "$data" | grep -i NetworkState
    echo "$data" | grep -i PulseNumber
    echo
done
