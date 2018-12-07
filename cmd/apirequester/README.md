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

    ./bin/apirequester -k=scripts/insolard/root_member_keys.json

### Options

        -k rootmemberkeys
                Path to file with RootMember keys.

        -u url
                API url for requests (default - http://localhost:19191/api).
