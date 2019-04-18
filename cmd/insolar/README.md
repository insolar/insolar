# Insolar

## Usage

### Build

    make insolar

### Start insolard

    ./scripts/insolard/launchnet.sh -g

### Send request example (RegisterNode)

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

    ./bin/insolar send_request --config=./scripts/insolard/configs/root_member_keys.json --root-caller --params=params.json

Check available commands: `./bin/insolar -h`

Help on any command: `./bin/insolar help COMMAND`

## how to generate certificate and keys for node

    ./bin/insolar certgen --root-conf=scripts/insolard/configs/root_member_keys.json
