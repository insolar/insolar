# Insolar

## Usage

### Build

    make insolar

### Start insolard

    ./insolar-scripts/insolard/launchnet.sh -g

### Send request example (RegisterNode)

You should have ```params.json``` with something like this:

    {
      "callParams": [
       "<public_key from node config>",
        <numberOfBootstrapNodes>,
        <majorityRule>,
        [<roles>],
        "<ip>"
      ],
      "callSite": "RegisterNode"
    }

Than use send_request command with this file:

    ./bin/insolar send-request --root-keys=.artifacts/launchnet/configs/ --root-caller --params=params.json
    ./bin/insolar send-request --root-keys=.artifacts/launchnet/configs/ --migration-admin-caller --params=params.json

Check available commands: `./bin/insolar -h`

Help on any command: `./bin/insolar help COMMAND`

## how to generate certificate and keys for node

    ./bin/insolar certgen --root-keys=scripts/insolard/configs/root_member_keys.json

### Options

        -a admin-url
                API url for requests (default - http://localhost:19001/admin-api/rpc).
        -u url
                API url for requests (default - http://localhost:19101/api/rpc).



### CreateMember Example
        
        params.json
        
        {
          "callSite": "member.migrationCreate",
          "callParams": {},
          "publicKey": "-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE3xcoC2lnprrVuc83K6b2R1gvA5kB\nfEUI7xBi1GX/LWtDzex5s47oEXPlXhysnrLOKL75kP8/5hRc3QJm12KuRw==\n-----END PUBLIC KEY-----\n"
        ​
        }

    ./bin/insolar send-request --root-keys=./scripts/insolard/configs/new_member_keys.json --params=params.json -u=http://localhost:19101/api/rpc -a=http://localhost:19001/admin-api/rpc

    for localhost:
    ./bin/insolar send-request --root-keys=./scripts/insolard/configs/new_member_keys.json --params=params.json 

### Migration example

        params.json
        
        {
          "callSite":"deposit.migration",
          "callParams": {
            "amount": "1000000000",
            "ethTxHash": "394578234932493486739856jfgd48756348563495846djf",
            "migrationAddress": "0x83274348763847632487326482346328462384632486234"
          }, 
          "reference":"1734543583274348763847632487326482346328462384632",
          "publicKey":"-----BEGIN PUBLIC KEY-----\nMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE3xcoC2lnprrVuc83K6b2R1gvA5kB\nfEUI7xBi1GX/LWtDzex5s47oEXPlXhysnrLOKL75kP8/5hRc3QJm12KuRw==\n-----END PUBLIC KEY-----\n"
        }
        ​
    ./bin/insolar send-request --root-keys=./scripts/insolard/configs/migration_daemon_keys.json --params=params.json -u=http://localhost:19001/admin-api/rpc -a=http://localhost:19001/admin-api/rpc

    for localhost:
    ./bin/insolar send-request --root-keys=./scripts/insolard/configs/migration_daemon_keys.json --params=params.json -u=http://localhost:19001/admin-api/rpc

## How to get count of free migration addresses in every shard

    ./bin/insolar free-migration-count --migration-admin-keys=.artifacts/launchnet/configs/ --alert-level=100 --shards-count=10

### Options

        -k migration-admin-keys
                Path to dir with config that contains public/private keys of migration admin.
        -l alert-level
                If one of shard have less free addresses than this value, command will print alert message.
        -s shards-count
                Count of shards at platform (must be a multiple of ten).

## How to add migration addresses from files to every shard

    ./bin/insolar add-migration-addresses --migration-admin-keys=.artifacts/launchnet/configs/ --shards-count=100 --addresses=../../migrationAddressGenerator/bin/addresses.json

### Options

        -k migration-admin-keys
                Dir with config that contains public/private keys of admin member.
        -g addresses
                Path to files with addresses. We expect files will be match generator utility output (from insolar/migrationAddressGenerator).
        -s shards-count
                Count of shards at platform (must be a multiple of ten).
