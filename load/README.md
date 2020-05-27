##### Exporter api load test

Install lib
```
go get github.com/skudasov/loadgen
```
Install cli
```
go install ~/go/pkg/mod/github.com/skudasov/loadgen\@v1.1.2/cmd/loadcli.go
```

Bootstrap grafana + grafite
```
docker run -d -p 8181:80 -p 8125:8125/udp -p 8126:8126 --publish=2003:2003 --name kamon-grafana-dashboard kamon/grafana_graphite
```

Create example local config in user home dir (~/generator_local.yaml):
```yaml
host:
  name: local_host
  network_iface: en0
generator:
  target: "127.0.0.1:5678"
  responseTimeoutSec: 20
  rampUpStrategy: linear
  verbose: true
execution_mode: parallel
grafana:
  url: http://0.0.0.0:8181
  login: "admin"
  password: "admin"
graphite:
  url: 0.0.0.0:2003
  flushDurationSec: 1
  loadGeneratorPrefix: exporter
checks:
  handle_threshold_percent: 1.20
root_package_name: insolar
load_scripts_dir: load
timezone: Europe/Moscow
logging: 
  level: info 
  encoding: console

```

Build & Run:
```go
loadcli b darwin
./load_suite -gen_config generator_local.yaml -config load/run_configs/get_pulses.yaml
./load_suite -gen_config generator_local.yaml -config load/run_configs/get_records.yaml
```
