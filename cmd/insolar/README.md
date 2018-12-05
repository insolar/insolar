Insolar
===============

Usage
----------
#### Build

    make insolar
   
#### Start insolard

    ./scripts/insolard/launchnet.sh -g
   
#### Send request example (RegisterNode)

You should have ```params.json``` with something like this:

    {
      "params": [
       "<public_key from node config>",
        <numberOfBootstrapNodes>,
        <majorityRule>,
        [<roles>],
        "<ip>"
      ],
      "method": "RegisterNode"
    }
Than use send_request command with this file:

    ./bin/insolar -c=send_request --config=./scripts/insolard/configs/root_member_keys.json --root_as_caller --params=params.json

### Options

        -c cmd
                Command. Available commands: default_config | random_ref | version | gen_keys | gen_certificate | send_request | gen_send_configs. 

        -v verbose
                Be verbose (default false).

        -o output
            Path to output file (use - for STDOUT).

        -u url
            API url (default http://localhost:19191/api).

        -g config
                Path to file with caller config or caller+params config.

        -p params
                Path to params file (default params.json).

        -r root_as_caller
                Do request from RootMember (default false).
