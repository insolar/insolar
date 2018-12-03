


for port in 53835 53837 53839 
do
    echo "Port $port. Status: "
    curl --header "Content-Type:application/json"  --data '{ "jsonrpc": "2.0", "method": "info.Get","id": "" }'  "localhost:$port/api/rpc" 2>/dev/null | grep -i status
done
