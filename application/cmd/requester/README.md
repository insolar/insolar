# API requester 

Insolar API requester is a simple CLI tool for sending [API requests](https://apidocs.insolar.io/platform/latest/#section/Request-Specification) to Insolar Platform.

The requester automates the following:

1. Gets a seed value from a node.
2. Extracts a public key from a provided ECDSA private key.
3. Forms a correct request from a provided JSON payload — replaces the seed and public key with the values acquired in previous steps.
4. Sends the formed request to a specified API endpoint.

## Usage

To use the requester, do the following:

1. Build it. In your `insolar/insolar/` directory, run:

   ```console
   make requester
   ```
   
   This builds a `.bin/requester` binary in the `insolar/insolar/` directory.

2. Generate a key pair. In the same directory, run:

   ```console
   ./bin/insolar gen-key-pair --target=user > /tmp/userkey
   ```
   
   This generates an ECDSA key pair and puts it in the specified file.
   
3. Copy-paste a JSON payload of any Platform request sample from the [API specification](https://apidocs.insolar.io/platform/latest/#operation/member-create) into a `payload.json` file.
   
   For example, the `member.create` payload sample:

   ```json
   {
     "jsonrpc": "2.0",
     "method": "contract.call",
     "id": 1,
     "params": {
       "seed": "<value>",
       "callSite": "member.create",
       "publicKey": "<value>"
     }
   }
   ```
   
4. To run the requester, specify:
 
   - Insolar Platform RPC endpoint as a URL parameter.
   - Path to the key pair as the `-k` option's value.
   - Path to `payload.json` as the `-r` option's value.
   
   For example:

   ```console
   ./bin/requester https://<endpoint>/api/rpc -k /tmp/userkey -r payload.json  
   ```

## Requester options

    Insolar API requester is a simple CLI tool for sending requests to Insolar Platform
    
    Usage:
      requester <insolar_endpoint> [flags]
    
    Examples:
    ./requester http://localhost:19101/api/rpc  -k /tmp/userkey  -r payload.json  -v
    
    Flags:
      -p, --autocompletekey     Extract a public key value from the specified key pair and replace the corresponding one in the request body with it (default true)
      -s, --autocompleteseed    Request a new seed value and replace the corresponding one in the request body with the new (default true)
      -h, --help                Help for requester
      -k, --memberkeys string   Path to a key pair
      -r, --request string      JSON request body or path to the file containing it
      -v, --verbose             Print request information
