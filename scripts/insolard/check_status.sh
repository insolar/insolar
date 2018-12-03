


for port in 19191 19195 19199
do
    echo "Port $port. Status: "
    curl --header "Content-Type:application/json"  --data '{ "jsonrpc": "2.0", "method": "status.Get","id": "" }'  "localhost:$port/api/rpc" 2>/dev/null |  python -m json.tool | grep -i NetworkState
    echo
done
