# The requester 

The utility helps send swagger-like requests to Insolar platform. 

## Usage

### Build

    make requester

### Options
    ./requester --help
    
    Usage of ./requester:
      -p, --autocompletekey     should replace publicKey to correct value
      -s, --autocompleteseed    should replace seed to correct value
      -k, --memberkeys string   path to member key paramsFile
      -r, --request string      the request body or path to request paramsFile
      -u, --url string          API URL. for example http://localhost:19101/api/rpc


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

`./bin/requester -k=/tmp/userkey -u http://localhost:19101/api/rpc -r params.json -p true -s true` <br>
or <br>
```./bin/requester -k=/tmp/userkey -u http://localhost:19101/api/rpc -r="`cat params.json`" -p true -s true```
   
