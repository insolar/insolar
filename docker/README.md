### Configuration passing

One volume is expected: -> `/opt/config`

It should have at least 2 files (1 optional):
* `cert.json` - certificate - information about discovery nodes and node itself,
   mandatory
* `keys.json` - keys - private/public keys, mandatory
* `insolar.yaml` - configuration file for insolar. It'll be backuped at first
  run and all needed fields would be overwrited (look at genconfig.go for more
  information), optional. Default one would be created if it's not passed.

### Configuration variables

* `INSOLARD_ROLE` - container role, possible values: "`ve`", "`insgorund`",
  "`ve+insgorund`" (default: "`ve+insgorund`", common)
* `INSOLARD_LOG_LEVEL` - logging level (default: info, common)
* `INSOLARD_TRANSPORT_LISTEN_PORT` - port to run transport on (default: 13831,
   applied to `ve*` roles)
* `INSOLARD_METRICS_LISTEN_PORT` - insolard port to run metrics listener on
   (default: 8001, applied to `ve*` roles)
* `INSGORUND_METRICS_LISTEN_PORT` - insgorund port to run metrics listener on
   (default: 8002, applied to `*insgorund` roles)
* `INSOLARD_RPC_LISTEN_PORT` - port to run rpc on (default 33001, applied to
   `ve*` roles)
* `INSGORUND_RPC_ENDPORINT` - URI for insgorund to connect to (should be passed
   explicitly, only `insgorund` role)
* `INSOLARD_JAEGER_ENDPOINT` - if exists, run node with `--trace` support and
  send all data to Jaeger provided here (default: off, applied to `ve*` roles)
* `INSGORUND_ENDPOINT` - endpoint where insgorund listens to requests (should be
   passed explicitly, only `ve` role)


### Example usages in docker-compose

In case you want to run insolard and insgorund in one container:
```yaml
insolard_insgorund:
  links:
    - "launchnet"
  environment:
    - INSOLARD_LOG_LEVEL=debug
    - INSOLARD_LOG_TO_FILE=1
  volumes:
    - type: bind
      source: ./config-01
      target: /opt/config
  build:
    context: insolard
    dockerfile: Dockerfile
  restart: always
```

In case you want to run insolard and insgorund in different containers
```yaml
insolard:
  links:
    - "launchnet"
  environment:
    - INSOLARD_ROLE=insolard
    - INSOLARD_LOG_LEVEL=debug
    - INSOLARD_LOG_TO_FILE=1
    - INSGORUND_ENDPOINT=insgorund:33002
  volumes:
    - type: bind
      source: ./config-02
      target: /opt/config
  build:
    context: insolard
    dockerfile: Dockerfile
  restart: always
insgorund:
  links:
    - "launchnet"
  environment:
    - INSOLARD_ROLE=insgorund
    - INSOLARD_LOG_LEVEL=debug
    - INSOLARD_LOG_TO_FILE=1
    - INSOLARD_RPC_ENDPOINT=insolard:33001
  build:
    context: insolard
    dockerfile: Dockerfile
  restart: always
```
