API requester
===============
   Makes several requests to API url and show results.

Usage
----------
#### Build

    make apirequester
   
#### Start insolard

    ./scripts/insolard/launchnet.sh -g
   
#### Start apirequester

    ./bin/apirequester -k=scripts/insolard/configs/

### Options

        -k path to members keys
                Path to dir with members keys.
        -a adminurl
                API url for requests (default - http://localhost:19001/admin-api/rpc).
        -p publicurl
                API url for requests (default - http://localhost:19101/api/rpc).
