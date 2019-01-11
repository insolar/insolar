Connecting node certificate generation
===============

Usage
----------
#### Build

    make bin/certgen
   
#### Start insolard

    ./scripts/insolard/launchnet.sh -g
   
#### Start benchmark

    ./bin/certgen --root-conf=scripts/insolard/configs/root_member_keys.json

### Options

    -c, --cert-file string   The OUT file the node certificate (default "cert.json")
    
    -k, --keys-file string   The OUT file for public/private keys of the node (default "keys.json")
    
    -r, --role string        The role of the new node (default "virtual")
    
        --root-conf string   Config that contains public/private keys of root member
    
    -h, --url string         Insolar API URL (default "http://localhost:19191/api")
    
    -v, --verbose            Be verbose (default false)
