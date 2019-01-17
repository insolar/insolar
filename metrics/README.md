Metrics component
-----------------

We use [Prometheus](https://prometheus.io) monitoring system and time series database for collecting and store metrics.

Package **metrics** is based on [Prometheus golang client](https://github.com/prometheus/client_golang).
It contains metrics collectors of entire project. Component starts http server on `http://0.0.0.0:8080/metrics` by default(can be changed in configuration)

If you want to add metrics in your component code, 
you need to describe it as global collector variable in this package. 
Each global collector must be registered in constructor `NewMetrics()`

When you creating collector, you need to fill `Opts` structure. 
You should to [read this guide](https://prometheus.io/docs/practices/naming/) before choosing `Opts.Name`

#### Collector types

 - [Counter](https://godoc.org/github.com/prometheus/client_golang/prometheus#Counter)
 - [Gauge](https://godoc.org/github.com/prometheus/client_golang/prometheus#Gauge)
 - [Summary](https://godoc.org/github.com/prometheus/client_golang/prometheus#Summary) 
 - [Histogram](https://godoc.org/github.com/prometheus/client_golang/prometheus#Histogram)
 
 
#### Labels

Labels is used to create query with filters. For example, when we count total number of sent packets, then
we can make query with filter by packet type in report. Generally, You don't need to use `Opts.ConstLabels`. 
This field is used for labels, which not be changed in runtime(e. g. for app version).

You should create a `metricVec` using specific method for particular collector type.

```go
// NetworkPacketSentTotal is total number of sent packets metric
var NetworkPacketSentTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name:      "packet_sent_total",
	Help:      "Total number of sent packets",
	Namespace: insolarNamespace,
	Subsystem: "network",
}, []string{"packetType"})
```


#### Using collectors in your code

Collectors are thread safe, you can manipulate with it from any goroutine.

```go
// labeled counter usage example
metrics.NetworkPacketSentTotal.WithLabelValues(packet.Type.String()).Inc()
```