# The requester 

The requester is a simple CLI for sending requests to Insolar Platform. 

## Usage

### Build

    make requester

### Options
    The requester is a simple CLI for sending requests to Insolar Platform
    
    Usage:
      requester <insolar url> [flags]
    
    Examples:
    ./requester http://localhost:19101/api/rpc  -k /tmp/userkey  -r params.json  -v
    
    Flags:
      -p, --autocompletekey     Should replace publicKey to correct value (default true)
      -s, --autocompleteseed    Should replace seed to correct value (default true)
      -h, --help                help for requester
      -k, --memberkeys string   Path to member key
      -r, --request string      The request body or path to request params file
      -v, --verbose             Print request information



### how to generate keypair 

    ./bin/insolar gen-key-pair --target=user > /tmp/userkey

### CreateMember Example
  ```
    params.json    
    {
      "jsonrpc": "2.0",
      "method": "contract.call",
      "id": 1,
      "params": {
        "seed": "fhDEwRRbSnYnbMnALKMh8gXdzaSvRv/nfsGC9S7kqik=",
        "callSite": "member.create",
        "publicKey": "-----BEGIN PUBLIC KEY-----\\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAEMSbA4KvO/jlwY+8WFDEdwhCLlsEC\\nF3/GYvu9iTWHwCctx1wTbGGjNLY03EjXyYxaf8coNbSbZeu+jXcWeMHG0A==\\n-----END PUBLIC KEY-----"
      }
    }
```      

`./bin/requester -k=/tmp/userkey -r params.json  http://localhost:19101/api/rpc` <br>
or <br>
```./bin/requester http://localhost:19101/api/rpc -k=/tmp/userkey -r="`cat params.json`" ```
   
