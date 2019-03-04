### Configuration passing

One volume is expected: -> `/opt/config`

It should have at least 2 files (1 optional):
* `cert.json` - certificate - information about discovery nodes and node itself,
   mandatory
* `keys.json` - keys - private/public keys, mandatory
* `insolar.yaml` - configuration file for insolar. It'll be backuped at first
  run and all needed fields will be overwritten (look at genconfig.go for more
  information), optional. Default one would be created if it's not passed.
* `genesis` (???)

### Configuration variables

#### Insolard

* `INSOLARD_LOG_LEVEL` - logging level (default: info)
* `INSOLARD_TRANSPORT_LISTEN_PORT` - port to run transport on (default: 7900)
* `INSOLARD_TRANSPORT_FIXED_ADDRESS` - URI for insolard to pretend (should be passed
   explicitly)
* `INSOLARD_JAEGER_ENDPOINT` - if exists, run node with `--trace` support and
  send all data to Jaeger provided here (default: off, applied to `ve*` roles)
* `INSGORUND_ENDPOINT` - endpoint where insgorund listens to requests (should be
   passed explicitly)

#### Insgorund

* `INSOLARD_LOG_LEVEL` - logging level (default: info)
* `INSOLARD_RPC_ENDPOINT` - endpoint where insolard listens to requests (should
   be passed explicitly)

### Example usages in docker-compose

Look at `docker-compose` folder in project root.
